# Reminder

Quran, hadith and names of Allah in one app and API.

## Overview

The Reminder is an app and API for the Quran, Hadith and Names of Allah. It provides search and query answers using an LLM. 
RAG (Retrieval Augmented Generation) and references ensures answers are grounded in authentic Islamic texts. 

The goal is to consolidate these texts
into a single app and API, providing helpful daily reminders using web push notifications and leveraging AI as a tool for search summaries.

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
# requires pnpm
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
