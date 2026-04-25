#!/usr/bin/env python3
"""Vault Ninja SDK — Python integration test.

Usage:
    VN_API_KEY=vn_org_... VN_API_URL=http://... python3 tests/test_python.py
"""
import os
import sys


def main():
    api_key  = os.environ.get("VN_API_KEY", "")
    base_url = os.environ.get("VN_API_URL", "")

    if not api_key:
        print("Error: VN_API_KEY required", file=sys.stderr)
        sys.exit(1)

    sys.path.insert(0, os.path.join(os.path.dirname(__file__), "../python"))
    from vaultninja import vn

    passed = 0
    failed = 0

    def ok(label):
        nonlocal passed
        print(f"  ✓ {label}")
        passed += 1

    def bad(label):
        nonlocal failed
        print(f"  ✗ {label}")
        failed += 1

    def results():
        print(f"\nResults: {passed} passed, {failed} failed")

    print("[python SDK]")

    kwargs = {"api_key": api_key}
    if base_url:
        kwargs["base_url"] = base_url
    client = vn(**kwargs)

    # 1. list_secrets
    try:
        secrets = client.list_secrets()
    except Exception as e:
        bad(f"list_secrets — {e}")
        results()
        sys.exit(1)

    if not secrets:
        bad("list_secrets — returned 0 secrets (need at least 1 to continue)")
        results()
        sys.exit(1)

    ok(f"list_secrets — found {len(secrets)} secret(s)")
    secret_id = secrets[0].id

    # 2. get_secret
    try:
        secret = client.get_secret(id=secret_id)
        ok(f'get_secret({secret_id}) — "{secret.title}"')
    except Exception as e:
        bad(f"get_secret({secret_id}) — {e}")
        results()
        sys.exit(1)

    # 3. get_field
    fields = secret.fields or []
    if not fields:
        print("  - get_field — skipped (no fields on secret)")
    else:
        field_id = fields[0].id
        try:
            field = client.get_field(id=secret_id, fid=field_id)
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
            client.get_file(id=secret_id, fid=file_id)
            ok(f"get_file({secret_id}, {file_id}) — {size_bytes} bytes")
        except Exception as e:
            bad(f"get_file({secret_id}, {file_id}) — {e}")

    results()
    if failed:
        sys.exit(1)


if __name__ == "__main__":
    main()
