const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const ts = require('../node_modules/typescript');
const code = ts.transpileModule(fs.readFileSync('web/app/utils/screen-wake-lock.ts', 'utf8'),
  { compilerOptions: { module: ts.ModuleKind.CommonJS } }).outputText;
const flush = () => new Promise(resolve => setImmediate(resolve));
function setup(request) {
  const document = new EventTarget();
  document.visibilityState = 'visible';
  const exports = {};
  vm.runInNewContext(code, { exports, document, navigator: request ? { wakeLock: { request } } : {} });
  return { start: exports.keepScreenAwake, visibility(value) {
    document.visibilityState = value;
    document.dispatchEvent(new Event('visibilitychange'));
  } };
}
function sentinel() {
  return { released: false, calls: 0, async release() { this.calls++; this.released = true; } };
}
(async () => {
  const locks = [];
  const h = setup(async type => { assert.equal(type, 'screen'); const lock = sentinel(); locks.push(lock); return lock; });
  const stop = h.start(); await flush(); assert.equal(locks.length, 1);
  h.visibility('visible'); await flush(); assert.equal(locks.length, 1);
  h.visibility('hidden'); await flush(); assert.equal(locks[0].calls, 1);
  h.visibility('visible'); await flush(); assert.equal(locks.length, 2);
  stop(); await flush(); assert.equal(locks[1].calls, 1);
  h.visibility('visible'); await flush(); assert.equal(locks.length, 2);

  let resolve;
  const late = sentinel();
  const delayed = setup(() => new Promise(done => { resolve = done; }));
  const cancel = delayed.start(); cancel(); resolve(late); await flush();
  assert.equal(late.calls, 1, 'pause/unmount must release a late grant');

  const denied = setup(async () => { throw new Error('denied'); });
  const cleanDenied = denied.start(); await flush(); cleanDenied();
  setup().start()(); // Unsupported browsers are a no-op.
  console.log('PASS: acquire, release, foreground reacquire, no duplicate requests, late-grant cleanup, denial and unsupported browsers.');
})().catch(error => { console.error(error); process.exitCode = 1; });
