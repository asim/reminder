const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const ts = require('typescript');

let slots = [], cursor = 0, continuous = false, streams = [];
const react = {
  useRef: value => slots[cursor++] ||= { current: value },
  useState: value => { const i = cursor++; if (!(i in slots)) slots[i] = value; return [slots[i], v => slots[i] = v]; },
  useEffect() {},
  useImperativeHandle: (ref, create) => { ref.current = create(); },
};
class ContinuousAudio {
  static supported() { return true; }
  constructor(audio, urls, first) { this.urls = urls; this.first = first; streams.push(this); }
  seek() { return false; }
  dispose() { this.disposed = true; }
}
const moduleObject = { exports: {} };
const code = ts.transpileModule(fs.readFileSync('app/components/quran/audio-player.tsx', 'utf8'), {
  compilerOptions: { module: ts.ModuleKind.CommonJS, jsx: ts.JsxEmit.ReactJSX },
}).outputText;
vm.runInNewContext(code, {
  exports: moduleObject.exports, module: moduleObject,
  require: name => name === 'react' ? react :
    name === 'react/jsx-runtime' ? { jsx: (type, props) => ({ type, props }), jsxs: (type, props) => ({ type, props }) } :
    name === './continuous-audio' ? { ContinuousAudio } : {},
  navigator: {}, location: { get search() { return continuous ? '?continuousAudio=1' : ''; } },
  URLSearchParams,
});
const { AudioPlayer } = moduleObject.exports;
const ref = { current: null };
const playlist = [1, 2, 3].flatMap(v => ['arabic', 'english'].map(language => ({
  url: v + '-' + language, language, label: 'Quran 2:' + v,
})));
let audioNode, tree, playCalls = 0;
const audio = {
  paused: true, ended: false,
  play() { this.paused = false; this.ended = false; playCalls++; return Promise.resolve(); },
  pause() { this.paused = true; },
};
function find(node, predicate) {
  if (!node || typeof node !== 'object') return;
  if (predicate(node)) return node;
  return [node.props?.children].flat(2).map(n => find(n, predicate)).find(Boolean);
}
function render() {
  cursor = 0;
  tree = AudioPlayer({ ref, playlist, verseLabel: 'Quran 2:1' });
  audioNode = find(tree, n => n.type === 'audio');
  audioNode.props.ref.current = audio;
}
function ended() { audio.ended = true; audio.paused = true; audioNode.props.onEnded(); render(); }
render();
ref.current.toggleVerse('Quran 2:2');
assert.equal(audio.src, '2-arabic');
render();
ref.current.toggleVerse('Quran 2:2');
assert.equal(audio.paused, true);
ref.current.toggleVerse('Quran 2:2');
assert.equal(audio.paused, false);
ended();
assert.equal(audio.src, '2-english');
const beforeEnd = playCalls;
ended();
assert.equal(playCalls, beforeEnd, 'Must stop after selected English verse');
ref.current.toggleVerse('Quran 2:2');
assert.equal(audio.src, '2-arabic', 'Replay starts in Arabic');
ref.current.toggleVerse('Quran 2:3');
assert.equal(audio.src, '3-arabic', 'Another verse replaces playback');
continuous = true;
ref.current.toggleVerse('Quran 2:1');
assert.equal(streams.at(-1).urls.join(','), '1-arabic,1-english');
const oldStream = streams.at(-1);
ref.current.toggleVerse('Quran 2:2');
assert.equal(oldStream.disposed, true);
assert.equal(streams.at(-1).first, 2);
assert.equal(streams.at(-1).urls.at(-1), '2-english');
slots = []; continuous = false; audio.paused = true; audio.ended = false; render();
find(tree, n => n.props?.['aria-label'] === 'Play audio').props.onClick();
assert.equal(audio.src, '1-arabic');
ended(); ended();
assert.equal(audio.src, '2-arabic', 'Default chapter playback still advances');
console.log('Verse playback: Arabic/English, pause/resume, stop, replay, switching, continuous boundaries, chapter playback passed');
