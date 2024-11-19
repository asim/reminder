# Reminder

A reminder to the whole world

## Overview

This is the reminder (dhikr) of Allah (God) through a programmable interface. Access the Quran in English, the names of Allah and more. 
Help spread the word of God and remind all why we're here. Use it via Go or as a HTTP server. Includes LLM summarisation for questions.

## Features

- Read Quran in English
- Names of Allah & Meaning
- Hadith (Bukhari) in English
- Index & Search using GPT 4o mini

## Install

Find the latest [release](https://github.com/asim/reminder/releases/latest)

Or Go install

```
go get github.com.com/asim/reminder@latest
```

## LLM Usage

Set the `OPENAI_API_KEY` value

```
export OPENAI_API_KEY=xxx
```

## Search

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
