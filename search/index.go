package search

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode"

	_ "modernc.org/sqlite"
)

// Index is a full-text search index backed by SQLite FTS5.
type Index struct {
	mu    sync.RWMutex
	db    *sql.DB
	count int
	built bool // true when the index was loaded from an existing file
}

// Result represents a single search result.
type Result struct {
	Text     string            `json:"text"`
	Score    float32           `json:"score"`
	Metadata map[string]string `json:"metadata"`
}

// synonyms maps common Islamic terms to related words so that a search
// for one concept also matches documents using alternative vocabulary.
// Only the canonical (map key) is looked up; values are extras to OR in.
var synonyms = map[string][]string{
	"mercy":        {"merciful", "compassion", "compassionate", "forgiveness", "pardon", "rahma"},
	"merciful":     {"mercy", "compassion", "compassionate", "rahma"},
	"forgiveness":  {"forgive", "pardon", "mercy", "repentance", "tawbah"},
	"forgive":      {"forgiveness", "pardon", "mercy", "repentance"},
	"repentance":   {"repent", "forgiveness", "tawbah", "pardon"},
	"repent":       {"repentance", "forgiveness", "tawbah"},
	"prayer":       {"salah", "salat", "worship", "pray"},
	"pray":         {"prayer", "salah", "salat", "worship"},
	"salah":        {"prayer", "salat", "worship"},
	"fasting":      {"fast", "sawm", "ramadan"},
	"fast":         {"fasting", "sawm", "ramadan"},
	"charity":      {"sadaqah", "zakat", "alms", "giving"},
	"zakat":        {"charity", "sadaqah", "alms"},
	"pilgrimage":   {"hajj", "umrah"},
	"hajj":         {"pilgrimage", "umrah"},
	"patience":     {"patient", "sabr", "perseverance", "steadfast"},
	"patient":      {"patience", "sabr", "perseverance"},
	"sabr":         {"patience", "patient", "perseverance"},
	"faith":        {"iman", "belief", "believe", "trust"},
	"iman":         {"faith", "belief", "believe"},
	"belief":       {"faith", "iman", "believe"},
	"righteous":    {"righteousness", "piety", "taqwa", "good"},
	"righteousness": {"righteous", "piety", "taqwa"},
	"taqwa":        {"piety", "righteousness", "righteous", "god-consciousness"},
	"piety":        {"taqwa", "righteous", "righteousness", "devout"},
	"worship":      {"prayer", "ibadah", "devotion"},
	"sin":          {"sins", "transgression", "wrongdoing", "evil"},
	"sins":         {"sin", "transgression", "wrongdoing"},
	"paradise":     {"jannah", "heaven", "garden", "gardens"},
	"jannah":       {"paradise", "heaven", "garden"},
	"heaven":       {"paradise", "jannah", "garden"},
	"hell":         {"jahannam", "hellfire", "fire", "punishment"},
	"jahannam":     {"hell", "hellfire", "fire"},
	"angel":        {"angels", "malaika"},
	"angels":       {"angel", "malaika"},
	"prophet":      {"prophets", "messenger", "messengers", "rasul"},
	"prophets":     {"prophet", "messenger", "messengers"},
	"messenger":    {"prophet", "messengers", "rasul"},
	"knowledge":    {"learn", "wisdom", "ilm", "understanding"},
	"wisdom":       {"knowledge", "wise", "hikma"},
	"death":        {"die", "dying", "grave", "hereafter", "akhira"},
	"hereafter":    {"akhira", "afterlife", "death", "judgment"},
	"judgment":     {"judgement", "reckoning", "hereafter", "account"},
	"justice":      {"just", "fairness", "equity"},
	"truth":        {"truthful", "honest", "honesty", "true"},
	"grateful":     {"gratitude", "thankful", "thanks", "shukr"},
	"gratitude":    {"grateful", "thankful", "shukr"},
}

// expandSynonyms adds synonyms for each query word. Returns a
// deduplicated list containing the original words plus any synonyms.
func expandSynonyms(words []string) []string {
	seen := make(map[string]struct{}, len(words)*2)
	expanded := make([]string, 0, len(words)*2)

	for _, w := range words {
		if _, ok := seen[w]; ok {
			continue
		}
		seen[w] = struct{}{}
		expanded = append(expanded, w)

		for _, syn := range synonyms[w] {
			if _, ok := seen[syn]; ok {
				continue
			}
			seen[syn] = struct{}{}
			expanded = append(expanded, syn)
		}
	}
	return expanded
}

// tokenize splits text into lowercase words suitable for FTS5 queries.
func tokenize(s string) []string {
	s = strings.ToLower(s)
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// New creates an SQLite FTS5 search index. If path is empty the index
// is created in memory; otherwise it is persisted to the given file.
// When a file already exists and contains a populated index, the
// returned Index has Built() == true so the caller can skip re-indexing.
func New(path ...string) *Index {
	dsn := ":memory:"
	if len(path) > 0 && path[0] != "" {
		dsn = path[0]
	}

	// Check whether the DB file already exists before opening it
	// (sql.Open will create the file if it doesn't exist).
	existed := false
	if dsn != ":memory:" {
		if _, err := os.Stat(dsn); err == nil {
			existed = true
		}
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(fmt.Sprintf("search: failed to open sqlite: %v", err))
	}

	idx := &Index{db: db}

	if existed {
		// Verify the existing DB has the expected tables.
		var cnt int
		if err := db.QueryRow(`SELECT COUNT(*) FROM docs`).Scan(&cnt); err == nil && cnt > 0 {
			idx.count = cnt
			idx.built = true
			return idx
		}
		// Table missing or empty — fall through and create it.
	}

	// Create tables for a fresh index.
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS docs (id INTEGER PRIMARY KEY AUTOINCREMENT, text TEXT NOT NULL, metadata TEXT NOT NULL)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS docs_fts USING fts5(text, content=docs, content_rowid=id)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			panic(fmt.Sprintf("search: failed to init db: %v", err))
		}
	}

	return idx
}

// Built reports whether the index was loaded from an existing file
// that already contained documents.
func (i *Index) Built() bool {
	return i.built
}

// Store adds documents with the given metadata into the search index.
func (i *Index) Store(md map[string]string, content ...string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	mdJSON, err := json.Marshal(md)
	if err != nil {
		return fmt.Errorf("search: marshal metadata: %w", err)
	}

	tx, err := i.db.Begin()
	if err != nil {
		return fmt.Errorf("search: begin tx: %w", err)
	}
	defer tx.Rollback()

	insertDoc, err := tx.Prepare(`INSERT INTO docs (text, metadata) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("search: prepare insert docs: %w", err)
	}
	defer insertDoc.Close()

	insertFTS, err := tx.Prepare(`INSERT INTO docs_fts (rowid, text) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("search: prepare insert fts: %w", err)
	}
	defer insertFTS.Close()

	for _, c := range content {
		if len(c) == 0 {
			continue
		}
		res, err := insertDoc.Exec(c, string(mdJSON))
		if err != nil {
			return fmt.Errorf("search: insert doc: %w", err)
		}
		id, _ := res.LastInsertId()
		if _, err := insertFTS.Exec(id, c); err != nil {
			return fmt.Errorf("search: insert fts: %w", err)
		}
		i.count++
	}

	return tx.Commit()
}

// Query searches the index using FTS5 MATCH with BM25 ranking.
// Synonyms for common Islamic terms are automatically included so that
// e.g. searching "forgiveness" also matches "mercy", "pardon", etc.
// Returns the top 25 results ordered by relevance.
func (i *Index) Query(q string) ([]*Result, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	words := tokenize(q)
	if len(words) == 0 {
		return nil, nil
	}

	// Expand with synonyms so conceptually related terms are included.
	words = expandSynonyms(words)

	// Build an FTS5 query: each word is matched with OR so partial matches
	// are included, and BM25 naturally ranks documents with more matches higher.
	ftsQuery := strings.Join(words, " OR ")

	rows, err := i.db.Query(
		`SELECT docs.text, docs.metadata, -rank AS score
		   FROM docs_fts
		   JOIN docs ON docs.id = docs_fts.rowid
		  WHERE docs_fts MATCH ?
		  ORDER BY rank
		  LIMIT 25`,
		ftsQuery,
	)
	if err != nil {
		return nil, fmt.Errorf("search: query: %w", err)
	}
	defer rows.Close()

	var results []*Result
	for rows.Next() {
		var text, mdJSON string
		var score float32
		if err := rows.Scan(&text, &mdJSON, &score); err != nil {
			return nil, fmt.Errorf("search: scan: %w", err)
		}

		var md map[string]string
		if err := json.Unmarshal([]byte(mdJSON), &md); err != nil {
			return nil, fmt.Errorf("search: unmarshal metadata: %w", err)
		}

		results = append(results, &Result{
			Text:     text,
			Score:    score,
			Metadata: md,
		})
	}

	return results, rows.Err()
}

// Count returns the total number of documents in the index.
func (i *Index) Count() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.count
}

// Close closes the underlying database.
func (i *Index) Close() error {
	if i.db != nil {
		return i.db.Close()
	}
	return nil
}
