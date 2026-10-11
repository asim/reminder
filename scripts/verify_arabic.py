#!/usr/bin/env python3
"""Compare every stored ayah with QuranEnc, without modifying either text.

Usage: python scripts/verify_arabic.py [directory-of-QuranEnc-responses]
Without a directory, downloads the 114 chapter responses using curl.
"""
import concurrent.futures
import hashlib
import json
from pathlib import Path
import subprocess
import sys

root = Path(__file__).resolve().parents[1]
path = root / 'quran/data/arabic.json'
stored = json.loads(path.read_text())


def compare(chapter):
    if len(sys.argv) > 1:
        raw = (Path(sys.argv[1]) / f'{chapter}.json').read_bytes()
    else:
        raw = subprocess.check_output([
            'curl', '--fail', '--silent', '--show-error', '--location',
            '--retry', '2', '--max-time', '45',
            f'https://quranenc.com/api/v1/translation/sura/english_mokhtasar/{chapter}',
        ])
    source = json.loads(raw)['result']
    local = stored[str(chapter)]
    assert len(local) == len(source), f'Chapter {chapter}: verse count differs'
    for number, (expected, actual) in enumerate(zip(source, local), 1):
        assert int(expected['sura']) == actual['chapter'] == chapter
        assert int(expected['aya']) == actual['verse'] == number
        assert expected['arabic_text'] == actual['text'], f'Arabic differs at {chapter}:{number}'
        assert '\ufffd' not in actual['text'], f'Replacement character at {chapter}:{number}'
    return len(local)


with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
    count = sum(pool.map(compare, range(1, 115)))
assert count == 6236, f'Expected 6236 numbered ayahs, found {count}'
print(f'{count} / 6236 ayahs match QuranEnc exactly (no Unicode normalization).')
print('arabic.json SHA-256:', hashlib.sha256(path.read_bytes()).hexdigest())
