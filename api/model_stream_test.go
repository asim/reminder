package api

import (
	"context"
	"encoding/json"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestModelHTTPStreamAndErrors(t *testing.T) {
	for _, fail := range []bool{false, true} {
		handler := modelHTTP(nil, func(r *http.Request, _ GenerationAccess) (*model.Result, error) {
			model.Emit(r.Context(), model.Event{Type: "delta", Text: "Draft"})
			if fail {
				return nil, invalidModelRequest("Invalid passage")
			}
			return cachedFixture(), nil
		})
		request := httptest.NewRequest("POST", "/api/explain", strings.NewReader("{}"))
		request.Header.Set("Accept", "application/x-ndjson")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, request)
		if !w.Flushed || w.Header().Get("X-Accel-Buffering") != "no" {
			t.Fatal("stream not flushed")
		}
		var events []model.Event
		for _, line := range strings.Split(strings.TrimSpace(w.Body.String()), "\n") {
			var e model.Event
			if err := json.Unmarshal([]byte(line), &e); err != nil {
				t.Fatal(err)
			}
			events = append(events, e)
		}
		last := events[len(events)-1]
		if fail {
			if last.Type != "error" || last.Result != nil {
				t.Fatal("failed stream completed")
			}
		} else if last.Type != "result" || last.Result.Answer != cachedFixture().Answer {
			t.Fatal("missing final result")
		}
	}
}
func TestStreamCacheHitDoesNotGenerate(t *testing.T) {
	cache, err := openExplainCache(t.TempDir(), "corpus")
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	req := model.ExplainRequest{Context: model.ExplainContext{Chapter: 1}}
	calls := 0
	generate := func() (*model.Result, error) { calls++; return cachedFixture(), nil }
	cache.result(context.Background(), cache.key(req), generate)
	ctx := model.WithEvents(context.Background(), func(model.Event) {})
	if _, err := cache.result(ctx, cache.key(req), generate); err != nil || calls != 1 {
		t.Fatal("stream missed existing JSON cache")
	}
}

type countedExplainIndex struct {
	modelIndex
	calls int
}

func (idx *countedExplainIndex) Query(q string) ([]*search.Result, error) {
	idx.calls++
	return idx.modelIndex.Query(q)
}
func TestDefaultExplainSkipsBroadSearch(t *testing.T) {
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"finish_reason":"stop","message":{"role":"assistant","content":"Explanation [1]."}}]}`))
	}))
	defer provider.Close()
	t.Setenv("OLLAMA_LLM_MODEL", "test")
	t.Setenv("OLLAMA_BASE_URL", provider.URL+"/v1")
	t.Setenv("AI_IP_DAILY_LIMIT", "0")
	t.Setenv("AI_DAILY_LIMIT", "20")
	if err := model.ConfigureLimits(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	ready := make(chan bool)
	close(ready)
	index := &countedExplainIndex{}
	service := &ModelService{Index: index, Names: &names.Names{}, Quran: quran.Load(), Ready: ready}
	req := model.ExplainRequest{Context: model.ExplainContext{Source: "quran", Chapter: 1, Verse: 2}}
	if _, err := service.Explain(context.Background(), req, nil); err != nil {
		t.Fatal(err)
	}
	if index.calls != 0 {
		t.Fatal("default explanation ran broad search")
	}
	req.Question = "What does this mean?"
	if _, err := service.Explain(context.Background(), req, nil); err != nil {
		t.Fatal(err)
	}
	if index.calls == 0 {
		t.Fatal("custom question lost retrieval")
	}
}
