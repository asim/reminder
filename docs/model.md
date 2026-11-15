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
budgets, source terms and the 828 verses whose token alignment remains unverified.
