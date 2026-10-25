# vn — Vault Ninja CLI & MCP server (zero-knowledge)

`vn` reads your organization's secrets and decrypts them **on your machine**.
It calls the zero-knowledge SDK API (`/api/sdk/v2`), which only ever returns
ciphertext. Your API key's private half never leaves the machine; vn sends
only the part before the `.`. Every read shows up in the org audit log as
`secret.read` with `via: cli` or `via: mcp`.

Needs a **`vn_org_2_…` key** (Settings → API Keys → Create). Older keys
keep working with the language SDKs, but vn can't use them.

## Install

```bash
brew install teleology-io/tap/vn                                     # macOS / Linux
go install github.com/teleology-io/vault-ninja-sdk/go/cmd/vn@latest  # Go 1.25+
```

Or download a binary from the [releases](https://github.com/teleology-io/vault-ninja-sdk/releases) (`vn/v…` tags).

## Credentials

First match wins:

1. `VN_TOKEN` (or `VN_API_KEY`, the variable the SDKs use)
2. `--token-file <path>`
3. OS keychain, via `vn login` (paste the key; it's checked, then stored)

`VN_API_URL` overrides the server. A host or any `/api/sdk/…` URL works.

## Commands

```bash
vn ls                                   # secrets this key can read
vn get "Prod DB" url                    # one field (label, id or type)
vn get "Prod DB" "Login#password"       # part of a login / card field
vn get "Prod DB" --json                 # whole secret
vn file get "TLS" cert.pem -o cert.pem  # decrypt a file (written 0600)
vn env --format dotenv                  # variables from .vaultninja.yaml → stdout
vn run -- npm test                      # run with secrets injected as env vars
vn run --env STRIPE_KEY="Stripe:secret key" -- ./deploy.sh
```

`vn run` never writes secrets to disk. The child gets only the variables you
mapped, not vn's own key, and its exit code is passed through.

### Project file

`.vaultninja.yaml` in the project (or any parent directory) holds references,
never values, so it's safe to commit:

```yaml
env:
  DATABASE_URL: { secret: "Prod DB", field: url }
  DB_PASSWORD:  { secret: "Prod DB", field: "Login#password" }
mcp:
  allow_reveal: ["Staging DB"]       # what get_secret_value may return (also needs --allow-reveal)
  allowed_commands: [npm, go, make]  # optional: limit what run_with_secrets may run
```

## CI/CD

```yaml
# GitHub Actions
- run: go install github.com/teleology-io/vault-ninja-sdk/go/cmd/vn@latest
- run: vn run -- ./deploy.sh
  env:
    VN_TOKEN: ${{ secrets.VN_TOKEN }}
```

Use a **secret-scoped** key with an **expiry** for each pipeline.

## AI agents (MCP)

```bash
claude mcp add vault-ninja -- vn mcp                 # Claude Code
```

```json
{ "mcpServers": { "vault-ninja": { "command": "vn", "args": ["mcp"] } } }
```

| Tool | What it does |
|---|---|
| `list_secrets` | ids, titles, tags. No values. |
| `describe_secret` | description, field labels/types, file names. No values. |
| `run_with_secrets` | runs a command (no shell) with secrets as env vars; returns exit code and output with every injected value redacted, including base64/hex/URL-encoded forms |
| `get_secret_value` | **off by default**. Only exists with `vn mcp --allow-reveal`, and only for secrets in `mcp.allow_reveal` |

**What this protects against, and what it doesn't.** Values stay out of the
model's context and transcripts, and redaction catches accidental leaks such
as a test printing its config. But `run_with_secrets` can run commands, so an
agent that's determined (or prompt-injected) could still move a value out,
for example by sending it somewhere over the network. The real controls are:

- a **narrowly scoped, expiring** key, one per agent
- `mcp.allowed_commands`
- the audit log (`secret.read via=mcp`)
- revoking the key when you're done

Secret titles, descriptions and command output are treated as untrusted text
(control characters stripped, length capped).
