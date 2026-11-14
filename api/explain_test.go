package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asim/reminder/model"
	"github.com/asim/reminder/quran"
)

func TestGroundedVerseExplanationAndLimits(t *testing.T) {
	q := quran.Load()
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Messages []model.Message `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		system := body.Messages[0].Content
		for _, text := range []string{q.Get(1).Verse(1).Arabic, "Saheeh International", "Muhammad Asad"} {
			if !strings.Contains(system, text) {
				t.Errorf("missing source %s", text)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"Explanation [1] [2] [3]."}}]}`))
	}))
	defer provider.Close()
	t.Setenv("OLLAMA_LLM_MODEL", "test")
	t.Setenv("OLLAMA_BASE_URL", provider.URL+"/v1")
	t.Setenv("AI_IP_DAILY_LIMIT", "1")
	t.Setenv("AI_DAILY_LIMIT", "20")
	if err := model.ConfigureLimits(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	handler := VerseExplainHandler(q)
	for _, tc := range []struct {
		method, ch, v string
		want          int
	}{{"POST", "1", "1", 200}, {"POST", "1", "1", 429}, {"POST", "1", "0", 404}, {"POST", "115", "1", 404}, {"GET", "1", "1", 405}} {
		r := httptest.NewRequest(tc.method, "/", nil)
		r.SetPathValue("chapter", tc.ch)
		r.SetPathValue("verse", tc.v)
		r = r.WithContext(model.WithClient(r.Context(), "192.0.2.1"))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatalf("%+v: %d %s", tc, w.Code, w.Body.String())
		}
		if w.Code == http.StatusOK {
			var result struct {
				Text       string
				Generated  bool
				References []*model.Reference
			}
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if !result.Generated || result.Text != "Explanation [1] [2] [3]." || len(result.References) < 3 {
				t.Fatal("missing generated answer or citations")
			}
		}
	}
}
