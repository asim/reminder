package quran

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

func TestArabicIntegrityAndVerseLookup(t *testing.T) {
	q := Load()
	count := 0
	for _, ch := range q.Chapters {
		if ch.Verse(0) != nil || ch.Verse(ch.VerseCount+1) != nil {
			t.Fatalf("invalid lookup in chapter %d", ch.Number)
		}
		if q.Index().Chapters[ch.Number-1].VerseCount != ch.VerseCount {
			t.Fatal("index count differs")
		}
		for _, v := range ch.Verses {
			if v.Number == 0 {
				if ch.Number == 1 || ch.Number == 9 || v.Arabic != Bismillah {
					t.Fatal("incorrect opening")
				}
				continue
			}
			count++
			if ch.Verse(v.Number) != v {
				t.Fatalf("wrong verse returned for %d:%d", ch.Number, v.Number)
			}
			var words []string
			for _, w := range v.Words {
				words = append(words, w.Arabic)
			}
			if strings.Join(words, " ") != v.Arabic {
				t.Fatalf("Arabic changes in word view: %d:%d", ch.Number, v.Number)
			}
		}
	}
	if count != 6236 {
		t.Fatalf("got %d ayahs", count)
	}
	if q.Get(1).Verse(1).Arabic != Bismillah {
		t.Fatal("Bismillah differs from 1:1")
	}
	// Pins the corpus verified against QuranEnc on 2026-09-25. Changes require
	// a fresh source comparison, not merely updating this hash.
	data, _ := files.ReadFile("data/arabic.json")
	if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != "1bf198a7ff2e948cf8c882755adb4b5142c06f12cb8f9e8ef0762c068997c195" {
		t.Fatalf("Arabic corpus changed: %s", got)
	}
}

func TestTranslationCoverage(t *testing.T) {
	q := Load()
	for _, ch := range q.Chapters {
		for _, v := range ch.Verses {
			if strings.TrimSpace(v.Translations["asad"]) == "" {
				t.Fatalf("missing Asad translation at %d:%d", ch.Number, v.Number)
			}
			if strings.Contains(v.Text, "Most Beneficent") {
				t.Fatalf("modified Khattab wording at %d:%d", ch.Number, v.Number)
			}
			if v.Number == 0 && v.Translations["asad"] != q.Get(1).Verse(1).Translations["asad"] {
				t.Fatalf("incorrect Asad opening in chapter %d", ch.Number)
			}
		}
	}
	data, _ := files.ReadFile("data/asad.json")
	if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != "6d6c7544e4d83009ed45f5ba2ee3b49881bfc32724516e89199173b065eaf2e3" {
		t.Fatalf("Asad source snapshot changed: %s", got)
	}
}
