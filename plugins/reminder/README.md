# Reminder plugin

A portable remote-MCP plugin for Quran, hadith, Names of Allah, daily reminders
and Arabic understanding. The first version has no custom ChatGPT UI and adds no
new navigation to Reminder. Public source tools require no account.

## Build and test

From the repository root:

```sh
python3 scripts/package-plugin.py /tmp/reminder-plugin.zip
python3 scripts/test-plugin.py http://127.0.0.1:8080/mcp
# Optional official-client check (tested with mcp 1.29.0):
python3 -m pip install "mcp==1.29.0"
REMINDER_MCP_SMOKE=1 sh scripts/test-http.sh
```

Start Reminder with `go run . --serve` first. The smoke test uses source retrieval;
it does not spend model credits. Unit tests exercise model generation against a
local fixture. Run `go test ./api` and `sh scripts/test-http.sh` as well.

For private ChatGPT testing, deploy this branch, then go to ChatGPT Plugins →
Add custom MCP server, name it Reminder, use `https://reminder.dev/mcp`, and select
No authentication. Install it and invoke it in a Work chat. This package does not
create a registered plugin ID or silently install anything in your account.
The public deployment must contain `explain_quran` before testing that workflow.

## What it exposes

- Daily reminder and date-specific reading (no subscription or scheduling writes).
- Quran chapters/verses, including stored translations and attributed tafsir.
- Bukhari and Muslim books, Names of Allah, hybrid source search.
- Lane dictionary entries and QAC/QUL morphology with alignment/conflict metadata.
- `explain_quran`: the same grounded explanation and shared default cache as the
  reader. Requires a chapter context; omit verse for a bounded chapter overview.
  Follow-ups pass complete user/assistant pairs with the same context. It returns
  answer, references, a generated flag, a passage URL and optional chapter coverage. Default explanations are
  cached; custom questions and history are not persisted in the shared cache.

Server initialization instructions and tool descriptions tell ChatGPT to cite
sources, distinguish tafsir from generated interpretation, and avoid treating a
root or dictionary sense as a definitive verse meaning. They also distinguish
ChatGPT's own synthesis from an explicitly requested Reminder-generated answer.
These are instructions to the client, not a guarantee of its output quality.

Model requests use the existing anonymous/global allowances. A ChatGPT connector
may share an upstream IP with other readers, so review rate limits before a public
launch; no unverified user-ID header bypass is introduced. Reading remains useful
when model generation is unavailable or rate-limited.

## Submission status

The ZIP is prepared for private connection/testing, not yet approved or published.
Before public submission:

1. Deploy and run the smoke test against the public endpoint; test the supplied
   positive and negative review prompts in ChatGPT.
2. Add real, approved `privacyPolicyURL` and `termsOfServiceURL` values under
   `extensions.com.openai.interface`. Neither URL is invented in this package.
   The privacy notice must cover provider processing of explanation questions and
   selected text/history, operational logging, and cached public default answers.
3. Confirm the listing category and publisher identity in the submission portal,
   complete developer/domain verification, and upload the ZIP using With MCP.
4. Complete portal review checks and publish after approval. Directory discovery
   and subscriber growth are not guaranteed by installing the private plugin.

The icon is reused from Reminder's existing web icon. Instructions are supplied
by the MCP server; no separate skill or privileged lifecycle hook is required.
Only explicit package files are included in the ZIP, never datasets, credentials,
server configuration or the rest of the repository.

Official references (checked 2026-11-16):
- https://developers.openai.com/plugins/build/plugins
- https://developers.openai.com/plugins/build/mcp-server
- https://developers.openai.com/plugins/deploy/submission

## Transport compatibility

The server is stateless Streamable HTTP with JSON responses, no session ID and no
SSE stream (GET returns 405). It negotiates 2025-03-26, 2025-06-18 and 2025-11-25,
accepts client notifications without JSON-RPC responses, and supplies structured
object results alongside legacy text. Tools are annotated read-only and
non-destructive; search/explain are marked open-world because optional generation
can contact the configured model provider. All remaining reading tools use the
bounded bundled corpus or Reminder's stored daily archive.

Requests with an Origin header must use the server's own origin, ChatGPT or
reminder.dev. Additional browser origins can be explicitly allowed through
`MCP_ALLOWED_ORIGINS` (comma-separated origins); ordinary server-to-server clients
without Origin are supported. `MCP-Protocol-Version` is allowed in CORS preflight.

Validation in this change: portable manifest/MCP JSON schemas, official Python
MCP client 1.29.0 over a live local HTTP server, HTTP route smoke tests and Go tests.
This is not a claim that the public deployment has been updated, that the plugin
has been installed in ChatGPT, or that directory review has completed.
