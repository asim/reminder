package search

import (
	"fmt"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"strings"
	"testing"
)

func TestSourceFilterBeforeLimit(t *testing.T) {
	idx := NewFTS()
	defer idx.Close()
	for i := 0; i < 40; i++ {
		if err := idx.Store(map[string]string{"source": "bukhari", "number": fmt.Sprint(i)}, "patience"); err != nil {
			t.Fatal(err)
		}
	}
	for _, source := range []string{"quran", "tafsir", "names", "muslim"} {
		if err := idx.Store(map[string]string{"source": source}, "patience and its rewards"); err != nil {
			t.Fatal(err)
		}
	}
	for _, source := range []string{"quran", "tafsir", "names", "hadith"} {
		refs, err := idx.Query(source + ":patience")
		if err != nil || len(refs) == 0 {
			t.Fatalf("%s buried: %v", source, err)
		}
		for _, ref := range refs {
			if !MatchesSource(source, ref.Metadata["source"]) {
				t.Fatal("source leaked")
			}
		}
	}
	if _, err := ScopedQuery("patience", "unknown"); err == nil {
		t.Fatal("invalid source accepted")
	}
}

func TestMentionInventoryIsCompleteAndDeduplicated(t *testing.T) {
	q := &quran.Quran{}
	for ch := 1; ch <= 30; ch++ {
		q.Chapters = append(q.Chapters, &quran.Chapter{Number: ch, Name: fmt.Sprint(ch), Verses: []*quran.Verse{
			{Number: 1, Text: "Moses spoke", Translations: map[string]string{"other": "Musa spoke"}},
			{Number: 2, Text: "Moses-like is not a test for a substring", Translations: map[string]string{}},
			{Number: 3, Text: "Mosesian is not a whole-word mention"},
		}})
	}
	refs, err := QuranMentions(q, []string{"Musa", "Moses"})
	if err != nil || len(refs) != 1 {
		t.Fatal(refs, err)
	}
	if refs[0].Metadata["matched_chapters"] != "30" || refs[0].Metadata["matched_verses"] != "60" || !strings.Contains(refs[0].Text, `"chapter":30`) {
		t.Fatal("truncated or double counted inventory", refs[0])
	}
	refs, err = QuranMentions(q, []string{"nonexistent"})
	if err != nil || refs[0].Metadata["matched_verses"] != "0" {
		t.Fatal("zero matches not represented")
	}
}

func TestMusaLookupAndScope(t *testing.T) {
	q := quran.Load()
	refs, err := QuranMentions(q, []string{"Musa", "Moses"})
	if err != nil || refs[0].Metadata["matched_chapters"] != "34" {
		t.Fatalf("Musa inventory: %+v %v", refs, err)
	}
	lookup := Lookups(evidenceIndex{}, &names.Names{}, q, nil, "hadith")
	refs, err = lookup(model.Lookup{Tool: "quran_mentions", Terms: []string{"Moses"}})
	if err == nil || len(refs) != 0 {
		t.Fatal("model widened source restriction or hid the restriction")
	}
}

func TestGeneralExplanationContextWithoutGrammar(t *testing.T) {
	q := quran.Load()
	refs, err := ExplanationPassage(q, model.ExplainContext{Chapter: 10, Verse: 106})
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, r := range refs {
		seen[r.Metadata["verse"]] = true
	}
	if !seen["105"] || !seen["106"] || !seen["107"] || refs[0].Metadata["verse"] != "106" {
		t.Fatal("missing neighbouring context or lost focus")
	}
	for _, tc := range []struct {
		query   string
		lexical bool
	}{{"Quran 10:106 What does it mean", false}, {"Quran 10:106 Explain the Arabic word", true}} {
		refs, err = ModelReferences(evidenceIndex{}, &names.Names{}, q, nil, tc.query)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, r := range refs {
			if r.Metadata["source"] == "morphology" {
				found = true
			}
		}
		if found != tc.lexical {
			t.Fatalf("grammar evidence for %q: %v", tc.query, found)
		}
	}
}

func TestArabicToolLookup(t *testing.T) {
	q := quran.Load()
	lookup := Lookups(evidenceIndex{}, &names.Names{}, q, nil, "quran")
	refs, err := lookup(model.Lookup{Tool: "arabic", Chapter: 10, Verse: 106})
	if err != nil || len(refs) == 0 {
		t.Fatal("missing Arabic tool evidence", err)
	}
	found := false
	for _, ref := range refs {
		if ref.Metadata["source"] == "morphology" {
			found = true
		}
	}
	if !found {
		t.Fatal("missing attributed verse analysis")
	}
	for _, call := range []model.Lookup{{Tool: "arabic"}, {Tool: "arabic", Chapter: 115, Verse: 1}, {Tool: "arabic", Chapter: 1, Verse: 99}} {
		if _, err := lookup(call); err == nil {
			t.Fatal("invalid Arabic lookup accepted")
		}
	}
	restricted := Lookups(evidenceIndex{}, &names.Names{}, q, nil, "hadith")
	if _, err := restricted(model.Lookup{Tool: "arabic", Chapter: 1, Verse: 2}); err == nil {
		t.Fatal("Arabic tool widened source restriction")
	}
}
