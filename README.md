# Vault Ninja SDK

Read-only programmatic access to your [Vault Ninja](https://vaultninja.org) organization secrets. Use this SDK in CI/CD pipelines, deployment scripts, and automation tools to pull credentials at runtime — no hardcoding required.

**OpenAPI spec:** [`openapi.yaml`](./openapi.yaml)

---

## Authentication

All requests require an `Authorization` header with a Bearer API key:

```
Authorization: Bearer vn_org_<your_key>
```

Create keys in your org settings → **API Keys**.

---

## Base URL

The default base URL is `https://api.vaultninja.org/api/sdk/v1`.

Override it by setting the `VAULT_API_URL` environment variable or passing it directly to the client constructor.

---

## Python

**Install:**
```bash
pip install "https://github.com/teleology-io/vault-ninja-sdk/releases/download/v1.0.0/vault_ninja-1.0.0.tar.gz"
```

**Usage:**
```python
from vault_ninja.client import VaultNinjaClient

client = VaultNinjaClient(api_key="vn_org_...")

# List all secrets
secrets = client.secrets.list_secrets()

# Get a specific secret with decrypted fields
secret = client.secrets.get_secret(id="<secret-id>")
for field in secret.fields:
    print(field.label, "=", field.value)

# Inject a single credential
field = client.secrets.get_field(id="<secret-id>", fid="<field-id>")
print(field.value)
```

---

## TypeScript / Node.js

**Install:**
```bash
npm install "https://github.com/teleology-io/vault-ninja-sdk/releases/download/v1.0.0/teleology-vn-1.0.0.tgz"
```

**Usage:**
```typescript
import { VaultNinjaClient } from '@teleology/vn';

const client = new VaultNinjaClient('vn_org_...');

const secrets = await client.listSecrets();
const secret  = await client.getSecret('<secret-id>');
const field   = await client.getField('<secret-id>', '<field-id>');
const file    = await client.getFile('<secret-id>', '<file-id>');

console.log(field.value);
```

---

## Go

**Install:**
```bash
go get github.com/teleology-io/vault-ninja-sdk/go@v1.0.0
```

**Usage:**
```go
import "github.com/teleology-io/vault-ninja-sdk/go/vaultninja"

client := vaultninja.New("vn_org_...")

secrets, err := client.ListSecrets(ctx)
secret,  err := client.GetSecret(ctx, "<secret-id>")
field,   err := client.GetField(ctx, "<secret-id>", "<field-id>")

fmt.Println(field.Value)
```

---

## Bash

**Source inline (zero install):**
```bash
source <(curl -fsSL https://raw.githubusercontent.com/teleology-io/vault-ninja-sdk/master/bash/vault-ninja.sh)

# List secrets
vn_list_secrets "$VN_API_KEY"

# Get a specific secret
vn_get_secret "$VN_API_KEY" "<secret-id>"

# Inject a field value (requires jq)
PASSWORD=$(vn_get_field "$VN_API_KEY" "<secret-id>" "<field-id>" | jq -r '.value')

# Download a file
vn_get_file "$VN_API_KEY" "<secret-id>" "<file-id>" > cert.pem
```


**GitHub Actions example:**
```yaml
- name: Fetch credentials from Vault Ninja
  run: |
    source <(curl -fsSL https://raw.githubusercontent.com/teleology-io/vault-ninja-sdk/master/bash/vault-ninja.sh)
    DB_PASS=$(vn_get_field "$VN_API_KEY" "$SECRET_ID" "$FIELD_ID" | jq -r '.value')
    echo "::add-mask::$DB_PASS"
    echo "DB_PASSWORD=$DB_PASS" >> "$GITHUB_ENV"
  env:
    VN_API_KEY: ${{ secrets.VN_API_KEY }}
```

---

## Endpoints

| Method | Path | Description |
|---|---|---|
| `GET` | `/secrets` | List all org secrets (no fields/files) |
| `GET` | `/secrets/{id}` | Get secret with decrypted fields + file metadata |
| `GET` | `/secrets/{id}/fields/{fid}` | Get a single decrypted field |
| `GET` | `/secrets/{id}/files/{fid}` | Download a file |

Full schema: see [`openapi.yaml`](./openapi.yaml)

---

## Error Responses

All errors return `{ "error": "message" }` with an appropriate HTTP status:

| Status | Meaning |
|---|---|
| `401` | Missing, invalid, or revoked API key |
| `404` | Secret/field/file not found or doesn't belong to this org |
| `429` | Rate limit exceeded — use exponential backoff |
| `500` | Server error |
