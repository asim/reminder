# Reminder

A reminder to the whole world

## Overview

The Reminder is an API and web app for accessing and searching the Quran, Hadith and Names of Allah. It provides contextual search and summarisation using 
OpenAI and includes referenced context usage. The API provides straightforward access to the Quran and ability to query the LLM. More to follow soon.

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

All queries are returned as JSON

- `/api/quran` - to get the entire quran
- `/api/names` - to get the list of names
- `/api/hadith` - to get the entire hadith
- `/api/search` - to get summarised answer
  * `q` param for the query
  * `POST` using `content-type` as `application/json`
  * `curl -d '{"q": "what is islam"}' http://localhost:8080/api/search`


## Notes

The Quran says in [6:90](https://quran.com/6:90)

```
Say, “I ask no reward of you for this ˹Quran˺—it is a reminder to the whole world.”
```

Therefore, we ask nothing in compensation. 
