package search

import (
	"path/filepath"
	"testing"
)

func TestSyncSourceRefreshesOnlyChangedNames(t *testing.T) {
	path := filepath.Join(t.TempDir(), "search.db")
	idx := NewFTS(path)
	if err := idx.Store(map[string]string{"source": "quran"}, "unchangedverse"); err != nil {
		t.Fatal(err)
	}
	if err := idx.Store(map[string]string{"source": "names"}, "supersededlabel"); err != nil {
		t.Fatal(err)
	}
	if err := idx.MarkBuilt(); err != nil {
		t.Fatal(err)
	}
	idx.Close()
	idx = NewFTS(path)
	defer idx.Close()
	var verseID int
	if err := idx.DB().QueryRow(`SELECT id FROM docs WHERE text = 'unchangedverse'`).Scan(&verseID); err != nil {
		t.Fatal(err)
	}
	documents := []*Result{{Text: "correctedlabel", Metadata: map[string]string{"source": "names", "number": "68"}}}
	changed, err := idx.SyncSource("names", documents)
	if err != nil || !changed {
		t.Fatalf("refresh: %v %v", changed, err)
	}
	if !idx.Built() || idx.Count() != 2 {
		t.Fatalf("completion/count lost: %v %d", idx.Built(), idx.Count())
	}
	for query, want := range map[string]int{"supersededlabel": 0, "correctedlabel": 1, "unchangedverse": 1} {
		results, err := idx.Query(query)
		if err != nil || len(results) != want {
			t.Fatalf("%s: %v %v", query, results, err)
		}
	}
	var unchangedID, nameID int
	idx.DB().QueryRow(`SELECT id FROM docs WHERE text = 'unchangedverse'`).Scan(&unchangedID)
	idx.DB().QueryRow(`SELECT id FROM docs WHERE text = 'correctedlabel'`).Scan(&nameID)
	if unchangedID != verseID {
		t.Fatal("unrelated source was rebuilt")
	}
	idx.Close()
	idx = NewFTS(path)
	defer idx.Close()
	if !idx.Built() {
		t.Fatal("completion marker not persisted")
	}
	changed, err = idx.SyncSource("names", documents)
	if err != nil || changed {
		t.Fatalf("unchanged content rewritten: %v %v", changed, err)
	}
	var afterID int
	idx.DB().QueryRow(`SELECT id FROM docs WHERE text = 'correctedlabel'`).Scan(&afterID)
	if afterID != nameID {
		t.Fatal("unchanged names rewritten")
	}
	if _, err = idx.SyncSource("names", []*Result{{Text: "bad", Metadata: map[string]string{"source": "quran"}}}); err == nil {
		t.Fatal("accepted wrong source")
	}
	results, err := idx.Query("correctedlabel")
	if err != nil || len(results) != 1 {
		t.Fatal("failed update damaged existing names")
	}
}

func TestSyncSourceDoesNotMarkPartialCorpusComplete(t *testing.T) {
	idx := NewFTS()
	defer idx.Close()
	if _, err := idx.SyncSource("names", []*Result{{Text: "name", Metadata: map[string]string{"source": "names"}}}); err != nil {
		t.Fatal(err)
	}
	if idx.Built() {
		t.Fatal("names alone marked entire corpus complete")
	}
}
