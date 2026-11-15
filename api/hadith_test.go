package api

import (
	"encoding/json"
	"github.com/asim/reminder/hadith"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHadithCollectionRoutes(t *testing.T) {
	mux := http.NewServeMux()
	RegisterHadithRoutes(mux, hadith.Load())
	for _, tc := range []struct {
		path                  string
		number, count, status int
	}{
		{"/api/hadith/1", 1, 0, 200}, {"/api/hadith/muslim/1", 1, 0, 200}, {"/api/hadith/muslim/0", 0, 91, 200},
		{"/api/hadith/0", 0, 0, 404}, {"/api/hadith/muslim/999", 0, 0, 404}, {"/api/hadith/muslim/invalid", 0, 0, 404},
	} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if w.Code != tc.status {
			t.Fatalf("%s: %d", tc.path, w.Code)
		}
		if tc.status == 200 {
			var b hadith.Book
			if err := json.Unmarshal(w.Body.Bytes(), &b); err != nil {
				t.Fatal(err)
			}
			if b.Number != tc.number || len(b.Hadiths) == 0 || (tc.count > 0 && len(b.Hadiths) != tc.count) {
				t.Fatalf("%s: wrong book", tc.path)
			}
		}
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/api/hadith/collections", nil))
	var c []map[string]string
	json.Unmarshal(w.Body.Bytes(), &c)
	if len(c) != 2 || c[1]["id"] != "muslim" {
		t.Fatal("missing collection")
	}
}
