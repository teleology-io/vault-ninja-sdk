"use strict";
/**
 * Vault Ninja SDK client wrapper.
 * Adds VAULT_API_URL env var support and a clean constructor on top of the
 * generated SecretsApi classes.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.VaultNinjaClient = void 0;
const index_1 = require("./index");
const DEFAULT_BASE_URL = 'https://api.vaultninja.org/api/sdk/v1';
class VaultNinjaClient {
    /**
     * @param apiKey   Your org API key (starts with `vn_org_`).
     * @param baseUrl  Override the API base URL. Falls back to the
     *                 `VAULT_API_URL` environment variable, then production.
     */
    constructor(apiKey, baseUrl) {
        const basePath = baseUrl ??
            (typeof process !== 'undefined' ? process.env.VAULT_API_URL : undefined) ??
            DEFAULT_BASE_URL;
        const config = new index_1.Configuration({ basePath, accessToken: apiKey });
        this.api = new index_1.SecretsApi(config);
    }
    listSecrets() {
        return this.api.listSecrets();
    }
    getSecret(id) {
        return this.api.getSecret({ id });
    }
    getField(id, fid) {
        return this.api.getField({ id, fid });
    }
    getFile(id, fid) {
        return this.api.getFileRaw({ id, fid }).then((r) => r.raw);
    }
}
exports.VaultNinjaClient = VaultNinjaClient;
