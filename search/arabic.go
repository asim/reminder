package search

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/asim/reminder/arabic"
	"github.com/asim/reminder/quran"
)

// LexicalReferences adds bounded, attributed evidence only when a question
// supplies Arabic or an explicit verse. Plain Search keeps its existing corpus.
func LexicalReferences(q *quran.Quran, query string) ([]*Result, error) {
	out := []*Result{}
	seenEntries := map[string]bool{}
	seenQueries := map[string]bool{}
	add := func(word string) error {
		key := arabic.Normalize(word)
		if key == "" || seenQueries[key] || len(seenQueries) >= 12 {
			return nil
		}
		seenQueries[key] = true
		entries, _, e := arabic.Lookup(word, 0, 5)
		if e != nil {
			return e
		}
		for _, entry := range entries {
			if seenEntries[entry.ID] || len(entry.Text) > 6000 || len(out) >= 10 {
				continue
			}
			seenEntries[entry.ID] = true
			out = append(out, &Result{Text: entry.Text, Metadata: map[string]string{"source": "lexicon", "title": "Lane’s Arabic-English Lexicon", "entry_id": entry.ID, "source_id": entry.SourceID, "root": entry.Root, "headword": entry.Word, "source_file": entry.SourceFile, "source_url": arabic.LaneURL, "licence": "CC BY-SA 3.0 US", "evidence_kind": "lexical senses, not a verse interpretation"}})
		}
		return nil
	}
	// Only exact Arabic tokens are looked up; no model-generated root guessing.
	for _, word := range strings.FieldsFunc(query, func(r rune) bool { return !unicode.In(r, unicode.Arabic) && !unicode.IsMark(r) }) {
		if e := add(word); e != nil {
			return nil, e
		}
	}
	seenVerses := map[string]bool{}
	for _, match := range explicitVerse.FindAllStringSubmatch(query, 2) {
		ch, _ := strconv.Atoi(match[1])
		v, _ := strconv.Atoi(match[2])
		if q == nil || ch < 1 || ch > len(q.Chapters) {
			continue
		}
		verse := q.Get(ch).Verse(v)
		if verse == nil {
			continue
		}
		key := fmt.Sprintf("%d:%d", ch, v)
		if seenVerses[key] {
			continue
		}
		seenVerses[key] = true
		analysis, e := arabic.Analyze(ch, v, verse.Arabic)
		if e != nil {
			return nil, e
		}
		// Do not present unverified word attachments as evidence for this edition.
		if !analysis.Aligned {
			continue
		}
		var evidence strings.Builder
		evidence.WriteString("Verified source word positions for Quran " + key + ". Roots and lemmas are corpus annotations; an absent root is unannotated, not inferred.\n")
		for _, word := range analysis.Words {
			fmt.Fprintf(&evidence, "Word %d: %s. ", word.Position, word.Arabic)
			for _, segment := range word.Segments {
				fmt.Fprintf(&evidence, "Segment %s; tag %s; lemma %s; root %s; grammar %s. ", segment.Arabic, segment.Tag, segment.Lemma, segment.Root, segment.Features)
			}
			evidence.WriteString("\n")
		}
		data := evidence.String()
		if len(data) <= 14000 {
			out = append(out, &Result{Text: data, Metadata: map[string]string{"source": "morphology", "title": "Quranic Arabic Corpus 0.4", "chapter": match[1], "verse": match[2], "source_url": arabic.CorpusURL, "alignment": "verified written tokens"}})
		}
		for _, word := range analysis.Words {
			for _, s := range word.Segments {
				if s.Lemma != "" {
					if e := add(s.Lemma); e != nil {
						return nil, e
					}
				}
				if s.Root != "" {
					if e := add(s.Root); e != nil {
						return nil, e
					}
				}
			}
		}
	}
	return out, nil
}
