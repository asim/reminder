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
