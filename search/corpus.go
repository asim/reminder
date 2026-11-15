package search

import (
	"fmt"
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"strconv"
)

func hadithDocument(c *hadith.Collection, book *hadith.Book, h *hadith.Hadith) *Result {
	source := c.ID
	if source == "" {
		source = "bukhari"
	}
	md := map[string]string{"source": source, "collection": source, "collection_name": c.Name, "book": book.Name, "book_num": fmt.Sprint(book.Number), "number": fmt.Sprint(h.Number), "narrator": h.Narrator, "arabic": h.Arabic}
	if source == "muslim" {
		md["source_url"] = "https://github.com/asim/hadith-json/blob/main/db/by_book/the_9_books/muslim.json"
		md["numbering"] = "dataset idInBook"
		md["source_id"] = fmt.Sprint(h.SourceID)
	}
	text := h.English
	if text == "" {
		text = h.Arabic
		md["english_available"] = "false"
	}
	return &Result{Text: text, Metadata: md}
}
func SyncHadith(fts *FTS, c *hadith.Collection) (bool, error) {
	var docs []*Result
	for _, book := range c.Books {
		for _, h := range book.Hadiths {
			doc := hadithDocument(c, book, h)
			if h.English != "" {
				doc.Text += "\n" + h.Arabic
			}
			docs = append(docs, doc)
		}
	}
	source := c.ID
	if source == "" {
		source = "bukhari"
	}
	return fts.SyncSource(source, docs)
}

// SyncCorpus refreshes each source independently, even on an already-built database.
func SyncCorpus(fts *FTS, q *quran.Quran, n *names.Names, collections ...*hadith.Collection) error {
	var docs []*Result
	for _, c := range q.Chapters {
		for _, v := range c.Verses {
			if v.Number > 0 {
				docs = append(docs, &Result{Text: v.Text, Metadata: map[string]string{"source": "quran", "chapter": fmt.Sprint(c.Number), "verse": fmt.Sprint(v.Number), "name": c.Name, "arabic": v.Arabic}})
			}
		}
	}
	if _, err := fts.SyncSource("quran", docs); err != nil {
		return err
	}
	if _, err := SyncNames(fts, n); err != nil {
		return err
	}
	if _, err := SyncTafsir(fts, q); err != nil {
		return err
	}
	for _, c := range collections {
		if _, err := SyncHadith(fts, c); err != nil {
			return err
		}
	}
	return fts.MarkBuilt()
}

// CurrentHadithResults returns the canonical record rather than the combined
// English/Arabic indexing text or a stale vector excerpt.
func CurrentHadithResults(results []*Result, b *hadith.Collection) []*Result {
	out := make([]*Result, 0, len(results))
	for _, result := range results {
		if result == nil {
			continue
		}
		source := result.Metadata["source"]
		if source != "bukhari" && source != "muslim" {
			out = append(out, result)
			continue
		}
		c := hadith.Select(source, b)
		if c == nil {
			continue
		}
		number, err := strconv.Atoi(result.Metadata["number"])
		if err != nil {
			continue
		}
		bookNumber, err := strconv.Atoi(result.Metadata["book_num"])
		if err != nil {
			continue
		}
		book := c.Get(bookNumber)
		if book == nil {
			continue
		}
		for _, h := range book.Hadiths {
			if h.Number == number {
				doc := hadithDocument(c, book, h)
				doc.Score = result.Score
				out = append(out, doc)
				break
			}
		}
	}
	return out
}
