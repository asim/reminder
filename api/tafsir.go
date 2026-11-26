package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asim/reminder/tafsir"
)

// RegisterTafsirRoutes exposes the bundled commentary independently of generation.
func RegisterTafsirRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tafsir", func(w http.ResponseWriter, r *http.Request) {
		editions := make([]tafsir.Edition, 0)
		for _, collection := range tafsir.Editions() {
			editions = append(editions, collection.Edition)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(editions)
	})
	mux.HandleFunc("GET /api/tafsir/{edition}/{chapter}/{verse}", func(w http.ResponseWriter, r *http.Request) {
		chapter, chErr := strconv.Atoi(r.PathValue("chapter"))
		verse, vErr := strconv.Atoi(r.PathValue("verse"))
		if chErr != nil || vErr != nil || chapter < 1 || chapter > 114 || verse < 1 {
			http.Error(w, "Invalid chapter or verse", http.StatusBadRequest)
			return
		}
		for _, collection := range tafsir.Editions() {
			if collection.Edition.ID != r.PathValue("edition") {
				continue
			}
			passage := collection.Get(chapter, verse)
			if passage == nil {
				break
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				Edition tafsir.Edition `json:"edition"`
				*tafsir.Passage
			}{collection.Edition, passage})
			return
		}
		http.Error(w, "Tafsir edition or passage not found", http.StatusNotFound)
	})
}
