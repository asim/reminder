package search

import (
	"strconv"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
)

// ModelReferences resolves search matches to the bundled source texts. Old
// indexes may contain excerpts or superseded translations; their identifiers,
// rather than their text, select the evidence supplied to the model.
func ModelReferences(idx Searcher, n *names.Names, q *quran.Quran, b *hadith.Collection, query string) ([]*Result, error) {
	refs, err := References(idx, n, query)
	if err != nil {
		return nil, err
	}
	out := make([]*Result, 0, len(refs))
	seen := make(map[string]bool)
	for _, ref := range refs {
		current := *ref
		current.Metadata = make(map[string]string, len(ref.Metadata)+1)
		for k, v := range ref.Metadata {
			current.Metadata[k] = v
		}
		md := current.Metadata
		var key string
		switch md["source"] {
		case "quran":
			chapter, _ := strconv.Atoi(md["chapter"])
			number, _ := strconv.Atoi(md["verse"])
			if q == nil || chapter < 1 || chapter > len(q.Chapters) || number < 1 {
				continue
			}
			ch := q.Get(chapter)
			if ch == nil {
				continue
			}
			verse := ch.Verse(number)
			if verse == nil {
				continue
			}
			current.Text = verse.Text
			md["arabic"] = verse.Arabic
			md["name"] = ch.Name
			key = "quran:" + strconv.Itoa(chapter) + ":" + strconv.Itoa(number)
		case "bukhari":
			book, _ := strconv.Atoi(md["book_num"])
			number, _ := strconv.Atoi(md["number"])
			if b == nil || book < 1 || number < 1 {
				continue
			}
			ch := b.Get(book)
			if ch == nil {
				continue
			}
			var passage *hadith.Hadith
			for _, h := range ch.Hadiths {
				if h.Number == number {
					passage = h
					break
				}
			}
			if passage == nil {
				continue
			}
			current.Text = passage.English
			md["arabic"] = passage.Arabic
			md["book"] = ch.Name
			md["narrator"] = passage.Narrator
			key = "bukhari:" + strconv.Itoa(book) + ":" + strconv.Itoa(number)
		case "names":
			// References already resolves Names against the current corpus, including
			// their Arabic identity and explanatory descriptions.
			key = "names:" + md["number"]
		default:
			continue
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, &current)
	}
	return out, nil
}
