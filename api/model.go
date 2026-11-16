package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
)

// ModelHandler exposes grounded prompting independently of search summaries.
func ModelHandler(keys *KeyStore, idx search.Searcher, corpus *names.Names, q *quran.Quran, b *hadith.Collection, indexed <-chan bool) http.Handler {
	return modelHandler(keys, idx, corpus, q, b, indexed, false, nil)
}

// ExplainHandler uses API key allowances when keys is non-nil, and the existing
// anonymous model limits for reader requests otherwise.
func ExplainHandler(keys *KeyStore, idx search.Searcher, corpus *names.Names, q *quran.Quran, b *hadith.Collection, indexed <-chan bool, caches ...*ExplainCache) http.Handler {
	var cache *ExplainCache
	if len(caches) > 0 {
		cache = caches[0]
	}
	return modelHandler(keys, idx, corpus, q, b, indexed, true, cache)
}

func modelHandler(keys *KeyStore, idx search.Searcher, corpus *names.Names, q *quran.Quran, b *hadith.Collection, indexed <-chan bool, explain bool, cache *ExplainCache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		fail := func(status int, message string) {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"message": message})
		}
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			fail(405, "POST required")
			return
		}
		var id string
		var err error
		if keys != nil {
			fields := strings.Fields(r.Header.Get("Authorization"))
			if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
				w.Header().Set("WWW-Authenticate", "Bearer")
				fail(401, "Requires authorisation: Authorization: Bearer <key>")
				return
			}
			id, err = keys.Authenticate(fields[1])
			if err != nil {
				if errors.Is(err, errKeyUnauthorized) {
					w.Header().Set("WWW-Authenticate", "Bearer")
					fail(401, "Invalid API key")
				} else {
					log.Printf("Model key lookup: %v", err)
					fail(503, "Model access unavailable")
				}
				return
			}
		}
		var request model.Request
		var explanation model.ExplainRequest
		var passage []*model.Reference
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, model.ModelBodyLimit()))
		decoder.DisallowUnknownFields()
		var body any = &request
		if explain {
			body = &explanation
		}
		if err := decoder.Decode(body); err != nil && !(explain && r.PathValue("chapter") != "" && err == io.EOF) {
			fail(400, "Invalid model request")
			return
		}
		if err := decoder.Decode(new(any)); err != io.EOF {
			fail(400, "Expected one JSON object")
			return
		}
		if explain {
			if r.PathValue("chapter") != "" {
				chapter, e1 := strconv.Atoi(r.PathValue("chapter"))
				verse, e2 := strconv.Atoi(r.PathValue("verse"))
				if e1 != nil || e2 != nil || chapter < 1 || verse < 1 {
					fail(400, "Invalid passage")
					return
				}
				if (explanation.Context.Source != "" && explanation.Context.Source != "quran") || (explanation.Context.Chapter != 0 && explanation.Context.Chapter != chapter) || (explanation.Context.Verse != 0 && explanation.Context.Verse != verse) {
					fail(400, "Passage context does not match URL")
					return
				}
				explanation.Context.Source = "quran"
				explanation.Context.Chapter = chapter
				explanation.Context.Verse = verse
			}
			if explanation.Context.Chapter > 0 {
				if explanation.Context.Source == "" {
					explanation.Context.Source = "quran"
				}
				if explanation.Context.Translation == "" {
					explanation.Context.Translation = "khattab"
				}
			}
			request, err = explanation.Request()
			if err != nil {
				fail(400, err.Error())
				return
			}
			passage, err = search.ExplanationPassage(q, explanation.Context)
			if err != nil {
				fail(404, err.Error())
				return
			}
		}
		if err := request.Validate(); err != nil {
			fail(400, err.Error())
			return
		}
		select {
		case <-indexed:
		default:
			fail(503, "Indexing content, please try again later")
			return
		}
		generate := func() (*model.Result, error) {
			if !model.Available() {
				return nil, &explainHTTPError{503, "No model configured"}
			}
			if keys != nil {
				if err := keys.Reserve(id, time.Now()); err != nil {
					if errors.Is(err, errKeyLimit) {
						return nil, &explainHTTPError{429, err.Error()}
					} else {
						log.Printf("Model key allowance: %v", err)
						return nil, &explainHTTPError{503, "Model access unavailable"}
					}
				}
			}
			// Key allowances replace the anonymous per-IP allowance, not the global cap.
			ctx := r.Context()
			if keys != nil {
				ctx = model.WithClient(ctx, "")
				ctx = model.WithTokenAllowance(ctx, func(units int) error { return keys.ReserveTokens(id, units, time.Now()) })
			}
			var focus *model.ExplainRequest
			if explain {
				focus = &explanation
			}
			return respondWithSources(ctx, request, focus, passage, idx, corpus, q, b)

		}
		var result *model.Result
		if explain {
			result, err = cache.result(r.Context(), cache.key(explanation), generate)
		} else {
			result, err = generate()
		}

		if err != nil {
			var status *explainHTTPError
			if errors.As(err, &status) {
				fail(status.code, status.message)
				return
			}
			var limit *model.LimitError
			if errors.As(err, &limit) {
				fail(429, limit.Error())
			} else {
				log.Printf("Model request failed: %v", err)
				fail(503, "Model unavailable. Please try again later.")
			}
			return
		}

		json.NewEncoder(w).Encode(result)
	})
}

type explainHTTPError struct {
	code    int
	message string
}

func (e *explainHTTPError) Error() string { return e.message }

// respondWithSources is shared by HTTP and MCP explanations and search prompting.
func respondWithSources(ctx context.Context, request model.Request, explanation *model.ExplainRequest, passage []*model.Reference, idx search.Searcher, corpus *names.Names, q *quran.Quran, b *hadith.Collection) (*model.Result, error) {
	retrieve := func(query string) ([]*model.Reference, error) {
		refs, err := search.ModelReferences(idx, corpus, q, b, query)
		if err != nil {
			return nil, err
		}
		// Dedupe exact records without merging distinct source annotations.
		out := make([]*model.Reference, 0, len(passage)+len(refs))
		seen := map[string]bool{}
		for _, ref := range append(append([]*model.Reference{}, passage...), refs...) {
			if ref == nil {
				continue
			}
			identity, _ := json.Marshal(struct {
				Text     string
				Metadata map[string]string
			}{ref.Text, ref.Metadata})
			if !seen[string(identity)] {
				seen[string(identity)] = true
				out = append(out, ref)
			}
		}
		return out, nil
	}
	var result *model.Result
	var err error
	if explanation != nil {
		result, err = model.Explain(ctx, *explanation, retrieve)
	} else {
		result, err = model.Respond(ctx, request, retrieve)
	}
	if err != nil {
		return nil, err
	}
	if explanation != nil && explanation.Context.Chapter > 0 && explanation.Context.Verse == 0 {
		seen := map[string]bool{}
		for _, ref := range result.References {
			md := ref.Metadata
			if md["chapter"] == strconv.Itoa(explanation.Context.Chapter) && (md["source"] == "quran" || md["source"] == "quran_translation") {
				seen[md["verse"]] = true
			}
		}
		total := 0
		for _, verse := range q.Get(explanation.Context.Chapter).Verses {
			if verse.Number > 0 {
				total++
			}
		}
		result.Coverage = map[string]int{"included_verses": len(seen), "total_verses": total}
	}

	return result, nil
}
