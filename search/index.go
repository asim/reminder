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
	built bool
}

// Result represents a single search result.
type Result struct {
	Text     string            `json:"text"`
	Score    float32           `json:"score"`
	Metadata map[string]string `json:"metadata"`
}

// synonyms maps common Islamic terms to related words so that a search
// for one concept also matches documents using alternative vocabulary.
var synonyms = map[string][]string{
	"mercy":         {"merciful", "compassion", "compassionate", "forgiveness", "pardon", "rahma"},
	"merciful":      {"mercy", "compassion", "compassionate", "rahma"},
	"forgiveness":   {"forgive", "pardon", "mercy", "repentance", "tawbah"},
	"forgive":       {"forgiveness", "pardon", "mercy", "repentance"},
	"repentance":    {"repent", "forgiveness", "tawbah", "pardon"},
	"repent":        {"repentance", "forgiveness", "tawbah"},
	"prayer":        {"salah", "salat", "worship", "pray"},
	"pray":          {"prayer", "salah", "salat", "worship"},
	"salah":         {"prayer", "salat", "worship"},
	"fasting":       {"fast", "sawm", "ramadan"},
	"fast":          {"fasting", "sawm", "ramadan"},
	"charity":       {"sadaqah", "zakat", "alms", "giving"},
	"zakat":         {"charity", "sadaqah", "alms"},
	"pilgrimage":    {"hajj", "umrah"},
	"hajj":          {"pilgrimage", "umrah"},
	"patience":      {"patient", "sabr", "perseverance", "steadfast"},
	"patient":       {"patience", "sabr", "perseverance"},
	"sabr":          {"patience", "patient", "perseverance"},
	"faith":         {"iman", "belief", "believe", "trust"},
	"iman":          {"faith", "belief", "believe"},
	"belief":        {"faith", "iman", "believe"},
	"righteous":     {"righteousness", "piety", "taqwa", "good"},
	"righteousness": {"righteous", "piety", "taqwa"},
	"taqwa":         {"piety", "righteousness", "righteous", "god-consciousness"},
	"piety":         {"taqwa", "righteous", "righteousness", "devout"},
	"worship":       {"prayer", "ibadah", "devotion"},
	"sin":           {"sins", "transgression", "wrongdoing", "evil"},
	"sins":          {"sin", "transgression", "wrongdoing"},
	"paradise":      {"jannah", "heaven", "garden", "gardens"},
	"jannah":        {"paradise", "heaven", "garden"},
	"heaven":        {"paradise", "jannah", "garden"},
	"hell":          {"jahannam", "hellfire", "fire", "punishment"},
	"jahannam":      {"hell", "hellfire", "fire"},
	"angel":         {"angels", "malaika"},
	"angels":        {"angel", "malaika"},
	"prophet":       {"prophets", "messenger", "messengers", "rasul"},
	"prophets":      {"prophet", "messenger", "messengers"},
	"messenger":     {"prophet", "messengers", "rasul"},
	"knowledge":     {"learn", "wisdom", "ilm", "understanding"},
	"wisdom":        {"knowledge", "wise", "hikma"},
	"death":         {"die", "dying", "grave", "hereafter", "akhira"},
	"hereafter":     {"akhira", "afterlife", "death", "judgment"},
	"judgment":      {"judgement", "reckoning", "hereafter", "account"},
	"justice":       {"just", "fairness", "equity"},
	"truth":         {"truthful", "honest", "honesty", "true"},
	"grateful":      {"gratitude", "thankful", "thanks", "shukr"},
	"gratitude":     {"grateful", "thankful", "shukr"},
}

// expandSynonyms adds synonyms for each query word, returning a deduplicated list.
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

// New creates an SQLite FTS5 search index. If path is non-empty the index
// is persisted to disk; otherwise it lives in memory. When a file already
// exists with data, Built() returns true so the caller can skip re-indexing.
func New(path ...string) *Index {
	dsn := ":memory:"
	if len(path) > 0 && path[0] != "" {
		dsn = path[0]
	}

	existed := false
	if dsn != ":memory:" {
		if _, err := os.Stat(dsn); err == nil {
			existed = true
		}
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(fmt.Sprintf("search: open sqlite: %v", err))
	}

	// Performance pragmas
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA synchronous=NORMAL")

	idx := &Index{db: db}

	if existed {
		var cnt int
		if err := db.QueryRow(`SELECT COUNT(*) FROM docs`).Scan(&cnt); err == nil && cnt > 0 {
			idx.count = cnt
			idx.built = true
			return idx
		}
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS docs (id INTEGER PRIMARY KEY AUTOINCREMENT, text TEXT NOT NULL, metadata TEXT NOT NULL)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS docs_fts USING fts5(text, content=docs, content_rowid=id)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			panic(fmt.Sprintf("search: init db: %v", err))
		}
	}

	return idx
}

// Built reports whether the index was loaded from an existing persisted file.
func (i *Index) Built() bool {
	return i.built
}

// Count returns the number of indexed documents.
func (i *Index) Count() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.count
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
		return fmt.Errorf("search: prepare insert: %w", err)
	}
	defer insertDoc.Close()

	insertFTS, err := tx.Prepare(`INSERT INTO docs_fts (rowid, text) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("search: prepare fts insert: %w", err)
	}
	defer insertFTS.Close()

	for _, c := range content {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}

		res, err := insertDoc.Exec(c, string(mdJSON))
		if err != nil {
			return fmt.Errorf("search: insert doc: %w", err)
		}
		rowid, _ := res.LastInsertId()

		if _, err := insertFTS.Exec(rowid, c); err != nil {
			return fmt.Errorf("search: insert fts: %w", err)
		}
		i.count++
	}

	return tx.Commit()
}

// Query performs a full-text search and returns up to 25 results ranked by BM25.
func (i *Index) Query(q string) ([]*Result, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	words := tokenize(q)
	if len(words) == 0 {
		return nil, nil
	}

	expanded := expandSynonyms(words)

	// Build FTS5 query: OR together all terms for broad recall
	for j, w := range expanded {
		expanded[j] = `"` + w + `"`
	}
	ftsQuery := strings.Join(expanded, " OR ")

	rows, err := i.db.Query(`
		SELECT d.text, d.metadata, bm25(docs_fts) AS rank
		FROM docs_fts f
		JOIN docs d ON d.id = f.rowid
		WHERE docs_fts MATCH ?
		ORDER BY rank
		LIMIT 25
	`, ftsQuery)
	if err != nil {
		return nil, fmt.Errorf("search: query: %w", err)
	}
	defer rows.Close()

	var results []*Result
	for rows.Next() {
		var text, mdJSON string
		var rank float64
		if err := rows.Scan(&text, &mdJSON, &rank); err != nil {
			return nil, fmt.Errorf("search: scan: %w", err)
		}

		md := make(map[string]string)
		json.Unmarshal([]byte(mdJSON), &md)

		// BM25 returns negative scores (lower = better match).
		// Convert to a 0-1 similarity score for API compatibility.
		score := float32(1.0 / (1.0 - rank))

		results = append(results, &Result{
			Text:     text,
			Score:    score,
			Metadata: md,
		})
	}

	return results, rows.Err()
}

// Close closes the underlying database connection.
func (i *Index) Close() error {
	return i.db.Close()
}
