# Reminder

A reminder to the whole world

## Overview

This is the english translation of the Quran made more accessible through a programmable interface. 

## Import

```go
import "github.com/asim/reminder/quran"
```

## Load

```go
q := quran.Load()

for _, chapter := range q {
  fmt.Println(chapter.Name)
  fmt.Println(strings.Join(chapter.Verses, "\n"))
}
```

## Markdown

To simply get the markdown

```go
md := quran.Markdown()

os.WriteFile("reminder.md", []byte(text), 0644)
```
