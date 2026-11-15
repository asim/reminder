import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

// Build-only data, exported with the Go corpus loaders, never the live API.
export function readingData(key: string, id?: string) {
  if (id !== undefined && !/^[1-9]\d*$/.test(id)) throw new Response('Not found', { status: 404 });
  try { return JSON.parse(readFileSync(resolve('.generated', `${key}${id ? `-${id}` : ''}.json`), 'utf8')); }
  catch { throw new Response('Not found', { status: 404 }); }
}
export function readingLoader(key: string, queryKey: unknown[], id?: string) {
  return { readingQueries: [{ queryKey, data: readingData(key, id) }] };
}
