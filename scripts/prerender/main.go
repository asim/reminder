// Export the same bundled records used by the API for deterministic web builds.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"os"
	"path/filepath"
)

func main() {
	root := ".generated"
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		panic(err)
	}
	write := func(name string, v any) {
		b, e := json.Marshal(v)
		if e != nil {
			panic(e)
		}
		if e = os.WriteFile(filepath.Join(root, name+".json"), b, 0644); e != nil {
			panic(e)
		}
	}
	q, b, n := quran.Load(), hadith.Load(), names.Load()
	write("chapters", q.Index().Chapters)
	write("books", b.Index().Books)
	write("names", n)
	embeds := []string{}
	paths := []string{"/", "/quran", "/hadith", "/names", "/api"}
	for _, c := range q.Chapters {
		write(fmt.Sprintf("quran-%d", c.Number), c)
		paths = append(paths, fmt.Sprintf("/quran/%d", c.Number))
		for _, v := range c.Verses {
			if v.Number > 0 {
				embeds = append(embeds, fmt.Sprintf("/embed/quran/%d/%d", c.Number, v.Number))
				paths = append(paths, fmt.Sprintf("/quran/%d/%d", c.Number, v.Number))
			}
		}
	}
	for _, book := range b.Books {
		write(fmt.Sprintf("hadith-%d", book.Number), book)
		paths = append(paths, fmt.Sprintf("/hadith/%d", book.Number))
	}
	for _, name := range *n {
		write(fmt.Sprintf("name-%d", name.Number), name)
		paths = append(paths, fmt.Sprintf("/names/%d", name.Number))
	}
	for _, book := range b.Books {
		for _, h := range book.Hadiths {
			if h.Number > 0 {
				embeds = append(embeds, fmt.Sprintf("/embed/hadith/%d/%d", book.Number, h.Number))
			}
		}
	}
	write("embeds", embeds)
	write("paths", paths)
}
