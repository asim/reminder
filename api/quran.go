package api

import (
	"net/http"
	"strconv"
)

// RegisterQuranVerseRoute registers the API endpoint for one Quran verse.
// lookup returns the JSON representation and whether the requested verse exists.
func RegisterQuranVerseRoute(mux *http.ServeMux, lookup func(chapter, verse int) ([]byte, bool)) {
	mux.HandleFunc("/api/quran/{chapter}/{verse}", func(w http.ResponseWriter, r *http.Request) {
		chapter, err := strconv.Atoi(r.PathValue("chapter"))
		if err != nil || chapter < 1 {
			writeQuranNotFound(w)
			return
		}

		verse, err := strconv.Atoi(r.PathValue("verse"))
		if err != nil || verse < 1 {
			writeQuranNotFound(w)
			return
		}

		body, ok := lookup(chapter, verse)
		if !ok {
			writeQuranNotFound(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	})
}

func writeQuranNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("{}"))
}
