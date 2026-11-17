package search

import (
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/quran"
	"testing"
)

func TestExplanationChapterPrioritizesEveryVerse(t *testing.T) {
	q := quran.Load()
	refs, err := ExplanationPassage(q, model.ExplainContext{Source: "quran", Chapter: 2, Translation: "sahih"})
	if err != nil {
		t.Fatal(err)
	}
	if len(refs) < 286 {
		t.Fatal("missing chapter text")
	}
	seen := map[string]bool{}
	for _, ref := range refs[:286] {
		if ref.Metadata["chapter"] != "2" || ref.Metadata["edition"] != "sahih" {
			t.Fatalf("wrong primary evidence: %+v", ref.Metadata)
		}
		seen[ref.Metadata["verse"]] = true
	}
	if len(seen) != 286 {
		t.Fatalf("chapter only has %d distinct verses", len(seen))
	}
	if _, err := ExplanationPassage(q, model.ExplainContext{Chapter: 2, Verse: 287}); err == nil {
		t.Fatal("accepted nonexistent verse")
	}
}

func TestHadithExplanationUsesExactRecord(t *testing.T) {
	b := hadith.Load()
	for _, corpus := range []*hadith.Collection{b, hadith.Muslim()} {
		book := corpus.Books[0]
		h := book.Hadiths[0]
		context := model.ExplainContext{Source: "hadith", Collection: corpus.ID, Book: book.Number, Hadith: h.Number}
		refs, err := HadithExplanationPassage(b, context)
		if err != nil || len(refs) != 1 {
			t.Fatalf("missing %s record: %v", corpus.ID, err)
		}
		if refs[0].Text != h.English || refs[0].Metadata["arabic"] != h.Arabic || refs[0].Metadata["collection"] != corpus.ID {
			t.Fatalf("wrong evidence: %+v", refs[0])
		}
		context.Hadith = 999999
		if _, err := HadithExplanationPassage(b, context); err == nil {
			t.Fatal("nonexistent record accepted")
		}
		context.Book = 999999
		if _, err := HadithExplanationPassage(b, context); err == nil {
			t.Fatal("nonexistent book accepted")
		}
	}
}
