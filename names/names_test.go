package names

import (
	"encoding/json"
	"fmt"
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
		if regexp.MustCompile(`\d+:\d+`).MatchString(name.Description + " " + name.Summary) {
			t.Fatalf("verse citation in reading text for name %d", name.Number)
		}
		for _, ref := range name.Location {
			if !strings.Contains(name.HTML(), ">"+ref+"</a>") {
				t.Fatalf("missing reference link %s for name %d", ref, name.Number)
			}
		}
	}
}
