import { createHash } from 'node:crypto';
import { readFileSync, readdirSync } from 'node:fs';
import { join } from 'node:path';
import { fileURLToPath } from 'node:url';
const root = fileURLToPath(new URL('../', import.meta.url));
const hash = createHash('sha256');
function add(path) {
  for (const entry of readdirSync(join(root, path), { withFileTypes: true }).sort((a, b) => a.name.localeCompare(b.name))) {
    const child = join(path, entry.name);
    if (entry.isDirectory()) add(child);
    else hash.update(child).update(readFileSync(join(root, child)));
  }
}
for (const path of ['app', 'public', 'scripts']) add(path);
for (const path of ['vite.config.ts', 'package.json', 'package-lock.json']) hash.update(path).update(readFileSync(join(root, path)));
export const appVersion = hash.digest('hex').slice(0, 16);
