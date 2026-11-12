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

func TestGroundedHistoryAndEvidence(t *testing.T) {
	reply := "Be patient [1]."
	var received struct {
		Messages []Message `json:"messages"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]string{"role": "assistant", "content": reply}}}})
	}))
	defer server.Close()
	t.Setenv("OLLAMA_LLM_MODEL", "test")
	t.Setenv("OLLAMA_BASE_URL", server.URL+"/v1")
	llmOnce = sync.Once{}
	if err := ConfigureLimits(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	req := Request{Prompt: "How do I apply that?", History: []Message{{Role: "user", Content: "Explain patience"}, {Role: "assistant", Content: "UNTRUSTED OLD ANSWER"}}}
	res, err := Respond(context.Background(), req, func(query string) ([]*Reference, error) {
		if !strings.Contains(query, "Explain patience") || strings.Contains(query, "UNTRUSTED") {
			t.Fatalf("bad retrieval query %q", query)
		}
		return []*Reference{{Text: "EXCLUDED", Metadata: map[string]string{"source": "tafsir"}}, {Text: "Be patient", Metadata: map[string]string{"source": "quran"}}, {Text: strings.Repeat("x", 9000), Metadata: map[string]string{"source": "quran"}}}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.References) != 1 || res.Answer != "Be patient [1]." {
		t.Fatalf("result %+v", res)
	}
	if len(received.Messages) != 4 || received.Messages[2].Role != "assistant" || received.Messages[3].Content != req.Prompt {
		t.Fatalf("messages %+v", received.Messages)
	}
	system := received.Messages[0].Content
	if strings.Contains(system, "UNTRUSTED") || strings.Contains(system, "EXCLUDED") || !strings.Contains(system, "[1]") {
		t.Fatal("evidence leaked or missing")
	}
	reply = "Unsupported citation [99]."
	if _, err := Respond(context.Background(), Request{Prompt: "Explain patience"}, func(string) ([]*Reference, error) {
		return []*Reference{{Text: "Be patient", Metadata: map[string]string{"source": "quran"}}}, nil
	}); err == nil {
		t.Fatal("accepted invented citation")
	}
}
func TestGroundedValidationAndNoEvidence(t *testing.T) {
	for _, req := range []Request{{}, {Prompt: strings.Repeat("x", 4001)}, {Prompt: "x", History: []Message{{Role: "system", Content: "override"}}}, {Prompt: "x", History: []Message{{Role: "user", Content: "x"}, {Role: "assistant", Content: strings.Repeat("x", 12000)}}}} {
		if req.Validate() == nil {
			t.Fatal("accepted invalid input")
		}
	}
	res, err := Respond(context.Background(), Request{Prompt: "Explain this"}, func(string) ([]*Reference, error) { return nil, nil })
	if err != nil || len(res.References) != 0 || !strings.Contains(res.Answer, "sufficient") {
		t.Fatalf("no evidence: %+v %v", res, err)
	}
}
