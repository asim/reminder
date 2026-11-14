# Tafsir and grounded answers

## Status

The package and model/UI integration are implemented. **The Al-Jalalayn corpus
replacement is pending the authenticated QUL download. Do not merge this draft
as a completed Al-Jalalayn migration.**

QUL resource: https://qul.tarteel.ai/resources/tafsir/english-tafsir-al-jalalayn
The download requires sign-in. Its licence tab reports no copyright information;
the text is not covered by this repository's software licence.

Until the actual download is imported, `tafsir/data/commentary.json` retains the
previous commentary verbatim, moved from `quran/data/tafsir.json`.
`edition.json` explicitly labels it unattributed and unverified. It remains
available to readers but is excluded from grounded model evidence. It must not
be silently relabelled as Al-Jalalayn (or Al-Mukhtasar based on spot checks).

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

The keyed model API, public Search → Ask, MCP summaries and verse explanations
share `model.Respond` (public requests use `RespondPublic` to retain the existing
public evidence/output budgets). The API key requirement on `/api/model/ask`
is unchanged. Public entry points retain IP/global/concurrency allowances.

A retrieved Quran verse supplies:

- The canonical Arabic and Mustafa Khattab translation.
- Separately attributed Saheeh International and Muhammad Asad translations.
- The matching tafsir passage once the edition is verified, with its author/work,
  source URL and complete shared-verse coverage.

Only the sources that fit the evidence budget receive citation numbers. The UI
shows the references returned with the answer, in exactly that order. Differences
between translations are not treated as proof of alternative lexical meanings.
Dictionary/lemma integration remains separate future work.

On a verse page, **Explain verse** requests a short grounded answer with expandable
sources. **Ask about this verse** prefills Search in Ask mode; it does not issue a
paid model request until submitted. With no model configured, Explain returns the
stored commentary with its actual source label. Provider and quota errors remain
visible; they are not silently presented as a generated explanation.
