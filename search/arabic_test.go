package search

import (
	"github.com/asim/reminder/quran"
	"strings"
	"testing"
)

func TestLexicalEvidence(t *testing.T) {
	q := quran.Load()
	refs, e := LexicalReferences(q, "Explain الحمد in 1:2")
	if e != nil {
		t.Fatal(e)
	}
	morphology, lexicon := false, false
	seen := map[string]bool{}
	for _, r := range refs {
		switch r.Metadata["source"] {
		case "morphology":
			morphology = true
			if !strings.Contains(r.Text, "حمد") {
				t.Fatal("missing root")
			}
		case "lexicon":
			lexicon = true
			id := r.Metadata["entry_id"]
			if id == "" || seen[id] {
				t.Fatal("missing/duplicate identity")
			}
			seen[id] = true
			if len(r.Text) > 6000 {
				t.Fatal("unbounded lexical record")
			}
		}
	}
	if !morphology || !lexicon {
		t.Fatalf("missing evidence %v %v", morphology, lexicon)
	}
	refs, e = LexicalReferences(q, "patience")
	if e != nil || len(refs) != 0 {
		t.Fatal("ordinary search gained unrelated lexical evidence")
	}
	refs, e = LexicalReferences(q, "صبر")
	if e != nil || len(refs) == 0 {
		t.Fatal("direct root lookup missing")
	}
	refs, e = LexicalReferences(q, "999:999")
	if e != nil || len(refs) != 0 {
		t.Fatal("invalid verse produced evidence")
	}
}
