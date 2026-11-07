package model

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestQuotaPersistenceAndRollover(t *testing.T) {
	dir := t.TempDir()
	l := &limits{path: filepath.Join(dir, "ai-usage.json"), daily: 3, perClient: 1, concurrent: 1}
	now := time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	ctx := WithClient(context.Background(), "192.0.2.1")
	release, err := l.acquire(ctx, now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = l.acquire(WithClient(context.Background(), "192.0.2.2"), now); err == nil {
		t.Fatal("concurrency bypass")
	}
	release()
	if _, err = l.acquire(ctx, now); err == nil {
		t.Fatal("IP quota bypass")
	}
	t.Setenv("AI_DAILY_LIMIT", "3")
	t.Setenv("AI_IP_DAILY_LIMIT", "1")
	old := budget
	defer func() { budget = old }()
	if err = ConfigureLimits(dir); err != nil {
		t.Fatal(err)
	}
	if _, err = budget.acquire(ctx, now); err == nil {
		t.Fatal("restart reset IP allowance")
	}
	for _, ip := range []string{"192.0.2.2", "192.0.2.3"} {
		done, err := budget.acquire(WithClient(context.Background(), ip), now)
		if err != nil {
			t.Fatal(err)
		}
		done()
	}
	if _, err = budget.acquire(WithClient(context.Background(), "192.0.2.4"), now); err == nil {
		t.Fatal("global quota bypass")
	}
	done, err := budget.acquire(ctx, now.Add(24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	done()
	data, _ := os.ReadFile(l.path)
	if strings.Contains(string(data), "192.0.2.") {
		t.Fatal("raw IP persisted")
	}
}

func TestQuotaConcurrentReservations(t *testing.T) {
	l := &limits{path: filepath.Join(t.TempDir(), "usage.json"), daily: 5, perClient: 100, concurrent: 100}
	var wg sync.WaitGroup
	var mu sync.Mutex
	accepted := 0
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			done, err := l.acquire(context.Background(), time.Now())
			if err == nil {
				defer done()
				mu.Lock()
				accepted++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if accepted != 5 {
		t.Fatalf("accepted %d, want 5", accepted)
	}
}

func TestQuotaFailsClosed(t *testing.T) {
	l := &limits{path: filepath.Join(t.TempDir(), "missing", "usage.json"), daily: 5, perClient: 5, concurrent: 1}
	if _, err := l.acquire(context.Background(), time.Now()); err == nil {
		t.Fatal("allowed without persistence")
	}
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "ai-usage.json"), []byte("broken"), 0600)
	if err := ConfigureLimits(dir); err == nil {
		t.Fatal("ignored corrupt counters")
	}
}

func TestInternalGenerationBypassesPublicLimits(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"reflection"}}]}`))
	}))
	defer server.Close()
	t.Setenv("OLLAMA_LLM_MODEL", "test")
	t.Setenv("OLLAMA_BASE_URL", server.URL)
	old := budget
	defer func() { budget = old }()
	budget = &limits{path: filepath.Join(t.TempDir(), "usage.json"), daily: 0, perClient: 0, concurrent: 0, output: 1024}
	if got := AskInternal(context.Background(), nil, "daily reflection"); got != "reflection" {
		t.Fatal(got)
	}
	// Public HTTP/MCP calls use Ask/AskSafe and cannot set the internal context key.
	func() {
		defer func() {
			if recover() == nil {
				t.Error("public request bypassed disabled quota")
			}
		}()
		Ask(context.Background(), nil, "public")
	}()
	if calls != 1 || budget.usage.Total != 0 {
		t.Fatalf("calls=%d usage=%d", calls, budget.usage.Total)
	}
}
