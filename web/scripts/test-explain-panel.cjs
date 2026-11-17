const assert = require('node:assert/strict'), fs = require('node:fs'), vm = require('node:vm'), ts = require('typescript');
let slots = [], cursor = 0, effects = [], available = false, requests = 0, lastBody, success = false;
const jsx = (type, props) => ({ type, props: props || {} });
const react = {
 useState(initial) { const i = cursor++; if (!(i in slots)) slots[i] = initial; return [slots[i], v => { slots[i] = typeof v === 'function' ? v(slots[i]) : v; }]; },
 useRef(initial) { const i = cursor++; return slots[i] ||= { current: initial }; },
 useEffect(fn, deps) { const i = cursor++, old = slots[i]; if (!old || deps.some((v,n) => old[n] !== v)) effects.push(fn); slots[i] = deps; },
};
const output = {};
vm.runInNewContext(ts.transpileModule(fs.readFileSync('app/components/quran/explain-button.tsx','utf8'), {
 compilerOptions: { module: ts.ModuleKind.CommonJS, jsx: ts.JsxEmit.ReactJSX },
}).outputText, {
 exports: output, AbortController, Error,
 require: name => name === 'react' ? react : name === 'react/jsx-runtime' ? { jsx, jsxs: jsx } : name === '~/utils/explanation-stream' ? { readExplanation: response => response.json() } : name === '~/queries/capabilities' ? { useModelAvailable: () => available } : {},
 window: { getSelection: () => null },
 fetch: async (url, options) => { requests++; lastBody = JSON.parse(options.body); return { ok: success, json: async () => success ? { answer: 'Explanation', references: [] } : { message: 'Temporarily unavailable' } }; },
});
const props = { chapter: 2, verse: 153, translation: 'sahih', commentary: [{ number: 153, comments: 'Stored tafsir text', commentary_source: 'Tafsir Al-Jalalayn' }] };
function flatten(n) {
 if (n == null || n === false) return [];
 if (Array.isArray(n)) return n.flatMap(flatten);
 return typeof n === 'object' ? [n,...flatten(n.props.children)] : [n];
}
function render(p = props) {
 cursor = 0; effects = []; const nodes = flatten(output.ExplainButton(p));
 for (const n of nodes) if (n.type === 'dialog') n.props.ref.current = { showModal() {}, close() {} };
 effects.forEach(fn => fn()); return nodes;
}
const button = (nodes,label) => nodes.find(n => n.type === 'button' && (n.props['aria-label'] === label || n.props.children === label));
function reset(model) { slots = []; available = model; requests = 0; }
(async () => {
 reset(false); let nodes = render();
 assert.equal(nodes.filter(n => n.type === 'button').length,1,'one outside action');
 button(nodes,'Tafsir Quran 2:153').props.onClick(); nodes = render();
 assert(nodes.includes('Stored tafsir text')); assert(nodes.includes('Tafsir Al-Jalalayn'));
 assert(!nodes.some(n => n.type === 'form')); assert.equal(requests,0);
 reset(false); assert(button(render({...props,commentary:[]}), 'Tafsir Quran 2:153'));
 reset(true); nodes = render(); assert.equal(nodes.filter(n => n.type === 'button').length,1);
 button(nodes,'Explain Quran 2:153').props.onClick(); nodes = render();
 nodes.find(n => n.type === 'textarea').props.onChange({target:{value:'What does this mean?'}}); nodes = render();
 nodes.find(n => n.type === 'form').props.onSubmit({preventDefault(){}});
 await new Promise(setImmediate); nodes = render();
 assert.equal(nodes.find(n => n.type === 'textarea').props.value,'What does this mean?');
 assert(nodes.some(n => n.props?.role === 'alert'));
 button(nodes,'Tafsir').props.onClick(); nodes = render(); assert(nodes.includes('Stored tafsir text')); assert.equal(requests,1);
 button(nodes,'Explain').props.onClick(); nodes = render();
 assert.equal(nodes.find(n => n.type === 'textarea').props.value,'What does this mean?');
 const hadith = { hadith: { collection: 'muslim', book: 0, number: 7369, label: 'Sahih Muslim record 7369' } };
 reset(false); assert.equal(render(hadith).length,0);
 reset(true); success = true; nodes = render(hadith);
 button(nodes,'Explain Sahih Muslim record 7369').props.onClick(); nodes = render(hadith);
 assert(!nodes.some(n => n.type === 'nav'), 'hadith has no tafsir tab');
 nodes.find(n => n.type === 'form').props.onSubmit({preventDefault(){}});
 await new Promise(setImmediate); nodes = render(hadith);
 assert.equal(lastBody.context.source,'hadith'); assert.equal(lastBody.context.collection,'muslim');
 assert.equal(lastBody.context.book,0); assert.equal(lastBody.context.hadith,7369);
 assert.equal(lastBody.context.chapter,undefined); assert.equal(lastBody.question,'');
 nodes.find(n => n.type === 'textarea').props.onChange({target:{value:'What does that mean?'}}); nodes = render(hadith);
 nodes.find(n => n.type === 'form').props.onSubmit({preventDefault(){}});
 await new Promise(setImmediate); nodes = render(hadith);
 assert.equal(lastBody.context.hadith,7369); assert.equal(lastBody.history.length,2);
 assert.equal(lastBody.question,'What does that mean?');
 render({hadith:{...hadith.hadith,collection:'bukhari',book:1,number:1}});
 nodes = render({hadith:{...hadith.hadith,collection:'bukhari',book:1,number:1}});
 assert(!nodes.some(n => n.type === 'dialog'), 'changing records clears the panel');
 console.log('Explain panel: single action, availability, attributed tafsir and recoverable failure passed');
})().catch(err => { console.error(err); process.exitCode = 1; });
