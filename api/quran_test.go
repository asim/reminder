package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuranVerseRouteReturnsJSONForFoundVerse(t *testing.T) {
	mux := http.NewServeMux()
	RegisterQuranVerseRoute(mux, func(chapter, verse int) ([]byte, bool) {
		if chapter == 1 && verse == 1 {
			return []byte(`{"synthetic":"verse"}`), true
		}
		return nil, false
	})

	request := httptest.NewRequest(http.MethodGet, "/api/quran/1/1", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if got := response.Body.String(); got != `{"synthetic":"verse"}` {
		t.Fatalf("body = %q, want synthetic verse JSON", got)
	}
}

func TestQuranVerseRouteReturnsNotFoundForInvalidOrMissingVerse(t *testing.T) {
	mux := http.NewServeMux()
	RegisterQuranVerseRoute(mux, func(chapter, verse int) ([]byte, bool) {
		if chapter == 1 && verse == 1 {
			return []byte(`{"synthetic":"verse"}`), true
		}
		return nil, false
	})

	for _, path := range []string{
		"/api/quran/not-a-chapter/1",
		"/api/quran/0/1",
		"/api/quran/1/not-a-verse",
		"/api/quran/1/0",
		"/api/quran/1/2",
	} {
		t.Run(path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)

			if response.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
			}
			if got := response.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", got)
			}
			if got := response.Body.String(); got != "{}" {
				t.Fatalf("body = %q, want {}", got)
			}
		})
	}
}
