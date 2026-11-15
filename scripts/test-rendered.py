#!/usr/bin/env python3
"""Check generated reading HTML independently of browser JavaScript."""
import gzip
import json
from html import unescape
from pathlib import Path
import re

root = Path(__file__).resolve().parents[1] / 'web'

def normal(s): return ' '.join(s.split())
paths = json.loads((root / '.generated/paths.json').read_text())
for route in paths:
    file = root / 'build/client' / route.lstrip('/') / 'index.html.gz'
    html = gzip.decompress(file.read_bytes()).decode()
    assert re.search(r'<title>[^<]+</title>', html), f'{route}: missing title'
    assert f'rel="canonical" href="https://reminder.dev{route}"' in html, f'{route}: missing canonical'
    visible = re.sub(r'<(script|style)\b[^>]*>.*?</\1>', '', html, flags=re.S)
    text = normal(unescape(re.sub(r'<[^>]*>', ' ', visible)))
    parts = route.strip('/').split('/')
    if parts[0] == 'quran' and len(parts) > 1:
        chapter = json.loads((root / f'.generated/quran-{parts[1]}.json').read_text())
        verses = [v for v in chapter['verses'] if v['number'] > 0 and (len(parts) == 2 or v['number'] == int(parts[2]))]
        for verse in verses:
            assert normal(verse['text']) in text, f'{route}: missing verse {verse["number"]}'
    elif parts[0] == 'names' and len(parts) == 2:
        name = json.loads((root / f'.generated/name-{parts[1]}.json').read_text())
        assert normal(name['summary']) in text, f'{route}: missing name summary'
    elif parts[0] == 'hadith' and len(parts) == 2:
        book = json.loads((root / f'.generated/hadith-{parts[1]}.json').read_text())
        assert normal(book['hadiths'][0]['english']) in text, f'{route}: missing hadith'
print(f'Validated content, titles and canonical URLs for {len(paths)} HTML pages without JavaScript.')

import xml.etree.ElementTree as ET
sitemap = ET.parse(root / 'build/client/sitemap.xml')
locations = {node.text for node in sitemap.findall('.//{http://www.sitemaps.org/schemas/sitemap/0.9}loc')}
assert locations == {'https://reminder.dev' + route for route in paths if route != '/api'}
assert 'Sitemap: https://reminder.dev/sitemap.xml' in (root / 'build/client/robots.txt').read_text()
print('Sitemap covers the rendered reading routes; robots.txt advertises it.')
