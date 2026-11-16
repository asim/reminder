import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import ts from 'typescript';

const source = await readFile(new URL('../app/utils/search-cache.ts', import.meta.url), 'utf8');
const js = ts.transpileModule(source, {compilerOptions: {module: ts.ModuleKind.ESNext}}).outputText;
const {cacheSearchResult, getCachedSearchResult, getAllCachedSearches} = await import(`data:text/javascript;base64,${Buffer.from(js).toString('base64')}`);

test('plain results survive reload and source/mode variants coexist', () => {
  const storage = new Map();
  globalThis.window = {};
  globalThis.localStorage = {getItem: key => storage.get(key) ?? null, setItem: (key,value) => storage.set(key,value)};
  const result = (source,mode) => ({q:'Musa',source,mode,answer:mode === 'ask' ? 'An answer' : '',references:[{text:source,metadata:{source}}]});
  cacheSearchResult('Musa',result('quran','search'));
  cacheSearchResult('Musa',result('hadith','search'));
  cacheSearchResult('Musa',result('quran','ask'));
  assert.equal(getAllCachedSearches().length,3);
  // Every read rehydrates the persisted JSON, as a fresh page does.
  assert.equal(getCachedSearchResult('Musa','quran','search').references[0].text,'quran');
  assert.equal(getCachedSearchResult('Musa','quran','search').answer,'');
  assert.equal(getCachedSearchResult('Musa','quran','ask').answer,'An answer');
  cacheSearchResult(' musa ',result('quran','search'));
  assert.equal(getAllCachedSearches().length,3);
  assert.equal(getCachedSearchResult('Musa','tafsir','search'),null);
});

test('unavailable storage does not break searching', () => {
  globalThis.localStorage = {getItem:()=>{throw Error('blocked')},setItem:()=>{throw Error('blocked')}};
  const previous=console.error; console.error=()=>{};
  try {
    assert.doesNotThrow(()=>cacheSearchResult('query',{q:'query',answer:'',references:[]}));
    assert.deepEqual(getAllCachedSearches(),[]);
  } finally {console.error=previous;}
});
