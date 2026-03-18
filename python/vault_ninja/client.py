"""
Vault Ninja SDK client wrapper.
Adds VAULT_API_URL env var support and a clean constructor on top of the generated API classes.
"""
import os

from vault_ninja.api.secrets_api import SecretsApi
from vault_ninja.api_client import ApiClient
from vault_ninja.configuration import Configuration

_DEFAULT_BASE_URL = "https://vaultninja.org/api/sdk/v1"


class VaultNinjaClient:
    """
    Thin wrapper around the generated SecretsApi that handles base URL
    resolution and bearer token injection.

    Args:
        api_key:  Your org API key (starts with ``vn_org_``).
        base_url: Override the API base URL. Falls back to the
                  ``VAULT_API_URL`` environment variable, then the
                  production URL.

    Example::

        from vault_ninja.client import VaultNinjaClient

        client = VaultNinjaClient(api_key="vn_org_...")
        secrets = client.secrets.list_secrets()
    """

    def __init__(self, api_key: str, base_url: str | None = None) -> None:
        resolved = base_url or os.environ.get("VAULT_API_URL") or _DEFAULT_BASE_URL
        config = Configuration(host=resolved, access_token=api_key)
        self._api_client = ApiClient(configuration=config)
        self.secrets = SecretsApi(self._api_client)

    def close(self) -> None:
        self._api_client.close()

    def __enter__(self) -> "VaultNinjaClient":
        return self

    def __exit__(self, *_) -> None:
        self.close()
