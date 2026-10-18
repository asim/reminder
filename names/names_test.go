package names

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestComposedNamesAndReferences(t *testing.T) {
	corpus := Load()
	if len(*corpus) != 99 {
		t.Fatalf("got %d names", len(*corpus))
	}
	var quran map[string][]struct {
		Verse int `json:"verse"`
	}
	data, err := os.ReadFile("../quran/data/arabic.json")
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(data, &quran); err != nil {
		t.Fatal(err)
	}
	refPattern := regexp.MustCompile(`^(\d+):(\d+)$`)
	seen := map[string]bool{}
	for i, name := range *corpus {
		if name.Number != i+1 || name.Meaning == "" || name.Description == "" || name.Summary == "" {
			t.Fatalf("incomplete name %d", i+1)
		}
		if seen[name.Arabic] {
			t.Fatalf("duplicate Arabic identity %d", name.Number)
		}
		seen[name.Arabic] = true
		data, err := files.ReadFile(fmt.Sprintf("data/%d.json", name.Number))
		if err != nil {
			t.Fatal(err)
		}
		var detail struct{ Name, Description, Summary string }
		if err = json.Unmarshal(data, &detail); err != nil {
			t.Fatal(err)
		}
		if detail.Name != name.English || detail.Description != name.Description || detail.Summary != name.Summary {
			t.Fatalf("composed fields drifted for %d", name.Number)
		}
		if len(name.Location) == 0 {
			t.Fatalf("missing references for %d", name.Number)
		}
		for _, ref := range name.Location {
			parts := refPattern.FindStringSubmatch(ref)
			if parts == nil {
				t.Fatalf("unsupported verse link %d: %q", name.Number, ref)
			}
			verse, _ := strconv.Atoi(parts[2])
			if verse < 1 || verse > len(quran[parts[1]]) {
				t.Fatalf("nonexistent verse %q", ref)
			}
		}
		if name.ReferenceNote != "" && !strings.Contains(name.HTML(), html.EscapeString(name.ReferenceNote)) {
			t.Fatalf("missing references in legacy page %d", name.Number)
		}
	}
	for _, number := range []int{20, 21, 22, 23, 24, 25, 29, 41, 49, 57, 58, 59, 61, 64, 65, 71, 72, 81, 86, 89, 90, 91, 92, 96, 98, 99} {
		if corpus.Get(number).ReferenceNote == "" {
			t.Errorf("name %d presents related evidence without qualification", number)
		}
	}
	if !strings.Contains(corpus.Markdown(), corpus.Get(64).ReferenceNote) {
		t.Fatal("markdown omitted reference distinction")
	}
}
