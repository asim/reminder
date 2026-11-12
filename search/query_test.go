package search

import (
	"encoding/json"
	"github.com/asim/reminder/model"
	"strings"
	"testing"
)

func TestContextKeepsLongHadithAndHonoursBudget(t *testing.T) {
	long := &Result{Text: strings.Repeat("a", 15881)}
	records := Context([]*Result{long, {Text: strings.Repeat("b", model.SearchEvidenceLimit()+1)}, {Text: "short"}})
	if len(records) != 2 {
		t.Fatalf("got %d records", len(records))
	}
	var first Result
	if err := json.Unmarshal([]byte(records[0]), &first); err != nil {
		t.Fatal(err)
	}
	if first.Text != long.Text {
		t.Fatal("long source truncated")
	}
	total := 0
	for _, r := range records {
		total += len(r)
	}
	if total > model.SearchEvidenceLimit() {
		t.Fatal("budget exceeded")
	}
}
