#!/usr/bin/env python3
"""Normalize the pinned asim/hadith-json Muslim export without renumbering it."""
import gzip, hashlib, json, sys
from pathlib import Path
raw = Path(sys.argv[1]).read_bytes()
blob = hashlib.sha1(b'blob ' + str(len(raw)).encode() + b'\0' + raw).hexdigest()
if blob != '7941fc29ab2eb780406d4d9142da623f89466a5d':
    raise SystemExit('Unexpected source revision; review the source before updating the pinned blob')
data = json.loads(raw)
books = {c['id']: {'number': c['id'], 'name': c['english'], 'hadiths': []} for c in data['chapters']}
assert books[0]['name'] == 'Introduction'
seen = set()
source_ids = set()
assert len(books) == len(data['chapters']) == 57
for h in data['hadiths']:
    number = h['idInBook']
    assert number not in seen and h['id'] not in source_ids and h['chapterId'] in books and h['bookId'] == 2
    source_ids.add(h['id'])
    assert h['arabic'].strip()  # Record 6292 has no English in the source; retain its Arabic.
    seen.add(number)
    books[h['chapterId']]['hadiths'].append({'number': number, 'source_id': h['id'], 'arabic': h['arabic'], 'narrator': h['english']['narrator'], 'english': h['english']['text']})
assert len(seen) == data['metadata']['length'] == 7459
out = {'id': 'muslim', 'name': 'Sahih Muslim', 'arabic': data['metadata']['arabic']['title'], 'books': [books[key] for key in sorted(books)]}
root = Path(__file__).resolve().parents[1]
(root/'hadith/data/muslim.json.gz').write_bytes(gzip.compress((json.dumps(out, ensure_ascii=False, separators=(',', ':'))+'\n').encode(),mtime=0))
(root/'hadith/muslim-source.json').write_text(json.dumps({'repository':'https://github.com/asim/hadith-json','path':'db/by_book/the_9_books/muslim.json','git_blob':blob,'sha256':hashlib.sha256(raw).hexdigest(),'records':len(seen),'numbering':'number preserves dataset idInBook; source_id preserves dataset id. Neither is asserted to be Sunnah.com numbering.','chapter_zero':'91 Introduction records retained as source chapter 0.'},indent=2)+'\n')
print('Imported',len(seen),'records in',len(books),'groups')
