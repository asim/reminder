package hadith

import "testing"

func TestMuslimPreservesSourceRecords(t *testing.T) {
	c := Muslim()
	if c != Muslim() || c.ID != "muslim" || len(c.Books) != 57 {
		t.Fatal("collection metadata")
	}
	if c.Get(0) == nil || c.Get(0).Name != "Introduction" || len(c.Get(0).Hadiths) != 91 {
		t.Fatal("introduction lost")
	}
	ids := map[int]bool{}
	sourceIDs := map[int]bool{}
	missing := 0
	for _, b := range c.Books {
		for _, h := range b.Hadiths {
			if ids[h.Number] || sourceIDs[h.SourceID] || h.Arabic == "" {
				t.Fatal("duplicate or missing source record")
			}
			ids[h.Number] = true
			sourceIDs[h.SourceID] = true
			if h.English == "" {
				missing++
				if h.Number != 6292 {
					t.Fatal("unexpected missing translation")
				}
			}
		}
	}
	if len(ids) != 7459 || missing != 1 {
		t.Fatalf("records=%d untranslated=%d", len(ids), missing)
	}
	if Select("unknown", Load()) != nil || Load().Get(0) != nil {
		t.Fatal("invalid collection/book accepted")
	}
}
