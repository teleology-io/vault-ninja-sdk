/**
 * Vault Ninja SDK client wrapper.
 * Adds VAULT_API_URL env var support and a clean constructor on top of the
 * generated SecretsApi classes.
 */
import type { Secret, SecretField, SecretSummary } from './index';
export declare class VaultNinjaClient {
    private api;
    /**
     * @param apiKey   Your org API key (starts with `vn_org_`).
     * @param baseUrl  Override the API base URL. Falls back to the
     *                 `VAULT_API_URL` environment variable, then production.
     */
    constructor(apiKey: string, baseUrl?: string);
    listSecrets(): Promise<SecretSummary[]>;
    getSecret(id: string): Promise<Secret>;
    getField(id: string, fid: string): Promise<SecretField>;
    getFile(id: string, fid: string): Promise<Response>;
}
