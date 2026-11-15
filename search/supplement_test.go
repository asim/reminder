package search

import (
	"fmt"
	"testing"
)

func TestHybridRankingAndCollection(t *testing.T) {
	record := func(source, number string) *Result {
		return &Result{Text: number, Metadata: map[string]string{"source": source, "book_num": "1", "number": number}}
	}
	lexical, shared, semantic := record("muslim", "1"), record("bukhari", "2"), record("bukhari", "3")
	idx := WithSupplement(evidenceIndex{results: []*Result{lexical, shared}}, evidenceIndex{results: []*Result{semantic, shared, shared}})
	got, err := idx.Query("patience")
	if err != nil || len(got) != 3 || got[0].Text != shared.Text || got[1].Text != lexical.Text || got[2].Text != semantic.Text {
		t.Fatalf("fusion lost or misranked evidence: %v %v", got, err)
	}
	got, err = idx.Query("muslim: patience")
	if err != nil || len(got) != 1 || got[0].Text != lexical.Text {
		t.Fatal("collection restriction leaked", got, err)
	}
}

func TestHybridLegacyMetadataAndFailure(t *testing.T) {
	current := &Result{Text: "current", Metadata: map[string]string{"source": "bukhari", "book_num": "1", "number": "1"}}
	legacy := &Result{Text: "stale", Metadata: map[string]string{"Source": "bukhari", "Book_Num": "1", "Number": "1"}}
	other := &Result{Text: "semantic only", Metadata: map[string]string{"Source": "bukhari", "Book_Num": "1", "Number": "2"}}
	idx := WithSupplement(evidenceIndex{results: []*Result{current}}, evidenceIndex{results: []*Result{legacy, other}})
	got, err := idx.Query("bukhari: patience")
	if err != nil || len(got) != 2 || got[0].Text != "current" || got[1].Metadata["number"] != "2" {
		t.Fatalf("legacy fusion: %#v %v", got, err)
	}
	if legacy.Metadata["source"] != "" || current.Score != 0 || got[0].Score <= got[1].Score {
		t.Fatal("fusion mutated inputs or lost ranking")
	}
	idx = WithSupplement(evidenceIndex{results: []*Result{current}}, evidenceIndex{err: fmt.Errorf("vector unavailable")})
	got, err = idx.Query("patience")
	if err != nil || len(got) != 1 || got[0] != current {
		t.Fatalf("lost lexical fallback: %v %v", got, err)
	}
}
