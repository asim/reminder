const assert = require('node:assert/strict'), fs = require('node:fs'), vm = require('node:vm'), ts = require('typescript');
let slots = [], cursor = 0, effects = [], available = false, requests = 0;
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
 require: name => name === 'react' ? react : name === 'react/jsx-runtime' ? { jsx, jsxs: jsx } : name === '~/queries/capabilities' ? { useModelAvailable: () => available } : {},
 window: { getSelection: () => null },
 fetch: async () => { requests++; return { ok: false, json: async () => ({ message: 'Temporarily unavailable' }) }; },
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
 reset(false); assert.equal(render({...props,commentary:[]}).length,0);
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
 console.log('Explain panel: single action, availability, attributed tafsir and recoverable failure passed');
})().catch(err => { console.error(err); process.exitCode = 1; });
