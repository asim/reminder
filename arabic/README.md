# Arabic learning evidence

The `arabic` package joins two separately attributed resources without treating
root meanings as verse interpretations:

- **Lane’s Arabic-English Lexicon:** historical dictionary entries, headwords
  and root families. The 49,473 unchanged entries remain in `dict/lane-xml.zip`.
  The original XML retains distinctions that a plain-text excerpt cannot display.
- **Quranic Arabic Corpus 0.4:** 77,429 words, 128,219 morphological segments,
  and all 6,236 numbered verses. Its roots, lemmas and grammar are linguistic
  annotations, not new translations. Particles and pronouns often have no root.

A *word* is the form in the verse. A *lemma* is its dictionary form. A *root*
connects related words; it is not a translation and does not determine which
sense applies. Prefixes, stems and suffixes retain their separate annotations.

## Use

- `GET /api/arabic/dictionary?q=صبر&limit=5&offset=0`: exact headword first,
  then entries under that root. Complete entries and total count; limit 1–25.
- `GET /api/arabic/roots?q=رح`: up to 100 roots beginning with the prefix.
- `GET /api/arabic/quran/1/2`: source words, segments, roots, lemmas and grammar,
  with the reader's canonical Arabic and an explicit `aligned` flag.
- MCP tools: `arabic_dictionary` and `quran_morphology`.
- Ask AI and `/api/model/ask`: explicit Arabic tokens or verse references add
  bounded lexical evidence through the shared retrieval implementation.
  For example: “Explain الحمد in 1:2: how do its root and lemma help?”

The model receives complete dictionary records up to 6,000 bytes each, at most
10 lexical records, with up to 12 distinct lookups and two verse analyses per
request. Longer entries remain available through the dictionary API. Selection
is deterministic and bounded, not a claim to have considered every possible
sense. Ordinary topical questions without Arabic or an explicit verse retain
existing retrieval. The reader Explain modal is not introduced by this change.

## Alignment and provenance

All 6,236 verses are available as source annotations. **5,408** have verified
written-token correspondence with the reader's edition. Verification compares
every word, not just counts. It ignores vocalization/recitation marks, alif wasla,
kasra-marked final undotted ya, and the documented hamza-plus-long-alif spelling
variants. It never alters either source or the reader's Arabic. It does not
conflate alif maqsura and ya in general, nor normalize dictionary letters.

The remaining **828** are flagged `aligned: false`; source positions must not
be attached to reader word positions. Their morphology is excluded from model
evidence. Matching source tokens are also not indexes into the existing grouped
English word-gloss array. That reader attachment remains a separate step.

The original morphology file is bundled byte-for-byte inside `morphology/qac.txt.gz`,
including its copyright blocks, and verified against the SHA-256 in `source.json`.
The mirror commit is pinned. See `morphology/ATTRIBUTION.txt`: the corpus identifies
GNU GPL with additional attribution/verbatim-copy terms and includes Tanzil text
with its own notice. Those data notices are retained separately from application
code. [Original corpus and terms](https://corpus.quran.com/download/).

QUL offers word-root/lemma resources but its pages currently say copyright
information is unavailable. This implementation uses the explicitly attributed
Quranic Arabic Corpus source instead. The earlier uploaded dataset has not been
identified or silently substituted.
