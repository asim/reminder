package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
)

func TestMCPSearchUsesGroundedModel(t *testing.T) {
	calls := 0
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var body struct {
			Messages []model.Message `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if len(body.Messages) == 0 || !strings.Contains(body.Messages[0].Content, "Sources (JSON records)") || !strings.Contains(body.Messages[0].Content, quran.Load().Get(1).Verse(2).Arabic) || !strings.Contains(body.Messages[0].Content, "jalalayn-en") {
			t.Error("MCP summary must use canonical Arabic, translations and attributed tafsir")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"Explanation [1]."}}]}`))
	}))
	defer provider.Close()
	t.Setenv("OLLAMA_LLM_MODEL", "test")
	t.Setenv("OLLAMA_BASE_URL", provider.URL+"/v1")
	t.Setenv("AI_IP_DAILY_LIMIT", "0")
	t.Setenv("AI_DAILY_LIMIT", "100")
	if err := model.ConfigureLimits(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	ready := make(chan bool)
	close(ready)
	server := NewMCP(quran.Load(), names.Load(), nil, nil, modelIndex{}, ready)
	for _, summary := range []bool{false, true} {
		text, err := server.handlers["search"](context.Background(), map[string]interface{}{"q": "Explain 1:2", "summarise": summary})
		if err != nil {
			t.Fatal(err)
		}
		var result struct {
			Answer     string             `json:"answer"`
			References []*model.Reference `json:"references"`
		}
		if err := json.Unmarshal([]byte(text), &result); err != nil {
			t.Fatal(err)
		}
		if summary {
			if calls != 1 || result.Answer != "Explanation [1]." || len(result.References) < 4 {
				t.Fatalf("ungrounded summary: %s", text)
			}
		} else if calls != 0 || result.Answer != "" {
			t.Fatal("plain search invoked a model")
		}
	}
}
