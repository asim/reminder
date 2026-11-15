package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/asim/reminder/arabic"
	"github.com/asim/reminder/quran"
)

func RegisterArabicRoutes(mux *http.ServeMux, q *quran.Quran) {
	send := func(w http.ResponseWriter, data any) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
	mux.HandleFunc("GET /api/arabic/dictionary", func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		if query == "" || len(query) > 200 {
			http.Error(w, "q must contain 1–200 bytes", 400)
			return
		}
		limit, offset := 10, 0
		var e error
		if v := r.URL.Query().Get("limit"); v != "" {
			limit, e = strconv.Atoi(v)
			if e != nil || limit < 1 || limit > 25 {
				http.Error(w, "limit must be 1–25", 400)
				return
			}
		}
		if v := r.URL.Query().Get("offset"); v != "" {
			offset, e = strconv.Atoi(v)
			if e != nil || offset < 0 || offset > 5000 {
				http.Error(w, "offset must be 0–5000", 400)
				return
			}
		}
		entries, total, e := arabic.Lookup(query, offset, limit)
		if e != nil {
			http.Error(w, "Dictionary unavailable", 503)
			return
		}
		send(w, map[string]any{"query": query, "entries": entries, "total": total, "offset": offset, "source": "Lane’s Arabic-English Lexicon", "source_url": arabic.LaneURL, "licence": "CC BY-SA 3.0 US"})
	})
	mux.HandleFunc("GET /api/arabic/roots", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if len(query) > 200 {
			http.Error(w, "q is too long", 400)
			return
		}
		roots, e := arabic.Roots(query, 100)
		if e != nil {
			http.Error(w, "Dictionary unavailable", 503)
			return
		}
		send(w, map[string]any{"roots": roots, "source_url": arabic.LaneURL, "limit": 100})
	})
	mux.HandleFunc("GET /api/arabic/quran/{chapter}/{verse}", func(w http.ResponseWriter, r *http.Request) {
		ch, _ := strconv.Atoi(r.PathValue("chapter"))
		v, _ := strconv.Atoi(r.PathValue("verse"))
		if ch < 1 || ch > len(q.Chapters) || v < 1 {
			http.NotFound(w, r)
			return
		}
		verse := q.Get(ch).Verse(v)
		if verse == nil {
			http.NotFound(w, r)
			return
		}
		a, e := arabic.Analyze(ch, v, verse.Arabic)
		if e != nil {
			http.Error(w, "Morphology unavailable", 503)
			return
		}
		send(w, a)
	})
}
