package search

import (
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
