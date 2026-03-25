# Reminder

Quran, hadith and names of Allah in one app and API.

## Overview

The Reminder is an app and API for the Quran, Hadith and Names of Allah. It provides full-text search across Islamic texts 
and optional LLM-powered answer summarisation using Fanar or Ollama. 

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
- Full-text search across all content
- Optional Fanar or Ollama LLM integration for answer summarisation
- API to query Quran, hadith, names
- Daily reminder web notifications

## Install

Find the latest [release](https://github.com/asim/reminder/releases/latest)

Or Go install

```
go get github.com.com/asim/reminder@latest
```

## Usage

**Quick Start** (no external dependencies required):

```bash
reminder --serve
```

Search works out of the box using full-text search with no external services needed.

**LLM Configuration** (optional - for AI-powered answer summarisation):

```bash
# Use Fanar API (Islamic-focused LLM)
export FANAR_API_KEY=xxx

# Or use local Ollama (default model: llama3.2)
# Install Ollama first: https://ollama.ai/
ollama pull llama3.2
export OLLAMA_LLM_MODEL=llama3.2

# Optional: use a different Ollama instance
export OLLAMA_BASE_URL=http://localhost:11434/v1
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
- `/api/search` - to search content and get results
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

File an issue to discuss 
