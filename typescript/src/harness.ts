/**
 * Vault Ninja SDK — TypeScript integration test harness.
 *
 * Usage (from the typescript/ directory):
 *   npm run harness -- --api-key=vn_org_... --endpoint=http://...
 *   VN_API_KEY=vn_org_... VAULT_API_URL=http://... npm run harness
 */

import { VaultNinjaClient } from './client';

let passed = 0;
let failed = 0;

function ok(label: string) {
  console.log(`  ✓ ${label}`);
  passed++;
}

function bad(label: string) {
  console.log(`  ✗ ${label}`);
  failed++;
}

function results() {
  console.log(`\nResults: ${passed} passed, ${failed} failed`);
}

function parseArgs(): { apiKey: string; endpoint: string } {
  const args = process.argv.slice(2);
  let apiKey = process.env.VN_API_KEY ?? '';
  let endpoint = process.env.VAULT_API_URL ?? '';

  for (const arg of args) {
    if (arg.startsWith('--api-key=')) apiKey = arg.slice('--api-key='.length);
    else if (arg.startsWith('--endpoint=')) endpoint = arg.slice('--endpoint='.length);
  }

  return { apiKey, endpoint };
}

async function main() {
  const { apiKey, endpoint } = parseArgs();

  if (!apiKey) {
    console.error('Error: --api-key or VN_API_KEY required');
    process.exit(1);
  }

  console.log('[typescript SDK]');

  const client = new VaultNinjaClient(apiKey, endpoint || undefined);

  // 1. listSecrets
  let secrets;
  try {
    secrets = await client.listSecrets();
  } catch (e) {
    bad(`listSecrets — ${e}`);
    results();
    process.exit(1);
  }

  if (secrets.length === 0) {
    bad('listSecrets — returned 0 secrets (need at least 1 to continue)');
    results();
    process.exit(1);
  }
  ok(`listSecrets — found ${secrets.length} secret(s)`);

  const secretId = secrets[0].id;

  // 2. getSecret
  let secret;
  try {
    secret = await client.getSecret(secretId);
    ok(`getSecret(${secretId}) — "${secret.title}"`);
  } catch (e) {
    bad(`getSecret(${secretId}) — ${e}`);
    results();
    process.exit(1);
  }

  // 3. getField
  const fields = secret.fields ?? [];
  if (fields.length === 0) {
    console.log('  - getField — skipped (no fields on secret)');
  } else {
    const fieldId = fields[0].id;
    try {
      const field = await client.getField(secretId, fieldId);
      ok(`getField(${secretId}, ${fieldId}) — label: "${field.label}"`);
    } catch (e) {
      bad(`getField(${secretId}, ${fieldId}) — ${e}`);
    }
  }

  // 4. getFile
  const files = secret.files ?? [];
  if (files.length === 0) {
    console.log('  - getFile — skipped (no files on secret)');
  } else {
    const fileId = files[0].id;
    const sizeBytes = files[0].sizeBytes;
    try {
      await client.getFile(secretId, fileId);
      ok(`getFile(${secretId}, ${fileId}) — ${sizeBytes} bytes`);
    } catch (e) {
      bad(`getFile(${secretId}, ${fileId}) — ${e}`);
    }
  }

  results();
  if (failed > 0) process.exit(1);
}

main();
