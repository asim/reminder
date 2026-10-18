// Run with: node web/scripts/test-reading-bookmarks.cjs
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const ts = require('typescript');
const source = fs.readFileSync(require('node:path').join(__dirname, '../app/hooks/use-reading-bookmark.ts'), 'utf8');
const code = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.CommonJS } }).outputText;
const key = 'reminder_reading_bookmarks';
const saved = new Map([[key, JSON.stringify({ quran: { url: '/quran/2#1', label: 'Old', timestamp: '2026-01-01' } })]]);
const listeners = new Set();
const unmount = [];
let reads = 0, getSnapshot, getServerSnapshot;
const context = {
  exports: {}, URL, console,
  window: {
    addEventListener(type, fn) { assert.equal(type, 'storage'); listeners.add(fn); },
    removeEventListener(type, fn) { assert.equal(type, 'storage'); listeners.delete(fn); },
  },
  localStorage: {
    getItem(k) { reads++; return saved.get(k) ?? null; },
    setItem(k, value) { saved.set(k, value); },
  },
  require(name) {
    assert.equal(name, 'react');
    return { useSyncExternalStore(subscribe, get, server) {
      getSnapshot = get; getServerSnapshot = server;
      unmount.push(subscribe(() => get()));
      return get();
    } };
  },
};
vm.runInNewContext(code, context);
const buttons = Array.from({ length: 286 }, () => context.exports.useReadingBookmark());
assert.equal(listeners.size, 1, 'one listener for a whole chapter');
assert.equal(reads, 1, 'initial storage read is shared');
assert.equal(getSnapshot(), getSnapshot(), 'stable snapshot between events');
assert.equal(getServerSnapshot(), getServerSnapshot(), 'stable server snapshot');
assert(buttons[285].isReadingBookmark('quran', '/quran/2?translation=asad#1'));
reads = 0;
buttons[0].addReadingBookmark('quran', 'Second', '/quran/2?translation=asad#2');
assert.equal(reads, 1, 'one fresh read per mutation, not one per subscriber');
assert.equal(getSnapshot().quran.length, 2);
buttons[1].removeReadingBookmark('quran', '/quran/2/1');
assert.equal(getSnapshot().quran.length, 1, 'stale consumer preserves the newer bookmark');
assert.equal(getSnapshot().quran[0].url, '/quran/2?translation=asad#2');
reads = 0;
saved.set(key, JSON.stringify({ quran: [{ url: '/quran/2#3' }], hadith: [], names: [] }));
listeners.forEach(fn => fn({ key }));
assert.equal(reads, 1, 'one parse/read per cross-tab event');
assert.equal(getSnapshot().quran[0].url, '/quran/2#3');
reads = 0;
listeners.forEach(fn => fn({ key: 'unrelated' }));
assert.equal(reads, 0);
saved.delete(key);
listeners.forEach(fn => fn({ key: null }));
assert.equal(getSnapshot().quran.length, 0, 'storage clear propagates');
unmount.forEach(fn => fn());
assert.equal(listeners.size, 0, 'last unmount removes the listener');
saved.set(key, JSON.stringify({ names: [{ url: '/names/68' }] }));
assert(context.exports.useReadingBookmark().isReadingBookmark('names', '/names/68'));
assert.equal(listeners.size, 1, 'remount refreshes from storage');
console.log('PASS: 286 consumers share one listener/read; mutations, cross-tab updates, migration and remount remain correct.');
