package model

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExplanationUsesOnePlainCall(t *testing.T) {
	for _, provider := range []string{"openai", "anthropic"} {
		for _, failure := range []bool{false, true} {
			t.Run(provider+map[bool]string{false: "-success", true: "-failure"}[failure], func(t *testing.T) {
				calls := 0
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls++
					var body map[string]any
					json.NewDecoder(r.Body).Decode(&body)
					if body["tools"] != nil || body["tool_choice"] != nil {
						t.Error("explanation sends tools")
					}
					data, _ := json.Marshal(body)
					for _, text := range []string{"EXACT PASSAGE", "TAFSIR", "previous question", "What does that mean?"} {
						if !strings.Contains(string(data), text) {
							t.Errorf("missing %s", text)
						}
					}
					if failure {
						w.WriteHeader(400)
						w.Write([]byte(`{"error":{"type":"invalid_request_error","message":"bad request"}}`))
						return
					}
					nativeResponse(w, provider, "Meaning [1].", nil)
				}))
				defer server.Close()
				configureResearchTest(t, provider, server.URL)
				request := ExplainRequest{Context: ExplainContext{Source: "quran", Chapter: 1, Verse: 2}, Question: "What does that mean?", History: []Message{{Role: "user", Content: "previous question"}, {Role: "assistant", Content: "previous answer"}}}
				result, err := Explain(context.Background(), request, func(string) ([]*Reference, error) {
					return []*Reference{{Text: "EXACT PASSAGE", Metadata: map[string]string{"source": "quran"}}, {Text: "TAFSIR", Metadata: map[string]string{"source": "tafsir"}}}, nil
				})
				if calls != 1 {
					t.Fatalf("expected one call, got %d", calls)
				}
				if failure {
					if err == nil || result != nil {
						t.Fatal("provider failure swallowed")
					}
				} else if err != nil || result.Answer != "Meaning [1]." || !result.Complete {
					t.Fatalf("%+v %v", result, err)
				}
			})
		}
	}
}
