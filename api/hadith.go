package api

import (
	"encoding/json"
	"github.com/asim/reminder/hadith"
	"net/http"
	"strconv"
)

func RegisterHadithRoutes(mux *http.ServeMux, b *hadith.Collection) {
	send := func(w http.ResponseWriter, data any) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	}
	mux.HandleFunc("GET /api/hadith/collections", func(w http.ResponseWriter, r *http.Request) {
		send(w, []map[string]string{{"id": "bukhari", "name": b.Name}, {"id": "muslim", "name": hadith.Muslim().Name}})
	})
	for _, entry := range []struct {
		prefix string
		corpus *hadith.Collection
	}{{"/api/hadith", b}, {"/api/hadith/muslim", hadith.Muslim()}} {
		c := entry.corpus
		mux.HandleFunc("GET "+entry.prefix, func(w http.ResponseWriter, r *http.Request) { send(w, c) })
		mux.HandleFunc("GET "+entry.prefix+"/books", func(w http.ResponseWriter, r *http.Request) { send(w, c.Index().Books) })
		mux.HandleFunc("GET "+entry.prefix+"/{book}", func(w http.ResponseWriter, r *http.Request) {
			number, err := strconv.Atoi(r.PathValue("book"))
			if err != nil {
				http.NotFound(w, r)
				return
			}
			book := c.Get(number)
			if book == nil {
				http.NotFound(w, r)
				return
			}
			send(w, book)
		})
	}
}
