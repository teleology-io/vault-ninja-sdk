#!/usr/bin/env node
'use strict';
// Vault Ninja SDK — TypeScript integration test.
// Run from repo root after building: node tests/test_typescript.js
//   VN_API_KEY=vn_org_... VN_API_URL=http://... node tests/test_typescript.js

const path = require('path');
const { vn, VaultNinjaError } = require(path.join(__dirname, '../typescript/dist/index.js'));

let passed = 0;
let failed = 0;
const ok  = (l) => { console.log(`  ✓ ${l}`); passed++; };
const bad = (l) => { console.log(`  ✗ ${l}`); failed++; };

async function main() {
  const apiKey  = process.env.VN_API_KEY  ?? '';
  const baseUrl = process.env.VN_API_URL  ?? '';

  if (!apiKey) {
    console.error('Error: VN_API_KEY required');
    process.exit(1);
  }

  console.log('[typescript SDK]');

  const client = new vn(apiKey, baseUrl || undefined);

  // 1. list_secrets
  let secrets;
  try {
    secrets = await client.list_secrets();
  } catch (e) {
    bad(`list_secrets — ${e}`);
    results(); process.exit(1);
  }
  if (!secrets.length) {
    bad('list_secrets — returned 0 secrets (need at least 1 to continue)');
    results(); process.exit(1);
  }
  ok(`list_secrets — found ${secrets.length} secret(s)`);

  const secretId = secrets[0].id;

  // 2. get_secret
  let secret;
  try {
    secret = await client.get_secret(secretId);
    ok(`get_secret(${secretId}) — "${secret.title}"`);
  } catch (e) {
    bad(`get_secret(${secretId}) — ${e}`);
    results(); process.exit(1);
  }

  // 3. get_field
  const fields = secret.fields ?? [];
  if (!fields.length) {
    console.log('  - get_field — skipped (no fields on secret)');
  } else {
    const fieldId = fields[0].id;
    try {
      const field = await client.get_field(secretId, fieldId);
      ok(`get_field(${secretId}, ${fieldId}) — label: "${field.label}"`);
    } catch (e) {
      bad(`get_field(${secretId}, ${fieldId}) — ${e}`);
    }
  }

  // 4. get_file
  const files = secret.files ?? [];
  if (!files.length) {
    console.log('  - get_file — skipped (no files on secret)');
  } else {
    const fileId = files[0].id;
    const sizeBytes = files[0].size_bytes;
    try {
      await client.get_file(secretId, fileId);
      ok(`get_file(${secretId}, ${fileId}) — ${sizeBytes} bytes`);
    } catch (e) {
      bad(`get_file(${secretId}, ${fileId}) — ${e}`);
    }
  }

  results();
  if (failed) process.exit(1);
}

function results() {
  console.log(`\nResults: ${passed} passed, ${failed} failed`);
}

main();
