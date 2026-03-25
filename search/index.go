package search

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
}

// Result represents a single search result.
type Result struct {
	Text     string            `json:"text"`
	Score    float32           `json:"score"`
	Metadata map[string]string `json:"metadata"`
}

// tokenize splits text into lowercase words suitable for FTS5 queries.
func tokenize(s string) []string {
	s = strings.ToLower(s)
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// New creates a new in-memory SQLite FTS5 search index.
func New() *Index {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(fmt.Sprintf("search: failed to open sqlite: %v", err))
	}

	// Create a regular table for metadata and an FTS5 virtual table for text.
	stmts := []string{
		`CREATE TABLE docs (id INTEGER PRIMARY KEY AUTOINCREMENT, text TEXT NOT NULL, metadata TEXT NOT NULL)`,
		`CREATE VIRTUAL TABLE docs_fts USING fts5(text, content=docs, content_rowid=id)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			panic(fmt.Sprintf("search: failed to init db: %v", err))
		}
	}

	return &Index{db: db}
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
// Returns the top 25 results ordered by relevance.
func (i *Index) Query(q string) ([]*Result, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	words := tokenize(q)
	if len(words) == 0 {
		return nil, nil
	}

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
