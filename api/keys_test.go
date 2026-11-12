package api

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestKeysAcrossProcessesAndRestart(t *testing.T) {
	dir := t.TempDir()
	a, err := OpenKeys(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	b, err := OpenKeys(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	id, secret, err := a.Create("test", 2, 1000000)
	if err != nil {
		t.Fatal(err)
	}
	got, err := b.Authenticate(secret)
	if err != nil || got != id {
		t.Fatalf("authenticate: %s %v", got, err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "model-keys.db"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte(secret)) {
		t.Fatal("plaintext secret persisted")
	}
	now := time.Now()
	var admitted atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if b.Reserve(id, now) == nil {
				admitted.Add(1)
			}
		}()
	}
	wg.Wait()
	if admitted.Load() != 2 {
		t.Fatalf("admitted %d", admitted.Load())
	}
	c, err := OpenKeys(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if err := c.Reserve(id, now); err != errKeyLimit {
		t.Fatalf("restart lost usage: %v", err)
	}
	if err := c.Reserve(id, now.Add(24*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := a.Revoke(id); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Authenticate(secret); err != errKeyUnauthorized {
		t.Fatalf("revoked: %v", err)
	}
	if _, err := b.Authenticate(secret + "x"); err != errKeyUnauthorized {
		t.Fatalf("invalid: %v", err)
	}
}
func TestKeyCLI(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer
	if err := RunKeys(dir, []string{"create", "--name", "demo", "--daily-limit", "1"}, &out); err != nil {
		t.Fatal(err)
	}
	parts := strings.Fields(out.String())
	id, secret := parts[1], parts[3]
	out.Reset()
	if err := RunKeys(dir, []string{"list"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), id) || strings.Contains(out.String(), secret) {
		t.Fatal("invalid listing")
	}
	if err := RunKeys(dir, []string{"revoke", id}, &out); err != nil {
		t.Fatal(err)
	}
	s, err := OpenKeys(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if _, err := s.Authenticate(secret); err != errKeyUnauthorized {
		t.Fatal("CLI revocation failed")
	}
}

func TestTokenAllowancePersistsAndResets(t *testing.T) {
	dir := t.TempDir()
	s, err := OpenKeys(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	id, _, err := s.Create("budget", 100, 1000)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	var allowed atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if s.ReserveTokens(id, 100, now) == nil {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if allowed.Load() != 10 {
		t.Fatalf("allowed %d reservations", allowed.Load())
	}
	other, err := OpenKeys(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	if err := other.ReserveTokens(id, 1, now); err == nil {
		t.Fatal("lost persisted token usage")
	}
	if err := other.ReserveTokens(id, 1000, now.Add(24*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := s.Revoke(id); err != nil {
		t.Fatal(err)
	}
	if err := other.ReserveTokens(id, 1, now.Add(48*time.Hour)); err == nil {
		t.Fatal("charged revoked key")
	}
}

func TestKeySchemaUpgrade(t *testing.T) {
	dir := t.TempDir()
	db, err := sql.Open("sqlite", filepath.Join(dir, "model-keys.db"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE model_keys (id TEXT PRIMARY KEY,name TEXT NOT NULL,hash TEXT UNIQUE NOT NULL,daily_limit INTEGER NOT NULL,day TEXT NOT NULL DEFAULT '',used INTEGER NOT NULL DEFAULT 0,revoked INTEGER NOT NULL DEFAULT 0); INSERT INTO model_keys(id,name,hash,daily_limit) VALUES('old','old','hash',20)`)
	db.Close()
	if err != nil {
		t.Fatal(err)
	}
	s, err := OpenKeys(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.ReserveTokens("old", 100, time.Now()); err != nil {
		t.Fatal(err)
	}
}
