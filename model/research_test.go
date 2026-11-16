package model

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestResearchRetrievesBeyondInitialResults(t *testing.T) {
	calls, lookups := 0, 0
	plan := `{"lookups":[{"tool":"quran_mentions","terms":["Musa","Moses"]}]}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var body struct {
			Messages []Message `json:"messages"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		reply := plan
		if calls%2 == 0 {
			if !strings.Contains(body.Messages[0].Content, "COMPLETE INVENTORY") {
				t.Error("lookup results not supplied to synthesis")
			}
			reply = "Musa is mentioned in these chapters [1]."
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"finish_reason": "stop", "message": map[string]string{"content": reply}}}})
	}))
	defer server.Close()
	t.Setenv("OLLAMA_LLM_MODEL", "test")
	t.Setenv("OLLAMA_BASE_URL", server.URL+"/v1")
	llmOnce = sync.Once{}
	if err := ConfigureLimits(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	retrieve := func(string) ([]*Reference, error) {
		return []*Reference{{Text: "one initial verse", Metadata: map[string]string{"source": "quran"}}}, nil
	}
	lookup := func(call Lookup) ([]*Reference, error) {
		lookups++
		if call.Tool != "quran_mentions" || len(call.Terms) != 2 {
			t.Fatal(call)
		}
		return []*Reference{{Text: "COMPLETE INVENTORY", Metadata: map[string]string{"source": "quran", "kind": "mention_inventory"}}}, nil
	}
	result, err := Research(context.Background(), Request{Prompt: "Which chapters is Musa mentioned in?"}, retrieve, lookup)
	if err != nil || calls != 2 || lookups != 1 || len(result.References) != 2 || !result.Complete {
		t.Fatalf("research failed: %+v %v calls=%d", result, err, calls)
	}
	// A malformed plan must not quietly produce an apparently complete answer.
	plan = "not JSON"
	if _, err = Research(context.Background(), Request{Prompt: "Which chapters?"}, retrieve, lookup); err == nil || lookups != 1 {
		t.Fatal("invalid plan executed")
	}
}
