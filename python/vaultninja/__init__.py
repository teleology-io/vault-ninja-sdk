"""Vault Ninja SDK — slim read-only client for org secrets.

stdlib only, no external dependencies.

Usage::

    from vaultninja import vn

    client = vn(api_key="vn_org_...")
    secrets = client.list_secrets()
    secret  = client.get_secret(id="<secret_id>")
    for field in secret.fields:
        print(field.label, "=", field.value)
    field = client.get_field(id="<secret_id>", fid="<field_id>")
    data  = client.get_file(id="<secret_id>", fid="<file_id>")

Env vars:
    VN_API_KEY — API key (used when api_key arg is omitted)
    VN_API_URL — override base URL
"""
from __future__ import annotations

import json
import os
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from typing import List

_DEFAULT_URL = "https://api.vaultninja.org/api/sdk/v1"


class VaultNinjaError(Exception):
    """Raised for non-2xx API responses."""

    def __init__(self, status: int, message: str) -> None:
        super().__init__(f"Vault Ninja API error {status}: {message}")
        self.status = status
        self.message = message


@dataclass
class SecretField:
    id: str
    secret_id: str
    position: int
    field_type: str
    label: str
    value: str
    created_at: str


@dataclass
class SecretFile:
    id: str
    secret_id: str
    original_name: str
    content_type: str
    size_bytes: int
    position: int
    created_at: str
    url: str


@dataclass
class Secret:
    id: str
    title: str
    description: str
    tags: List[str]
    org_id: str
    created_by: str
    created_at: str
    updated_at: str
    fields: List[SecretField] = field(default_factory=list)
    files: List[SecretFile] = field(default_factory=list)


def _parse_field(d: dict) -> SecretField:
    return SecretField(
        id=d["id"],
        secret_id=d["secret_id"],
        position=d["position"],
        field_type=d["field_type"],
        label=d["label"],
        value=d["value"],
        created_at=d["created_at"],
    )


def _parse_file(d: dict) -> SecretFile:
    return SecretFile(
        id=d["id"],
        secret_id=d["secret_id"],
        original_name=d["original_name"],
        content_type=d["content_type"],
        size_bytes=d["size_bytes"],
        position=d["position"],
        created_at=d["created_at"],
        url=d["url"],
    )


def _parse_secret(d: dict) -> Secret:
    return Secret(
        id=d["id"],
        title=d["title"],
        description=d["description"],
        tags=d.get("tags") or [],
        org_id=d["org_id"],
        created_by=d["created_by"],
        created_at=d["created_at"],
        updated_at=d["updated_at"],
        fields=[_parse_field(f) for f in d.get("fields") or []],
        files=[_parse_file(f) for f in d.get("files") or []],
    )


class vn:
    """Vault Ninja SDK client.

    Args:
        api_key:  Your org API key (starts with ``vn_org_``).
                  Falls back to the ``VN_API_KEY`` environment variable.
        base_url: Override the API base URL.
                  Falls back to ``VN_API_URL``, then the production URL.
    """

    def __init__(self, api_key: str = "", base_url: str = "") -> None:
        self._api_key = api_key or os.environ.get("VN_API_KEY", "")
        self._base_url = (base_url or os.environ.get("VN_API_URL", "") or _DEFAULT_URL).rstrip("/")
        if not self._api_key:
            raise ValueError("api_key is required (or set VN_API_KEY)")

    def _get(self, path: str) -> bytes:
        req = urllib.request.Request(
            self._base_url + path,
            headers={"Authorization": f"Bearer {self._api_key}"},
        )
        try:
            with urllib.request.urlopen(req) as resp:
                return resp.read()
        except urllib.error.HTTPError as exc:
            body = exc.read()
            try:
                msg = json.loads(body).get("error", body.decode())
            except Exception:
                msg = body.decode(errors="replace")
            raise VaultNinjaError(exc.code, msg) from None

    def list_secrets(self) -> List[Secret]:
        return [_parse_secret(d) for d in json.loads(self._get("/secrets"))]

    def get_secret(self, id: str) -> Secret:
        return _parse_secret(json.loads(self._get(f"/secrets/{id}")))

    def get_field(self, id: str, fid: str) -> SecretField:
        return _parse_field(json.loads(self._get(f"/secrets/{id}/fields/{fid}")))

    def get_file(self, id: str, fid: str) -> bytes:
        return self._get(f"/secrets/{id}/files/{fid}")
