package search

import (
	"errors"
	"testing"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
)

type evidenceIndex struct {
	results []*Result
	err     error
}

func (i evidenceIndex) Store(map[string]string, ...string) error { return nil }
func (i evidenceIndex) Query(string) ([]*Result, error)          { return i.results, i.err }

func TestModelReferencesResolveBundledArabic(t *testing.T) {
	q, b, n := quran.Load(), hadith.Load(), names.Load()
	verse := &Result{Text: "stale excerpt", Score: 0.8, Metadata: map[string]string{"source": "quran", "chapter": "2", "verse": "1", "arabic": "wrong"}}
	var h *hadith.Hadith
	for _, candidate := range b.Get(1).Hadiths {
		if candidate.Number > 0 {
			h = candidate
			break
		}
	}
	if h == nil {
		t.Fatal("no numbered hadith")
	}
	refs := []*Result{verse, verse, {Text: "old hadith", Metadata: map[string]string{"source": "bukhari", "book_num": "1", "number": "1"}}, {Metadata: map[string]string{"source": "names", "arabic": (*n)[0].Arabic}}}
	out, err := ModelReferences(evidenceIndex{results: refs}, n, q, b, "test")
	if err != nil || len(out) != 3 {
		t.Fatalf("results %v, %v", out, err)
	}
	if out[0].Text != q.Get(2).Verse(1).Text || out[0].Metadata["arabic"] != q.Get(2).Verse(1).Arabic || out[0].Score != verse.Score {
		t.Fatal("Quran evidence differs from numbered verse")
	}
	if out[1].Text != h.English || out[1].Metadata["arabic"] != h.Arabic || h.Arabic == "" {
		t.Fatal("hadith evidence differs from source")
	}
	if out[2].Metadata["arabic"] != (*n)[0].Arabic {
		t.Fatal("missing Arabic name")
	}
	if verse.Text != "stale excerpt" || verse.Metadata["arabic"] != "wrong" {
		t.Fatal("enrichment mutated index result")
	}
}

func TestModelReferencesRejectUnresolvedSources(t *testing.T) {
	q, b, n := quran.Load(), hadith.Load(), names.Load()
	var refs []*Result
	for _, md := range []map[string]string{
		{"source": "quran", "chapter": "999", "verse": "1"},
		{"source": "quran", "chapter": "1", "verse": "999"},
		{"source": "quran", "chapter": "1", "verse": "0"},
		{"source": "quran", "chapter": "invalid", "verse": "1"},
		{"source": "quran", "chapter": "1"},
		{"source": "bukhari", "book_num": "999", "number": "1"},
		{"source": "bukhari", "book_num": "1", "number": "999999"},
		{"source": "bukhari", "book_num": "1", "number": "0"},
		{"source": "tafsir", "chapter": "1", "verse": "1"},
	} {
		refs = append(refs, &Result{Text: "unresolved", Metadata: md})
	}
	out, err := ModelReferences(evidenceIndex{results: refs}, n, q, b, "test")
	if err != nil || len(out) != 0 {
		t.Fatalf("unresolved evidence included: %v, %v", out, err)
	}
	want := errors.New("index unavailable")
	if _, err := ModelReferences(evidenceIndex{err: want}, n, q, b, "test"); !errors.Is(err, want) {
		t.Fatalf("lost retrieval error: %v", err)
	}
}
