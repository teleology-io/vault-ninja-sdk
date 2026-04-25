#!/usr/bin/env bash
# Vault Ninja SDK — run all SDK test harnesses.
#
# Usage:
#   VN_API_KEY=vn_org_... VN_API_URL=http://localhost:8080/api/sdk/v1 ./tests/run_all.sh

set -uo pipefail

VN_API_KEY="${VN_API_KEY:-}"
VN_API_URL="${VN_API_URL:-}"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TESTS_DIR="${REPO_ROOT}/tests"

if [[ -z "${VN_API_KEY}" ]]; then
  echo "Error: VN_API_KEY is required" >&2
  exit 1
fi

OVERALL_PASS=0
OVERALL_FAIL=0

run_sdk() {
  local name="$1"; shift
  if "$@"; then
    ((OVERALL_PASS++)) || true
  else
    ((OVERALL_FAIL++)) || true
  fi
  echo ""
}

# ── Bash ──────────────────────────────────────────────────────────────────────
run_sdk "bash" bash "${TESTS_DIR}/test_bash.sh"

# ── Go ────────────────────────────────────────────────────────────────────────
run_sdk "go" bash -c "
  cd '${REPO_ROOT}/go'
  go run ./cmd/harness -api-key='${VN_API_KEY}' -endpoint='${VN_API_URL}'
"

# ── Python ────────────────────────────────────────────────────────────────────
run_sdk "python" python3 "${TESTS_DIR}/test_python.py"

# ── TypeScript ────────────────────────────────────────────────────────────────
run_sdk "typescript" bash -c "
  cd '${REPO_ROOT}/typescript' && npm run build --silent
  node '${TESTS_DIR}/test_typescript.js'
"

# ── Summary ───────────────────────────────────────────────────────────────────
echo "══════════════════════════════════════"
echo "SDKs passed: ${OVERALL_PASS}  failed: ${OVERALL_FAIL}"
echo "══════════════════════════════════════"

[[ $OVERALL_FAIL -eq 0 ]]
