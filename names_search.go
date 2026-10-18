package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/asim/reminder/names"
	"github.com/asim/reminder/search"
)

func nameDocument(name *names.Name) *search.Result {
	return &search.Result{
		Text: strings.Join([]string{name.Meaning, name.English, name.Arabic, name.Description, name.Summary}, " - "),
		Metadata: map[string]string{
			"source": "names", "number": fmt.Sprint(name.Number),
			"meaning": name.Meaning, "english": name.English, "arabic": name.Arabic,
		},
	}
}

func syncNames(fts *search.FTS, corpus *names.Names) (bool, error) {
	documents := make([]*search.Result, 0, len(*corpus))
	for _, name := range *corpus {
		documents = append(documents, nameDocument(name))
	}
	return fts.SyncSource("names", documents)
}

// Old vector indexes have no stable name ID, but retain the Arabic identity.
// Resolve their matches to the current corpus before showing or summarising
// them. Never send a superseded description back to a reader or model.
func currentNameResults(results []*search.Result, corpus *names.Names) []*search.Result {
	identity := func(s string) string {
		return strings.Map(func(r rune) rune {
			if unicode.IsMark(r) || unicode.IsSpace(r) || r == '\u0640' {
				return -1
			}
			return r
		}, s)
	}
	byArabic := make(map[string]*names.Name, len(*corpus))
	for _, name := range *corpus {
		byArabic[identity(name.Arabic)] = name
	}
	out := make([]*search.Result, 0, len(results))
	for _, result := range results {
		metadata := make(map[string]string, len(result.Metadata))
		for key, value := range result.Metadata {
			metadata[strings.ToLower(key)] = value
		}
		result.Metadata = metadata
		if result.Metadata["source"] != "names" {
			out = append(out, result)
			continue
		}
		if name := byArabic[identity(result.Metadata["arabic"])]; name != nil {
			current := nameDocument(name)
			current.Score = result.Score
			out = append(out, current)
		}
	}
	return out
}
