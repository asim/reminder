package search

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
)

// Lookups exposes only bounded, read-only corpus operations. The caller's scope
// is authoritative and cannot be widened by the model's proposed source.
func Lookups(idx Searcher, n *names.Names, q *quran.Quran, b *hadith.Collection, scope string) model.LookupFunc {
	return func(call model.Lookup) ([]*model.Reference, error) {
		switch call.Tool {
		case "search":
			if strings.TrimSpace(call.Query) == "" || len(call.Query) > 1000 {
				return nil, fmt.Errorf("invalid lookup query")
			}
			source := call.Source
			if scope != "" && scope != "all" {
				source = scope
			}
			query, err := ScopedQuery(call.Query, source)
			if err != nil {
				return nil, err
			}
			return ModelReferences(idx, n, q, b, query)
		case "passage":
			if scope != "" && scope != "all" && scope != "quran" && scope != "tafsir" {
				return nil, nil
			}
			if call.Verse < 1 {
				return nil, fmt.Errorf("passage lookup requires a verse")
			}
			refs, err := ExplanationPassage(q, model.ExplainContext{Chapter: call.Chapter, Verse: call.Verse})
			if err != nil {
				return nil, err
			}
			return FilterSource(refs, scope), nil
		case "quran_mentions":
			if scope != "" && scope != "all" && scope != "quran" {
				return nil, nil
			}
			return QuranMentions(q, call.Terms)
		default:
			return nil, fmt.Errorf("unknown lookup tool")
		}
	}
}

// QuranMentions scans the complete bundled Quran, not a top-N search. Each verse
// is counted once even when aliases match several translations of that verse.
func QuranMentions(q *quran.Quran, terms []string) ([]*Result, error) {
	if q == nil || len(terms) == 0 || len(terms) > 8 {
		return nil, fmt.Errorf("provide 1–8 mention terms")
	}
	needles := []string{}
	for _, term := range terms {
		words := mentionWords(term)
		if len(words) == 0 || len(term) > 100 {
			return nil, fmt.Errorf("invalid mention term")
		}
		needles = append(needles, " "+strings.Join(words, " ")+" ")
	}
	type chapterMatches struct {
		Chapter int    `json:"chapter"`
		Name    string `json:"name"`
		Verses  []int  `json:"verses"`
	}
	var chapters []chapterMatches
	total := 0
	for _, ch := range q.Chapters {
		matches := chapterMatches{Chapter: ch.Number, Name: ch.Name}
		for _, verse := range ch.Verses {
			if verse.Number < 1 {
				continue
			}
			texts := []string{verse.Text, verse.Arabic}
			for _, translation := range verse.Translations {
				texts = append(texts, translation)
			}
			found := false
			for _, text := range texts {
				haystack := " " + strings.Join(mentionWords(text), " ") + " "
				for _, needle := range needles {
					if strings.Contains(haystack, needle) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				matches.Verses = append(matches.Verses, verse.Number)
				total++
			}
		}
		if len(matches.Verses) > 0 {
			chapters = append(chapters, matches)
		}
	}
	data, _ := json.Marshal(chapters)
	text := fmt.Sprintf("Complete scan of the bundled Quran Arabic and English translations for explicit whole-word mentions of any of %q. Found %d distinct verses in %d chapters. Counts are verses, not word occurrences. This does not count unnamed references or thematic allusions. Matches (chapter, name, verse numbers): %s", terms, total, len(chapters), data)
	return []*Result{{Text: text, Score: 1, Metadata: map[string]string{"source": "quran", "kind": "mention_inventory", "title": "Quran mention lookup", "coverage": "complete corpus scan for supplied terms", "matched_verses": strconv.Itoa(total), "matched_chapters": strconv.Itoa(len(chapters))}}}, nil
}

// Strip vocalisation marks before tokenisation, retaining Arabic word boundaries.
func mentionWords(text string) []string {
	return tokenize(strings.Map(func(r rune) rune {
		if unicode.IsMark(r) || r == '\u0640' {
			return -1
		}
		return r
	}, text))
}
