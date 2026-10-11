package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/asim/reminder/quran"
)

func registerExplanationRoute(mux *http.ServeMux, q *quran.Quran) {
	mux.HandleFunc("POST /api/quran/{chapter}/{verse}/explain", func(w http.ResponseWriter, r *http.Request) {
		chapter, _ := strconv.Atoi(r.PathValue("chapter"))
		number, _ := strconv.Atoi(r.PathValue("verse"))
		if chapter < 1 || chapter > len(q.Chapters) {
			http.Error(w, "Chapter not found", http.StatusNotFound)
			return
		}
		verse := q.Get(chapter).Verse(number)
		if verse == nil {
			http.Error(w, "Verse not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Commentary remains available without a configured model.
		answer, generated := verse.Comments, false
		if LLMAvailable() {
			ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
			defer cancel()
			context := []string{fmt.Sprintf("Quran %d:%d\nArabic: %s\nTranslation: %s\nStored commentary: %s", chapter, number, verse.Arabic, verse.Text, verse.Comments)}
			question := fmt.Sprintf("Explain Quran %d:%d in plain English using only the supplied verse, translation and commentary. Cite the verse number. Distinguish the translation from interpretation. Do not invent historical circumstances, scholarly attributions or rulings. Do not generate or rewrite Arabic. Use plain text, at most three short paragraphs.", chapter, number)
			if result, err := askLLMSafe(ctx, context, question); err == nil && result != "" {
				answer, generated = result, true
			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"text": answer, "generated": generated, "reference": fmt.Sprintf("%d:%d", chapter, number)})
	})
}
