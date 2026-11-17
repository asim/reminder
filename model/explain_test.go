package model

import (
	"context"
	"strings"
	"testing"
)

func TestExplainFollowupRetrievalKeepsPassage(t *testing.T) {
	req := ExplainRequest{Context: ExplainContext{Source: "quran", Chapter: 2, Verse: 153, Translation: "sahih"}, Question: "What does that mean?", History: []Message{{Role: "user", Content: "What is patience?"}, {Role: "assistant", Content: "UNTRUSTED"}}}
	_, err := Explain(context.Background(), req, func(query string) ([]*Reference, error) {
		if !strings.Contains(query, "2:153") || !strings.Contains(query, "What is patience?") || strings.Contains(query, "UNTRUSTED") || strings.Contains(query, "Keep this passage") {
			t.Fatalf("bad query: %s", query)
		}
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, req := range []ExplainRequest{{}, {Context: ExplainContext{Source: "quran"}}, {Context: ExplainContext{Chapter: 115}}, {Question: "x", Q: "y"}, {Question: "x", Context: ExplainContext{Translation: "invented"}}} {
		if _, err := req.Request(); err == nil {
			t.Fatalf("accepted invalid request: %+v", req)
		}
	}
}

func TestHadithExplanationContext(t *testing.T) {
	for _, collection := range []string{"bukhari", "muslim"} {
		req := ExplainRequest{Context: ExplainContext{Source: "hadith", Collection: collection, Book: 0, Hadith: 7369}}
		got, err := req.Request()
		if err != nil || !strings.Contains(got.Prompt, collection+" book 0, record 7369") || strings.Contains(got.Prompt, "surrounding verses") {
			t.Fatalf("incorrect hadith prompt: %+v, %v", got, err)
		}
		req.Question = "What does that mean?"
		req.History = []Message{{Role: "user", Content: "What is the lesson?"}, {Role: "assistant", Content: "UNTRUSTED"}}
		_, err = Explain(context.Background(), req, func(query string) ([]*Reference, error) {
			if !strings.Contains(query, collection+" book 0 hadith 7369") || !strings.Contains(query, "What is the lesson?") || strings.Contains(query, "UNTRUSTED") {
				t.Fatalf("lost hadith follow-up context: %s", query)
			}
			return nil, nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, c := range []ExplainContext{
		{Source: "hadith", Collection: "invented", Book: 1, Hadith: 1},
		{Source: "hadith", Collection: "bukhari", Book: -1, Hadith: 1},
		{Source: "hadith", Collection: "bukhari", Book: 1},
		{Source: "hadith", Collection: "bukhari", Book: 1, Hadith: 1, Chapter: 2},
		{Source: "quran", Chapter: 2, Collection: "muslim"},
	} {
		if _, err := (ExplainRequest{Context: c}).Request(); err == nil {
			t.Fatalf("accepted mixed or invalid context: %+v", c)
		}
	}
}
