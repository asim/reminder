#!/usr/bin/env python3
"""Read-only MCP smoke test. --sdk also checks the official Python MCP client."""
import argparse
import asyncio
import json
import urllib.request

parser = argparse.ArgumentParser()
parser.add_argument('endpoint')
parser.add_argument('--sdk', action='store_true')
args = parser.parse_args()
sequence = 0

def rpc(method, params=None, notification=False):
    global sequence
    sequence += 1
    payload = {'jsonrpc': '2.0', 'method': method}
    if not notification:
        payload['id'] = sequence
    if params is not None:
        payload['params'] = params
    request = urllib.request.Request(args.endpoint, data=json.dumps(payload).encode(), headers={
        'Content-Type': 'application/json', 'Accept': 'application/json, text/event-stream',
        'MCP-Protocol-Version': '2025-06-18',
    })
    with urllib.request.urlopen(request, timeout=45) as response:
        if notification:
            assert response.status == 202
            return
        result = json.load(response)
    assert result['id'] == sequence and 'error' not in result, result
    return result['result']

init = rpc('initialize', {'protocolVersion': '2025-06-18', 'capabilities': {}, 'clientInfo': {'name': 'reminder-smoke', 'version': '1.0'}})
assert init['protocolVersion'] == '2025-06-18' and init['instructions']
rpc('notifications/initialized', notification=True)
tools = {tool['name']: tool for tool in rpc('tools/list')['tools']}
for name in ['get_daily', 'get_quran_verse', 'search', 'arabic_dictionary', 'quran_morphology', 'explain_quran']:
    assert tools[name]['annotations']['readOnlyHint'] is True
assert tools['explain_quran']['outputSchema']['required'] == ['answer', 'references']
verse = rpc('tools/call', {'name': 'get_quran_verse', 'arguments': {'chapter': 2, 'verse': 153}})
record = verse['structuredContent']
assert not verse.get('isError') and record['arabic'] and record['comments']
assert record['url'] == 'https://reminder.dev/quran/2/153'
invalid = rpc('tools/call', {'name': 'explain_quran', 'arguments': {'context': {'chapter': 115}}})
assert invalid['isError'], 'invalid chapter must fail before generation'
print(f'MCP smoke passed: {len(tools)} tools, source retrieval and invalid-context handling')

async def sdk_test():
    from mcp import ClientSession
    from mcp.client.streamable_http import streamablehttp_client
    from urllib.parse import urlparse
    options = {}
    if urlparse(args.endpoint).hostname in {'127.0.0.1', 'localhost', '::1'}:
        import httpx
        def local_client(**kwargs):
            return httpx.AsyncClient(**kwargs, trust_env=False, follow_redirects=True)
        options['httpx_client_factory'] = local_client
    async with streamablehttp_client(args.endpoint, **options) as (read, write, _):
        async with ClientSession(read, write) as session:
            initialized = await session.initialize()
            listing = await session.list_tools()
            assert any(tool.name == 'explain_quran' for tool in listing.tools)
            result = await session.call_tool('get_quran_verse', {'chapter': 1, 'verse': 2})
            assert not result.isError
            assert json.loads(result.content[0].text)['number'] == 2
            print(f'Official Python MCP client passed: {initialized.protocolVersion}')
if args.sdk:
    asyncio.run(sdk_test())
