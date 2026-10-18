package search

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// SyncSource atomically replaces one source when its content changes. Other
// sources and their completion status are preserved; no embeddings are used.
func (i *FTS) SyncSource(source string, documents []*Result) (bool, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, document := range documents {
		if document == nil || document.Metadata["source"] != source || document.Text == "" {
			return false, fmt.Errorf("search: invalid document for source %q", source)
		}
	}
	data, err := json.Marshal(documents)
	if err != nil {
		return false, err
	}
	fingerprint := fmt.Sprintf("%x", sha256.Sum256(data))
	key := "source:" + source
	var previous string
	if err := i.db.QueryRow(`SELECT value FROM meta WHERE key = ?`, key).Scan(&previous); err == nil && previous == fingerprint {
		return false, nil
	}
	tx, err := i.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	// FTS5 external-content tables need the original text when removing tokens.
	if _, err = tx.Exec(`INSERT INTO docs_fts(docs_fts, rowid, text)
		SELECT 'delete', id, text FROM docs WHERE json_extract(metadata, '$.source') = ?`, source); err != nil {
		return false, err
	}
	if _, err = tx.Exec(`DELETE FROM docs WHERE json_extract(metadata, '$.source') = ?`, source); err != nil {
		return false, err
	}
	for _, document := range documents {
		metadata, err := json.Marshal(document.Metadata)
		if err != nil {
			return false, err
		}
		result, err := tx.Exec(`INSERT INTO docs(text, metadata) VALUES (?, ?)`, document.Text, string(metadata))
		if err != nil {
			return false, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return false, err
		}
		if _, err = tx.Exec(`INSERT INTO docs_fts(rowid, text) VALUES (?, ?)`, id, document.Text); err != nil {
			return false, err
		}
	}
	var count int
	if err = tx.QueryRow(`SELECT COUNT(*) FROM docs`).Scan(&count); err != nil {
		return false, err
	}
	if _, err = tx.Exec(`INSERT OR REPLACE INTO meta(key, value) VALUES (?, ?)`, key, fingerprint); err != nil {
		return false, err
	}
	if i.built {
		if _, err = tx.Exec(`INSERT OR REPLACE INTO meta(key, value) VALUES (?, ?)`, metaDocCount, count); err != nil {
			return false, err
		}
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}
	i.count = count
	return true, nil
}
