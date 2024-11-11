# Reminder

A reminder to the whole world

## Overview

This is the reminder (dhikr) of Allah through a programmable interface. Access the Quran in English, the names of Allah and more. 
Help spread the word of God and remind all of why we're here. 

Integrate it into any piece of software you have.

## Import

For the Quran in English

```go
import "github.com/asim/reminder/quran"
```

For the names of Allah

```go
import "github.com/asim/reminder/names"
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

## Render

Render in markdown

```go
md := quran.Markdown()

os.WriteFile("reminder.md", []byte(text), 0644)
```

For the names

```go
md := names.Markdown()
```

## Serve

Run the http server on :8080 

```
go run main.go
```
