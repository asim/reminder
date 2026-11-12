package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/search"
)

// ModelHandler exposes grounded prompting independently of search summaries.
func ModelHandler(keys *KeyStore, idx search.Searcher, corpus *names.Names, indexed <-chan bool) http.Handler {
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
		fields := strings.Fields(r.Header.Get("Authorization"))
		if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
			w.Header().Set("WWW-Authenticate", "Bearer")
			fail(401, "Requires authorisation: Authorization: Bearer <key>")
			return
		}
		id, err := keys.Authenticate(fields[1])
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
		var request model.Request
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16384))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil {
			fail(400, "Invalid model request")
			return
		}
		if err := decoder.Decode(new(any)); err != io.EOF {
			fail(400, "Expected one JSON object")
			return
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
		if !model.Available() {
			fail(503, "No model configured")
			return
		}
		if err := keys.Reserve(id, time.Now()); err != nil {
			if errors.Is(err, errKeyLimit) {
				fail(429, err.Error())
			} else {
				log.Printf("Model key allowance: %v", err)
				fail(503, "Model access unavailable")
			}
			return
		}
		// Key allowances replace the anonymous per-IP allowance, not the global cap.
		ctx := model.WithClient(r.Context(), "")
		result, err := model.Respond(ctx, request, func(query string) ([]*model.Reference, error) { return search.References(idx, corpus, query) })
		if err != nil {
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
