import { appVersion } from './app-version.mjs';
import { createHash } from 'node:crypto';
import { readFile, readdir, writeFile } from 'node:fs/promises';
import { join } from 'node:path';

const output = new URL('../build/client/', import.meta.url);
const source = new URL('../../', import.meta.url);
async function files(root, prefix = '') {
  const result = [];
  for (const entry of await readdir(new URL(prefix, root), { withFileTypes: true })) {
    const path = join(prefix, entry.name);
    if (entry.isDirectory()) result.push(...await files(root, path + '/'));
    else result.push(path);
  }
  return result.sort();
}
const assets = (await files(output)).filter(path =>
  path.startsWith('assets/') || path.startsWith('fonts/') ||
  ['__spa-fallback.html', 'manifest.webmanifest', 'icon-192.png', 'icon-96.png', 'reminder.png', 'favicon.ico'].includes(path));
if (!assets.includes('__spa-fallback.html')) throw new Error('Missing offline SPA shell');
const shellHash = createHash('sha256');
for (const path of assets) shellHash.update(path).update(await readFile(new URL(path, output)));
const textHash = createHash('sha256');
for (const dir of ['quran/', 'hadith/', 'names/']) {
  for (const path of await files(new URL(dir, source))) {
    if (path.endsWith('.json') || path.endsWith('.go')) textHash.update(dir + path).update(await readFile(new URL(dir + path, source)));
  }
}
const worker = await readFile(new URL('../public/reminder.js', import.meta.url), 'utf8');
const build = { appVersion, version: shellHash.update(worker).digest('hex').slice(0, 16), textVersion: textHash.digest('hex').slice(0, 16), assets: assets.map(path => '/' + path) };
await writeFile(new URL('reminder.js', output), `self.__REMINDER_BUILD__ = ${JSON.stringify(build)};\n${worker}`);
console.log(`Offline shell: ${assets.length} files; text revision ${build.textVersion}`);
