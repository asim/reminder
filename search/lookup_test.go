package search

import (
	"fmt"
	"github.com/asim/reminder/quran"
	"testing"
)

type lookupFallback struct{ calls int }

func (f *lookupFallback) Store(map[string]string, ...string) error { return nil }
func (f *lookupFallback) Query(string) ([]*Result, error) {
	f.calls++
	return []*Result{{Text: "topic result"}}, nil
}
func TestQuranNavigationLookup(t *testing.T) {
	corpus := quran.Load()
	fallback := &lookupFallback{}
	lookup := WithQuran(fallback, corpus)
	for _, query := range []string{"Al hashr", "Al-Hashr", "alhashr", "Hashr", "Surah Al Hashr", "The Exile", "59", "chapter 59"} {
		results, err := lookup.Query(query)
		if err != nil || len(results) != 1 || results[0].Metadata["chapter"] != "59" || results[0].Metadata["verse"] != "" {
			t.Fatalf("%q: %+v %v", query, results, err)
		}
	}
	for _, query := range []string{"59:22", "59/22", "Al Hashr 22", "Surah Al Hashr:22", "Quran 59:22"} {
		results, err := lookup.Query(query)
		if err != nil || len(results) != 1 || results[0].Metadata["verse"] != "22" || results[0].Text != corpus.Get(59).Verse(22).Text {
			t.Fatalf("%q: %+v %v", query, results, err)
		}
	}
	for _, query := range []string{"1:1", "1/1", "1 1"} {
		results, err := lookup.Query(query)
		if err != nil || len(results) != 1 || results[0].Metadata["chapter"] != "1" || results[0].Metadata["verse"] != "1" {
			t.Fatalf("reference mistaken for chapter: %q %+v", query, results)
		}
	}
	for _, query := range []string{"0", "115", "59:0", "59:25", "115:1", "Al Hashr 25"} {
		results, err := lookup.Query(query)
		if err != nil || len(results) != 0 {
			t.Fatalf("invalid reference %q: %+v %v", query, results, err)
		}
	}
	for _, ch := range corpus.Chapters {
		results, err := lookup.Query(ch.Name)
		if err != nil || len(results) != 1 || results[0].Metadata["chapter"] != fmt.Sprint(ch.Number) {
			t.Fatalf("chapter name %q not resolved", ch.Name)
		}
	}
	if fallback.calls != 0 {
		t.Fatalf("navigation called semantic index %d times", fallback.calls)
	}
	results, err := lookup.Query("how to be patient")
	if err != nil || len(results) != 1 || results[0].Text != "topic result" || fallback.calls != 1 {
		t.Fatal("topic fallback broken")
	}
}
