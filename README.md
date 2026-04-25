# Vault Ninja SDK

Read-only programmatic access to your [Vault Ninja](https://vaultninja.org) organization secrets. Use this SDK in CI/CD pipelines, deployment scripts, and automation tools to pull credentials at runtime — no hardcoding required.

---

## Authentication

All requests require an `Authorization` header with a Bearer API key:

```
Authorization: Bearer vn_org_<your_key>
```

Create keys in your org settings → **API Keys**.

---

## Environment Variables

| Variable | Description |
|---|---|
| `VN_API_KEY` | Your org API key |
| `VN_API_URL` | Override the base URL (default: `https://api.vaultninja.org/api/sdk/v1`) |

---

## Python

**Install:**
```bash
pip install "https://github.com/teleology-io/vault-ninja-sdk/releases/latest/download/teleology-vn-python.tgz"
```

**Usage:**
```python
from vaultninja import vn

client = vn(api_key="vn_org_...")

secrets = client.list_secrets()

secret = client.get_secret(id="<secret-id>")
for field in secret.fields:
    print(field.label, "=", field.value)

field = client.get_field(id="<secret-id>", fid="<field-id>")
print(field.value)

data = client.get_file(id="<secret-id>", fid="<file-id>")
```

Zero dependencies — stdlib only.

---

## TypeScript

**Install:**
```bash
npm install "https://github.com/teleology-io/vault-ninja-sdk/releases/latest/download/teleology-vn-node.tgz"
```

**Usage:**
```typescript
import { vn } from '@teleology/vn';

const client = new vn('vn_org_...');

const secrets = await client.list_secrets();
const secret  = await client.get_secret('<secret-id>');
const field   = await client.get_field('<secret-id>', '<field-id>');
const file    = await client.get_file('<secret-id>', '<file-id>');

console.log(field.value);
```

No runtime dependencies — native `fetch`, full TypeScript types.

---

## Go

**Install:**
```bash
go get github.com/teleology-io/vault-ninja-sdk/go@latest
```

**Usage:**
```go
import "github.com/teleology-io/vault-ninja-sdk/go/vaultninja"

client := vaultninja.New("vn_org_...")

secrets, err := client.ListSecrets(ctx)

secret, err := client.GetSecret(ctx, "<secret-id>")
for _, f := range secret.Fields {
    fmt.Println(f.Label, "=", f.Value)
}

field, err := client.GetField(ctx, "<secret-id>", "<field-id>")
fmt.Println(field.Value)

data, err := client.GetFile(ctx, "<secret-id>", "<file-id>")
```

Zero dependencies — stdlib only.

---

## Bash

**Source inline (zero install):**
```bash
source <(curl -fsSL https://raw.githubusercontent.com/teleology-io/vault-ninja-sdk/master/bash/vault-ninja.sh)
export VN_API_KEY="vn_org_..."

# List secrets
vn list

# Get a specific secret
vn secret "<secret-id>"

# Inject a field value (requires jq)
PASSWORD=$(vn field "<secret-id>" "<field-id>" | jq -r '.value')

# Download a file
vn file "<secret-id>" "<file-id>" > cert.pem
```

**GitHub Actions example:**
```yaml
- name: Fetch credentials from Vault Ninja
  run: |
    source <(curl -fsSL https://raw.githubusercontent.com/teleology-io/vault-ninja-sdk/master/bash/vault-ninja.sh)
    DB_PASS=$(vn field "$SECRET_ID" "$FIELD_ID" | jq -r '.value')
    echo "::add-mask::$DB_PASS"
    echo "DB_PASSWORD=$DB_PASS" >> "$GITHUB_ENV"
  env:
    VN_API_KEY: ${{ secrets.VN_API_KEY }}
```

---

## Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/secrets` | List all org secrets |
| `GET` | `/secrets/{id}` | Get secret with decrypted fields + file metadata |
| `GET` | `/secrets/{id}/fields/{fid}` | Get a single decrypted field |
| `GET` | `/secrets/{id}/files/{fid}` | Download a file |

---

## Error Responses

All errors return `{ "error": "message" }` with an appropriate HTTP status:

| Status | Meaning |
|---|---|
| `401` | Missing, invalid, or revoked API key |
| `404` | Secret/field/file not found or doesn't belong to this org |
| `429` | Rate limit exceeded — use exponential backoff |
| `500` | Server error |
