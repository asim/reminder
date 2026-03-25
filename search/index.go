package search

import (
	"sort"
	"strings"
	"sync"
	"unicode"
)

type document struct {
	text     string
	lower    string
	metadata map[string]string
}

type Index struct {
	mu   sync.RWMutex
	docs []document
}

type Result struct {
	Text     string            `json:"text"`
	Score    float32           `json:"score"`
	Metadata map[string]string `json:"metadata"`
}

// tokenize splits text into lowercase words for matching.
func tokenize(s string) []string {
	s = strings.ToLower(s)
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func (i *Index) Store(md map[string]string, content ...string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, c := range content {
		if len(c) == 0 {
			continue
		}
		i.docs = append(i.docs, document{
			text:     c,
			lower:    strings.ToLower(c),
			metadata: md,
		})
	}
	return nil
}

func (i *Index) Query(q string) ([]*Result, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	queryWords := tokenize(q)
	if len(queryWords) == 0 {
		return nil, nil
	}

	queryLower := strings.ToLower(q)

	type scored struct {
		doc   document
		score float32
	}

	var matches []scored

	for _, doc := range i.docs {
		// Count how many query words appear in the document
		matchCount := 0
		for _, word := range queryWords {
			if strings.Contains(doc.lower, word) {
				matchCount++
			}
		}

		if matchCount == 0 {
			continue
		}

		// Base score: fraction of query words matched
		score := float32(matchCount) / float32(len(queryWords))

		// Boost for exact phrase match
		if strings.Contains(doc.lower, queryLower) {
			score = 1.0
		}

		matches = append(matches, scored{doc: doc, score: score})
	}

	// Sort by score descending
	sort.Slice(matches, func(a, b int) bool {
		return matches[a].score > matches[b].score
	})

	// Return top 25 results
	limit := 25
	if len(matches) < limit {
		limit = len(matches)
	}

	var results []*Result
	for _, m := range matches[:limit] {
		results = append(results, &Result{
			Text:     m.doc.text,
			Score:    m.score,
			Metadata: m.doc.metadata,
		})
	}

	return results, nil
}

func (i *Index) Count() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return len(i.docs)
}

func New(name string) *Index {
	return &Index{}
}
