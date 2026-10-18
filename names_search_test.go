package main

import (
	"strings"
	"testing"

	"github.com/asim/reminder/names"
	"github.com/asim/reminder/search"
)

func TestOldNameSearchResultsUseCurrentCorpus(t *testing.T) {
	corpus := names.Load()
	name := corpus.Get(68)
	results := currentNameResults([]*search.Result{
		{Text: "old meaning", Score: 0.8, Metadata: map[string]string{"Source": "names", "Arabic": name.Arabic, "Meaning": "old meaning"}},
		{Text: "unrelated", Metadata: map[string]string{"source": "quran"}},
		{Text: "unknown name", Metadata: map[string]string{"source": "names", "arabic": "unknown"}},
	}, corpus)
	if len(results) != 2 {
		t.Fatalf("got %d results", len(results))
	}
	if results[0].Metadata["meaning"] != name.Meaning || !strings.Contains(results[0].Text, name.Summary) || strings.Contains(results[0].Text, "old meaning") {
		t.Fatal("stale name wording survived")
	}
	if results[0].Score != 0.8 || results[1].Text != "unrelated" {
		t.Fatal("changed ranking or unrelated result")
	}
}
