package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/asim/reminder/tafsir"
)

func TestTafsirRoutes(t *testing.T) {
	mux := http.NewServeMux()
	RegisterTafsirRoutes(mux)
	request := func(method, path string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest(method, path, nil))
		return w
	}
	w := request("GET", "/api/tafsir")
	var editions []tafsir.Edition
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &editions) != nil || len(editions) != len(tafsir.Editions()) {
		t.Fatalf("editions: %s", w.Body.String())
	}
	for _, collection := range tafsir.Editions() {
		// Exercise a shared-passage pointer where available, not only its first verse.
		passage := collection.Passages()[0]
		for _, p := range collection.Passages() {
			if len(p.Ayahs) > 1 {
				passage = p
				break
			}
		}
		ref := passage.Ayahs[len(passage.Ayahs)-1]
		ch, v := tafsir.Numbers(ref)
		w = request("GET", "/api/tafsir/"+collection.Edition.ID+"/"+strconv.Itoa(ch)+"/"+strconv.Itoa(v))
		var result struct {
			Edition tafsir.Edition `json:"edition"`
			tafsir.Passage
		}
		if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &result) != nil {
			t.Fatalf("passage: %s", w.Body.String())
		}
		if result.Edition != collection.Edition || result.Text != passage.Text || result.Ref != passage.Ref || len(result.Ayahs) != len(passage.Ayahs) {
			t.Fatal("lost attribution or shared passage")
		}
	}
	for _, tc := range []struct {
		method, path string
		code         int
	}{
		{"POST", "/api/tafsir", 405},
		{"GET", "/api/tafsir/missing/1/1", 404},
		{"GET", "/api/tafsir/jalalayn-en/1/999", 404},
		{"GET", "/api/tafsir/jalalayn-en/no/1", 400},
		{"GET", "/api/tafsir/jalalayn-en/115/1", 400},
	} {
		if w := request(tc.method, tc.path); w.Code != tc.code {
			t.Errorf("%s: %d want %d", tc.path, w.Code, tc.code)
		}
	}
}
