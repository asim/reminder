package api

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
)

func TestEveryDocumentedEndpointHasExamples(t *testing.T) {
	docs := Load()
	seen := map[string]bool{}
	for _, endpoint := range docs.Endpoints {
		key := endpoint.Method + " " + endpoint.Path
		if seen[key] {
			t.Fatalf("duplicate endpoint %s", key)
		}
		seen[key] = true
		if endpoint.Example == nil || endpoint.Example.Request == "" || !json.Valid(endpoint.Example.Response) {
			t.Fatalf("missing/invalid example for %s", key)
		}
		if !strings.Contains(endpoint.Example.Request, "https://reminder.dev/") {
			t.Fatalf("missing request URL: %s", key)
		}
		if endpoint.Method == "POST" {
			parts := strings.SplitN(endpoint.Example.Request, "-d '", 2)
			if len(parts) != 2 || !json.Valid([]byte(strings.TrimSuffix(parts[1], "'"))) {
				t.Fatalf("invalid request JSON: %s", key)
			}
		}
	}
	if len(seen) != len(examples) {
		t.Fatal("unattached examples")
	}
	markdown := docs.Markdown()
	if strings.Count(markdown, "### Example request") != len(seen) || strings.Count(markdown, "### Example response") != len(seen) {
		t.Fatal("missing rendered examples")
	}
	if !strings.Contains(markdown, "Authorization: Bearer $REMINDER_API_KEY") {
		t.Fatal("missing bearer example")
	}
	// The React page imports the same JSON rather than duplicating examples.
}

// Examples intentionally omit fields and array entries; every value shown must
// still match the loaded corpus and serialization used by the handlers.
func exampleSubset(example, actual any) bool {
	switch e := example.(type) {
	case map[string]any:
		a, ok := actual.(map[string]any)
		if !ok {
			return false
		}
		for k, v := range e {
			value, ok := a[k]
			if !ok || !exampleSubset(v, value) {
				return false
			}
		}
		return true
	case []any:
		a, ok := actual.([]any)
		if !ok {
			return false
		}
		for _, v := range e {
			found := false
			for _, value := range a {
				if exampleSubset(v, value) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(example, actual)
	}
}
func TestContentExamplesMatchCorpus(t *testing.T) {
	q := quran.Load()
	h := hadith.Load()
	n := names.Load()
	refs, err := search.WithQuran(nil, q).Query("Al Hashr")
	if err != nil {
		t.Fatal(err)
	}
	responses := map[string]any{
		"GET /api/quran": q, "GET /api/quran/chapters": q.Index().Chapters,
		"GET /api/quran/{chapter}": q.Get(112), "GET /api/quran/{chapter}/{verse}": q.Get(112).Verse(1),
		"GET /api/hadith": h, "GET /api/hadith/{book}": h.Get(1), "GET /api/names": n,
		"POST /api/search": map[string]any{"q": "Al Hashr", "answer": "", "references": refs},
	}
	for key, response := range responses {
		bytes, err := json.Marshal(response)
		if err != nil {
			t.Fatal(err)
		}
		var actual, example any
		if err := json.Unmarshal(bytes, &actual); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(examples[key].Response, &example); err != nil {
			t.Fatal(err)
		}
		if !exampleSubset(example, actual) {
			t.Errorf("example differs from corpus: %s", key)
		}
	}
}
func TestMCPInitializeExample(t *testing.T) {
	example := examples["POST /mcp"]
	parts := strings.SplitN(example.Request, "-d '", 2)
	r := httptest.NewRequest("POST", "/mcp", strings.NewReader(strings.TrimSuffix(parts[1], "'")))
	w := httptest.NewRecorder()
	NewServer("reminder", "1.0.0").ServeHTTP(w, r)
	var want, got any
	if err := json.Unmarshal(example.Response, &want); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("MCP example mismatch: %s", w.Body.String())
	}
}
