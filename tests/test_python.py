#!/usr/bin/env python3
# Run from the python/ directory: python3 ../tests/test_python.py
"""
Vault Ninja SDK — Python integration test harness.

Usage (from the python/ directory, after pip install -e .):
    python ../tests/test_python.py --api-key vn_org_... --endpoint http://...
    VN_API_KEY=vn_org_... VAULT_API_URL=http://... python ../tests/test_python.py
"""
import argparse
import os
import sys


def main():
    parser = argparse.ArgumentParser(description="Vault Ninja Python SDK harness")
    parser.add_argument("--api-key", default=os.environ.get("VN_API_KEY", ""))
    parser.add_argument("--endpoint", default=os.environ.get("VAULT_API_URL", ""))
    args = parser.parse_args()

    if not args.api_key:
        print("Error: --api-key or VN_API_KEY required", file=sys.stderr)
        sys.exit(1)

    from vault_ninja.client import VaultNinjaClient

    passed = 0
    failed = 0

    def ok(label):
        nonlocal passed
        print(f"  \u2713 {label}")
        passed += 1

    def bad(label):
        nonlocal failed
        print(f"  \u2717 {label}")
        failed += 1

    print("[python SDK]")

    kwargs = {"api_key": args.api_key}
    if args.endpoint:
        kwargs["base_url"] = args.endpoint

    with VaultNinjaClient(**kwargs) as client:
        # 1. list_secrets
        try:
            secrets = client.secrets.list_secrets()
        except Exception as e:
            bad(f"list_secrets — {e}")
            _results(passed, failed)
            sys.exit(1)

        if not secrets:
            bad("list_secrets — returned 0 secrets (need at least 1 to continue)")
            _results(passed, failed)
            sys.exit(1)

        ok(f"list_secrets — found {len(secrets)} secret(s)")
        secret_id = secrets[0].id

        # 2. get_secret
        try:
            secret = client.secrets.get_secret(secret_id)
            ok(f'get_secret({secret_id}) — "{secret.title}"')
        except Exception as e:
            bad(f"get_secret({secret_id}) — {e}")
            _results(passed, failed)
            sys.exit(1)

        # 3. get_field
        fields = secret.fields or []
        if not fields:
            print("  - get_field — skipped (no fields on secret)")
        else:
            field_id = fields[0].id
            try:
                field = client.secrets.get_field(secret_id, field_id)
                ok(f'get_field({secret_id}, {field_id}) — label: "{field.label}"')
            except Exception as e:
                bad(f"get_field({secret_id}, {field_id}) — {e}")

        # 4. get_file
        files = secret.files or []
        if not files:
            print("  - get_file — skipped (no files on secret)")
        else:
            file_id = files[0].id
            size_bytes = files[0].size_bytes
            try:
                client.secrets.get_file(secret_id, file_id)
                ok(f"get_file({secret_id}, {file_id}) — {size_bytes} bytes")
            except Exception as e:
                bad(f"get_file({secret_id}, {file_id}) — {e}")

    _results(passed, failed)
    if failed:
        sys.exit(1)


def _results(passed, failed):
    print(f"\nResults: {passed} passed, {failed} failed")


if __name__ == "__main__":
    main()
