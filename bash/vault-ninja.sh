#!/usr/bin/env bash
# Vault Ninja SDK — Bash
# source this file, then call: vn <subcommand> [args]
#
# Usage:
#   source vault-ninja.sh
#   export VN_API_KEY="vn_org_..."
#   vn list
#   vn secret <secret_id>
#   vn field  <secret_id> <field_id>
#   vn file   <secret_id> <file_id>
#
# Environment:
#   VN_API_KEY  API key
#   VN_API_URL  Override base URL (default: https://api.vaultninja.org/api/sdk/v1)

VN_API_URL="${VN_API_URL:-https://api.vaultninja.org/api/sdk/v1}"

vn() {
  local key="${VN_API_KEY}"
  if [[ -z "$key" ]]; then
    echo "vn: VN_API_KEY is not set" >&2
    return 1
  fi

  local cmd="$1"; shift
  case "$cmd" in
    list)
      curl -fsSL \
        -H "Authorization: Bearer ${key}" \
        "${VN_API_URL}/secrets"
      ;;
    secret)
      local id="$1"
      curl -fsSL \
        -H "Authorization: Bearer ${key}" \
        "${VN_API_URL}/secrets/${id}"
      ;;
    field)
      local id="$1" fid="$2"
      curl -fsSL \
        -H "Authorization: Bearer ${key}" \
        "${VN_API_URL}/secrets/${id}/fields/${fid}"
      ;;
    file)
      local id="$1" fid="$2"
      curl -fsSL \
        -H "Authorization: Bearer ${key}" \
        "${VN_API_URL}/secrets/${id}/files/${fid}"
      ;;
    *)
      echo "vn: unknown command '${cmd}'" >&2
      echo "Usage: vn list | secret <id> | field <id> <fid> | file <id> <fid>" >&2
      return 1
      ;;
  esac
}
