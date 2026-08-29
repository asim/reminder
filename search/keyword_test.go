package search

import (
	"os"
	"path/filepath"
	"testing"
)

// verses is a small slice of the real corpus, in the archaic register the
// translations actually use.
var verses = []string{
	"O you who have believed, seek help through patience and prayer. Indeed, Allah is with the patient.",
	"And We will surely test you with something of fear and hunger and a loss of wealth and lives and fruits, but give good tidings to the patient.",
	"Do not grieve; indeed Allah is with us.",
	"So do not weaken and do not grieve, and you will be superior if you are believers.",
	"Unquestionably, by the remembrance of Allah hearts are assured.",
	"And whoever fears Allah, He will make for him a way out.",
	"Indeed, with hardship will be ease.",
	"And He found you lost and guided you.",
}

func seedVerses(t *testing.T, idx *Index) {
	t.Helper()
	for i, v := range verses {
		md := map[string]string{"source": "quran", "verse": string(rune('1' + i))}
		if err := idx.Store(md, v); err != nil {
			t.Fatalf("store %d: %v", i, err)
		}
	}
}

// TestKeywordSearchWithoutEmbedder is the guarantee that matters most: with no
// API key and no embedder, search still returns results. Before the move to
// FTS5 an install without OPENAI_API_KEY failed every query outright.
func TestKeywordSearchWithoutEmbedder(t *testing.T) {
	idx := New()
	defer idx.Close()
	seedVerses(t, idx)

	if idx.Embedder != nil {
		t.Fatal("expected no embedder on a bare index")
	}

	for _, q := range []string{"patience", "hardship", "hearts", "guided"} {
		res, err := idx.Query(q)
		if err != nil {
			t.Fatalf("query %q: %v", q, err)
		}
		if len(res) == 0 {
			t.Errorf("query %q returned no results", q)
		}
	}
}

// TestPorterStemming checks that inflected forms match. "grieving" and
// "grieved" should both find a verse containing "grieve".
func TestPorterStemming(t *testing.T) {
	idx := New()
	defer idx.Close()
	seedVerses(t, idx)

	for _, q := range []string{"grieving", "grieved", "believing", "fearing"} {
		res, err := idx.Query(q)
		if err != nil {
			t.Fatalf("query %q: %v", q, err)
		}
		if len(res) == 0 {
			t.Errorf("stemming failed: %q returned no results", q)
		}
	}
}

// TestEmotionalVocabulary covers the gap between how people search and how the
// translations are worded. Someone typing "anxiety" is looking for the verses
// about grief and fear, but that word does not appear in the text.
func TestEmotionalVocabulary(t *testing.T) {
	idx := New()
	defer idx.Close()
	seedVerses(t, idx)

	for _, q := range []string{
		"anxiety", "anxious", "depressed", "sad", "worried",
		"grief", "stress", "despair", "afraid", "lost",
	} {
		res, err := idx.Query(q)
		if err != nil {
			t.Fatalf("query %q: %v", q, err)
		}
		if len(res) == 0 {
			t.Errorf("query %q returned no results; the synonym table should reach the archaic wording", q)
		}
	}
}

// TestStaleSchemaRebuild verifies that an index built with an older tokenizer
// is discarded rather than silently reused, which would leave stemming off.
func TestStaleSchemaRebuild(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "search.db")

	// Build an index using the old schema (no porter tokenizer).
	old := New(path)
	if _, err := old.db.Exec(`DROP TABLE docs_fts`); err != nil {
		t.Fatalf("drop: %v", err)
	}
	if _, err := old.db.Exec(
		`CREATE VIRTUAL TABLE docs_fts USING fts5(text, content=docs, content_rowid=id)`,
	); err != nil {
		t.Fatalf("create legacy fts: %v", err)
	}
	seedVerses(t, old)
	old.Close()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected db at %s: %v", path, err)
	}

	// Reopening must not treat the stale index as built.
	reopened := New(path)
	defer reopened.Close()

	if reopened.Built() {
		t.Error("stale index was reused; it should be dropped and rebuilt")
	}

	// And it must be usable again once repopulated, with stemming active.
	seedVerses(t, reopened)
	res, err := reopened.Query("grieving")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(res) == 0 {
		t.Error("expected stemmed match after rebuild")
	}
}

// TestPartialIndexNotMarkedBuilt covers an interrupted first build. Previously
// any non-empty docs table counted as built, so a process killed midway
// through indexing came back up and skipped the rest of the corpus forever —
// permanently missing whole sections. This matters more under systemd, which
// restarts the service automatically after a crash.
func TestPartialIndexNotMarkedBuilt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "search.db")

	// Simulate a build that stored some documents and then died before
	// finishing: no MarkBuilt call.
	partial := New(path)
	seedVerses(t, partial)
	if partial.Built() {
		t.Fatal("index should not report built before MarkBuilt")
	}
	partial.Close()

	// Restarting must not treat the partial corpus as complete.
	reopened := New(path)
	defer reopened.Close()
	if reopened.Built() {
		t.Fatal("partial index was treated as built; the rest of the corpus would never be indexed")
	}
}

// TestMarkBuiltPersists is the companion: a completed build is reused.
func TestMarkBuiltPersists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "search.db")

	first := New(path)
	seedVerses(t, first)
	if err := first.MarkBuilt(); err != nil {
		t.Fatalf("mark built: %v", err)
	}
	want := first.Count()
	first.Close()

	reopened := New(path)
	defer reopened.Close()
	if !reopened.Built() {
		t.Fatal("completed index should be reused, not rebuilt")
	}
	if reopened.Count() != want {
		t.Errorf("expected %d documents, got %d", want, reopened.Count())
	}

	// And it must still be queryable without re-indexing.
	res, err := reopened.Query("patience")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(res) == 0 {
		t.Error("expected results from the reused index")
	}
}
