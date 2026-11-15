package arabic

import (
	"github.com/asim/reminder/quran"
	"testing"
)

func TestSources(t *testing.T) {
	entries, total, e := Lookup("صَبْر", 0, 3)
	if e != nil || len(entries) == 0 || total < 1 {
		t.Fatalf("Lane lookup %v %d %v", entries, total, e)
	}
	if entries[0].Root != "صبر" {
		t.Fatal(entries[0])
	}
	a, e := Analyze(1, 2, quran.Load().Get(1).Verse(2).Arabic)
	if e != nil {
		t.Fatal(e)
	}
	if !a.Aligned || a.Words[0].Segments[1].Root != "حمد" {
		t.Fatalf("bad analysis %+v", a)
	}
	entries[0].Text = "changed"
	again, _, _ := Lookup("صبر", 0, 3)
	if again[0].Text == "changed" {
		t.Fatal("mutated shared entry")
	}
	a.Words[0].Segments[0].Root = "changed"
	againA, _ := Analyze(1, 2, a.CanonicalArabic)
	if againA.Words[0].Segments[0].Root == "changed" {
		t.Fatal("mutated shared annotation")
	}
}
func TestAlignmentCoverage(t *testing.T) {
	q := quran.Load()
	aligned, words := 0, 0
	for _, ch := range q.Chapters {
		for _, v := range ch.Verses {
			if v.Number < 1 {
				continue
			}
			a, e := Analyze(ch.Number, v.Number, v.Arabic)
			if e != nil {
				t.Fatal(e)
			}
			words += len(a.Words)
			if a.Aligned {
				aligned++
			}
		}
	}
	t.Logf("Verified token alignment: %d / 6236 verses; %d source words", aligned, words)
	if aligned != 5832 {
		t.Fatal("unexpected alignment regression")
	}
	a, e := Analyze(1, 2, "different text")
	if e != nil || a.Aligned || a.AlignmentNote == "" {
		t.Fatal("unverified alignment accepted")
	}
}

func TestAlignmentNormalizationIsNarrow(t *testing.T) {
	if alignmentKey("فِى") != alignmentKey("فِي") || alignmentKey("أٓ") != alignmentKey("ءَا") {
		t.Fatal("known orthographic equivalence missing")
	}
	if alignmentKey("عَلَى") == alignmentKey("عَلَي") {
		t.Fatal("alif maqsura conflated with ya")
	}
	if Normalize("أمل") != Normalize("أمل") {
		t.Fatal("canonical Unicode equivalents diverge")
	}
	if Normalize("أمل") == Normalize("امل") || Normalize("على") == Normalize("علي") {
		t.Fatal("lexical letters conflated")
	}
}

func TestCorpusTransliterationAndMissingRoots(t *testing.T) {
	loadCorpus()
	if corpusErr != nil {
		t.Fatal(corpusErr)
	}
	for ref, words := range corpus {
		for _, word := range words {
			for _, s := range word.Segments {
				for _, r := range s.Arabic + s.Lemma + s.Root {
					if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
						t.Fatalf("unconverted Buckwalter %s: %+v", ref, s)
					}
				}
			}
		}
	}
	a, e := Analyze(1, 5, quran.Load().Get(1).Verse(5).Arabic)
	if e != nil {
		t.Fatal(e)
	}
	if a.Words[0].Segments[0].Root != "" {
		t.Fatal("root invented for pronoun")
	}
}
