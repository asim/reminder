package search

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/asim/reminder/quran"
)

// WithQuran resolves navigation queries from the corpus before consulting an
// embedding or keyword index. Indexing itself still uses the original Searcher.
func WithQuran(index Searcher, corpus *quran.Quran) Searcher {
	return &quranLookup{Searcher: index, corpus: corpus}
}

type quranLookup struct {
	Searcher
	corpus *quran.Quran
}

var verseReference = regexp.MustCompile(`^(.+?)\s*[:/]\s*([0-9]+)$`)
var namedVerse = regexp.MustCompile(`^(.+?)\s+([0-9]+)$`)

func chapterQuery(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, prefix := range []string{"surah ", "surat ", "sura ", "chapter ", "quran "} {
		s = strings.TrimPrefix(s, prefix)
	}
	return s
}
func chapterKey(s string) string {
	s = chapterQuery(s)
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, s)
}
func (l *quranLookup) chapter(s string) *quran.Chapter {
	key := chapterKey(s)
	if number, err := strconv.Atoi(chapterQuery(s)); err == nil {
		for _, ch := range l.corpus.Chapters {
			if ch.Number == number {
				return ch
			}
		}
		return nil
	}
	aliases := map[string]int{"fatiha": 1, "alfatiha": 1, "yaseen": 36, "yaaseen": 36, "yasin": 36, "ikhlas": 112, "alikhlas": 112}
	for _, ch := range l.corpus.Chapters {
		name := chapterKey(ch.Name)
		if key == name || key == strings.TrimPrefix(name, "al") || key == chapterKey(ch.English) || aliases[key] == ch.Number {
			return ch
		}
	}
	return nil
}
func (l *quranLookup) Query(query string) ([]*Result, error) {
	query = strings.TrimSpace(query)
	if ch := l.chapter(query); ch != nil {
		return []*Result{{Text: fmt.Sprintf("%s — %s. Chapter %d, %d verses.", ch.Name, ch.English, ch.Number, ch.VerseCount), Score: 1, Metadata: map[string]string{"source": "quran", "kind": "chapter", "chapter": strconv.Itoa(ch.Number), "name": ch.Name}}}, nil
	}
	match := verseReference.FindStringSubmatch(query)
	if match == nil {
		match = namedVerse.FindStringSubmatch(query)
	}
	if match != nil {
		if ch := l.chapter(match[1]); ch != nil {
			number, err := strconv.Atoi(match[2])
			if err != nil {
				return []*Result{}, nil
			}
			verse := ch.Verse(number)
			if verse == nil {
				return []*Result{}, nil
			}
			return []*Result{{Text: verse.Text, Score: 1, Metadata: map[string]string{"source": "quran", "chapter": strconv.Itoa(ch.Number), "verse": strconv.Itoa(verse.Number), "name": ch.Name}}}, nil
		}
	}
	// Invalid numeric references must not produce unrelated semantic matches.
	if _, err := strconv.Atoi(chapterQuery(query)); err == nil {
		return []*Result{}, nil
	}
	if match != nil {
		if _, err := strconv.Atoi(chapterKey(match[1])); err == nil {
			return []*Result{}, nil
		}
	}
	return l.Searcher.Query(query)
}
