#!/usr/bin/env python3
"""Import the six supplied QUL SQLite archives without guessing or repairing annotations.
Usage: python3 scripts/import_qul.py /path/to/uploads
Outputs lossless field values grouped by chapter, plus provenance/checksums.
"""
import hashlib,json,sqlite3,sys,tempfile,zipfile
from pathlib import Path

ROOT=Path(__file__).resolve().parents[1]
EXPECTED={
 'ayah-lemma':'59c88d4bdf7d174cbae1330ae0357ff31e854f8d4b64db5775b6349ff409397d',
 'ayah-root':'d9ef78e0c5f5a5c5c5b496aad879e0a23bc0b6fab9f0ec92eadaed1674e7fb2d',
 'ayah-stem':'3b70021da6904109c03c2c56444b0586964b9cdb1a09d384ee0ebddae00b3634',
 'word-lemma':'8fd431f0b765d9a66e8ad584edbfc19bb2c62091bb1809a97666e47648f4a356',
 'word-root':'a85c325c669fdf3295f8c150f9fbe780391ec9fc3c29b2c8b37b93083038db69',
 'word-stem':'57992661ebd776d2400fb129562ab931a7ea2211a8cd5b757d72e59556924784'}

def main(folder):
 chapters={};manifest={'source':'Quranic Universal Library (user-supplied exports)','source_url':'https://qul.tarteel.ai/resources/morphology','licence':'Not specified in the supplied archives; no new licence asserted.','archives':{},'chapters':{}}
 with tempfile.TemporaryDirectory() as scratch:
  for name,digest in EXPECTED.items():
   archive=folder/(name+'.db.zip');raw=archive.read_bytes()
   if hashlib.sha256(raw).hexdigest()!=digest:raise ValueError('Unexpected source revision: '+name)
   with zipfile.ZipFile(archive) as z:
    if z.namelist()!=[name+'.db']:raise ValueError('Unexpected archive members')
    data=z.read(name+'.db')
   dbpath=Path(scratch)/(name+'.db');dbpath.write_bytes(data)
   with sqlite3.connect(dbpath.as_uri()+'?mode=ro',uri=True) as db:
    if db.execute('PRAGMA integrity_check').fetchone()[0]!='ok':raise ValueError('Corrupt SQLite')
    level,kind=name.split('-');table=kind+'s';index=['root','lemma','stem'].index(kind)
    if level=='word':
     col='arabic_trilateral' if kind=='root' else 'text'
     orphan=db.execute(f'SELECT count(*) FROM {kind}_words w LEFT JOIN {table} t ON t.id=w.{kind}_id WHERE t.id IS NULL').fetchone()[0]
     if orphan:raise ValueError('Orphan links: '+name)
     rows=db.execute(f'SELECT w.word_location,t.{col} FROM {kind}_words w JOIN {table} t ON t.id=w.{kind}_id').fetchall()
    else:rows=db.execute(f'SELECT verse_key,text FROM {table}').fetchall()
    seen=set()
    for location,value in rows:
     if location in seen:raise ValueError('Duplicate location: '+location)
     seen.add(location);parts=location.split(':')
     if len(parts)!=(3 if level=='word' else 2) or any(not p.isdigit() or int(p)<1 for p in parts):raise ValueError('Invalid location')
     ch,verse=map(int,parts[:2])
     if ch>114:raise ValueError('Invalid chapter')
     record=chapters.setdefault(ch,{}).setdefault(str(verse),{'ayah':[None,None,None],'words':{}})
     target=record['ayah'] if level=='ayah' else record['words'].setdefault(str(int(parts[2])),[None,None,None])
     target[index]=value
    manifest['archives'][archive.name]={'sha256':digest,'database_sha256':hashlib.sha256(data).hexdigest(),'rows':len(rows)}
 out=ROOT/'arabic/qul';out.mkdir(exist_ok=True)
 if len(chapters)!=114:raise ValueError('Incomplete chapters')
 for ch,verses in sorted(chapters.items()):
  data=(json.dumps(verses,ensure_ascii=False,sort_keys=True,separators=(',',':'))+'\n').encode()
  name=f'{ch}.json';(out/name).write_bytes(data);manifest['chapters'][name]=hashlib.sha256(data).hexdigest()
 (out/'source.json').write_text(json.dumps(manifest,ensure_ascii=False,indent=2)+'\n')
 print('Imported all six archives into',len(chapters),'chapters; original values preserved.')
if __name__=='__main__':main(Path(sys.argv[1]))
