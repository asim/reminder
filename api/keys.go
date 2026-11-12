package api

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// KeyStore is shared by the running server and the operator's CLI. SQLite
// transactions allow issuance and revocation without restarting the server.
type KeyStore struct{ db *sql.DB }

func OpenKeys(dir string) (*KeyStore, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "model-keys.db")
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
	if _, err = db.Exec(`PRAGMA busy_timeout=5000;
 CREATE TABLE IF NOT EXISTS model_keys (
 id TEXT PRIMARY KEY, name TEXT NOT NULL, hash TEXT UNIQUE NOT NULL,
 daily_limit INTEGER NOT NULL, day TEXT NOT NULL DEFAULT '', used INTEGER NOT NULL DEFAULT 0,
 revoked INTEGER NOT NULL DEFAULT 0
 )`); err != nil {
		db.Close()
		return nil, err
	}
	return &KeyStore{db}, nil
}
func (s *KeyStore) Close() error { return s.db.Close() }
func keyHash(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}
func (s *KeyStore) Create(name string, limit int) (id, secret string, err error) {
	if strings.TrimSpace(name) == "" || len(name) > 100 || strings.ContainsAny(name, "\r\n\t") || limit < 1 {
		return "", "", fmt.Errorf("a name (1–100 bytes) and positive daily limit are required")
	}
	var random [32]byte
	if _, err = rand.Read(random[:]); err != nil {
		return
	}
	secret = "rm_" + hex.EncodeToString(random[:])
	id = keyHash(secret)[:16]
	_, err = s.db.Exec("INSERT INTO model_keys(id,name,hash,daily_limit) VALUES(?,?,?,?)", id, name, keyHash(secret), limit)
	return
}
func (s *KeyStore) Revoke(id string) error {
	res, err := s.db.Exec("UPDATE model_keys SET revoked=1 WHERE id=?", id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("key not found")
	}
	return nil
}

var errKeyUnauthorized = errors.New("invalid API key")
var errKeyLimit = errors.New("API key daily limit reached")

func (s *KeyStore) Authenticate(secret string) (string, error) {
	if len(secret) != 67 || !strings.HasPrefix(secret, "rm_") {
		return "", errKeyUnauthorized
	}
	var id string
	err := s.db.QueryRow("SELECT id FROM model_keys WHERE hash=? AND revoked=0", keyHash(secret)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errKeyUnauthorized
	}
	return id, err
}
func (s *KeyStore) Reserve(id string, now time.Time) error {
	day := now.UTC().Format("2006-01-02")
	res, err := s.db.Exec(`UPDATE model_keys SET used=CASE WHEN day=? THEN used+1 ELSE 1 END, day=? WHERE id=? AND revoked=0 AND (day<>? OR used<daily_limit)`, day, day, id, day)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errKeyLimit
	}
	return nil
}

// RunKeys executes operator commands before corpus indexing or model setup.
func RunKeys(dir string, args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: reminder keys create|list|revoke")
	}
	fs := flag.NewFlagSet("keys "+args[0], flag.ContinueOnError)
	fs.SetOutput(out)
	name := fs.String("name", "", "Application name")
	limit := fs.Int("daily-limit", 20, "Daily request allowance (UTC)")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	switch args[0] {
	case "create":
		if fs.NArg() != 0 {
			return fmt.Errorf("unexpected arguments")
		}
	case "list":
		if len(args) != 1 {
			return fmt.Errorf("usage: reminder keys list")
		}
	case "revoke":
		if len(args) != 2 {
			return fmt.Errorf("usage: reminder keys revoke <id>")
		}
	default:
		return fmt.Errorf("unknown keys command")
	}
	store, err := OpenKeys(dir)
	if err != nil {
		return err
	}
	defer store.Close()
	switch args[0] {
	case "create":
		id, secret, err := store.Create(*name, *limit)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(out, "ID: %s\nKey: %s\nSave this key now; it will not be shown again.\n", id, secret)
		return err
	case "revoke":
		return store.Revoke(fs.Arg(0))
	case "list":
		rows, err := store.db.Query("SELECT id,name,daily_limit,revoked FROM model_keys ORDER BY name,id")
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id, name string
			var limit, revoked int
			if err := rows.Scan(&id, &name, &limit, &revoked); err != nil {
				return err
			}
			status := "active"
			if revoked != 0 {
				status = "revoked"
			}
			if _, err := fmt.Fprintf(out, "%s\t%s\t%d/day\t%s\n", id, name, limit, status); err != nil {
				return err
			}
		}
		return rows.Err()
	}
	return nil
}
