#!/usr/bin/env python3
"""Validate a QUL English Al-Jalalayn JSON/ZIP download before replacing the corpus."""
import argparse
import hashlib
import json
from html.parser import HTMLParser
from pathlib import Path
import zipfile

ROOT = Path(__file__).resolve().parents[1]
MAX_BYTES = 64 * 1024 * 1024

class PlainText(HTMLParser):
    def __init__(self):
        super().__init__(convert_charrefs=True)
        self.parts = []
        self.skip = 0
    def handle_starttag(self, tag, attrs):
        if tag in {'script', 'style', 'iframe', 'object', 'svg'}:
            self.skip += 1
        if tag in {'p', 'div', 'br', 'li', 'h2', 'h3', 'h4', 'blockquote'} and not self.skip:
            self.parts.append('\n')
    def handle_endtag(self, tag):
        if tag in {'script', 'style', 'iframe', 'object', 'svg'} and self.skip:
            self.skip -= 1
        if tag in {'p', 'div', 'li', 'h2', 'h3', 'h4', 'blockquote'} and not self.skip:
            self.parts.append('\n')
    def handle_data(self, value):
        if not self.skip:
            self.parts.append(value)

def unique_object(pairs):
    result = {}
    for key, value in pairs:
        if key in result:
            raise ValueError(f'duplicate key: {key}')
        result[key] = value
    return result

def convert(source):
    if len(source) > MAX_BYTES:
        raise ValueError('input exceeds 64 MiB')
    raw = json.loads(source, object_pairs_hook=unique_object)
    quran = json.loads((ROOT / 'quran/data/arabic.json').read_text())
    expected = {f'{ch}:{v["verse"]}' for ch, verses in quran.items() for v in verses}
    if not isinstance(raw, dict) or set(raw) != expected:
        raise ValueError('download must contain exactly the 6,236 bundled verse references')
    result = {}
    for key, value in raw.items():
        if isinstance(value, str):
            result[key] = value
            continue
        if not isinstance(value, dict) or not isinstance(value.get('text'), str):
            raise ValueError(f'invalid passage {key}')
        parser = PlainText()
        parser.feed(value['text'])
        text = '\n\n'.join(p.strip() for p in ''.join(parser.parts).splitlines() if p.strip())
        if not text:
            raise ValueError(f'empty passage {key}')
        result[key] = {'text': text}
        if 'ayah_keys' in value:
            if not isinstance(value['ayah_keys'], list) or any(not isinstance(v, str) for v in value['ayah_keys']):
                raise ValueError(f'invalid passage group {key}')
            result[key]['ayah_keys'] = value['ayah_keys']
    resolved = {}
    for key in result:
        target, seen = key, set()
        while isinstance(result.get(target), str):
            if target in seen:
                raise ValueError(f'cyclic pointer {key}')
            seen.add(target)
            target = result[target]
        if target not in result:
            raise ValueError(f'missing target {key}')
        resolved[key] = target
    for key, value in result.items():
        if isinstance(value, dict):
            for ref in value.get('ayah_keys', []):
                if resolved.get(ref) != key:
                    raise ValueError(f'inconsistent passage group {key}: {ref}')
    return result

def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('input', type=Path)
    args = ap.parse_args()
    if args.input.stat().st_size > MAX_BYTES:
        raise ValueError('input exceeds 64 MiB')
    downloaded = args.input.read_bytes()
    source = downloaded
    if zipfile.is_zipfile(args.input):
        with zipfile.ZipFile(args.input) as archive:
            files = [f for f in archive.infolist() if not f.is_dir()]
            if len(files) != 1 or not files[0].filename.endswith('.json') or files[0].file_size > MAX_BYTES:
                raise ValueError('expected one JSON file, at most 64 MiB')
            source = archive.read(files[0])
    result = convert(source)
    edition = {
        'id': 'jalalayn-en', 'title': 'Tafsir Al-Jalalayn (English)',
        'language': 'en', 'verified': True,
        'source_url': 'https://qul.tarteel.ai/resources/tafsir/english-tafsir-al-jalalayn',
        'download_sha256': hashlib.sha256(downloaded).hexdigest(),
        'rights': 'QUL reports no copyright information for this resource; not covered by the code licence.',
    }
    # All validation completes before either tracked file is changed. Review and
    # commit both files together; this is a build-time importer, never a server job.
    destination = ROOT / 'tafsir/data'
    for name, content in [('commentary.json', result), ('edition.json', edition)]:
        (destination / name).write_text(json.dumps(content, ensure_ascii=False, indent=2) + '\n')
    print('Imported Al-Jalalayn for all 6,236 verses. Run go test ./tafsir ./quran ./search.')

if __name__ == '__main__':
    main()
