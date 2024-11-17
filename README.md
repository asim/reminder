# Reminder

A reminder to the whole world

## Overview

This is the reminder (dhikr) of Allah (God) through a programmable interface. Access the Quran in English, the names of Allah and more. 
Help spread the word of God and remind all why we're here. Use it via Go or as a HTTP server. Includes LLM summarisation for questions.

## Features

- Go Library
- Quran in English
- Names of Allah & Meaning
- Hadith (Bukhari) in English
- Index & Search using GPT 4o mini

## Import in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/asim/reminder.svg)](https://pkg.go.dev/github.com/asim/reminder)

For the Quran in English

```go
import "github.com/asim/reminder/quran"
```

For the names of Allah

```go
import "github.com/asim/reminder/names"
```

For the hadith (bukhari) of the prophet Muhammad (PBUH)

```go
import "github.com/asim/reminder/hadith"
```

### Load

```go
q := quran.Load()

for _, chapter := range q {
  fmt.Println(chapter.Name)
  fmt.Println(chapter.Verses)
}
```

Load the names

```go
n := names.Load()

for _, name := range n {
  fmt.Println(name.English)
  fmt.Println(name.Meaning)
}
```

## Build the Binary

To build the binary `reminder`

```
go install
```

## Index & Search

To index the content in a vector DB and create embeddings you'll need the `OPENAI_API_KEY` variable set

```
reminder --index
```

You can then export this via the `--export` flag to `$HOME/reminder.idx.gob.gz` and import via `--import`.

The `OPENAI_API_KEY` will be required for any search queries to craete embeddings for your query 

Note: By default an embedded index is included. You will need to replace ./index/data/reminder.idx.gob.gz and 
rebuild the binary if you want to replace the built-in version or use `--import` at runtime.

## LLM Answers

Using GPT 4o mini we're able to provide a summarised answer to each search query. The results will return 
the answer first and then the references from Quran, Hadith and Names used for the context of the query.

Provide `OPENAI_API_KEY` to make use of this in your own server.

## Serve HTTP

Run the http server on :8080 

```
reminder --serve
```

## Generate HTML

To regenerate the HTML files

```
reminder --generate
```

Make sure to rebuild the binary after as the files are embedded by default.
