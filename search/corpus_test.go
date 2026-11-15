package search

import (
	"errors"
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"path/filepath"
	"testing"
)

func TestMuslimAddedToExistingSQLiteAndModel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "search.db")
	idx := NewFTS(path)
	defer idx.Close()
	if err := idx.Store(map[string]string{"source": "bukhari"}, "old text"); err != nil {
		t.Fatal(err)
	}
	if err := idx.MarkBuilt(); err != nil {
		t.Fatal(err)
	}
	b := hadith.Load()
	m := hadith.Muslim()
	q := quran.Load()
	n := names.Load()
	if err := SyncCorpus(idx, q, n, b, m); err != nil {
		t.Fatal(err)
	}
	for _, c := range []*hadith.Collection{b, m} {
		var count int
		err := idx.DB().QueryRow(`SELECT COUNT(*) FROM docs WHERE json_extract(metadata,'$.source')=?`, c.ID).Scan(&count)
		want := 0
		for _, book := range c.Books {
			want += len(book.Hadiths)
		}
		if err != nil || count != want {
			t.Fatalf("%s got %d want %d: %v", c.ID, count, want, err)
		}
		if changed, err := SyncHadith(idx, c); err != nil || changed {
			t.Fatalf("unchanged %s reindexed: %v", c.ID, err)
		}
	}
	got, err := idx.Query("muslim: decree")
	if err != nil || len(got) == 0 {
		t.Fatalf("search: %v", err)
	}
	for _, result := range got {
		if result.Metadata["source"] != "muslim" {
			t.Fatal("collection filter leaked another source")
		}
	}
	stale := []*Result{{Text: "wrong", Metadata: map[string]string{"source": "muslim", "book_num": "0", "number": "7369"}}}
	canonical := CurrentHadithResults(stale, b)
	if len(canonical) != 1 || canonical[0].Text != m.Get(0).Hadiths[0].English {
		t.Fatal("public summaries received indexing text")
	}
	refs := EnrichReferences(stale, n, q, b)
	if len(refs) != 1 || refs[0].Text != m.Get(0).Hadiths[0].English || refs[0].Metadata["numbering"] != "dataset idInBook" || refs[0].Metadata["arabic"] == "" {
		t.Fatal("model lost Muslim identity or canonical text")
	}
}
func TestSemanticSupplementCannotHideSQLite(t *testing.T) {
	p := &Result{Text: "current", Metadata: map[string]string{"source": "muslim", "book_num": "1", "number": "1"}}
	query := WithSupplement(evidenceIndex{results: []*Result{p}}, evidenceIndex{err: errors.New("embedding unavailable")})
	refs, err := query.Query("faith")
	if err != nil || len(refs) != 1 || refs[0].Text != p.Text || refs[0].Metadata["number"] != p.Metadata["number"] {
		t.Fatal("lost SQLite source")
	}
	extra := &Result{Text: "older duplicate", Metadata: map[string]string{"source": "muslim", "book_num": "1", "number": "1"}}
	refs, err = WithSupplement(evidenceIndex{results: []*Result{p}}, evidenceIndex{results: []*Result{extra}}).Query("faith")
	if err != nil || len(refs) != 1 || refs[0].Text != p.Text || refs[0].Metadata["number"] != p.Metadata["number"] {
		t.Fatal("duplicate vector result displaced current source")
	}
}
