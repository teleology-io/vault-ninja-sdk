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

# --fail-with-body (used below to surface API error bodies, e.g. the 409 an
# end-to-end encrypted secret returns when this key has no usable data key)
# needs curl >= 7.76.0 (Mar 2021). Checked once per shell, cached in
# _VN_CURL_OK so `vn` calls after the first don't re-spawn `curl --version`.
_vn_check_curl() {
  if [[ -z "${_VN_CURL_OK:-}" ]]; then
    local ver major minor
    ver=$(curl --version 2>/dev/null | head -1 | awk '{print $2}')
    IFS='.' read -r major minor _ <<< "$ver"
    if [[ -n "$major" && -n "$minor" ]] && { (( major > 7 )) || { (( major == 7 )) && (( minor >= 76 )); }; }; then
      _VN_CURL_OK=1
    else
      _VN_CURL_OK=0
      echo "vn: curl ${ver:-<unknown>} is too old (need >= 7.76.0) — --fail-with-body isn't supported, so API error details (e.g. a 409 on an end-to-end encrypted secret) can't be shown. Upgrade curl to use this SDK." >&2
    fi
  fi
  [[ "$_VN_CURL_OK" == 1 ]]
}

# GETs a path, streaming the body straight to stdout (must stay a stream, not
# a captured shell variable — `file` downloads raw binary bytes, which a
# bash string can't hold safely). --fail-with-body behaves like -f (nonzero
# exit on non-2xx) but still writes the body instead of discarding it — the
# plain `-f` this used to use swallowed the server's actual explanation on a
# 409 (e.g. an end-to-end encrypted secret with no usable key on this token).
_vn_get() {
  curl -sSL --fail-with-body -H "Authorization: Bearer ${2}" "${VN_API_URL}${1}"
}

vn() {
  _vn_check_curl || return 1

  local key="${VN_API_KEY}"
  if [[ -z "$key" ]]; then
    echo "vn: VN_API_KEY is not set" >&2
    return 1
  fi

  local cmd="$1"; shift
  case "$cmd" in
    list)
      _vn_get "/secrets" "$key"
      ;;
    secret)
      local id="$1"
      _vn_get "/secrets/${id}" "$key"
      ;;
    field)
      local id="$1" fid="$2"
      _vn_get "/secrets/${id}/fields/${fid}" "$key"
      ;;
    file)
      local id="$1" fid="$2"
      _vn_get "/secrets/${id}/files/${fid}" "$key"
      ;;
    *)
      echo "vn: unknown command '${cmd}'" >&2
      echo "Usage: vn list | secret <id> | field <id> <fid> | file <id> <fid>" >&2
      return 1
      ;;
  esac
}
