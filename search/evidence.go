package search

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/tafsir"
)

var explicitVerse = regexp.MustCompile(`\b([0-9]{1,3}):([0-9]{1,3})\b`)

// ModelReferences resolves identifiers against bundled texts and adds attributed
// translations and verified tafsir. Stale index excerpts are never evidence.
func ModelReferences(idx Searcher, n *names.Names, q *quran.Quran, b *hadith.Collection, query string) ([]*Result, error) {
	var refs []*Result
	// Questions such as "Explain Quran 112:1" must retrieve that exact verse.
	for _, match := range explicitVerse.FindAllStringSubmatch(query, 8) {
		refs = append(refs, &Result{Score: 1, Metadata: map[string]string{"source": "quran", "chapter": match[1], "verse": match[2]}})
	}
	if len(refs) == 0 {
		var err error
		refs, err = References(WithQuran(idx, q), n, query)
		if err != nil {
			return nil, err
		}
	}
	return EnrichReferences(refs, n, q, b), nil
}

func tafsirDocument(c *tafsir.Collection, p *tafsir.Passage, score float32) *Result {
	ch, v := tafsir.Numbers(p.Ref)
	return &Result{Text: p.Text, Score: score, Metadata: map[string]string{
		"source": "tafsir", "edition": c.Edition.ID, "title": c.Edition.Title,
		"source_url": c.Edition.SourceURL, "chapter": strconv.Itoa(ch), "verse": strconv.Itoa(v),
		"ayahs": strings.Join(p.Ayahs, ", "), "ref": p.Ref,
	}}
}

// VerseReferences is also used by the reading page without a search/index call.
func VerseReferences(q *quran.Quran, chapter, number int, score float32) []*Result {
	if q == nil || chapter < 1 || chapter > len(q.Chapters) || number < 1 {
		return nil
	}
	ch := q.Get(chapter)
	if ch == nil {
		return nil
	}
	v := ch.Verse(number)
	if v == nil {
		return nil
	}
	makeRef := func(source, edition, title, text string) *Result {
		return &Result{Text: text, Score: score, Metadata: map[string]string{
			"source": source, "edition": edition, "translator": title, "chapter": strconv.Itoa(chapter), "verse": strconv.Itoa(number), "arabic": v.Arabic, "name": ch.Name,
		}}
	}
	refs := []*Result{makeRef("quran", "khattab", "Mustafa Khattab — The Clear Quran", v.Text)}
	for _, edition := range []struct{ ID, Title string }{{"sahih", "Saheeh International"}, {"asad", "Muhammad Asad — The Message of the Qur’an"}} {
		if text := v.Translations[edition.ID]; text != "" {
			refs = append(refs, makeRef("quran_translation", edition.ID, edition.Title, text))
		}
	}
	if q.Tafsir != nil && q.Tafsir.Edition.Verified {
		if p := q.Tafsir.Get(chapter, number); p != nil {
			refs = append(refs, tafsirDocument(q.Tafsir, p, score))
		}
	}
	return refs
}

func EnrichReferences(refs []*Result, n *names.Names, q *quran.Quran, b *hadith.Collection) []*Result {
	out := make([]*Result, 0, len(refs))
	seen := map[string]bool{}
	add := func(r *Result) {
		md := r.Metadata
		key := strings.Join([]string{md["source"], md["edition"], md["chapter"], md["verse"], md["book_num"], md["number"]}, ":")
		if !seen[key] {
			seen[key] = true
			out = append(out, r)
		}
	}
	for _, ref := range refs {
		if ref == nil {
			continue
		}
		md := ref.Metadata
		switch md["source"] {
		case "quran", "quran_translation", "tafsir":
			ch, _ := strconv.Atoi(md["chapter"])
			v, _ := strconv.Atoi(md["verse"])
			if md["kind"] == "chapter" && q != nil && ch >= 1 && ch <= len(q.Chapters) {
				for _, verse := range q.Get(ch).Verses {
					for _, r := range VerseReferences(q, ch, verse.Number, ref.Score) {
						add(r)
					}
				}
				continue
			}
			for _, r := range VerseReferences(q, ch, v, ref.Score) {
				add(r)
			}
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
			for _, h := range ch.Hadiths {
				if h.Number == number {
					add(&Result{Text: h.English, Score: ref.Score, Metadata: map[string]string{"source": "bukhari", "book_num": strconv.Itoa(book), "number": strconv.Itoa(number), "arabic": h.Arabic, "book": ch.Name, "narrator": h.Narrator}})
					break
				}
			}
		case "names":
			if n == nil {
				continue
			}
			for _, r := range CurrentNameResults([]*Result{ref}, n) {
				add(r)
			}
		}
	}
	return out
}

// CurrentTafsirResults refreshes search results even for prebuilt vector indexes.
func CurrentTafsirResults(refs []*Result, q *quran.Quran) []*Result {
	out := make([]*Result, 0, len(refs))
	seen := map[string]bool{}
	for _, r := range refs {
		if r.Metadata["source"] != "tafsir" {
			out = append(out, r)
			continue
		}
		ch, _ := strconv.Atoi(r.Metadata["chapter"])
		v, _ := strconv.Atoi(r.Metadata["verse"])
		if q == nil || q.Tafsir == nil {
			continue
		}
		p := q.Tafsir.Get(ch, v)
		if p == nil {
			continue
		}
		key := fmt.Sprintf("%s:%s", q.Tafsir.Edition.ID, p.Ref)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, tafsirDocument(q.Tafsir, p, r.Score))
	}
	return out
}

func SyncTafsir(fts *FTS, q *quran.Quran) (bool, error) {
	var docs []*Result
	for _, p := range q.Tafsir.Passages() {
		docs = append(docs, tafsirDocument(q.Tafsir, p, 0))
	}
	return fts.SyncSource("tafsir", docs)
}
