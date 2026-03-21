/**
 * Vault Ninja SDK client wrapper.
 * Adds VAULT_API_URL env var support and a clean constructor on top of the
 * generated SecretsApi classes.
 */

import { Configuration, SecretsApi } from './index';
import type { Secret, SecretField, SecretSummary } from './index';

const DEFAULT_BASE_URL = 'https://api.vaultninja.org/api/sdk/v1';

export class VaultNinjaClient {
  private api: SecretsApi;

  /**
   * @param apiKey   Your org API key (starts with `vn_org_`).
   * @param baseUrl  Override the API base URL. Falls back to the
   *                 `VAULT_API_URL` environment variable, then production.
   */
  constructor(apiKey: string, baseUrl?: string) {
    const basePath =
      baseUrl ??
      (typeof process !== 'undefined' ? process.env.VAULT_API_URL : undefined) ??
      DEFAULT_BASE_URL;

    const config = new Configuration({ basePath, accessToken: apiKey });
    this.api = new SecretsApi(config);
  }

  listSecrets(): Promise<SecretSummary[]> {
    return this.api.listSecrets();
  }

  getSecret(id: string): Promise<Secret> {
    return this.api.getSecret({ id });
  }

  getField(id: string, fid: string): Promise<SecretField> {
    return this.api.getField({ id, fid });
  }

  getFile(id: string, fid: string): Promise<Response> {
    return this.api.getFileRaw({ id, fid }).then((r) => r.raw);
  }
}
