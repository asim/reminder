package arabic

import (
	"github.com/asim/reminder/quran"
	"testing"
)

func TestQULSourcesAndConflicts(t *testing.T) {
	q := quran.Load()
	a, e := Analyze(1, 2, q.Get(1).Verse(2).Arabic)
	if e != nil {
		t.Fatal(e)
	}
	if !a.QUL.PositionsVerified || a.QUL.Words["1"].Root.Status != "corroborated" || a.QUL.AyahRoot == nil {
		t.Fatalf("missing QUL corroboration: %+v", a.QUL)
	}
	for _, ref := range [][2]int{{5, 41}, {9, 124}} {
		a, e = Analyze(ref[0], ref[1], q.Get(ref[0]).Verse(ref[1]).Arabic)
		if e != nil {
			t.Fatal(e)
		}
		position := "13"
		if ref[0] == 9 {
			position = "15"
		}
		f := a.QUL.Words[position].Lemma
		if f.Status != "conflict" || f.Value == nil || *f.Value != "ءَامَنَ" {
			t.Fatalf("source disagreement hidden: %+v", f)
		}
	}
	a, e = Analyze(2, 181, q.Get(2).Verse(181).Arabic)
	if e != nil || a.QUL.PositionsVerified {
		t.Fatal("shifted source positions accepted")
	}
	a, e = Analyze(2, 1, q.Get(2).Verse(1).Arabic)
	if e != nil || a.QUL.AyahRoot != nil {
		t.Fatal("missing root summary inferred")
	}
	// Neither model/API clients can modify the cached source records.
	a, _ = Analyze(1, 2, q.Get(1).Verse(2).Arabic)
	*a.QUL.Words["1"].Root.Value = "changed"
	*a.QUL.AyahRoot = "changed"
	again, _ := Analyze(1, 2, q.Get(1).Verse(2).Arabic)
	if *again.QUL.Words["1"].Root.Value == "changed" || *again.QUL.AyahRoot == "changed" {
		t.Fatal("source cache mutated")
	}
}
func TestQULCoverage(t *testing.T) {
	q := quran.Load()
	positions, both := 0, 0
	statuses := map[string]int{}
	for _, ch := range q.Chapters {
		for _, v := range ch.Verses {
			if v.Number < 1 {
				continue
			}
			a, e := Analyze(ch.Number, v.Number, v.Arabic)
			if e != nil {
				t.Fatal(e)
			}
			if a.QUL.PositionsVerified {
				positions++
				if a.Aligned {
					both++
				}
			}
			for _, w := range a.QUL.Words {
				statuses[w.Lemma.Status]++
			}
		}
	}
	t.Logf("QUL position-verified verses: %d; also reader-aligned: %d; lemma statuses: %v", positions, both, statuses)
	if both < 5000 {
		t.Fatal("unexpected QUL alignment loss")
	}
}
