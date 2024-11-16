# Reminder

A reminder to the whole world

## Overview

This is the reminder (dhikr) of Allah (God) through a programmable interface. Access the Quran in English, the names of Allah and more. 
Help spread the word of God and remind all why we're here. Use it via Go or as a HTTP server.

## Import

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

## Load

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

Load the hadith

```go
hadith.Load()
```

## Render

Render in markdown

```go
md := quran.Markdown()
```

For the names

```go
md := names.Markdown()
```

For the hadith

```go
md := hadith.Markdown()
```

## Build

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
