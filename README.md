# Reminder

Quran, hadith and names of Allah all in one app and API.

## Overview

The Reminder is an API and app for the Quran, Hadith (Bukhari) and Names of Allah. It provides search summarisation using OpenAI and RAG context referencing. The goal is to consolidate these texts and 
information into a single API and app and 
leverage LLMs as a tool for searching. We
do not offload reasoning to LLMs but they 
are a new form of useful indexing for search.

## Contents

- [Features](#features)
- [Install](#install)
- [OpenAI Usage](#openai-usage)
- [Server](#server)
- [API](#api)
- [Web](#web)
- [Notes](#notes)
- [Community](#community)
- [Sources](#sources)

## Features

- Quran in English & Arabic
- Names of Allah & Meaning
- Hadith (Bukhari) in English
- Index & Search using GPT 4o
- RAG contextual referencing
- API to query LLM or Quran

## Install

Find the latest [release](https://github.com/asim/reminder/releases/latest)

Or Go install

```
go get github.com.com/asim/reminder@latest
```

## OpenAI Usage

Set the `OPENAI_API_KEY` value

```
export OPENAI_API_KEY=xxx
```

## Server

Run the http server 

```
reminder --serve
```

Go to [localhost:8080](https://localhost:8080)

## API

All queries are returned as JSON

- `/api/quran` - to get the entire quran
- `/api/names` - to get the list of names
- `/api/hadith` - to get the entire hadith
- `/api/search` - to get summarised answer
  * `q` param for the query
  * `POST` using `content-type` as `application/json`
  * `curl -d '{"q": "what is islam"}' http://localhost:8080/api/search`

See [`/api`](https://reminder.dev/api) for more details 

## App

The reminder bakes in a "lite" app by default. This can be replaced by a more featureful react app.

To build the react app

```
# requires pnpm
make setup
```

Build the app
```
make build
```

Pass the additional `--web` flag which replaces the lite app with the react app
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

## Community

Come chat on [Discord](https://discord.gg/F3xXRGbp9d)

## Sources

We have been requested to verify the sources of data

- [Arabic Quran](https://github.com/asim/quran-json-arabic)
- [English Quran](https://github.com/asim/quran-json)
- [Names of Allah](https://github.com/asim/99-Names-Of-Allah)
- [Sahih Bukhari](https://github.com/asim/bukhari)

Summarisation is provided by OpenAI but all sources of truth are authentic
