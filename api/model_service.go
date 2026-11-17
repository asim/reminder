package api

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
)

// ModelService owns grounding and explanation policy independently of transport.
type ModelService struct {
	Index  search.Searcher
	Names  *names.Names
	Quran  *quran.Quran
	Hadith *hadith.Collection
	Ready  <-chan bool
	Cache  *ExplainCache
}

// GenerationAccess reserves a caller's allowance only when generation is needed.
// A nil callback retains the anonymous limits already attached to the context.
type GenerationAccess func(context.Context) (context.Context, error)

type modelServiceError struct {
	kind    string
	message string
}

func (e *modelServiceError) Error() string     { return e.message }
func invalidModelRequest(message string) error { return &modelServiceError{"invalid", message} }

func (s *ModelService) ready() error {
	select {
	case <-s.Ready:
		return nil
	default:
		return &modelServiceError{"unavailable", "Indexing content, please try again later"}
	}
}
func generationContext(ctx context.Context, access GenerationAccess) (context.Context, error) {
	if !model.Available() {
		return nil, &modelServiceError{"unavailable", "No model configured"}
	}
	if access != nil {
		return access(ctx)
	}
	return ctx, nil
}

func (s *ModelService) Ask(ctx context.Context, request model.Request, access GenerationAccess) (*model.Result, error) {
	if err := request.Validate(); err != nil {
		return nil, invalidModelRequest(err.Error())
	}
	if err := s.ready(); err != nil {
		return nil, err
	}
	ctx, err := generationContext(ctx, access)
	if err != nil {
		return nil, err
	}
	return model.Research(ctx, request, s.references(nil), search.Lookups(s.Index, s.Names, s.Quran, s.Hadith, ""))
}

func (s *ModelService) Explain(ctx context.Context, request model.ExplainRequest, access GenerationAccess) (*model.Result, error) {
	if request.Context.Chapter > 0 {
		if request.Context.Source == "" {
			request.Context.Source = "quran"
		}
		if request.Context.Translation == "" {
			request.Context.Translation = "khattab"
		}
	}
	if _, err := request.Request(); err != nil {
		return nil, invalidModelRequest(err.Error())
	}
	passage, err := search.ExplanationPassage(s.Quran, request.Context)
	if request.Context.Source == "hadith" {
		passage, err = search.HadithExplanationPassage(s.Hadith, request.Context)
	}
	if err != nil {
		return nil, &modelServiceError{"not_found", err.Error()}
	}
	if err := s.ready(); err != nil {
		return nil, err
	}
	model.Emit(ctx, model.Event{Type: "status", Text: "Loading explanation…"})
	return s.Cache.result(ctx, s.Cache.key(request), func() (*model.Result, error) {
		generation, err := generationContext(ctx, access)
		if err != nil {
			return nil, err
		}
		retrieve := s.references(passage)
		if len(passage) > 0 && request.Question == "" && request.Q == "" && request.Context.Selection == "" && len(request.History) == 0 {
			// Default Explain already has an exact anchor and context. Tools remain
			// available for gaps; avoid an unconditional embedding/ranked search first.
			retrieve = func(string) ([]*model.Reference, error) { return passage, nil }
		}
		result, err := model.Explain(generation, request, retrieve, search.Lookups(s.Index, s.Names, s.Quran, s.Hadith, ""))
		if err != nil {
			return nil, err
		}
		s.addCoverage(result, request.Context)
		return result, nil
	})
}

func (s *ModelService) references(passage []*model.Reference) func(string) ([]*model.Reference, error) {
	return func(query string) ([]*model.Reference, error) {
		refs, err := search.ModelReferences(s.Index, s.Names, s.Quran, s.Hadith, query)
		if err != nil {
			return nil, err
		}
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
}

func (s *ModelService) addCoverage(result *model.Result, focus model.ExplainContext) {
	if focus.Chapter == 0 || focus.Verse != 0 {
		return
	}
	seen := map[string]bool{}
	for _, ref := range result.References {
		md := ref.Metadata
		if md["chapter"] == strconv.Itoa(focus.Chapter) && (md["source"] == "quran" || md["source"] == "quran_translation") {
			seen[md["verse"]] = true
		}
	}
	total := 0
	for _, verse := range s.Quran.Get(focus.Chapter).Verses {
		if verse.Number > 0 {
			total++
		}
	}
	result.Coverage = map[string]int{"included_verses": len(seen), "total_verses": total}
}
