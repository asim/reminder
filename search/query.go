package search

import (
	"context"
	"encoding/json"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"sync"
)

// References resolves search hits against the current names corpus.
func References(index Searcher, corpus *names.Names, question string) ([]*Result, error) {
	results, err := index.Query(question)
	if err != nil {
		return nil, err
	}
	return CurrentNameResults(results, corpus), nil
}

// Context preserves the existing 8000-byte context budget for summaries.
func Context(results []*Result) []string {
	var size int
	var contexts []string
	for _, result := range results {
		if size >= 8000 {
			break
		}
		data, _ := json.Marshal(result)
		size += len(data)
		contexts = append(contexts, string(data))
	}
	return contexts
}

func Summarise(ctx context.Context, results []*Result, question string) (string, error) {
	return model.AskSafe(ctx, Context(results), question)
}

// History owns the in-memory question/answer history for search sessions.
type History struct {
	mu       sync.RWMutex
	sessions map[string][]string
}

func (h *History) Get(session string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return append([]string{}, h.sessions[session]...)
}

func (h *History) Add(session, question, answer string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.sessions == nil {
		h.sessions = make(map[string][]string)
	}
	h.sessions[session] = append([]string{question, answer}, h.sessions[session]...)
}
