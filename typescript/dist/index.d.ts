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
export type FieldType = 'PASSWORD' | 'CARD' | 'TOTP' | 'URL' | 'NOTE' | 'TEXT';
export interface SecretField {
    id: string;
    secret_id: string;
    position: number;
    field_type: FieldType;
    label: string;
    value: string;
    created_at: string;
}
export interface SecretFile {
    id: string;
    secret_id: string;
    original_name: string;
    content_type: string;
    size_bytes: number;
    position: number;
    created_at: string;
    url: string;
}
export interface Secret {
    id: string;
    title: string;
    description: string;
    tags: string[];
    org_id: string;
    created_by: string;
    created_at: string;
    updated_at: string;
    fields: SecretField[];
    files: SecretFile[];
}
export declare class VaultNinjaError extends Error {
    readonly status: number;
    constructor(status: number, message: string);
}
export declare class vn {
    private readonly baseUrl;
    private readonly apiKey;
    constructor(api_key?: string, base_url?: string);
    private _get;
    list_secrets(): Promise<Secret[]>;
    get_secret(id: string): Promise<Secret>;
    get_field(id: string, fid: string): Promise<SecretField>;
    get_file(id: string, fid: string): Promise<ArrayBuffer>;
}
