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

// FTS is a keyword search index backed by SQLite FTS5. It is used when no
// OpenAI key is configured, so that search works without any external service.
type FTS struct {
	mu       sync.RWMutex
	db       *sql.DB
	count    int
	built    bool
}

// synonyms maps common Islamic terms to related words so that a search
// for one concept also matches documents using alternative vocabulary.
// ftsTokenizer is the FTS5 tokenizer. Porter stemming lets "grieving",
// "grieved" and "grieves" all match a query for "grieve". Changing this
// invalidates existing indexes, which are detected and rebuilt in New.
const ftsTokenizer = "porter unicode61"

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

	// Modern conversational vocabulary mapped onto the archaic register the
	// translations actually use. Someone typing "anxiety" or "I feel
	// depressed" finds nothing without this, and those are among the most
	// common things people search for.
	"anxiety":    {"anxious", "fear", "afraid", "grieve", "distress", "worry", "despair"},
	"anxious":    {"anxiety", "fear", "afraid", "grieve", "distress", "worry"},
	"depressed":  {"depression", "grieve", "grief", "sorrow", "despair", "distress", "sad"},
	"depression": {"depressed", "grieve", "grief", "sorrow", "despair", "distress"},
	"sad":        {"sadness", "grieve", "grief", "sorrow", "distress", "weep"},
	"sadness":    {"sad", "grieve", "grief", "sorrow", "distress"},
	"grief":      {"grieve", "sorrow", "sad", "mourn", "distress"},
	"grieve":     {"grief", "sorrow", "sad", "mourn", "distress"},
	"worried":    {"worry", "anxious", "fear", "afraid", "grieve"},
	"worry":      {"worried", "anxious", "fear", "afraid", "grieve"},
	"lonely":     {"loneliness", "alone", "forsaken", "abandoned", "solitude"},
	"loneliness": {"lonely", "alone", "forsaken", "abandoned"},
	"stress":     {"hardship", "distress", "burden", "difficulty", "affliction"},
	"despair":    {"hopeless", "grieve", "sorrow", "distress", "anxiety"},
	"hopeless":   {"despair", "grieve", "sorrow", "distress"},
	"hope":       {"hopeful", "tidings", "glad", "expectation", "mercy"},
	"afraid":     {"fear", "fearful", "frightened", "terror", "dread"},
	"fear":       {"afraid", "fearful", "frightened", "terror", "dread", "taqwa"},
	"suffering":  {"hardship", "affliction", "trial", "test", "distress", "pain"},
	"pain":       {"hardship", "affliction", "suffering", "distress", "hurt"},
	"hardship":   {"difficulty", "affliction", "trial", "adversity", "ease"},
	"struggle":   {"strive", "hardship", "jihad", "effort", "perseverance"},
	"angry":      {"anger", "wrath", "rage", "fury"},
	"anger":      {"angry", "wrath", "rage", "restrain"},
	"illness":    {"sick", "sickness", "disease", "healing", "cure", "ill"},
	"sick":       {"illness", "sickness", "disease", "healing", "cure"},
	"healing":    {"heal", "cure", "remedy", "shifa", "illness"},
	"poverty":    {"poor", "needy", "destitute", "hunger", "provision"},
	"poor":       {"poverty", "needy", "destitute", "hunger"},
	"wealth":     {"riches", "provision", "sustenance", "rizq", "money"},
	"money":      {"wealth", "riches", "provision", "sustenance", "rizq"},
	"lost":       {"astray", "guidance", "misguided", "straying"},
	"guidance":   {"guide", "guided", "straight", "path", "huda"},
	"doubt":      {"doubtful", "uncertain", "waver", "suspicion"},
	"parents":    {"mother", "father", "kindred", "family"},
	"family":     {"kindred", "relatives", "parents", "children", "household"},
	"marriage":   {"marry", "spouse", "wives", "wife", "husband", "nikah"},
	"children":   {"child", "offspring", "sons", "daughters", "progeny"},
	"trial":      {"test", "tested", "affliction", "tribulation", "hardship"},
	"test":       {"trial", "tested", "affliction", "tribulation"},
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
func NewFTS(path ...string) *FTS {
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

	idx := &FTS{db: db}

	if existed {
		// Reuse the existing index only if it was built with the current
		// schema and the build actually ran to completion. Older indexes used
		// the default unicode61 tokenizer with no stemming.
		var sqlText string
		err := db.QueryRow(
			`SELECT sql FROM sqlite_master WHERE type='table' AND name='docs_fts'`,
		).Scan(&sqlText)

		if err == nil && strings.Contains(sqlText, ftsTokenizer) {
			if cnt, ok := completedDocCount(db); ok {
				idx.count = cnt
				idx.built = true
				return idx
			}
		}

		// Stale schema, or a build that was interrupted partway through.
		// Either way the contents cannot be trusted, so start over rather
		// than serving a corpus that is silently missing whole sections.
		db.Exec(`DROP TABLE IF EXISTS docs_fts`)
		db.Exec(`DROP TABLE IF EXISTS docs`)
		db.Exec(`DROP TABLE IF EXISTS meta`)
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS docs (id INTEGER PRIMARY KEY AUTOINCREMENT, text TEXT NOT NULL, metadata TEXT NOT NULL)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS docs_fts USING fts5(text, content=docs, content_rowid=id, tokenize = '` + ftsTokenizer + `')`,
		`CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			panic(fmt.Sprintf("search: init db: %v", err))
		}
	}

	return idx
}

// Built reports whether the index was loaded from an existing persisted file.
// metaDocCount is the key under which a completed build records how many
// documents it wrote. Its presence is what marks an index as finished.
const metaDocCount = "indexed_docs"

// completedDocCount reports the document count of a finished build. It returns
// false when the marker is absent or disagrees with the rows actually present,
// which is what an interrupted build looks like.
func completedDocCount(db *sql.DB) (int, bool) {
	var marked int
	if err := db.QueryRow(
		`SELECT value FROM meta WHERE key = ?`, metaDocCount,
	).Scan(&marked); err != nil {
		return 0, false
	}

	var actual int
	if err := db.QueryRow(`SELECT COUNT(*) FROM docs`).Scan(&actual); err != nil {
		return 0, false
	}

	if marked == 0 || marked != actual {
		return 0, false
	}
	return actual, true
}

// Built reports whether the index holds a complete corpus. It is false for a
// build that was interrupted partway through, so that indexing runs again
// instead of permanently serving a partial corpus.
func (i *FTS) Built() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.built
}

// MarkBuilt records that indexing finished. Call it only once every source has
// been stored; until then a restart must redo the work.
func (i *FTS) MarkBuilt() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, err := i.db.Exec(
		`INSERT INTO meta (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		metaDocCount, i.count,
	); err != nil {
		return fmt.Errorf("search: mark built: %w", err)
	}

	i.built = true
	return nil
}

// Count returns the number of indexed documents.
func (i *FTS) Count() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.count
}

// Store adds documents with the given metadata into the search index.
func (i *FTS) Store(md map[string]string, content ...string) error {
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

// queryFTS performs a full-text search and returns up to N results ranked by BM25.
func (i *FTS) queryFTS(q string, limit int) ([]*Result, error) {
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
		LIMIT ?
	`, ftsQuery, limit)
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

// Query returns up to 25 keyword matches ranked by BM25.
func (i *FTS) Query(q string) ([]*Result, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.queryFTS(q, 25)
}

// DB returns the underlying database connection for direct queries.
func (i *FTS) DB() *sql.DB {
	return i.db
}

// Close closes the underlying database connection.
func (i *FTS) Close() error {
	return i.db.Close()
}
