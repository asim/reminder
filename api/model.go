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

	"github.com/asim/reminder/model"
)

func (s *ModelService) AskHandler(keys *KeyStore) http.Handler {
	return modelHTTP(keys, func(r *http.Request, access GenerationAccess) (*model.Result, error) {
		var request model.Request
		if err := decodeModelRequest(r, &request, false); err != nil {
			return nil, err
		}
		return s.Ask(r.Context(), request, access)
	})
}

func (s *ModelService) ExplainHandler(keys *KeyStore) http.Handler {
	return modelHTTP(keys, func(r *http.Request, access GenerationAccess) (*model.Result, error) {
		var request model.ExplainRequest
		if err := decodeModelRequest(r, &request, r.PathValue("chapter") != ""); err != nil {
			return nil, err
		}
		if r.PathValue("chapter") != "" {
			if err := applyPassageURL(&request, r); err != nil {
				return nil, err
			}
		}
		return s.Explain(r.Context(), request, access)
	})
}

func decodeModelRequest(r *http.Request, request any, allowEmpty bool) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(request); err != nil && !(allowEmpty && err == io.EOF) {
		return invalidModelRequest("Invalid model request")
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return invalidModelRequest("Expected one JSON object")
	}
	return nil
}

func applyPassageURL(request *model.ExplainRequest, r *http.Request) error {
	chapter, e1 := strconv.Atoi(r.PathValue("chapter"))
	verse, e2 := strconv.Atoi(r.PathValue("verse"))
	if e1 != nil || e2 != nil || chapter < 1 || verse < 1 {
		return invalidModelRequest("Invalid passage")
	}
	c := &request.Context
	if (c.Source != "" && c.Source != "quran") || (c.Chapter != 0 && c.Chapter != chapter) || (c.Verse != 0 && c.Verse != verse) {
		return invalidModelRequest("Passage context does not match URL")
	}
	c.Source = "quran"
	c.Chapter = chapter
	c.Verse = verse
	return nil
}

// modelHTTP owns HTTP concerns only: method, authentication, decoding bounds and errors.
func modelHTTP(keys *KeyStore, operation func(*http.Request, GenerationAccess) (*model.Result, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		var emit func(model.Event)
		fail := func(status int, message string) {
			if emit != nil {
				emit(model.Event{Type: "error", Text: message})
				return
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"message": message})
		}
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			fail(405, "POST required")
			return
		}
		access, err := authenticateModel(r, keys)
		if err != nil {
			if errors.Is(err, errKeyUnauthorized) {
				w.Header().Set("WWW-Authenticate", "Bearer")
				fail(401, "Requires a valid Authorization: Bearer <key>")
			} else {
				log.Printf("Model authentication: %v", err)
				fail(503, "Model access unavailable")
			}
			return
		}
		if r.Header.Get("Accept") == "application/x-ndjson" {
			ctx, cancel := context.WithCancel(r.Context())
			defer cancel()
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("X-Accel-Buffering", "no")
			encoder := json.NewEncoder(w)
			controller := http.NewResponseController(w)
			emit = func(event model.Event) {
				if ctx.Err() != nil {
					return
				}
				if err := encoder.Encode(event); err != nil {
					cancel()
					return
				}
				if err := controller.Flush(); err != nil {
					cancel()
				}
			}
			r = r.WithContext(model.WithEvents(ctx, emit))
		}
		r.Body = http.MaxBytesReader(w, r.Body, model.ModelBodyLimit())
		result, err := operation(r, access)
		if err != nil {
			var service *modelServiceError
			var limit *model.LimitError
			switch {
			case errors.As(err, &service):
				status := 503
				switch service.kind {
				case "invalid":
					status = 400
				case "not_found":
					status = 404
				}
				fail(status, service.Error())
			case errors.Is(err, errKeyLimit), errors.As(err, &limit):
				fail(429, err.Error())
			default:
				log.Printf("Model request failed: %v", err)
				fail(503, "Model unavailable. Please try again later.")
			}
			return
		}
		if emit != nil {
			emit(model.Event{Type: "result", Result: result})
		} else {
			json.NewEncoder(w).Encode(result)
		}
	})
}

func authenticateModel(r *http.Request, keys *KeyStore) (GenerationAccess, error) {
	if keys == nil {
		return nil, nil
	}
	fields := strings.Fields(r.Header.Get("Authorization"))
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		return nil, errKeyUnauthorized
	}
	id, err := keys.Authenticate(fields[1])
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context) (context.Context, error) {
		if err := keys.Reserve(id, time.Now()); err != nil {
			return nil, err
		}
		ctx = model.WithClient(ctx, "")
		return model.WithTokenAllowance(ctx, func(units int) error { return keys.ReserveTokens(id, units, time.Now()) }), nil
	}, nil
}
