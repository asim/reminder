# Model API

`POST /api/model/ask` requires authorisation: `Authorization: Bearer <key>`.
For reminder.dev, request access through a [GitHub issue](https://github.com/asim/reminder/issues).
An operator approves access and shares the key privately. Never post keys in an issue.
Self-hosted installations issue their own keys; no accounts or admin UI are required.

## Issue and revoke keys

Run as the same operating-system user as the server (same `HOME`):

```sh
reminder keys create --name "my-app" --daily-limit 20 --daily-tokens 1000000
reminder keys list
reminder keys revoke <id>
```

Creation displays the key once. Only its SHA-256 hash is stored in
`$HOME/.reminder/model-keys.db`, alongside its ID, label, allowance and usage.
Keep this directory private and persistent. Commands work without starting the
server or indexing content. New keys and revocations apply without restarting;
revocation prevents new requests but does not cancel an already admitted one.
To rotate a key, issue a replacement, update the caller, then revoke the old ID.
The key authorises this endpoint only. It is separate from the server's model-provider key.

## Request

```sh
curl https://reminder.dev/api/model/ask \
  -H "Authorization: Bearer $REMINDER_API_KEY" \
  -H 'Content-Type: application/json' \
  -d '{"prompt":"Explain patience","history":[]}'
```

Returns `{"answer":"...","references":[...]}`. The answer is plain text/Markdown,
not rendered HTML; numbered citations refer to the returned reference array.
Each reference contains `text`, `score` and `metadata`, as in search. Quran and
hadith references resolve to the bundled English and Arabic texts; the Arabic is
in `metadata.arabic` and is included in the model’s evidence alongside the English.

For follow-ups, supply earlier turns in chronological order:

```json
{
  "prompt": "How would that apply to a disagreement?",
  "history": [
    {"role": "user", "content": "Explain patience"},
    {"role": "assistant", "content": "The previous answer"}
  ]
}
```

The server does not store conversations. History helps interpret the prompt;
it is not evidence. Each call retrieves sources from the Quran, Bukhari, Muslim and
Names material. Quran evidence includes three attributed translations and matching
English Al-Jalalayn tafsir passages from QUL. Unrelated index records and
unverified editions are excluded. See [tafsir.md](tafsir.md). The model is
instructed to distinguish sources from interpretation and acknowledge missing
evidence. This is retrieval-grounded generation, not a trained religious model
or a guarantee of correct interpretation. The reusable operation is `model.Respond`;
it remains an internal integration capability exposed through the existing keyed
API, without a new Ask UI. Search and MCP retain their existing summarisation
path. Exact verse references supplement topical retrieval rather than replacing
it, including when earlier conversation turns contain verse references.

## Limits and operation

Defaults (configurable on the server):

| Setting | Default |
| --- | --- |
| `MODEL_MAX_PROMPT_BYTES` | 32 KiB |
| `MODEL_MAX_INPUT_BYTES` (prompt + history) | 128 KiB |
| `MODEL_MAX_EVIDENCE_BYTES` | 64 KiB |
| `MODEL_MAX_OUTPUT_TOKENS` | 4,096 tokens |
| `SEARCH_MAX_EVIDENCE_BYTES` | 32 KiB |

History must contain complete user/assistant pairs, with no fixed turn count.
The JSON body allows six times the input budget plus 64 KiB for escaping and
structure. Evidence includes complete records only; records that do not fit
are skipped. Public search summaries retain `AI_MAX_OUTPUT_TOKENS` (1,024).
Choose lower limits if your configured model has a smaller context window.

Each key defaults to 20 admitted requests and 1,000,000 **budgeted tokens** per
UTC day. `--daily-tokens` sets the latter when issuing a key. This is conservative
preflight accounting, not provider billing: one unit per UTF-8 input byte
(including sources), framing overhead, plus the maximum output allowance.
Longer conversations consume more allowance. Reservations are atomic, persist
across restarts and are not refunded after failures or shorter responses.
Existing key databases receive the default token allowance automatically.

Keys replace the anonymous per-IP allowance for this endpoint, but share
`AI_DAILY_LIMIT` and `AI_CONCURRENCY_LIMIT` with public search summaries.
Admitted requests count even when later blocked by token/global limits or
provider failures. Global accounting supports one server process.
Scheduled internal generation remains exempt from public limits.

Missing/invalid/revoked keys return 401; malformed requests 400; exhausted
allowances or busy generation 429; unavailable indexing/model/storage 503.
Errors use `{"message":"..."}`. Use HTTPS and keep keys in a server-side client,
not embedded in a public website or app bundle. Configure the existing model
provider on the server as usual; keys grant access but do not supply inference.

## Hadith collections

Both collections are available through the complete SQLite index, independent of
provider keys. `SEARCH_SEMANTIC=true` enables optional additional vector matches;
a missing Muslim embedding never excludes Muslim from retrieval. Model evidence
resolves each result by collection, book and record number against the bundled
corpus, including Muslim's Introduction (book 0). Muslim references identify the
source dataset numbering rather than claiming standard Sunnah.com references.

## Arabic lexical evidence

Explicit Arabic words and verse references now retrieve bounded Lane dictionary
entries and verified Quranic Arabic Corpus morphology through the shared model
path. For example, ask `Explain الحمد in 1:2, including its root and lemma`.
References distinguish `lexicon` from `morphology`, Quran, translations and tafsir.
Lexical entries retain stable IDs and attribution. The model must distinguish
root families, dictionary senses and contextual interpretation; it must not infer
a verse's meaning from a root alone or invent an unannotated root.

See [Arabic data and API](../arabic/README.md) for lookup endpoints, evidence
budgets, source terms and the 404 verses whose token alignment remains unverified.

QUL exports now provide separately attributed comparisons of roots, lemmas and
stems. The model receives them only when their source positions and the reader
text both align. It must acknowledge relevant conflicts; a conflicting field
cannot silently select a dictionary entry. See `arabic/README.md` for coverage.

## Contextual Explain

`POST /api/model/explain` uses the same API-key authentication and allowances as
`/api/model/ask`. It shares `model.Explain` with the reader's `POST /api/explain`
and `POST /api/quran/{chapter}/{verse}/explain` routes. Reader requests use the
anonymous per-IP and global model allowances. All return `{answer, references}`
with numbered citations. No conversation is stored on the server.

```json
{
  "context": {
    "source": "quran",
    "chapter": 2,
    "verse": 153,
    "translation": "sahih",
    "selection": "الصبر"
  },
  "question": "What does patience mean here?",
  "history": []
}
```

The context is resolved from the bundled Quran, independently of search ranking.
Translations, verified tafsir, Arabic morphology and Lane evidence supplement it.
Follow-ups resend the same context and complete user/assistant pairs. History and
selected text are user input, not source evidence. The reader keeps the most recent
four pairs; changing passage or translation resets the discussion.

Omit `verse` for a chapter overview. Chapter text is placed before supplementary
commentary, within the existing evidence budget. Long chapters may be only partly
represented; the model is instructed to state that limitation and explain only
what the supplied verses support. Chapter responses also include `coverage` with
`included_verses` and `total_verses`, computed from the source text actually sent
to the model; the reader displays these counts. This is not an exhaustive chapter
analysis.

The verse URL supplies chapter/verse automatically and accepts an empty body;
conflicting body context is rejected. `/api/explain` also accepts legacy `{ "q":
"Explain Quran 2:153" }` input, now grounded in retrieved sources. Its old ungrounded
JSON output and the verse endpoint's stored-tafsir response are replaced by the
consistent `{answer, references}` response. Missing model configuration returns
503 rather than presenting stored tafsir as a generated explanation. Stored tafsir
remains in the verse's `comments` and `commentary_source` fields and the reader's
commentary toggle.

Reader Explain is available at chapter and verse level, including selected text
within a single verse. The panel labels its passage, keeps follow-up questions
there, and exposes the evidence behind numbered citations. This first reader
integration covers Quran; hadith and Names retain their existing reading UI.

### Shared default-explanation cache

The three Explain routes share a persistent SQLite cache at
`~/.reminder/explain-cache.db`. Only requests with an explicit passage and empty
`question`, `q`, selection and history qualify. The reader sends an empty question
when its default Explain button is used. Custom questions and all follow-ups
bypass the cache; their text and conversation history are never stored there.

The cache stores the answer, numbered reference records and chapter coverage as
one JSON snapshot. Its key includes passage, translation, model/provider endpoint,
model limits, prompt and a SHA-256 fingerprint of the running executable. Because
the Quran, translations, tafsir, QUL and Lane data are embedded, replacing any of
those sources or explanation code invalidates cached answers automatically. This
is deliberately conservative: other binary changes can invalidate the cache too.
The same binary keeps its cache across restarts. Set `EXPLAIN_CACHE_EPOCH` to a new
value to refresh answers after a provider changes a model behind an unchanged name.
Entries expire after 30 days.

Concurrent requests for the same default explanation share one generation within
the server process. If the initiating reader cancels, active waiters retry through
one new shared generation using the new owner's request context and allowance.
A waiting reader can still cancel independently. Cache hits avoid retrieval and model generation and do not
consume model-generation allowances; external API key authentication is still
required on every request. Only nonempty answers with evidence and a provider
completion marker (`stop` / `end_turn`) are persisted. Errors, truncated responses,
missing completion markers and cancelled generations are not cached. A cache
read/write failure is logged; successful generation can still be returned.
