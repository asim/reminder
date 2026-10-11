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
- Hadith (Bukhari) in English
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

**Quick Start** (requires OpenAI API key for fast embeddings):

```bash
export OPENAI_API_KEY=xxx
reminder --serve
```

**Alternative: Fully Local** (slower, requires Ollama):

```bash
# Install and start Ollama (https://ollama.ai/)
ollama pull nomic-embed-text  # For embeddings
ollama pull gemma3:1b          # For LLM responses

# Don't set OPENAI_API_KEY - will use Ollama for both
reminder --serve
```

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

**Embedding Configuration**:

For **best performance**, use OpenAI embeddings (fast & cheap - $0.02/1M tokens):

```bash
export OPENAI_API_KEY=xxx  # Uses text-embedding-3-small (fast, 1536 dims)
```

For **local/offline** (slower, but no API costs):

```bash
# Install Ollama first: https://ollama.ai/
ollama pull nomic-embed-text

# Optional: use different model or instance
export OLLAMA_EMBED_MODEL=nomic-embed-text  # default
export OLLAMA_BASE_URL=http://localhost:11434/api  # default
```

**Migration Note**: If switching embedding providers, delete the old index cache:

```bash
rm -rf ~/.reminder/data/reminder.idx.gob.gz
```

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
- [Word by word](https://github.com/asim/quranwbw)

## Development

File an issue to discuss 

## Offline reading

The web app automatically saves the Quran (including Arabic, translation,
word-by-word text, transliteration and stored commentary), Bukhari and names
on this device in the background while online. The initial download needs a
connection and resumes after an interruption. Texts and the app shell survive reloads and browser restarts, subject
to the browser's storage limits and eviction policy.

Saved reading modes, bookmarks and unopened chapters work offline. Audio,
server search and AI features need a connection. Unsaved content shows a reconnect message
instead of waiting indefinitely. Previously viewed daily reminders are cached,
but new reminders require a connection.

`npm run build` generates the versioned service worker after building all route
bundles. Deploy the complete `web/build/client` output (as `make web`/`make build`
already do). Development mode does not install an offline cache. Production must
use HTTPS (or localhost) for service workers.
