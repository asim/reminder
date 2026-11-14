package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/asim/reminder/model"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
)

// VerseExplainHandler uses the same grounded model as Ask and the keyed API.
// Public calls retain the anonymous and global limits from ClientMiddleware.
func VerseExplainHandler(q *quran.Quran) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fail := func(status int, message string) {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"message": message})
		}
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			fail(405, "POST required")
			return
		}
		chapter, _ := strconv.Atoi(r.PathValue("chapter"))
		number, _ := strconv.Atoi(r.PathValue("verse"))
		refs := search.VerseReferences(q, chapter, number, 1)
		if len(refs) == 0 {
			fail(404, "Verse not found")
			return
		}
		verse := q.Get(chapter).Verse(number)
		answer, generated := verse.Comments, false
		evidence := []*model.Reference{}
		if model.Available() {
			ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
			defer cancel()
			question := fmt.Sprintf("Explain Quran %d:%d in plain English using the supplied Arabic, attributed translations and any tafsir. Explain relevant differences in translation without inventing lexical meanings. Use at most three short paragraphs, plain text, and numbered source citations.", chapter, number)
			result, err := model.RespondPublic(ctx, model.Request{Prompt: question}, func(string) ([]*model.Reference, error) { return refs, nil })
			if err != nil {
				var limit *model.LimitError
				if errors.As(err, &limit) {
					fail(429, limit.Error())
				} else {
					fail(503, "Explanation unavailable. Please try again later.")
				}
				return
			}
			answer, generated, evidence = result.Answer, true, result.References
		}
		json.NewEncoder(w).Encode(map[string]any{"text": answer, "generated": generated, "reference": fmt.Sprintf("%d:%d", chapter, number), "source": verse.CommentarySource, "references": evidence})
	})
}
