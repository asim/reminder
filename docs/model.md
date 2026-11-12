# Model API

`POST /api/model/ask` requires authorisation: `Authorization: Bearer <key>`.
For reminder.dev, request access through a [GitHub issue](https://github.com/asim/reminder/issues).
An operator approves access and shares the key privately. Never post keys in an issue.
Self-hosted installations issue their own keys; no accounts or admin UI are required.

## Issue and revoke keys

Run as the same operating-system user as the server (same `HOME`):

```sh
reminder keys create --name "my-app" --daily-limit 20
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
Each reference contains `text`, `score` and `metadata`, as in search.

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
it is not evidence. Each call retrieves sources from the Quran, Bukhari and
Names material. Tafsir and unrelated index records are excluded. The model is
instructed to distinguish sources from interpretation and acknowledge missing
evidence. This is retrieval-grounded generation, not a trained religious model
or a guarantee of correct interpretation. The reusable operation is `model.Respond`;
existing search summaries continue to use `model.AskSafe`.

## Limits and operation

The prompt is limited to 4,000 UTF-8 bytes; prompt plus history to 12,000 bytes;
history to six complete user/assistant pairs; and the JSON body to 16 KiB.
Retrieved evidence is bounded to 8,000 bytes. Output uses `AI_MAX_OUTPUT_TOKENS`
(default 1,024). Longer conversations must be shortened by the caller.

Each key defaults to 20 admitted requests per UTC day. Attempts that fail after
admission still count, including provider failures or a full global allowance.
Keys replace the anonymous per-IP allowance for this endpoint, but share
`AI_DAILY_LIMIT` and `AI_CONCURRENCY_LIMIT` with public search summaries.
Key counters persist across restarts; global accounting supports one server process.
Scheduled internal generation remains exempt from public limits.

Missing/invalid/revoked keys return 401; malformed requests 400; exhausted
allowances or busy generation 429; unavailable indexing/model/storage 503.
Errors use `{"message":"..."}`. Use HTTPS and keep keys in a server-side client,
not embedded in a public website or app bundle. Configure the existing model
provider on the server as usual; keys grant access but do not supply inference.
