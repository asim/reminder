package api

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/asim/reminder/model"
)

type explainFlight struct {
	done chan struct{}
	data []byte
	err  error
}

// ExplainCache stores only default, public passage explanations. Flights are
// shared by all three routes in this server process; SQLite survives restarts.
type ExplainCache struct {
	db      *sql.DB
	version string
	mu      sync.Mutex
	flights map[string]*explainFlight
}

func OpenExplainCache(dir string) (*ExplainCache, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(executable)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	_, err = io.Copy(h, f)
	f.Close()
	if err != nil {
		return nil, err
	}
	// The executable includes prompts and all embedded Quran/tafsir/Arabic data.
	return openExplainCache(dir, hex.EncodeToString(h.Sum(nil)))
}

func openExplainCache(dir, version string) (*ExplainCache, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "explain-cache.db")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	f.Close()
	if err := os.Chmod(path, 0600); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(`PRAGMA busy_timeout=5000; CREATE TABLE IF NOT EXISTS explanations (key TEXT PRIMARY KEY, response BLOB NOT NULL, expires INTEGER NOT NULL); DELETE FROM explanations WHERE expires < unixepoch();`); err != nil {
		db.Close()
		return nil, err
	}
	return &ExplainCache{db: db, version: version, flights: map[string]*explainFlight{}}, nil
}
func (c *ExplainCache) Close() error { return c.db.Close() }

func (c *ExplainCache) key(r model.ExplainRequest) string {
	if c == nil || r.Context.Chapter < 1 || r.Question != "" || r.Q != "" || r.Context.Selection != "" || len(r.History) > 0 {
		return ""
	}
	r.Context.Source = "quran"
	if r.Context.Translation == "" {
		r.Context.Translation = "khattab"
	}
	data, _ := json.Marshal([]any{c.version, model.ExplanationVersion(), r})
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func (c *ExplainCache) result(ctx context.Context, key string, generate func() (*model.Result, error)) (*model.Result, error) {
	if key == "" {
		return generate()
	}
	// Coordinate both cache reads and misses, avoiding a lookup/generation race.
	c.mu.Lock()
	if flight := c.flights[key]; flight != nil {
		c.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-flight.done:
		}
		return decodeExplanation(flight.data, flight.err)
	}
	flight := &explainFlight{done: make(chan struct{})}
	c.flights[key] = flight
	c.mu.Unlock()
	defer func() { c.mu.Lock(); delete(c.flights, key); close(flight.done); c.mu.Unlock() }()
	var data []byte
	err := c.db.QueryRowContext(ctx, "SELECT response FROM explanations WHERE key=? AND expires>?", key, time.Now().Unix()).Scan(&data)
	if err == nil {
		if result, e := decodeExplanation(data, nil); e == nil && result.Answer != "" && len(result.References) > 0 {
			flight.data = data
			return result, nil
		}
	} else if err != sql.ErrNoRows && ctx.Err() == nil {
		log.Printf("Explain cache read: %v", err)
	}
	if ctx.Err() != nil {
		flight.err = ctx.Err()
		return nil, flight.err
	}
	result, err := generate()
	if err != nil {
		flight.err = err
		return nil, err
	}
	data, err = json.Marshal(result)
	if err != nil {
		flight.err = err
		return nil, err
	}
	flight.data = data
	if result.Complete && result.Answer != "" && len(result.References) > 0 && ctx.Err() == nil {
		if _, err := c.db.ExecContext(ctx, "DELETE FROM explanations WHERE expires<?", time.Now().Unix()); err != nil {
			log.Printf("Explain cache expiry: %v", err)
		}
		if _, err := c.db.ExecContext(ctx, "INSERT OR REPLACE INTO explanations(key,response,expires) VALUES(?,?,?)", key, data, time.Now().Add(30*24*time.Hour).Unix()); err != nil {
			log.Printf("Explain cache write: %v", err)
		}
	}
	return result, nil
}

func decodeExplanation(data []byte, err error) (*model.Result, error) {
	if err != nil {
		return nil, err
	}
	var result model.Result
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid cached explanation: %w", err)
	}
	return &result, nil
}
