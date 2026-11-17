package model

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ExplainContext identifies the passage independently of the user's question.
type ExplainContext struct {
	Source      string `json:"source"`
	Collection  string `json:"collection,omitempty"`
	Book        int    `json:"book,omitempty"`
	Hadith      int    `json:"hadith,omitempty"`
	Chapter     int    `json:"chapter"`
	Verse       int    `json:"verse,omitempty"`
	Translation string `json:"translation,omitempty"`
	Selection   string `json:"selection,omitempty"`
}

type ExplainRequest struct {
	Context  ExplainContext `json:"context"`
	Question string         `json:"question,omitempty"`
	Q        string         `json:"q,omitempty"` // Legacy /api/explain input.
	History  []Message      `json:"history,omitempty"`
}

func (r ExplainRequest) Request() (Request, error) {
	if r.Question != "" && r.Q != "" {
		return Request{}, fmt.Errorf("provide question or q, not both")
	}
	c := r.Context
	if c.Source != "" && c.Source != "quran" && c.Source != "hadith" {
		return Request{}, fmt.Errorf("unsupported explanation source")
	}
	if c.Source == "hadith" {
		if (c.Collection != "bukhari" && c.Collection != "muslim") || c.Book < 0 || c.Hadith < 1 {
			return Request{}, fmt.Errorf("collection, book and hadith are required")
		}
		if c.Chapter != 0 || c.Verse != 0 || c.Translation != "" {
			return Request{}, fmt.Errorf("Quran fields are not valid for hadith context")
		}
	} else if c.Collection != "" || c.Book != 0 || c.Hadith != 0 {
		return Request{}, fmt.Errorf("hadith fields require hadith source")
	}
	if c.Chapter < 0 || c.Chapter > 114 || c.Verse < 0 || (c.Verse > 0 && c.Chapter == 0) {
		return Request{}, fmt.Errorf("invalid Quran context")
	}
	if c.Source == "quran" && c.Chapter == 0 {
		return Request{}, fmt.Errorf("chapter is required")
	}
	if c.Translation != "" && c.Translation != "khattab" && c.Translation != "sahih" && c.Translation != "asad" {
		return Request{}, fmt.Errorf("unknown translation")
	}
	if len(c.Selection) > 2000 {
		return Request{}, fmt.Errorf("selection is too long")
	}
	question := strings.TrimSpace(r.Question)
	if question == "" {
		question = strings.TrimSpace(r.Q)
	}
	if question == "" && c.Chapter == 0 && c.Source != "hadith" {
		return Request{}, fmt.Errorf("a passage context or question is required")
	}
	if question == "" {
		question = "Explain the overall meaning and message of this passage, using its surrounding verses and tafsir."
	}
	anchor := "Explain the supplied question using retrieved evidence."
	if c.Chapter > 0 {
		anchor = fmt.Sprintf("Explain Quran chapter %d", c.Chapter)
		if c.Verse > 0 {
			anchor = fmt.Sprintf("Explain Quran %d:%d", c.Chapter, c.Verse)
		} else {
			anchor += ". Give an overview limited to the verses actually supplied as evidence; explicitly state any incomplete coverage. Do not imply a complete chapter analysis."
		}
	}
	prompt := fmt.Sprintf("%s\nUse the supplied passage evidence first. Answer directly when this evidence is sufficient; use tools for specific missing evidence, not to re-fetch records already supplied. Keep this passage central. Explain its meaning and message in accessible English, using surrounding context and tafsir where helpful. Use Arabic analysis when it helps answer the question. For a general question, explain the whole passage even if a fragment is selected.\nSelected translation: %s\nSelected text (user selection): %s\nQuestion: %s", anchor, c.Translation, c.Selection, question)
	if c.Source == "hadith" {
		if strings.TrimSpace(r.Question+r.Q) == "" {
			question = "Explain the meaning and message of this hadith."
		}
		prompt = fmt.Sprintf("Explain %s book %d, record %d. Keep this exact narration central, including for follow-up questions. Explain its meaning in accessible English. Use the supplied Arabic and translation and look up related evidence when useful. No stored hadith commentary is supplied: distinguish your explanation from attributed scholarly commentary, and verify any attribution with retrieved sources. Record numbers identify this dataset; do not invent canonical numbering, authenticity grades or rulings.\nSelected text (user selection): %s\nQuestion: %s", c.Collection, c.Book, c.Hadith, c.Selection, question)
	}
	request := Request{Prompt: prompt, History: r.History}
	return request, request.Validate()
}

// Explain shares grounding, limits and citation validation with search summaries.
func Explain(ctx context.Context, request ExplainRequest, retrieve func(string) ([]*Reference, error), lookups ...LookupFunc) (*Result, error) {
	r, err := request.Request()
	if err != nil {
		return nil, err
	}
	var lookup LookupFunc
	if len(lookups) > 0 {
		lookup = lookups[0]
	}
	return Research(ctx, r, func(_ string) ([]*Reference, error) {
		query := strings.TrimSpace(request.Question + " " + request.Q + " " + request.Context.Selection)
		if request.Context.Chapter > 0 {
			if request.Context.Verse > 0 {
				query = fmt.Sprintf("Quran %d:%d %s", request.Context.Chapter, request.Context.Verse, query)
			} else {
				query = fmt.Sprintf("Quran chapter %d %s", request.Context.Chapter, query)
			}
		}
		if request.Context.Source == "hadith" {
			query = fmt.Sprintf("%s book %d hadith %d %s", request.Context.Collection, request.Context.Book, request.Context.Hadith, query)
		}
		for i := len(request.History) - 2; i >= 0 && i >= len(request.History)-4; i -= 2 {
			if len(query)+len(request.History[i].Content)+1 <= contentLimits.Prompt {
				query += "\n" + request.History[i].Content
			}
		}
		return retrieve(query)
	}, lookup)
}

// ExplanationVersion excludes credentials; deployment fingerprints additionally
// cover prompt implementation and embedded data. Epoch allows operator refresh.
func ExplanationVersion() string {
	provider, name, endpoint := "anthropic", "claude-haiku-4-5", ""
	if os.Getenv("OLLAMA_LLM_MODEL") != "" || (os.Getenv("OPENAI_API_KEY") == "" && os.Getenv("ANTHROPIC_API_KEY") == "") {
		provider = "ollama"
		name = os.Getenv("OLLAMA_LLM_MODEL")
		if name == "" {
			name = "gemma3:1b"
		}
		endpoint = os.Getenv("OLLAMA_BASE_URL")
		if endpoint == "" {
			endpoint = "http://localhost:11434/v1"
		}
	} else if os.Getenv("OPENAI_API_KEY") != "" {
		provider = "openai"
		name = os.Getenv("OPENAI_MODEL")
		if name == "" {
			name = "gpt-6-luna"
		}
	}
	// Bump the policy revision when explanation/retrieval behavior changes.
	data, _ := json.Marshal([]any{provider, name, endpoint, contentLimits, groundedPrompt, "explanation-policy-v2", os.Getenv("EXPLAIN_CACHE_EPOCH")})
	return string(data)
}
