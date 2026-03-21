#!/usr/bin/env bash
# Bash SDK integration test harness.
# Usage:
#   VN_API_KEY=vn_org_... VAULT_API_URL=http://... bash tests/test_bash.sh
# Run from the repo root.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=../bash/vault-ninja.sh
source "${SCRIPT_DIR}/../bash/vault-ninja.sh"

KEY="${VN_API_KEY:-}"
if [[ -z "$KEY" ]]; then
  echo "Error: VN_API_KEY is required" >&2
  exit 1
fi

PASS=0
FAIL=0

pass() { echo "  ✓ $1"; ((PASS++)) || true; }
fail() { echo "  ✗ $1"; ((FAIL++)) || true; }

echo "[bash SDK]"

# 1. listSecrets
LIST=$(vn_list_secrets "$KEY" 2>&1) || { fail "listSecrets — curl error: $LIST"; exit 1; }
COUNT=$(echo "$LIST" | jq 'length' 2>/dev/null) || { fail "listSecrets — invalid JSON: $LIST"; exit 1; }
SECRET_ID=$(echo "$LIST" | jq -r '.[0].id // empty')
if [[ -z "$SECRET_ID" ]]; then
  fail "listSecrets — returned 0 secrets (need at least 1 to continue)"
  echo ""
  echo "Results: ${PASS} passed, ${FAIL} failed"
  exit 1
fi
pass "listSecrets — found ${COUNT} secret(s)"

# 2. getSecret
SECRET=$(vn_get_secret "$KEY" "$SECRET_ID" 2>&1) || { fail "getSecret(${SECRET_ID}) — curl error: $SECRET"; exit 1; }
SECRET_TITLE=$(echo "$SECRET" | jq -r '.title // empty') || { fail "getSecret(${SECRET_ID}) — invalid JSON"; exit 1; }
pass "getSecret(${SECRET_ID}) — \"${SECRET_TITLE}\""

# 3. getField (use first field)
FIELD_ID=$(echo "$SECRET" | jq -r '.fields[0].id // empty')
if [[ -z "$FIELD_ID" ]]; then
  echo "  - getField — skipped (no fields on secret)"
else
  FIELD=$(vn_get_field "$KEY" "$SECRET_ID" "$FIELD_ID" 2>&1) || { fail "getField(${SECRET_ID}, ${FIELD_ID}) — curl error: $FIELD"; }
  FIELD_LABEL=$(echo "$FIELD" | jq -r '.label // empty') || { fail "getField — invalid JSON"; }
  if [[ -n "$FIELD_LABEL" ]]; then
    pass "getField(${SECRET_ID}, ${FIELD_ID}) — label: \"${FIELD_LABEL}\""
  else
    fail "getField(${SECRET_ID}, ${FIELD_ID}) — missing label in response"
  fi
fi

# 4. getFile (use first file)
FILE_ID=$(echo "$SECRET" | jq -r '.files[0].id // empty')
if [[ -z "$FILE_ID" ]]; then
  echo "  - getFile — skipped (no files on secret)"
else
  FILE_SIZE=$(echo "$SECRET" | jq -r '.files[0].size_bytes // "?"')
  FILE_BYTES=$(vn_get_file "$KEY" "$SECRET_ID" "$FILE_ID" 2>&1) || { fail "getFile(${SECRET_ID}, ${FILE_ID}) — curl error"; }
  BYTES_RECEIVED=${#FILE_BYTES}
  pass "getFile(${SECRET_ID}, ${FILE_ID}) — ${FILE_SIZE} bytes (server), ${BYTES_RECEIVED} bytes received"
fi

echo ""
echo "Results: ${PASS} passed, ${FAIL} failed"
[[ $FAIL -eq 0 ]]
