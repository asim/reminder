# Reminder

Quran, hadith and names of Allah in one app and API.

## Overview

Reminder brings daily reminders, the Quran, hadith and the Names of Allah together
in one app and API. Reading, listening and returning to these texts remain the
core experience.

Search helps readers find passages; AI can add source-backed summaries and
explanations. Arabic language resources, translations and tafsir support deeper
understanding of the text. They should enrich reading without creating a separate
general chat or broad knowledge destination. See [project purpose](docs/product.md).

## Contents

- [Features](#features)
- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [Web](#web)
- [Notes](#notes)
- [Sources](#sources)

## Features

- Quran in English & Arabic
- Quran Audio Recitation (Arabic & English)
- Names of Allah & Meaning
- Sahih Bukhari and Sahih Muslim in Arabic and English
- Index & Search with RAG and LLMs
- Optional Claude, Ollama or OpenAI integration
- API to query Quran, hadith, names
- Daily reminder web notifications

## Install

Find the latest [release](https://github.com/asim/reminder/releases/latest)

Or Go install

```
go install github.com/asim/reminder@latest
```

## Usage

**Quick Start** (reading and SQLite search need no model API key):

```bash
reminder --serve --web
```

Search refreshes the bundled Quran, tafsir, Names, Bukhari and Muslim sources in
`~/.reminder/search.db` at startup. Collection updates do not require embeddings.

**LLM Configuration** (optional) - used for answer summaries:

```bash
# OpenAI (default: gpt-6-luna)
export OPENAI_API_KEY=xxx
export OPENAI_MODEL=gpt-6-luna  # optional override

# Use a specific Ollama model (default: gemma3:1b)
export OLLAMA_LLM_MODEL=gemma3:1b

# Use Claude when OpenAI and an explicit Ollama model are not configured
export ANTHROPIC_API_KEY=xxx
```

An explicit Ollama model takes priority, then OpenAI, then Anthropic.
OpenAI API credentials are required; a ChatGPT/Codex subscription is not an API key.

The reader links back to the full chapter at the selected verse.

Search offers All, Quran, Hadith, Tafsir and Names filters, applied before result
limits. `/api/search` accepts `source: "all" | "quran" | "hadith" | "tafsir" | "names"`.
The browser retains both searches and AI answers in Recent, separately by source
and mode; refresh restores the query and filter without regenerating a saved answer.

Ask AI and verse explanations use native tool calling with OpenAI-compatible
(including tool-capable Ollama) and Anthropic models. The model can answer directly
or repeatedly inspect lookup results and refine its investigation. Available tools
search the texts, retrieve a verse and its neighbours, scan Quran mentions, and
look up Arabic dictionary/morphology evidence. There is no mandatory planning call.

The prompt allows language knowledge, context and reasoning to support explanation.
Quotations, references, counts, exhaustive lists and attributed scholarly claims
must be verified against retrieved sources. Source numbers remain stable across
all tool rounds. Mention inventories count explicit whole-word matches for supplied
names/aliases, not thematic or unnamed references. General verse explanations
start with the passage's meaning, supported by context and tafsir where useful.

Research is bounded by a 120-second deadline, 1 MiB of conservative cumulative
request/output allowance, and an eight-turn loop backstop (the last turn is
answer-only). Existing per-client, key and global usage limits apply to every
provider call, including tool schemas and results. Local models must support
native tools; provider errors are surfaced rather than silently switching to an
unverified answer. No external web browsing is enabled.

Summaries are optional. With no model configured the summarise controls are
hidden and search still returns references.

**Optional semantic search supplement**:

```bash
export OPENAI_API_KEY=xxx
 # Enabled automatically; set SEARCH_SEMANTIC=false to disable
```

This combines the existing vector and current SQLite rankings using reciprocal
rank fusion. It is enabled automatically when an OpenAI key is configured. SQLite coverage is retained if embeddings are unavailable;
new collections do not need a vector rebuild to appear in search or model evidence.
`--index`, `--export` and `--import` remain explicit vector maintenance commands.
See [hadith collections](docs/hadith.md) for source numbering and collection URLs.

Run the server 

```
reminder --serve
```

Go to [localhost:8080](https://localhost:8080)

## API

All queries are returned as JSON

- `/api/latest` - for the latest reminder with LLM-generated contextual message
  * Returns: verse, hadith, name of Allah, links, and an AI-generated spiritual message
  * The `message` field contains an LLM-generated reflection (2-3 sentences) based on the verse, hadith, and name
  * Falls back to default message ("In the Name of Allah—the Most Beneficent, Most Merciful") if LLM is unavailable
- `/api/quran` - to get the entire quran
- `/api/names` - to get the list of names
- `/api/hadith` - to get the entire hadith
- `/api/search` - to get summarised answer
  * `q` param for the query
  * `POST` using `content-type` as `application/json`
  * `curl -d '{"q": "what is islam"}' http://localhost:8080/api/search`

See [`/api`](https://reminder.dev/api) for more details 

MCP clients can connect to `https://reminder.dev/mcp` using HTTP POST for
Quran, hadith, names, reminders and search tools. Search summaries are opt-in
with `summarise: true` (default: `false`).

## App

The reminder bakes in a "lite" app by default. This can be replaced by a featureful react app.

To build the react app

```
# requires Node.js 22 and npm
make setup
```

Build the app
```
make build
```

Pass the `--web` flag to use the react app
```
reminder --web --serve
```

## Notes

The Quran says in [6:90](https://reminder.dev/quran/6/90)

```
Say,
“I ask no reward of you for this (Quran) —
it is a reminder to the whole world.”
```

## Sources

See [the Arabic source audit](docs/arabic-audit.md) for verification scope and limitations.

The sources used for data

- [Arabic Quran](https://github.com/asim/quran-json-arabic)
- [English Quran](https://github.com/asim/quran-json)
- [Names of Allah](https://github.com/asim/99-Names-Of-Allah)
- [Sahih Bukhari](https://github.com/asim/bukhari)
- [Sahih Muslim JSON](https://github.com/asim/hadith-json)
- [Word by word](https://github.com/asim/quranwbw)

## Development

File an issue to discuss 

## Offline reading

Quran, Bukhari, Muslim and names are saved for offline reading after the initial download.
Audio, search and AI require a connection. Use `make build` to include the service worker.

## Model API

`POST /api/model/ask` requires authorisation: `Authorization: Bearer <key>`.
Operators issue keys on the CLI; see [model access and usage](docs/model.md).

## AI limits

Summaries require `summarise: true`. Public search defaults: 100 requests/day,
20 per IP, two concurrent generations and 1,024 output tokens. Configure with
`AI_DAILY_LIMIT`, `AI_IP_DAILY_LIMIT`, `AI_CONCURRENCY_LIMIT` and `AI_MAX_OUTPUT_TOKENS`.
Daily/hourly generation is exempt. Counters persist across restarts for a single server.
Nginx's `X-Real-IP` is trusted from loopback; set `TRUSTED_PROXY_CIDRS` for other proxy addresses.
