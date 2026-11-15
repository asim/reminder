package search

import (
	"encoding/json"
	"sort"
	"strings"
)

type supplemented struct {
	Searcher
	extra Searcher
}

// WithSupplement combines independent lexical and semantic rankings using
// reciprocal rank fusion. Raw similarity and SQLite scores are not comparable.
func WithSupplement(primary, extra Searcher) Searcher { return &supplemented{primary, extra} }
func (s *supplemented) Query(query string) ([]*Result, error) {
	primary, err := s.Searcher.Query(query)
	if err != nil {
		return nil, err
	}
	source, semanticQuery := collectionQuery(query)
	if source == "" {
		semanticQuery = query
	}
	extra, err := s.extra.Query(semanticQuery)
	if err != nil {
		return primary, nil
	}
	type candidate struct {
		result *Result
		rank   float64
	}
	candidates := []*candidate{}
	seen := map[string]*candidate{}
	for _, set := range [][]*Result{primary, extra} {
		ranked := map[string]bool{}
		for _, r := range set {
			if r == nil {
				continue
			}
			changed := false
			m := make(map[string]string, len(r.Metadata))
			for key, value := range r.Metadata {
				m[strings.ToLower(key)] = value
				changed = changed || strings.ToLower(key) != key
			}
			if source != "" && m["source"] != source {
				continue
			}
			if changed {
				canonical := *r
				canonical.Metadata = m
				r = &canonical
			}
			identity := []string{m["source"], m["edition"], m["chapter"], m["verse"], m["book_num"], m["number"]}
			if m["source"] == "names" {
				identity = []string{"names", m["arabic"]}
			}
			keyBytes, _ := json.Marshal(identity)
			key := string(keyBytes)
			if ranked[key] {
				continue
			}
			ranked[key] = true
			c := seen[key]
			if c == nil {
				c = &candidate{result: r}
				seen[key] = c
				candidates = append(candidates, c)
			}
			c.rank += 1.0 / float64(60+len(ranked))
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool { return candidates[i].rank > candidates[j].rank })
	results := make([]*Result, 0, len(candidates))
	for _, c := range candidates {
		result := *c.result
		result.Score = float32(c.rank)
		results = append(results, &result)
	}
	return results, nil
}
