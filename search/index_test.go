package search

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewAndClose(t *testing.T) {
	idx := New()
	if idx == nil {
		t.Fatal("expected non-nil index")
	}
	if err := idx.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}

func TestStoreAndCount(t *testing.T) {
	idx := New()
	defer idx.Close()

	md := map[string]string{"source": "test"}
	if err := idx.Store(md, "hello world", "foo bar"); err != nil {
		t.Fatalf("store: %v", err)
	}
	if got := idx.Count(); got != 2 {
		t.Fatalf("count = %d, want 2", got)
	}
}

func TestStoreSkipsEmpty(t *testing.T) {
	idx := New()
	defer idx.Close()

	md := map[string]string{"source": "test"}
	if err := idx.Store(md, "", "hello", ""); err != nil {
		t.Fatalf("store: %v", err)
	}
	if got := idx.Count(); got != 1 {
		t.Fatalf("count = %d, want 1", got)
	}
}

func TestQueryEmptyIndex(t *testing.T) {
	idx := New()
	defer idx.Close()

	results, err := idx.Query("hello")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestQueryEmptyString(t *testing.T) {
	idx := New()
	defer idx.Close()

	idx.Store(map[string]string{}, "hello world")

	results, err := idx.Query("")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if results != nil {
		t.Fatalf("expected nil results for empty query, got %v", results)
	}
}

func TestQueryBasicMatch(t *testing.T) {
	idx := New()
	defer idx.Close()

	idx.Store(map[string]string{"source": "quran"}, "In the name of Allah the most merciful")
	idx.Store(map[string]string{"source": "hadith"}, "The best of you are those who learn the Quran")
	idx.Store(map[string]string{"source": "names"}, "Ar-Rahman the most gracious")

	results, err := idx.Query("merciful")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least 1 result")
	}
	if results[0].Metadata["source"] != "quran" {
		t.Fatalf("expected source=quran, got %s", results[0].Metadata["source"])
	}
}

func TestQueryMultipleWords(t *testing.T) {
	idx := New()
	defer idx.Close()

	idx.Store(map[string]string{"id": "1"}, "prayer and fasting are pillars of Islam")
	idx.Store(map[string]string{"id": "2"}, "charity and kindness are beloved")
	idx.Store(map[string]string{"id": "3"}, "prayer is a duty for every believer")

	results, err := idx.Query("prayer Islam")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	// First result should be the one containing both words
	if results[0].Metadata["id"] != "1" {
		t.Fatalf("expected id=1 as top result (has both words), got %s", results[0].Metadata["id"])
	}
}

func TestQueryReturnsMetadata(t *testing.T) {
	idx := New()
	defer idx.Close()

	md := map[string]string{
		"source":  "quran",
		"chapter": "1",
		"verse":   "1",
	}
	idx.Store(md, "In the name of Allah the most gracious the most merciful")

	results, err := idx.Query("Allah")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.Metadata["source"] != "quran" {
		t.Errorf("source = %s, want quran", r.Metadata["source"])
	}
	if r.Metadata["chapter"] != "1" {
		t.Errorf("chapter = %s, want 1", r.Metadata["chapter"])
	}
	if r.Score <= 0 {
		t.Errorf("score = %f, want > 0", r.Score)
	}
}

func TestQueryLimit25(t *testing.T) {
	idx := New()
	defer idx.Close()

	for i := 0; i < 50; i++ {
		idx.Store(map[string]string{}, "the word of truth")
	}

	results, err := idx.Query("truth")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) > 25 {
		t.Fatalf("expected at most 25 results, got %d", len(results))
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello world", 2},
		{"Hello, World!", 2},
		{"", 0},
		{"   ", 0},
		{"one-two three", 3},
	}

	for _, tt := range tests {
		got := tokenize(tt.input)
		if len(got) != tt.want {
			t.Errorf("tokenize(%q) = %v (len %d), want len %d", tt.input, got, len(got), tt.want)
		}
	}
}

func TestSynonymExpansion(t *testing.T) {
	idx := New()
	defer idx.Close()

	// Only the word "mercy" appears in the document.
	idx.Store(map[string]string{"id": "1"}, "Allah is full of mercy towards His servants")
	idx.Store(map[string]string{"id": "2"}, "the sun rises in the east")

	// Search for "forgiveness" — should match via synonym expansion.
	results, err := idx.Query("forgiveness")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected synonym expansion to match 'mercy' via 'forgiveness'")
	}
	if results[0].Metadata["id"] != "1" {
		t.Fatalf("expected id=1, got %s", results[0].Metadata["id"])
	}
}

func TestExpandSynonymsDedup(t *testing.T) {
	words := expandSynonyms([]string{"mercy", "forgiveness"})
	seen := map[string]int{}
	for _, w := range words {
		seen[w]++
		if seen[w] > 1 {
			t.Fatalf("duplicate word %q in expanded list", w)
		}
	}
	// Both original words must be present.
	if seen["mercy"] != 1 {
		t.Error("missing 'mercy'")
	}
	if seen["forgiveness"] != 1 {
		t.Error("missing 'forgiveness'")
	}
}

func TestPersistentIndex(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	// Create and populate an index.
	idx := New(dbPath)
	if idx.Built() {
		t.Fatal("new DB should not be marked as built")
	}
	idx.Store(map[string]string{"source": "quran"}, "In the name of Allah the most merciful")
	if idx.Count() != 1 {
		t.Fatalf("count = %d, want 1", idx.Count())
	}
	idx.Close()

	// Verify the file was created.
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected DB file at %s: %v", dbPath, err)
	}

	// Re-open — should detect existing data and skip rebuild.
	idx2 := New(dbPath)
	defer idx2.Close()

	if !idx2.Built() {
		t.Fatal("reopened DB should be marked as built")
	}
	if idx2.Count() != 1 {
		t.Fatalf("reopened count = %d, want 1", idx2.Count())
	}

	// Queries should still work.
	results, err := idx2.Query("merciful")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results from persisted index")
	}
}

func TestInMemoryNotBuilt(t *testing.T) {
	idx := New()
	defer idx.Close()

	if idx.Built() {
		t.Fatal("in-memory index should not be marked as built")
	}
}
