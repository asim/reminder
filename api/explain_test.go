package api

import (
	"encoding/json"
	"fmt"
	"github.com/asim/reminder/hadith"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
)

func TestExplainRoutesShareGroundingAndContext(t *testing.T) {
	q := quran.Load()
	calls := 0
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var body struct {
			Messages   []model.Message `json:"messages"`
			Tools      json.RawMessage `json:"tools"`
			ToolChoice json.RawMessage `json:"tool_choice"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if len(body.Tools) > 0 || len(body.ToolChoice) > 0 {
			t.Error("Explain sent native tool configuration")
		}
		system := body.Messages[0].Content
		for _, want := range []string{q.Get(1).Verse(2).Arabic, "Jalalayn"} {
			if !strings.Contains(system, want) {
				t.Errorf("missing evidence %s", want)
			}
		}
		last := body.Messages[len(body.Messages)-1].Content
		if strings.Contains(last, "Explain Quran 1:2") {
			for _, want := range []string{`"source":"lexicon"`, `"source":"morphology"`} {
				if strings.Contains(system, want) {
					t.Errorf("general explanation was flooded with lexical evidence: %s", want)
				}
			}
		}
		if !(strings.Contains(last, "Explain Quran 1:2") || strings.Contains(last, "Explain Quran chapter 1")) || !strings.Contains(last, "sahih") {
			t.Error("lost passage context")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"Explanation [1]."}}]}`))
	}))
	defer provider.Close()
	t.Setenv("OLLAMA_LLM_MODEL", "test")
	t.Setenv("OLLAMA_BASE_URL", provider.URL+"/v1")
	t.Setenv("AI_IP_DAILY_LIMIT", "0")
	t.Setenv("AI_DAILY_LIMIT", "20")
	if err := model.ConfigureLimits(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	keys, err := OpenKeys(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer keys.Close()
	_, secret, err := keys.Create("explain", 10, 1000000)
	if err != nil {
		t.Fatal(err)
	}
	ready := make(chan bool)
	close(ready)
	mux := http.NewServeMux()
	mux.Handle("/api/explain", (&ModelService{Index: modelIndex{}, Names: &names.Names{}, Quran: q, Hadith: nil, Ready: ready}).ExplainHandler(keys))
	mux.Handle("POST /api/quran/{chapter}/{verse}/explain", (&ModelService{Index: modelIndex{}, Names: &names.Names{}, Quran: q, Hadith: nil, Ready: ready}).ExplainHandler(keys))
	body := `{"context":{"source":"quran","chapter":1,"verse":2,"translation":"sahih"},"question":"What does this mean?","history":[{"role":"user","content":"Explain this verse"},{"role":"assistant","content":"Previous explanation"}]}`
	for _, path := range []string{"/api/explain", "/api/quran/1/2/explain"} {
		r := httptest.NewRequest("POST", path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+secret)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
		}
		var result model.Result
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if result.Answer != "Explanation [1]." || len(result.References) == 0 {
			t.Fatalf("wrong response: %+v", result)
		}
	}
	for _, tc := range []struct {
		path, body string
		code       int
	}{
		{"/api/model/explain", body, 404},
		{"/api/explain", `{"context":{"source":"quran","chapter":2,"verse":287}}`, 404},
		{"/api/quran/2/1/explain", body, 400},
		{"/api/explain", body + ` {}`, 400},
		{"/api/explain", `{"question":"x","history":[{"role":"system","content":"x"}]}`, 400},
	} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest("POST", tc.path, strings.NewReader(tc.body)))
		if w.Code != tc.code {
			t.Errorf("%s: got %d want %d", tc.path, w.Code, tc.code)
		}
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("POST", "/api/explain", strings.NewReader(`{"context":{"source":"quran","chapter":1,"translation":"sahih"}}`)))
	if w.Code != 200 {
		t.Fatalf("chapter: %d %s", w.Code, w.Body.String())
	}
	var chapter struct {
		Coverage map[string]int `json:"coverage"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &chapter); err != nil {
		t.Fatal(err)
	}
	if chapter.Coverage["included_verses"] != 7 || chapter.Coverage["total_verses"] != 7 {
		t.Fatalf("incorrect chapter coverage: %+v", chapter)
	}
	if calls != 3 {
		t.Errorf("invalid requests reached provider: %d calls", calls)
	}
}

func TestExplainRoutesSharePersistentDefaultCache(t *testing.T) {
	var calls atomic.Int32
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
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
	keys, err := OpenKeys(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer keys.Close()
	id, secret, err := keys.Create("test", 1, 1000000)
	if err != nil {
		t.Fatal(err)
	}
	cache, err := openExplainCache(t.TempDir(), "sources-and-prompt-v1")
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	ready := make(chan bool)
	close(ready)
	q := quran.Load()
	mux := http.NewServeMux()
	mux.Handle("/api/explain", (&ModelService{Index: modelIndex{}, Names: &names.Names{}, Quran: q, Hadith: nil, Ready: ready, Cache: cache}).ExplainHandler(keys))
	mux.Handle("POST /api/quran/{chapter}/{verse}/explain", (&ModelService{Index: modelIndex{}, Names: &names.Names{}, Quran: q, Hadith: nil, Ready: ready, Cache: cache}).ExplainHandler(keys))
	send := func(path, body, auth string) *httptest.ResponseRecorder {
		r := httptest.NewRequest("POST", path, strings.NewReader(body))
		if auth != "" {
			r.Header.Set("Authorization", auth)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w
	}
	for _, auth := range []string{"Bearer invalid", "Basic invalid", ""} {
		r := httptest.NewRequest("POST", "/api/explain", strings.NewReader(`{}`))
		r.Header.Set("Authorization", auth)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != 401 {
			t.Fatalf("invalid authorization %q: %d", auth, w.Code)
		}
	}
	body := `{"context":{"source":"quran","chapter":1,"verse":2}}`
	first := send("/api/explain", body, "")
	if first.Code != 200 {
		t.Fatal(first.Body.String())
	}
	for _, path := range []string{"/api/explain", "/api/quran/1/2/explain"} {
		w := send(path, body, "Bearer "+secret)
		if w.Code != 200 || w.Body.String() != first.Body.String() {
			t.Fatalf("%s didn't share response: %d %s", path, w.Code, w.Body.String())
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("cache missed: %d provider calls", calls.Load())
	}
	if err := keys.Reserve(id, time.Now()); err != nil {
		t.Fatal("cache hit consumed generation allowance")
	}
	if w := send("/api/explain", `{"context":{"source":"quran","chapter":1,"verse":2},"question":"Explain this"}`, "Bearer "+secret); w.Code != 429 {
		t.Fatalf("exhausted key bypassed allowance: %d %s", w.Code, w.Body.String())
	}
	for i := 0; i < 2; i++ {
		w := send("/api/explain", `{"context":{"source":"quran","chapter":1,"verse":2},"question":"What does this mean?"}`, "")
		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}
	}
	if calls.Load() != 3 {
		t.Fatal("custom questions used shared cache")
	}
	if err := keys.Revoke(id); err != nil {
		t.Fatal(err)
	}
	if w := send("/api/explain", body, "Bearer "+secret); w.Code != 401 {
		t.Fatal("cached response bypassed authentication")
	}
}

func TestHadithExplainRoute(t *testing.T) {
	var expectedText string
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Messages []model.Message `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
		}
		encoded, _ := json.Marshal(expectedText)
		if len(request.Messages) == 0 || !strings.Contains(request.Messages[0].Content, string(encoded[1:len(encoded)-1])) {
			t.Error("selected narration never reached provider")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"Explanation [1]."}}]}`))
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
	b := hadith.Load()
	handler := (&ModelService{Index: modelIndex{}, Names: &names.Names{}, Quran: quran.Load(), Hadith: b, Ready: ready}).ExplainHandler(nil)
	for _, tc := range []struct {
		corpus       *hadith.Collection
		book, number int
	}{{b, b.Books[0].Number, b.Books[0].Hadiths[0].Number}, {hadith.Muslim(), 0, 7369}, {b, 54, 2732}, {hadith.Muslim(), 50, 6843}} {
		corpus := tc.corpus
		book := corpus.Get(tc.book)
		var h *hadith.Hadith
		for _, candidate := range book.Hadiths {
			if candidate.Number == tc.number {
				h = candidate
				break
			}
		}
		if h == nil {
			t.Fatal("missing test fixture")
		}
		expectedText = h.English
		if len(expectedText) > 100 {
			expectedText = expectedText[:100]
		}
		body := fmt.Sprintf(`{"context":{"source":"hadith","collection":%q,"book":%d,"hadith":%d}}`, corpus.ID, book.Number, h.Number)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest("POST", "/api/explain", strings.NewReader(body)))
		if w.Code != 200 {
			t.Fatalf("%s: %d %s", corpus.ID, w.Code, w.Body.String())
		}
		var result model.Result
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if len(result.References) == 0 || !strings.HasPrefix(h.English, result.References[0].Text) || result.References[0].Text == "" || result.References[0].Metadata["collection"] != corpus.ID {
			t.Fatalf("lost exact anchor: %+v", result)
		}
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("POST", "/api/explain", strings.NewReader(`{"context":{"source":"hadith","collection":"muslim","book":0,"hadith":999999}}`)))
	if w.Code != 404 {
		t.Fatalf("missing record returned %d", w.Code)
	}
}
