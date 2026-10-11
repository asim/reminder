package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/asim/reminder/quran"
)

func TestExplanationUsesRequestedVerseWithoutModel(t *testing.T) {
	// Avoid contacting any provider in this test.
	llmOnce = sync.Once{}
	llmOnce.Do(func() { llmReachable = false })
	t.Cleanup(func() { llmOnce = sync.Once{}; llmReachable = false })
	q := quran.Load()
	mux := http.NewServeMux()
	registerExplanationRoute(mux, q)
	for _, path := range []string{"/api/quran/0/1/explain", "/api/quran/2/0/explain", "/api/quran/2/287/explain"} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest("POST", path, nil))
		if w.Code != http.StatusNotFound {
			t.Fatalf("%s: %d", path, w.Code)
		}
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("POST", "/api/quran/2/255/explain", nil))
	var result struct {
		Text      string
		Generated bool
		Reference string
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Generated || result.Reference != "2:255" || result.Text != q.Get(2).Verse(255).Comments {
		t.Fatalf("wrong explanation: %+v", result)
	}
}
