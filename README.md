# Reminder

A reminder to the whole world

## Overview

This is the reminder (dhikr) of Allah (God) through a programmable interface. Access the Quran in English, the names of Allah and more. 
Help spread the word of God and remind all why we're here. Use it via Go or as a HTTP server. Includes LLM summarisation for questions.

## Features

- Quran in English
- Names of Allah & Meaning
- Hadith (Bukhari) in English
- Index & Search using GPT 4o mini

## Download && Install

Find the latest [release](https://github.com/asim/reminder/releases/latest)

Or Go install

To build the binary `reminder`

```
go get github.com.com/asim/reminder@latest
```

## Search

Set the `OPENAI_API_KEY` value

```
export OPENAI_API_KEY=xxx
```

Make a query

```
reminder search what did the prophet
```

## Server

Run the http server 

```
reminder --serve
```

Go to [localhost:8080](https://localhost:8080)
