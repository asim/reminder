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
  - High-quality Arabic recitation (Mishary Alafasy)
  - Modern English translation audio
  - Sequential playback (Arabic → English)
  - Full playback controls with accessibility support
- Names of Allah & Meaning
- Hadith (Bukhari) in English
- Index & Search with RAG and LLMs
- Optional Fanar or OpenAI integration
- API to query Quran, hadith, names
- Daily reminder web notifications

## Install

Find the latest [release](https://github.com/asim/reminder/releases/latest)

Or Go install

```
go get github.com.com/asim/reminder@latest
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
ollama pull llama3.2           # For LLM responses

# Don't set OPENAI_API_KEY - will use Ollama for both
reminder --serve
```

**LLM Configuration** (optional - defaults to local Ollama):

The app now uses **local Ollama by default** for LLM responses. You can override with:

```bash
# Use Fanar API (takes priority over Ollama/OpenAI)
export FANAR_API_KEY=xxx

# Use a specific Ollama model (default: llama3.2)
export OLLAMA_LLM_MODEL=llama3.1

# Use OpenAI as fallback (only if no Fanar key and OLLAMA_LLM_MODEL not set)
export OPENAI_API_KEY=xxx
```

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

- `/api/latest` - for the latest reminder
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

The sources used for data

- [Arabic Quran](https://github.com/asim/quran-json-arabic)
- [English Quran](https://github.com/asim/quran-json)
- [Names of Allah](https://github.com/asim/99-Names-Of-Allah)
- [Sahih Bukhari](https://github.com/asim/bukhari)
- [Word by word](https://github.com/asim/quranwbw)

## Development

Find me on [Discord](https://discord.gg/jwTYuUVAGh)
