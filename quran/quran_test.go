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
