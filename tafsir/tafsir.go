// Package tafsir owns the commentary corpus independently of Quran translations.
package tafsir

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

//go:embed data/*.json
var files embed.FS

type Edition struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Language  string `json:"language"`
	SourceURL string `json:"source_url"`
	Verified  bool   `json:"verified"`
}
type Passage struct {
	Ref   string   `json:"ref"`
	Text  string   `json:"text"`
	Ayahs []string `json:"ayahs"`
}
type Collection struct {
	Edition  Edition
	byVerse  map[string]*Passage
	passages []*Passage
}

var once sync.Once
var bundled *Collection

// Load parses the embedded corpus once. No network, extraction or model call.
func Load() *Collection {
	once.Do(func() {
		data, err := files.ReadFile("data/commentary.json")
		if err != nil {
			panic(err)
		}
		meta, err := files.ReadFile("data/edition.json")
		if err != nil {
			panic(err)
		}
		var edition Edition
		if err := json.Unmarshal(meta, &edition); err != nil {
			panic(err)
		}
		bundled, err = Parse(data, edition)
		if err != nil {
			panic(err)
		}
	})
	return bundled
}

// Parse accepts QUL's verse-keyed objects and shared-passage pointers.
func Parse(data []byte, edition Edition) (*Collection, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	if len(raw) == 0 || edition.ID == "" || edition.Title == "" {
		return nil, fmt.Errorf("empty tafsir or edition")
	}
	c := &Collection{Edition: edition, byVerse: map[string]*Passage{}}
	pointers := map[string]string{}
	groups := map[string][]string{}
	for ref, value := range raw {
		if ch, v := Numbers(ref); ch < 1 || ch > 114 || v < 1 || ref != fmt.Sprintf("%d:%d", ch, v) {
			return nil, fmt.Errorf("invalid reference %q", ref)
		}
		var target string
		if json.Unmarshal(value, &target) == nil {
			pointers[ref] = target
			continue
		}
		var entry struct {
			Text  string   `json:"text"`
			Ayahs []string `json:"ayah_keys"`
		}
		if err := json.Unmarshal(value, &entry); err != nil || strings.TrimSpace(entry.Text) == "" {
			return nil, fmt.Errorf("invalid passage %s", ref)
		}
		p := &Passage{Ref: ref, Text: entry.Text}
		c.byVerse[ref] = p
		c.passages = append(c.passages, p)
		groups[ref] = entry.Ayahs
	}
	for ref := range pointers {
		target := ref
		seen := map[string]bool{}
		for c.byVerse[target] == nil {
			if seen[target] {
				return nil, fmt.Errorf("cyclic reference %s", ref)
			}
			seen[target] = true
			next, ok := pointers[target]
			if !ok {
				return nil, fmt.Errorf("missing passage %s", target)
			}
			target = next
		}
		c.byVerse[ref] = c.byVerse[target]
	}
	for ref, p := range c.byVerse {
		p.Ayahs = append(p.Ayahs, ref)
	}
	for ref, keys := range groups {
		for _, key := range keys {
			if c.byVerse[key] != c.byVerse[ref] {
				return nil, fmt.Errorf("inconsistent passage group %s", ref)
			}
		}
	}
	less := func(a, b string) bool { ac, av := Numbers(a); bc, bv := Numbers(b); return ac*1000+av < bc*1000+bv }
	sort.Slice(c.passages, func(i, j int) bool { return less(c.passages[i].Ref, c.passages[j].Ref) })
	for _, p := range c.passages {
		sort.Slice(p.Ayahs, func(i, j int) bool { return less(p.Ayahs[i], p.Ayahs[j]) })
	}
	return c, nil
}
func Numbers(ref string) (int, int) {
	parts := strings.Split(ref, ":")
	if len(parts) != 2 {
		return 0, 0
	}
	ch, _ := strconv.Atoi(parts[0])
	v, _ := strconv.Atoi(parts[1])
	return ch, v
}
func (c *Collection) Get(chapter, verse int) *Passage {
	if c == nil {
		return nil
	}
	return c.byVerse[fmt.Sprintf("%d:%d", chapter, verse)]
}
func (c *Collection) Passages() []*Passage { return append([]*Passage(nil), c.passages...) }
