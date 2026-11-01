const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const ts = require('typescript');
let slots = [], cursor = 0, effects = [], playing = null, reduced = false, calls = [];
const react = {
  useRef(value) { return slots[cursor++] ||= { current: value }; },
  useState() { cursor++; return [playing, () => {}]; },
  useEffect(fn, deps) {
    const i = cursor++, previous = slots[i];
    if (!previous || deps.some((v, n) => v !== previous[n])) effects.push(fn);
    slots[i] = deps;
  },
};
const exportsObject = {};
const chapter = { number: 2, name: '', english: '', verses: [] };
const code = ts.transpileModule(fs.readFileSync('app/routes/_app.quran.$chapterNumber._index.tsx', 'utf8'), {
  compilerOptions: { module: ts.ModuleKind.CommonJS, jsx: ts.JsxEmit.ReactJSX },
}).outputText;
vm.runInNewContext(code, {
  exports: exportsObject,
  require: name => name === 'react' ? react :
    name === 'react/jsx-runtime' ? { jsx: () => null, jsxs: () => null } :
    name === 'react-router' ? { useLocation: () => ({ hash: '' }) } :
    name === '~/queries/quran' ? { useTranslatedChapter: () => chapter } :
    name === '~/components/quran/reading-settings' ? { useQuranSettings: () => ({ settings: {} }) } :
    name === '~/utils/quran-share' ? { buildQuranShareUrl: () => '' } : {},
  window: { location: { hash: '' }, matchMedia: () => ({ matches: reduced }) },
  document: { getElementById: id => ({ scrollIntoView: options => calls.push({ id, ...options }) }) },
});
function render(label, number = '2') {
  playing = label; cursor = 0; effects = [];
  exportsObject.default({ params: { chapterNumber: number } });
  effects.forEach(fn => fn());
}
render(null);
assert.equal(calls.length, 0, 'No scrolling before playback');
render('Quran 2:5');
assert.equal(calls.at(-1).id, '5');
assert.equal(calls.at(-1).behavior, 'smooth');
render(null); render('Quran 2:5'); render('Quran 2:5');
assert.equal(calls.length, 1, 'Pause/resume and language changes do not scroll again');
render('Quran 2:6');
assert.equal(calls.at(-1).id, '6', 'Follow the next verse');
reduced = true;
render('Quran 2:7');
assert.equal(calls.at(-1).behavior, 'instant', 'Respect reduced motion');
render('Quran 2:7', '3');
assert.equal(calls.length, 3, 'Ignore stale playback from the previous chapter');
render('Quran 3:1', '3');
assert.equal(calls.at(-1).id, '1');
console.log('Verse following: start, advancement, pause/resume, reduced motion and chapter changes passed');
