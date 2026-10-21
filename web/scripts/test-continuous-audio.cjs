const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const ts = require('../node_modules/typescript');
const source = fs.readFileSync('web/app/components/quran/continuous-audio.ts', 'utf8');
const code = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2022 } }).outputText;
const flush = () => new Promise(resolve => setImmediate(resolve));
const ranges = (start, end) => ({ length: end > start ? 1 : 0, start: () => start, end: () => end });

function harness(fetchOverride) {
  const requests = [], positions = [], failures = [], revoked = [];
  let media;
  class Buffer extends EventTarget {
    buffered = ranges(0, 0);
    updating = false;
    appendBuffer(bytes) {
      assert.equal(this.updating, false);
      this.updating = true;
      queueMicrotask(() => {
        this.buffered = ranges(this.buffered.start(0), this.buffered.end(0) + new Uint8Array(bytes)[0]);
        this.updating = false;
        this.dispatchEvent(new Event('updateend'));
      });
    }
    remove(start, end) {
      assert.equal(this.updating, false);
      this.updating = true;
      queueMicrotask(() => {
        this.buffered = ranges(end, this.buffered.end(0));
        this.updating = false;
        this.dispatchEvent(new Event('updateend'));
      });
    }
  }
  class MediaSource extends EventTarget {
    static isTypeSupported(type) { return type === 'audio/mpeg'; }
    readyState = 'closed';
    constructor() { super(); media = this; }
    addSourceBuffer(type) { assert.equal(type, 'audio/mpeg'); return this.buffer = new Buffer(); }
    endOfStream() { this.readyState = 'ended'; }
  }
  const audio = new EventTarget();
  Object.assign(audio, { currentTime: 0, paused: false, ended: false, sources: [], pause() { this.paused = true; } });
  Object.defineProperty(audio, 'src', { set(value) { this.sources.push(value); } });
  const exported = {};
  vm.runInNewContext(code, { exports: exported, MediaSource, URL: { createObjectURL: () => 'blob:stream', revokeObjectURL: url => revoked.push(url) },
    AbortController, setTimeout, clearTimeout, fetch: async (url, options) => {
      requests.push(url);
      return fetchOverride ? fetchOverride(url, options) : { ok: true, arrayBuffer: async () => Uint8Array.of(48).buffer };
    } });
  const stream = new exported.ContinuousAudio(audio, Array.from({ length: 8 }, (_, i) => `track-${i}`), 0,
    (segment, elapsed) => positions.push({ ...segment, elapsed }), message => failures.push(message));
  const open = () => { media.readyState = 'open'; media.dispatchEvent(new Event('sourceopen')); };
  const tick = time => { audio.currentTime = time; audio.dispatchEvent(new Event('timeupdate')); };
  return { stream, audio, requests, positions, failures, revoked, open, tick, media: () => media };
}

(async () => {
  const h = harness(); h.open(); await flush();
  assert.deepEqual(h.requests, ['track-0', 'track-1']); // 96 seconds, bounded prefetch.
  assert.equal(h.media().buffer.mode, 'sequence');
  h.tick(49); await flush();
  assert.equal(h.positions.at(-1).index, 1); assert.equal(h.positions.at(-1).elapsed, 1);
  assert.equal(h.stream.seek(1, 10), true); assert.equal(h.audio.currentTime, 58);
  assert.equal(h.stream.seek(7), false); // Caller may rebuild only on an explicit skip.
  h.audio.pause(); h.tick(60); await flush(); assert.equal(h.audio.paused, true);
  h.audio.paused = false;
  for (const time of [97, 145, 193, 241, 289, 337]) { h.tick(time); await flush(); }
  assert.equal(h.media().readyState, 'ended');
  assert.deepEqual(h.requests, Array.from({ length: 8 }, (_, i) => `track-${i}`));
  assert.deepEqual(h.audio.sources, ['blob:stream']); // No source reset across Arabic/English/verse boundaries.
  assert.equal(h.audio.paused, false);
  assert(h.media().buffer.buffered.start(0) > 0); // Played audio is removed.
  assert.equal(h.stream.seek(0), false);
  h.stream.dispose(); h.tick(350); await flush(); assert.deepEqual(h.revoked, ['blob:stream']);
  assert.deepEqual(h.failures, []);

  const bad = harness(async () => ({ ok: false })); bad.open(); await flush();
  assert.equal(bad.stream.hasFailed, true); assert.equal(bad.audio.paused, true); assert.equal(bad.failures.length, 1); bad.stream.dispose();
  let aborted = false;
  const cancelled = harness((url, { signal }) => new Promise((resolve, reject) => signal.addEventListener('abort', () => { aborted = true; reject(new Error('aborted')); })));
  cancelled.open(); await flush(); cancelled.stream.dispose(); await flush();
  assert.equal(aborted, true); assert.deepEqual(cancelled.failures, []);
  console.log('PASS: ordered continuous timeline, bounded lookahead, seeking, pause preservation, buffer eviction, failure and cancellation.');
})().catch(error => { console.error(error); process.exitCode = 1; });
