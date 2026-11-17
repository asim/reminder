package search

import (
	"fmt"
	"strings"
	"unicode"
)

func ValidSource(source string) bool {
	switch source {
	case "", "all", "quran", "hadith", "tafsir", "names", "bukhari", "muslim":
		return true
	}
	return false
}

func ScopedQuery(query, source string) (string, error) {
	if !ValidSource(source) {
		return "", fmt.Errorf("unknown search source")
	}
	if source == "" || source == "all" {
		return query, nil
	}
	_, text := collectionQuery(query)
	return source + ":" + text, nil
}

func MatchesSource(filter, source string) bool {
	return filter == "" || filter == "all" || filter == source || filter == "hadith" && (source == "bukhari" || source == "muslim") || filter == "quran" && source == "quran_translation"
}

func FilterSource(refs []*Result, source string) []*Result {
	out := make([]*Result, 0, len(refs))
	for _, ref := range refs {
		if ref != nil && MatchesSource(source, ref.Metadata["source"]) {
			out = append(out, ref)
		}
	}
	return out
}

// Grammar aids should answer a wording question, not displace a verse explanation.
func WantsLexical(query string) bool {
	for _, word := range tokenize(query) {
		switch word {
		case "arabic", "word", "words", "wording", "root", "lemma", "grammar", "grammatical", "morphology", "lexical", "conjunction", "particle", "translate", "translation":
			return true
		}
	}
	return strings.ContainsFunc(query, func(r rune) bool { return unicode.In(r, unicode.Arabic) })
}
