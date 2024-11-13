# Reminder

A reminder to the whole world

## Overview

This is the reminder (dhikr) of Allah through a programmable interface. Access the Quran in English, the names of Allah and more. 
Help spread the word of God and remind all why we're here. 

## Import

For the Quran in English

```go
import "github.com/asim/reminder/quran"
```

For the names of Allah

```go
import "github.com/asim/reminder/names"
```

For the hadith of the prophet Muhammad (PBUH)

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

## Serve

Run the http server on :8080 

```
go run main.go
```
