package main

import (
	"fmt"
	"strings"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
)

func indexContent(idx *search.Index, md map[string]string, text string) {
	// index the documents
	// TODO: use original json
	lines := strings.Split(text, "\n")

	fmt.Println("Indexing: ", md["source"], md["chapter"], md["verse"])

	if err := idx.Store(md, lines...); err != nil {
		fmt.Println("Error indexing", err)
	}
}

func indexQuran(idx *search.Index, q *quran.Quran) {
	fmt.Println("Indexing Quran")

	for _, chapter := range q.Chapters {
		for _, verse := range chapter.Verses {
			indexContent(idx, map[string]string{
				"source":  "quran",
				"chapter": fmt.Sprintf("%v", chapter.Number),
				"verse":   fmt.Sprintf("%v", verse.Number),
				"name":    chapter.Name,
			}, verse.Text)
		}
	}
}

func indexNames(idx *search.Index, n *names.Names) {
	fmt.Println("Indexing Names")

	for _, name := range *n {
		indexContent(idx, map[string]string{
			"source":  "names",
			"meaning": name.Meaning,
			"english": name.English,
			"arabic":  name.Arabic,
		}, strings.Join([]string{name.Meaning, name.English, name.Description}, " - "))
	}
}

func indexTafsir(idx *search.Index, q *quran.Quran) {
	fmt.Println("Indexing Tafsir")

	for _, comment := range q.Commentary {
		indexContent(idx, map[string]string{
			"source":  "tafsir",
			"chapter": fmt.Sprintf("%v", comment.Chapter),
			"verse":   fmt.Sprintf("%v", comment.Verse),
		}, comment.Text)
	}
}

func indexHadith(idx *search.Index, b *hadith.Collection) {
	fmt.Println("Indexing Hadith")

	for _, book := range b.Books {
		for _, h := range book.Hadiths {
			indexContent(idx, map[string]string{
				"source":   "bukhari",
				"book":     book.Name,
				"narrator": h.Narrator,
				"number":   fmt.Sprintf("%d", h.Number),
			}, h.English)
		}
	}
}
