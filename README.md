# Reminder

A reminder to the whole world

## Overview

Reminder is a search for knowledge with the Quran, Hadith and names of Allah used for contextual reference when querying an LLM. We're trying to uphold Islamic values while in the pursuit of knowledge.

## Features

- Quran in English & Arabic
- Names of Allah & Meaning
- Hadith (Bukhari) in English
- Index & Search using GPT 4o
- Contextual query referencing
- API to query LLM or get Quran

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

- `/api/quran` - to get the entire quran as JSON
- `/api/search` - to get a summarised answer to a question
  * `q` param for the query
  * `POST` using `content-type` as `application/json`
  * `curl -d '{"q": "what is islam"}' http://localhost:8080/api/search`


## Notes

The Quran says in [6:90](https://quran.com/6:90)

```
Say, “I ask no reward of you for this ˹Quran˺—it is a reminder to the whole world.”
```

Therefore, we ask nothing in compensation. 
