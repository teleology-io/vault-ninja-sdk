"use strict";
/**
 * Vault Ninja SDK — slim read-only client for org secrets.
 * No runtime dependencies. Uses native fetch.
 *
 * Usage:
 *   import { vn } from '@teleology/vn';
 *   const client = new vn('vn_org_...');
 *   const secrets = await client.list_secrets();
 *   const secret  = await client.get_secret('<id>');
 *   const field   = await client.get_field('<id>', '<fid>');
 *   const file    = await client.get_file('<id>', '<fid>');
 *
 * Env vars:
 *   VN_API_KEY — API key (used when api_key is omitted)
 *   VN_API_URL — override base URL
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.vn = exports.VaultNinjaError = void 0;
const DEFAULT_URL = 'https://api.vaultninja.org/api/sdk/v1';
class VaultNinjaError extends Error {
    constructor(status, message) {
        super(`Vault Ninja API error ${status}: ${message}`);
        this.status = status;
        this.name = 'VaultNinjaError';
    }
}
exports.VaultNinjaError = VaultNinjaError;
const env = (key) => typeof process !== 'undefined' ? process.env[key] : undefined;
class vn {
    constructor(api_key, base_url) {
        this.apiKey = api_key ?? env('VN_API_KEY') ?? '';
        this.baseUrl = (base_url ?? env('VN_API_URL') ?? DEFAULT_URL).replace(/\/$/, '');
        if (!this.apiKey)
            throw new Error('api_key is required (or set VN_API_KEY)');
    }
    async _get(path) {
        const res = await fetch(this.baseUrl + path, {
            headers: { Authorization: `Bearer ${this.apiKey}` },
        });
        if (!res.ok) {
            const body = await res.json().catch(() => ({ error: res.statusText }));
            throw new VaultNinjaError(res.status, body.error ?? res.statusText);
        }
        return res;
    }
    async list_secrets() {
        return (await this._get('/secrets')).json();
    }
    async get_secret(id) {
        return (await this._get(`/secrets/${id}`)).json();
    }
    async get_field(id, fid) {
        return (await this._get(`/secrets/${id}/fields/${fid}`)).json();
    }
    async get_file(id, fid) {
        return (await this._get(`/secrets/${id}/files/${fid}`)).arrayBuffer();
    }
}
exports.vn = vn;
