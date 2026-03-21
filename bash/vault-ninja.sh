#!/usr/bin/env bash
# Vault Ninja SDK — Bash
# source this file, then call vn_* functions with your API key.
#
# Usage:
#   source vault-ninja.sh
#   secrets=$(vn_list_secrets "$VN_API_KEY")
#   field=$(vn_get_field "$VN_API_KEY" "$SECRET_ID" "$FIELD_ID" | jq -r '.value')
#
# Environment:
#   VAULT_API_URL  Override the API base URL (default: https://api.vaultninja.org/api/sdk/v1)
#   VN_API_KEY     Conventionally used to hold your API key

VAULT_API_URL="${VAULT_API_URL:-https://api.vaultninja.org/api/sdk/v1}"

# List all secrets in the org (no fields or files).
# Usage: vn_list_secrets <api_key>
vn_list_secrets() {
  local key="$1"
  curl -fsSL \
    -H "Authorization: Bearer ${key}" \
    "${VAULT_API_URL}/secrets"
}

# Get a single secret with all decrypted fields and file metadata.
# Usage: vn_get_secret <api_key> <secret_id>
vn_get_secret() {
  local key="$1" id="$2"
  curl -fsSL \
    -H "Authorization: Bearer ${key}" \
    "${VAULT_API_URL}/secrets/${id}"
}

# Get a single decrypted field value.
# Usage: vn_get_field <api_key> <secret_id> <field_id>
vn_get_field() {
  local key="$1" id="$2" fid="$3"
  curl -fsSL \
    -H "Authorization: Bearer ${key}" \
    "${VAULT_API_URL}/secrets/${id}/fields/${fid}"
}

# Download a file's raw bytes.
# Usage: vn_get_file <api_key> <secret_id> <file_id>
vn_get_file() {
  local key="$1" id="$2" fid="$3"
  curl -fsSL \
    -H "Authorization: Bearer ${key}" \
    "${VAULT_API_URL}/secrets/${id}/files/${fid}"
}
