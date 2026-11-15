# Lane’s Arabic-English Lexicon

`lane-xml.zip` contains 38 original XML files (49,473 entries) from
[Alpheios](https://github.com/alpheios-project/lan), pinned to the commit in
`source.json`. The XML is unchanged. The manifest records every source-file
SHA-256 and the archive checksum. `ATTRIBUTION.txt` preserves the upstream notice:
Perseus Digital Library digitization, with corrections by the Alpheios Project.
The dataset is licensed **CC BY-SA 3.0 United States**; see the
[licence](https://creativecommons.org/licenses/by-sa/3.0/us/).
These dataset terms remain separate from the application code.

Build an optional local SQLite index, using only the Python standard library:

```sh
python3 scripts/lane_dictionary.py --build
python3 scripts/lane_dictionary.py --lookup كتب --limit 3
python3 scripts/lane_dictionary.py --lookup صبر --limit 3
```

The generated `lane.sqlite` is ignored by Git. `--db PATH` chooses another
location. Building checks the archive, each file, entry identities and total
count before accepting the index. Lookup normalizes diacritics, whitespace and
tatweel only; it does not guess roots or equate different Arabic letters.
The index retains source filename, entry ID, root, headword, text and XML.
The returned plain text is a lookup aid; the original XML retains typography
and sense boundaries needed for a future reader interface.

The Go `arabic` package now loads and verifies this archive for dictionary API/MCP lookup and bounded model evidence. Plain Search remains unchanged.
Lane gives lexical senses and historical examples; those do not establish the
meaning of a particular verse. Model evidence carries entry IDs and must
retain distinctions between senses and consider grammar, verse context,
translations and attributed tafsir. Generated English should be labelled an
explanation or interpretation, never a replacement for the Quranic Arabic.

## Quranic word mapping

[QUL morphology](https://qul.tarteel.ai/resources/morphology) supplies word and
ayah roots, lemmas and stems. Those annotations provide the bridge from a
specific Quran word to dictionary candidates. They are not an English dictionary.
The earlier user-uploaded morphology dataset has not been identified in this
checkout. The current integration instead uses the attributed Quranic Arabic
Corpus 0.4 with per-token verification; see [Arabic evidence](../README.md).
Do not infer word-position equivalence merely from matching counts.
