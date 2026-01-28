package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
)

// Checkpoint tracks indexing progress for resumption
type Checkpoint struct {
	QuranChapter int  `json:"quran_chapter"`
	QuranVerse   int  `json:"quran_verse"`
	QuranDone    bool `json:"quran_done"`
	NamesDone    bool `json:"names_done"`
	HadithBook   int  `json:"hadith_book"`
	HadithNum    int  `json:"hadith_num"`
	HadithDone   bool `json:"hadith_done"`
	TafsirDone   bool `json:"tafsir_done"`
}

func getCheckpointPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "reminder-index-checkpoint.json")
}

func loadCheckpoint() *Checkpoint {
	cp := &Checkpoint{}
	data, err := os.ReadFile(getCheckpointPath())
	if err != nil {
		return cp
	}
	json.Unmarshal(data, cp)
	return cp
}

func saveCheckpoint(cp *Checkpoint) {
	data, _ := json.MarshalIndent(cp, "", "  ")
	os.WriteFile(getCheckpointPath(), data, 0644)
}

func clearCheckpoint() {
	os.Remove(getCheckpointPath())
}

func indexContent(idx *search.Index, md map[string]string, text string) {
	// index the documents
	lines := strings.Split(text, "\n")

	fmt.Println("Indexing: ", md["source"], md["chapter"], md["verse"])

	if err := idx.Store(md, lines...); err != nil {
		fmt.Println("Error indexing", err)
	}
}

func indexQuran(idx *search.Index, q *quran.Quran) {
	cp := loadCheckpoint()
	if cp.QuranDone {
		fmt.Println("Quran already indexed, skipping")
		return
	}

	fmt.Println("Indexing Quran")
	if cp.QuranChapter > 0 {
		fmt.Printf("Resuming from chapter %d, verse %d\n", cp.QuranChapter, cp.QuranVerse)
	}

	for _, chapter := range q.Chapters {
		// Skip already indexed chapters
		if chapter.Number < cp.QuranChapter {
			continue
		}

		for _, verse := range chapter.Verses {
			// Skip already indexed verses in resumed chapter
			if chapter.Number == cp.QuranChapter && verse.Number <= cp.QuranVerse {
				continue
			}

			indexContent(idx, map[string]string{
				"source":  "quran",
				"chapter": fmt.Sprintf("%v", chapter.Number),
				"verse":   fmt.Sprintf("%v", verse.Number),
				"name":    chapter.Name,
			}, verse.Text)

			// Save checkpoint every 50 verses
			if verse.Number%50 == 0 {
				cp.QuranChapter = chapter.Number
				cp.QuranVerse = verse.Number
				saveCheckpoint(cp)
			}
		}
	}

	cp.QuranDone = true
	saveCheckpoint(cp)
}

func indexNames(idx *search.Index, n *names.Names) {
	cp := loadCheckpoint()
	if cp.NamesDone {
		fmt.Println("Names already indexed, skipping")
		return
	}

	fmt.Println("Indexing Names")

	for _, name := range *n {
		indexContent(idx, map[string]string{
			"source":  "names",
			"meaning": name.Meaning,
			"english": name.English,
			"arabic":  name.Arabic,
		}, strings.Join([]string{name.Meaning, name.English, name.Description}, " - "))
	}

	cp.NamesDone = true
	saveCheckpoint(cp)
}

func indexTafsir(idx *search.Index, q *quran.Quran) {
	cp := loadCheckpoint()
	if cp.TafsirDone {
		fmt.Println("Tafsir already indexed, skipping")
		return
	}

	fmt.Println("Indexing Tafsir")

	for _, comment := range q.Commentary {
		indexContent(idx, map[string]string{
			"source":  "tafsir",
			"chapter": fmt.Sprintf("%v", comment.Chapter),
			"verse":   fmt.Sprintf("%v", comment.Verse),
		}, comment.Text)
	}

	cp.TafsirDone = true
	saveCheckpoint(cp)
}

func indexHadith(idx *search.Index, b *hadith.Collection) {
	cp := loadCheckpoint()
	if cp.HadithDone {
		fmt.Println("Hadith already indexed, skipping")
		return
	}

	fmt.Println("Indexing Hadith")
	if cp.HadithBook > 0 {
		fmt.Printf("Resuming from book %d, hadith %d\n", cp.HadithBook, cp.HadithNum)
	}

	for bookIdx, book := range b.Books {
		// Skip already indexed books
		if bookIdx+1 < cp.HadithBook {
			continue
		}

		for _, h := range book.Hadiths {
			// Skip already indexed hadiths in resumed book
			if bookIdx+1 == cp.HadithBook && h.Number <= cp.HadithNum {
				continue
			}

			indexContent(idx, map[string]string{
				"source":   "bukhari",
				"book":     book.Name,
				"book_num": fmt.Sprintf("%d", book.Number),
				"narrator": h.Narrator,
				"number":   fmt.Sprintf("%d", h.Number),
			}, h.English)

			// Save checkpoint every 100 hadiths
			if h.Number%100 == 0 {
				cp.HadithBook = bookIdx + 1
				cp.HadithNum = h.Number
				saveCheckpoint(cp)
			}
		}
	}

	cp.HadithDone = true
	saveCheckpoint(cp)
	// Clear checkpoint when all done
	clearCheckpoint()
}
