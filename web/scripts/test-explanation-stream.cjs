const assert = require('node:assert/strict'), fs = require('node:fs'), vm = require('node:vm'), ts = require('typescript');
const output = {};
vm.runInNewContext(ts.transpileModule(fs.readFileSync('app/utils/explanation-stream.ts','utf8'), {compilerOptions:{module:ts.ModuleKind.CommonJS}}).outputText,{exports:output,TextDecoder,Error});
const encoder = new TextEncoder();
function response(lines, bytewise=false) {
 const bytes=encoder.encode(lines);
 return new Response(new ReadableStream({start(c){if(bytewise) for(const b of bytes)c.enqueue(new Uint8Array([b]));else c.enqueue(bytes);c.close();}}),{headers:{'Content-Type':'application/x-ndjson'}});
}
(async()=>{
 const events=[];
 const result=await output.readExplanation(response('{"type":"delta","text":"حديث"}\n{"type":"reset"}\n{"type":"result","result":{"answer":"Final [1].","references":[]}}\n',true),e=>events.push(e));
 assert.equal(events[0].text,'حديث');assert.equal(events[1].type,'reset');assert.equal(result.answer,'Final [1].');
 await assert.rejects(output.readExplanation(response('{"type":"delta","text":"Partial"}\n'),()=>{}),/interrupted/);
 await assert.rejects(output.readExplanation(response('{"type":"delta","text":"Partial"}\n{"type":"error","text":"Invalid citation"}\n'),()=>{}),/Invalid citation/);
 const fallback=await output.readExplanation(new Response('{"answer":"Cached","references":[]}',{headers:{'Content-Type':'application/json'}}),()=>{});
 assert.equal(fallback.answer,'Cached');
 // Progress arrives while the stream is still open, not after its final result.
 let controller;
 const live=new Response(new ReadableStream({start(c){controller=c;}}),{headers:{'Content-Type':'application/x-ndjson'}});
 let seen;
 const progress=new Promise(resolve=>seen=resolve);
 const pending=output.readExplanation(live,e=>seen(e));
 controller.enqueue(encoder.encode('{"type":"delta","text":"Early"}\n'));
 assert.equal((await progress).text,'Early');
 controller.enqueue(encoder.encode('{"type":"result","result":{"answer":"Done","references":[]}}\n'));controller.close();
 assert.equal((await pending).answer,'Done');
 console.log('Explanation stream: incremental text, UTF-8 boundaries, resets, failures and JSON fallback passed');
})().catch(err=>{console.error(err);process.exitCode=1;});
