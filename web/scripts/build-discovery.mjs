import { readFile, writeFile, readdir, unlink } from 'node:fs/promises';
import { gzipSync } from 'node:zlib';
const root = new URL('../build/client/', import.meta.url);
const paths = JSON.parse(await readFile(new URL('../.generated/paths.json', import.meta.url)));
const embeds = JSON.parse(await readFile(new URL('../.generated/embeds.json', import.meta.url)));
await writeFile(new URL('robots.txt', root), 'User-agent: *\nAllow: /\n\nSitemap: https://reminder.dev/sitemap.xml\n');
await writeFile(new URL('sitemap.xml', root), '<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n' + paths.filter(p => p !== '/api').map(p => `<url><loc>https://reminder.dev${p}</loc></url>`).join('\n') + '\n</urlset>\n');
await writeFile(new URL('embed-paths.txt', root), embeds.join('\n') + '\n');
// Keep thousands of reading pages compact in the embedded Go binary. The Go
// server transparently decompresses them; the browser still receives normal HTML.
async function compress(dir) {
 for (const entry of await readdir(dir, { withFileTypes: true })) {
  const file = new URL(entry.name, dir);
  if (entry.isDirectory()) { await compress(new URL(entry.name + '/', dir)); continue; }
  if ((entry.name.endsWith('.html') && entry.name !== '__spa-fallback.html') || entry.name.endsWith('.data')) {
   const data = await readFile(file);
   await writeFile(new URL(entry.name + '.gz', dir), gzipSync(data));
   await unlink(file);
  }
 }
}
await compress(root);
console.log(`Discovery: ${paths.length} rendered URLs, robots.txt and sitemap.xml`);
