#!/usr/bin/env python3
"""Build and query an offline Lane index from the pinned, unmodified XML archive."""
import argparse
import hashlib
import json
from pathlib import Path
import sqlite3
import unicodedata
import xml.etree.ElementTree as ET
import zipfile

DIRECTORY = Path(__file__).resolve().parents[1] / 'arabic' / 'dict'


def normalize(text):
    return ''.join(c for c in unicodedata.normalize('NFC', text)
                   if not unicodedata.category(c).startswith('M')
                   and not c.isspace() and c != '\u0640')


def build(destination):
    manifest = json.loads((DIRECTORY / 'source.json').read_text())
    archive = DIRECTORY / 'lane-xml.zip'
    if hashlib.sha256(archive.read_bytes()).hexdigest() != manifest['archive_sha256']:
        raise ValueError('Lane archive checksum mismatch')
    if destination.exists():
        raise ValueError(f'{destination} already exists; choose a new --db path')
    db = sqlite3.connect(destination)
    try:
        db.execute('CREATE TABLE entries (id TEXT PRIMARY KEY, source_id TEXT, root TEXT, word TEXT, '
                   'root_key TEXT, word_key TEXT, text TEXT, source_file TEXT, xml TEXT)')
        with zipfile.ZipFile(archive) as source, db:
            if set(source.namelist()) != set(manifest['files']):
                raise ValueError('Unexpected archive contents')
            for filename, digest in manifest['files'].items():
                data = source.read(filename)
                if hashlib.sha256(data).hexdigest() != digest:
                    raise ValueError(f'Checksum mismatch: {filename}')
                tree = ET.fromstring(data)
                ordinal = 0
                for group in tree.findall('.//div2'):
                    head = group.find('head')
                    root = ''.join(head.itertext()).strip() if head is not None else ''
                    for entry in group.findall('entryFree'):
                        ordinal += 1
                        orth = entry.find('form/orth')
                        word = entry.get('key') or (''.join(orth.itertext()).strip() if orth is not None else '')
                        text = ' '.join(''.join(entry.itertext()).split())
                        db.execute('INSERT INTO entries VALUES (?,?,?,?,?,?,?,?,?)',
                                   (f'{filename}:{ordinal}', entry.get('id', ''), root, word, normalize(root),
                                    normalize(word), text, filename, ET.tostring(entry, encoding='unicode')))
            count = db.execute('SELECT count(*) FROM entries').fetchone()[0]
            if count != manifest['entries']:
                raise ValueError(f'Expected {manifest["entries"]} entries, found {count}')
            db.execute('CREATE INDEX roots ON entries(root_key)')
            db.execute('CREATE INDEX words ON entries(word_key)')
        print(json.dumps({'database': str(destination), 'entries': count}))
    except Exception:
        db.close()
        destination.unlink(missing_ok=True)
        raise
    finally:
        db.close()


def lookup(database, query, limit):
    db = sqlite3.connect(database.resolve().as_uri() + '?mode=ro', uri=True)
    db.row_factory = sqlite3.Row
    try:
        rows = db.execute('SELECT id,source_id,root,word,text,source_file FROM entries '
                          'WHERE word_key=? OR root_key=? ORDER BY id LIMIT ?',
                          (normalize(query), normalize(query), limit)).fetchall()
        print(json.dumps([dict(row) for row in rows], ensure_ascii=False, indent=2))
    finally:
        db.close()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--db', type=Path, default=DIRECTORY / 'lane.sqlite')
    parser.add_argument('--build', action='store_true')
    parser.add_argument('--lookup', help='Arabic headword or root; diacritics optional')
    parser.add_argument('--limit', type=int, default=10)
    args = parser.parse_args()
    if not 1 <= args.limit <= 100:
        parser.error('--limit must be between 1 and 100')
    if not args.build and not args.lookup:
        parser.error('choose --build or --lookup')
    if args.build:
        build(args.db)
    if args.lookup:
        lookup(args.db, args.lookup, args.limit)
