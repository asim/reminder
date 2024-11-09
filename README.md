# Reminder

A reminder to the whole world

## Overview

This is the english translation of the Quran made more accessible through a programmable interface. 

## Load

```
import "github.com/asim/reminder/quran"

q := quran.Load()

for _, chapter := range q {
  fmt.Println(chapter.Name)
  fmt.Println(strings.Join(chapter.Verses, "\n"))
}
```

## Markdown

To compile it to markdown

```
go run main.go
```

An `index.md` file will be generated
