package search

import (
	"errors"
	"testing"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/tafsir"
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
	if err != nil || len(out) != 5 {
		t.Fatalf("results %v, %v", out, err)
	}
	if out[0].Text != q.Get(2).Verse(1).Text || out[0].Metadata["arabic"] != q.Get(2).Verse(1).Arabic || out[0].Score != verse.Score {
		t.Fatal("Quran evidence differs from numbered verse")
	}
	if out[3].Text != h.English || out[3].Metadata["arabic"] != h.Arabic || h.Arabic == "" {
		t.Fatal("hadith evidence differs from source")
	}
	if out[4].Metadata["arabic"] != (*n)[0].Arabic {
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
		{"source": "tafsir", "chapter": "999", "verse": "1"},
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

func TestAttributedTranslationsAndTafsir(t *testing.T) {
	q := quran.Load()
	var err error
	q.Tafsir, err = tafsir.Parse([]byte(`{"1:1":{"text":"Verified test commentary","ayah_keys":["1:1","1:2"]},"1:2":"1:1"}`), tafsir.Edition{ID: "test", Title: "Test author", Verified: true})
	if err != nil {
		t.Fatal(err)
	}
	refs, err := ModelReferences(evidenceIndex{err: errors.New("must not query index")}, names.Load(), q, nil, "Compare Quran 1:1 and 1:2")
	if err != nil {
		t.Fatal(err)
	}
	if len(refs) != 7 {
		t.Fatalf("want six translations and one shared tafsir, got %d", len(refs))
	}
	if refs[1].Text != q.Get(1).Verse(1).Translations["sahih"] || refs[2].Metadata["translator"] == "" {
		t.Fatal("missing attributed translations")
	}
	if refs[3].Metadata["ayahs"] != "1:1, 1:2" || refs[3].Metadata["title"] != "Test author" {
		t.Fatal("lost tafsir attribution or coverage")
	}
	stale := []*Result{{Text: "old commentary", Metadata: map[string]string{"source": "tafsir", "chapter": "1", "verse": "2"}}}
	current := CurrentTafsirResults(stale, q)
	if len(current) != 1 || current[0].Text != "Verified test commentary" {
		t.Fatal("stale tafsir survived")
	}
	fts := NewFTS()
	if _, err := SyncTafsir(fts, q); err != nil {
		t.Fatal(err)
	}
	if changed, err := SyncTafsir(fts, q); err != nil || changed {
		t.Fatalf("unchanged corpus reindexed: %v", err)
	}
	q.Tafsir.Edition.Verified = false
	if got := VerseReferences(q, 1, 1, 1); len(got) != 3 {
		t.Fatal("unverified tafsir included as model evidence")
	}
}
