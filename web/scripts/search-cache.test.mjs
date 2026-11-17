import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import ts from 'typescript';

const source = await readFile(new URL('../app/utils/search-cache.ts', import.meta.url), 'utf8');
const js = ts.transpileModule(source, {compilerOptions: {module: ts.ModuleKind.ESNext}}).outputText;
const {cacheSearchResult, getCachedSearchResult, getAllCachedSearches, restoreCachedSearches} = await import(`data:text/javascript;base64,${Buffer.from(js).toString('base64')}`);

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

test('restoring Ask AI never seeds or overwrites ordinary-search evidence', () => {
  const storage = new Map();
  globalThis.window = {};
  globalThis.localStorage = {getItem: key => storage.get(key) ?? null, setItem: (key,value) => storage.set(key,value)};
  const ai = {q:'Musa', source:'quran', mode:'ask', answer:'Answer', references:[{text:'Complete mention inventory', metadata:{source:'quran',kind:'mention_inventory'}}]};
  cacheSearchResult('Musa', ai);
  const restored = new Map();
  restoreCachedSearches('Musa', 'quran', 'ask', (mode,result) => restored.set(mode,result));
  assert.equal(restored.has('search'),false, 'Search must fetch its own results when none were saved');
  assert.deepEqual(restored.get('ask'),ai);

  const plain = {q:'Musa', source:'quran', mode:'search', answer:'', references:[{text:'Ranked verse result',metadata:{source:'quran'}}]};
  cacheSearchResult('Musa',plain);
  restored.clear();
  restoreCachedSearches('Musa','quran','ask',(mode,result) => restored.set(mode,result));
  assert.deepEqual(restored.get('search'),plain);
  assert.deepEqual(restored.get('ask'),ai);
  // Switching back to Search and persisting the restored results retains its own evidence.
  cacheSearchResult('Musa',restored.get('search'));
  assert.deepEqual(getCachedSearchResult('Musa','quran','search'),plain);
  restored.clear();
  restoreCachedSearches('Musa','quran','search',(mode,result) => restored.set(mode,result));
  assert.deepEqual([...restored.keys()],['search']);
});
