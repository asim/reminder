# Tafsir and grounded answers

## Status

English Tafsir Al-Jalalayn from QUL is bundled for all 6,236 verses. It remains the default
reader commentary and is available as attributed
evidence to the model through the shared tafsir package.

QUL resource: https://qul.tarteel.ai/resources/tafsir/english-tafsir-al-jalalayn
The download requires sign-in. Its licence tab reports no copyright information;
the text is not covered by this repository's software licence.

The imported archive is `tafsir-al-jalalayn.json.zip`; its SHA-256 is recorded in
`tafsir/data/edition.json`. The imported 73:4 passage matches the resource preview.

## Import

Download the English Al-Jalalayn JSON/ZIP from the resource above, then run:

```sh
python3 scripts/import_tafsir.py /path/to/download.json.zip
go test ./tafsir ./quran ./search ./model ./api
```

The build-time importer validates all 6,236 verse keys against the bundled Arabic
corpus, pointer resolution, nonempty passages and supplied group membership. It
rejects duplicate keys, missing verses and cycles before writing either file.
HTML is converted to plain text; scripts and active content are discarded.
Review the imported edition, passages and attribution against QUL before
committing both corpus and manifest together. The download checksum is retained
in the manifest. `verified` means this source import has been accepted, not that
every interpretation is independently authenticated.

## Loading and indexing

`tafsir.Load()` parses the embedded corpus once at startup. This compact commentary
uses an in-memory verse-to-passage map and shared passage objects; it does not
need a second SQLite database or disk extraction. Quran's existing `comments`
field is populated through this package, with additive `commentary_source` and
`commentary_url` fields. The Arabic and translation corpora are unchanged.

The existing SQLite keyword index replaces its tafsir records transactionally
when their content hash changes. One record is indexed per shared passage.
Prebuilt vector hits are resolved back to current verse/passage data; their
semantic ranking still needs a separate embedding rebuild. Offline text hashing
includes the tafsir package, so a new corpus refreshes saved chapters.

## Model and UI

The keyed model API uses `model.Respond`. It remains an internal integration
capability while its user-facing placement is decided. Search retains its
**Summarise with AI** checkbox and existing summarisation operation; MCP summaries
use that same existing generation path. Neither is converted into an Ask feature.

A retrieved Quran verse supplies:

- The canonical Arabic and Mustafa Khattab translation.
- Separately attributed Saheeh International and Muhammad Asad translations.
- The matching Al-Jalalayn tafsir passage, with its author/work,
  source URL and complete shared-verse coverage.

Only the sources that fit the model evidence budget receive citation numbers.
Exact verse matches are merged with topical search results, then canonicalised
and deduplicated. This preserves relevant hadith and other evidence in mixed
questions and conversation follow-ups.

The existing reader commentary displays Al-Jalalayn directly with attribution.
There is no new Explain or Ask button. The existing verse explanation endpoint
returns stored tafsir with `generated: false`, without calling a model.

Differences between translations are not proof of alternative lexical meanings.
Dictionary/lemma integration remains separate future work.

## Reading editions

Reading settings and the Tafsir panel offer Al-Jalalayn, Al-Mukhtasar and Ibn Kathir
in English. The choice is independent of the Commentary visibility toggle,
saved locally, and carried by shared reading URLs. All three editions travel
in chapter/verse payloads and work offline after refreshing the chapter.

The original commentary is English Al-Mukhtasar, QUL resource
https://qul.tarteel.ai/resources/tafsir/266. The supplied
`abridged-explanation-of-the-quran.json` matches all 6,236 original entries exactly.
Its existing `original-en` ID is retained for saved settings and shared links.
The supplied `en-tafisr-ibn-kathir.json` is imported as `ibn-kathir-en`, preserving
shared-passage pointers and converting HTML to plain text. Its exact upstream
resource URL was not supplied, so no URL is guessed. Both import manifests
record the archive SHA-256; source texts are not relicensed by the code licence.

Al-Jalalayn remains the default and model grounding source. Selecting another
reading edition does not change model evidence or explanation caches.

The importer accepts `--stem`, `--id`, `--title` and `--source-url` for additional
editions. Use stem `original` and ID `original-en` for Al-Mukhtasar, or stem
`ibn-kathir` and ID `ibn-kathir-en` for Ibn Kathir. Pass an empty source URL
when the exact resource address is not known.

Quran API verse records now expose a `commentaries` map keyed by edition ID.
Each value contains `text`, `source`, `url` (when known), and `verified`.
Existing `comments`, `commentary_source`, and `commentary_url` retain their
Al-Jalalayn values for compatibility. Missing editions are shown as unavailable,
without substituting another source.

To add a reading edition, register its parsed collection in `tafsir.Editions()`
and its display option in `web/app/utils/tafsir-editions.ts`. Record attribution
and rights before presenting a source as verified.

## HTTP API

`GET /api/tafsir` lists the bundled editions with `id`, `title`, `language`,
`source_url` and `verified`. No API key or model generation is needed.

`GET /api/tafsir/{edition}/{chapter}/{verse}` returns `edition` metadata and
`ref`, `text`, and `ayahs`. Shared commentary retains the full passage and all
covered ayahs, even when requesting a later verse in that passage. Newlines in
`text` are preserved. Unknown editions or missing passages return 404; malformed
chapter or verse numbers return 400.

Example: `/api/tafsir/jalalayn-en/2/153`. Use edition IDs from the listing;
`original-en` identifies Al-Mukhtasar (English).
