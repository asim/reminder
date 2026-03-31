package search

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreAndQuery(t *testing.T) {
	idx := New()
	defer idx.Close()

	err := idx.Store(map[string]string{"source": "quran", "chapter": "1"}, "In the name of Allah the most gracious the most merciful")
	if err != nil {
		t.Fatal(err)
	}
	err = idx.Store(map[string]string{"source": "quran", "chapter": "2"}, "This is the book about which there is no doubt a guidance for the righteous")
	if err != nil {
		t.Fatal(err)
	}

	results, err := idx.Query("merciful")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected results for 'merciful'")
	}
	if results[0].Metadata["source"] != "quran" {
		t.Errorf("expected source=quran, got %s", results[0].Metadata["source"])
	}
}

func TestSynonymExpansion(t *testing.T) {
	idx := New()
	defer idx.Close()

	idx.Store(map[string]string{"source": "test"}, "Allah is full of compassion and rahma towards believers")

	// "mercy" should match via synonym expansion to "compassion" and "rahma"
	results, err := idx.Query("mercy")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("synonym expansion should match 'mercy' to 'compassion'/'rahma'")
	}
}

func TestExpandSynonymsDedup(t *testing.T) {
	words := []string{"mercy", "compassion"}
	expanded := expandSynonyms(words)

	seen := make(map[string]bool)
	for _, w := range expanded {
		if seen[w] {
			t.Errorf("duplicate word in expansion: %s", w)
		}
		seen[w] = true
	}
}

func TestPersistentIndex(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	// Create and populate
	idx := New(dbPath)
	idx.Store(map[string]string{"source": "test"}, "patience and perseverance in the face of hardship")
	if idx.Built() {
		t.Error("fresh index should not report Built()")
	}
	if idx.Count() != 1 {
		t.Errorf("expected count=1, got %d", idx.Count())
	}
	idx.Close()

	// Reopen — should detect existing data
	idx2 := New(dbPath)
	defer idx2.Close()
	if !idx2.Built() {
		t.Error("reopened index should report Built()")
	}
	if idx2.Count() != 1 {
		t.Errorf("expected count=1 after reopen, got %d", idx2.Count())
	}

	results, err := idx2.Query("patience")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected results from persisted index")
	}
}

func TestInMemoryNotBuilt(t *testing.T) {
	idx := New()
	defer idx.Close()
	if idx.Built() {
		t.Error("in-memory index should not report Built()")
	}
}

func TestEmptyQuery(t *testing.T) {
	idx := New()
	defer idx.Close()
	results, err := idx.Query("")
	if err != nil {
		t.Fatal(err)
	}
	if results != nil {
		t.Error("empty query should return nil results")
	}
}

func TestPersistentFileCreated(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	idx := New(dbPath)
	idx.Store(map[string]string{}, "test content")
	idx.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("expected database file to be created on disk")
	}
}
