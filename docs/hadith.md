# Hadith collections and search

The reader offers Bukhari and Muslim through Hadith settings. The selected
collection is part of the URL, so reloading, sharing and bookmarking preserve it.
Bukhari keeps `/hadith/{book}#number` and its existing bookmark keys. Muslim uses
`/hadith/muslim/{book}#number` and `muslim:book:number` bookmark keys. Switching
collections returns to its book list instead of treating unrelated book numbers
as equivalents. Narration references and bookmark/share controls sit below the
text, matching the Quran reader.

## Muslim source

Imported from `asim/hadith-json`, `db/by_book/the_9_books/muslim.json`, Git blob
`7941fc29ab2eb780406d4d9142da623f89466a5d`. The source README attributes its data to
Sunnah.com. The archive's checksum and numbering conventions are retained in
`hadith/muslim-source.json`; the import is reproduced with:

```sh
python3 scripts/import_muslim.py /path/to/muslim.json
```

All 7,459 records are retained in 57 groups, including 91 Introduction records in
chapter 0. The source's `idInBook` becomes `number`; `id` becomes `source_id`.
These are dataset identifiers, not asserted to be standard Sunnah.com references.
The reader labels them Records and model evidence carries the numbering scheme.
Record 6292 has Arabic but no English in the source; its Arabic is retained and
indexed, and the reader explicitly indicates the missing English translation.
No text-based deduplication or renumbering is applied to Muslim.

The normalized JSON is gzip-compressed deterministically and embedded. It is
loaded and decompressed once. The importer verifies the pinned source blob,
unique record IDs, chapter membership and full record coverage before writing.
Source text remains attributed to its original providers; this import does not
make a new licensing claim over the corpus.

## APIs and integrations

Existing `/api/hadith`, `/api/hadith/books` and `/api/hadith/{book}` remain Bukhari.
New routes expose `/api/hadith/collections`, `/api/hadith/muslim`,
`/api/hadith/muslim/books` and `/api/hadith/muslim/{book}`. Invalid book IDs return
404. MCP `get_hadith_books` and `get_hadith_book` accept optional `collection`
(`bukhari` by default, or `muslim`). Daily reminders remain Bukhari.

Both collections are included in offline text downloads, content revision
hashing, pre-rendered HTML and the sitemap. The private user bookmark state is
not included in generated pages.

## Search

SQLite FTS5 is always populated from the complete bundled corpus before serving
requests, with per-source fingerprints and transactional replacement. This also
updates installations whose database was already marked built before Muslim was
added. Hadith indexing includes both Arabic and English.

Search and summaries cover both collections regardless of which one the reader
is displaying. An explicit prefix, such as `muslim: fasting` or
`bukhari: fasting`, restricts retrieval to that collection. Provider credentials control generation, not corpus coverage.
An OpenAI key enables embedding retrieval alongside SQLite by default. Rankings
are combined with reciprocal rank fusion, with duplicate references merged and
SQLite retained if vector retrieval fails. `SEARCH_SEMANTIC=false` disables vectors. Existing vector maintenance commands remain available, but this import
does not regenerate the legacy embedding artifact. Keyword relevance still needs
answer-quality evaluation; corpus coverage alone does not establish answer quality.

## Explanations and questions

When a model is configured, each narration has the same Explain icon and panel
as Quran. Readers can request an overview or ask follow-up questions. The panel
has no Tafsir tab: the bundled Hadith records contain Arabic, translation and
narrator details, but no scholarly commentary collection. Answers are labelled
AI-generated and show their retrieved sources.

`POST /api/explain` and `/api/model/explain` accept a Hadith context:

```json
{"context":{"source":"hadith","collection":"muslim","book":0,"hadith":7369},"question":"What does this mean?"}
```

The server resolves the exact corpus record before ranked retrieval, retaining
collection, book and record identity on follow-ups. The existing research loop
can look up additional evidence. Attributed commentary must be verified against
retrieved sources; generated explanation is not presented as stored commentary.
Muslim book 0 is supported. Invalid records return 404 instead of substituting
search results.
