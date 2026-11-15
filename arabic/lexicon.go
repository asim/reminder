// Package arabic provides attributed lexical and morphological evidence.
// It never infers a verse's meaning from a root or guesses missing roots.
package arabic

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"io"
	"sort"
	"strings"
	"sync"
	"unicode"
)

//go:embed dict/lane-xml.zip dict/source.json morphology/qac.txt.gz morphology/source.json
var files embed.FS

const LaneURL = "https://github.com/alpheios-project/lan"
const CorpusURL = "https://corpus.quran.com"

type Entry struct {
	ID         string `json:"id"`
	SourceID   string `json:"source_id,omitempty"`
	Root       string `json:"root"`
	Word       string `json:"word"`
	Text       string `json:"text"`
	SourceFile string `json:"source_file"`
}

type lexicon struct {
	entries []*Entry
	words   map[string][]*Entry
	roots   map[string][]*Entry
}

var laneOnce sync.Once
var lane *lexicon
var laneErr error

// Normalize removes marks and spacing, preserving distinctions between letters.
func Normalize(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsMark(r) || unicode.IsSpace(r) || r == '\u0640' {
			return -1
		}
		return r
	}, norm.NFC.String(s))
}

func plain(s string) (string, error) {
	d := xml.NewDecoder(strings.NewReader("<text>" + s + "</text>"))
	var b strings.Builder
	for {
		t, e := d.Token()
		if e == io.EOF {
			break
		}
		if e != nil {
			return "", e
		}
		if c, ok := t.(xml.CharData); ok {
			b.Write(c)
		}
	}
	return strings.Join(strings.Fields(b.String()), " "), nil
}

func loadLane() (*lexicon, error) {
	laneOnce.Do(func() { lane, laneErr = parseLane() })
	return lane, laneErr
}
func parseLane() (*lexicon, error) {
	data, e := files.ReadFile("dict/lane-xml.zip")
	if e != nil {
		return nil, e
	}
	var manifest struct {
		Files   map[string]string `json:"files"`
		Entries int               `json:"entries"`
		SHA     string            `json:"archive_sha256"`
	}
	meta, _ := files.ReadFile("dict/source.json")
	if e = json.Unmarshal(meta, &manifest); e != nil {
		return nil, e
	}
	if fmt.Sprintf("%x", sha256.Sum256(data)) != manifest.SHA {
		return nil, fmt.Errorf("Lane archive checksum mismatch")
	}
	z, e := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if e != nil {
		return nil, e
	}
	out := &lexicon{words: map[string][]*Entry{}, roots: map[string][]*Entry{}}
	seen := map[string]bool{}
	for _, f := range z.File {
		if seen[f.Name] || manifest.Files[f.Name] == "" {
			return nil, fmt.Errorf("unexpected Lane file")
		}
		seen[f.Name] = true
		r, e := f.Open()
		if e != nil {
			return nil, e
		}
		raw, e := io.ReadAll(r)
		r.Close()
		if e != nil {
			return nil, e
		}
		if fmt.Sprintf("%x", sha256.Sum256(raw)) != manifest.Files[f.Name] {
			return nil, fmt.Errorf("Lane file checksum mismatch")
		}
		d := xml.NewDecoder(bytes.NewReader(raw))
		ordinal := 0
		for {
			tok, e := d.Token()
			if e == io.EOF {
				break
			}
			if e != nil {
				return nil, e
			}
			start, ok := tok.(xml.StartElement)
			if !ok || start.Name.Local != "div2" {
				continue
			}
			var group struct {
				Head struct {
					XML string `xml:",innerxml"`
				} `xml:"head"`
				Entries []struct {
					ID   string `xml:"id,attr"`
					Key  string `xml:"key,attr"`
					XML  string `xml:",innerxml"`
					Form struct {
						Orth string `xml:"orth"`
					} `xml:"form"`
				} `xml:"entryFree"`
			}
			if e = d.DecodeElement(&group, &start); e != nil {
				return nil, e
			}
			root, e := plain(group.Head.XML)
			if e != nil {
				return nil, e
			}
			for _, v := range group.Entries {
				ordinal++
				text, e := plain(v.XML)
				if e != nil {
					return nil, e
				}
				word := v.Key
				if word == "" {
					word = v.Form.Orth
				}
				entry := &Entry{ID: fmt.Sprintf("%s:%d", f.Name, ordinal), SourceID: v.ID, Root: root, Word: word, Text: text, SourceFile: f.Name}
				out.entries = append(out.entries, entry)
				out.words[Normalize(word)] = append(out.words[Normalize(word)], entry)
				if root != "" {
					out.roots[Normalize(root)] = append(out.roots[Normalize(root)], entry)
				}
			}
		}
	}
	if len(seen) != len(manifest.Files) || len(out.entries) != manifest.Entries {
		return nil, fmt.Errorf("incomplete Lane corpus")
	}
	return out, nil
}

// Lookup prefers exact headwords before related entries under the same root.
// Returned entries are copies; callers cannot mutate the shared index.
func Lookup(query string, offset, limit int) ([]Entry, int, error) {
	if limit < 1 || limit > 25 || offset < 0 || offset > 5000 || len(query) > 200 {
		return nil, 0, fmt.Errorf("invalid dictionary lookup")
	}
	l, e := loadLane()
	if e != nil {
		return nil, 0, e
	}
	key := Normalize(query)
	refs := append(append([]*Entry{}, l.words[key]...), l.roots[key]...)
	out := []Entry{}
	seen := map[string]bool{}
	for _, r := range refs {
		if !seen[r.ID] {
			seen[r.ID] = true
			out = append(out, *r)
		}
	}
	total := len(out)
	if offset >= total {
		return []Entry{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return out[offset:end], total, nil
}

func Roots(query string, limit int) ([]string, error) {
	if limit < 1 || limit > 100 || len(query) > 200 {
		return nil, fmt.Errorf("invalid root lookup")
	}
	l, e := loadLane()
	if e != nil {
		return nil, e
	}
	key := Normalize(query)
	out := []string{}
	for root := range l.roots {
		if strings.HasPrefix(root, key) {
			out = append(out, root)
		}
	}
	sort.Strings(out)
	if len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
