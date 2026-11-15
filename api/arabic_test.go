package api

import (
	"encoding/json"
	"github.com/asim/reminder/quran"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArabicRoutes(t *testing.T) {
	mux := http.NewServeMux()
	RegisterArabicRoutes(mux, quran.Load())
	for _, tc := range []struct {
		path   string
		status int
	}{{"/api/arabic/dictionary?q=صبر", 200}, {"/api/arabic/dictionary?q=صبر&limit=26", 400}, {"/api/arabic/dictionary", 400}, {"/api/arabic/dictionary?q=صبر&offset=-1", 400}, {"/api/arabic/roots?q=رح", 200}, {"/api/arabic/quran/1/2", 200}, {"/api/arabic/quran/0/1", 404}, {"/api/arabic/quran/1/99", 404}} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if w.Code != tc.status {
			t.Fatalf("%s: %d %s", tc.path, w.Code, w.Body.String())
		}
		if tc.status == 200 && !json.Valid(w.Body.Bytes()) {
			t.Fatal("invalid JSON")
		}
	}
}
