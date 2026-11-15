package arabic

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Data is grouped by verse and source word location. Array order is root, lemma,
// stem. Nil denotes absence, never an inferred value. Original field strings stay intact.
type qulVerse struct {
	Ayah  [3]*string            `json:"ayah"`
	Words map[string][3]*string `json:"words"`
}
type QULField struct {
	Value  *string `json:"value"`
	Status string  `json:"status"`
}
type QULWord struct {
	Root  QULField `json:"root"`
	Lemma QULField `json:"lemma"`
	Stem  QULField `json:"stem"`
}
type QULAnalysis struct {
	Source            string             `json:"source"`
	SourceURL         string             `json:"source_url"`
	PositionsVerified bool               `json:"positions_verified"`
	AyahRoot          *string            `json:"ayah_root"`
	AyahLemma         *string            `json:"ayah_lemma"`
	AyahStem          *string            `json:"ayah_stem"`
	Words             map[string]QULWord `json:"words"`
	Note              string             `json:"note"`
}

var qulOnce sync.Once
var qulChapters map[int]map[string]qulVerse
var qulErr error

func loadQUL() {
	qulOnce.Do(func() {
		qulChapters = map[int]map[string]qulVerse{}
		raw, e := files.ReadFile("qul/source.json")
		if e != nil {
			qulErr = e
			return
		}
		var manifest struct {
			Chapters map[string]string `json:"chapters"`
		}
		if e = json.Unmarshal(raw, &manifest); e != nil {
			qulErr = e
			return
		}
		if len(manifest.Chapters) != 114 {
			qulErr = fmt.Errorf("incomplete QUL manifest")
			return
		}
		for ch := 1; ch <= 114; ch++ {
			name := strconv.Itoa(ch) + ".json"
			raw, e := files.ReadFile("qul/" + name)
			if e != nil {
				qulErr = e
				return
			}
			if fmt.Sprintf("%x", sha256.Sum256(raw)) != manifest.Chapters[name] {
				qulErr = fmt.Errorf("QUL chapter checksum mismatch")
				return
			}
			var verses map[string]qulVerse
			if e = json.Unmarshal(raw, &verses); e != nil {
				qulErr = e
				return
			}
			qulChapters[ch] = verses
		}
	})
}
func sourceKey(s string) string { return Normalize(strings.ReplaceAll(s, "ٱ", "ا")) }
func assess(value *string, expected []string) QULField {
	if value == nil {
		return QULField{Status: "missing"}
	}
	// Copy pointers so callers cannot mutate cached source data.
	copyValue := *value
	f := QULField{Value: &copyValue, Status: "unverified"}
	if len(expected) == 0 {
		return f
	}
	unique := map[string]bool{}
	for _, e := range expected {
		unique[sourceKey(e)] = true
	}
	if len(unique) == 1 && unique[sourceKey(*value)] {
		f.Status = "corroborated"
	} else {
		f.Status = "conflict"
	}
	return f
}
func compareQUL(ch, verse int, source []Word) (*QULAnalysis, error) {
	loadQUL()
	if qulErr != nil {
		return nil, qulErr
	}
	v, ok := qulChapters[ch][strconv.Itoa(verse)]
	if !ok {
		return nil, fmt.Errorf("QUL verse unavailable")
	}
	out := &QULAnalysis{Source: "Quranic Universal Library", SourceURL: "https://qul.tarteel.ai/resources/morphology", Words: map[string]QULWord{}, PositionsVerified: len(v.Words) == len(source), Note: "Corroborated means agreement with the corpus at verified source positions, not independent scholarly validation. Ayah fields are source summaries, not positional word arrays. Missing or conflicting values are not repaired."}
	copyString := func(s *string) *string {
		if s == nil {
			return nil
		}
		v := *s
		return &v
	}
	out.AyahRoot = copyString(v.Ayah[0])
	out.AyahLemma = copyString(v.Ayah[1])
	out.AyahStem = copyString(v.Ayah[2])
	for position, values := range v.Words {
		n, e := strconv.Atoi(position)
		var expected [3][]string
		if e == nil && n > 0 && n <= len(source) {
			for _, s := range source[n-1].Segments {
				if s.Root != "" {
					expected[0] = append(expected[0], s.Root)
				}
				if s.Lemma != "" {
					expected[1] = append(expected[1], s.Lemma)
				}
				if strings.HasPrefix(s.Features, "STEM") {
					expected[2] = append(expected[2], s.Arabic)
				}
			}
		}
		word := QULWord{Root: assess(values[0], expected[0]), Lemma: assess(values[1], expected[1]), Stem: assess(values[2], expected[2])}
		if word.Stem.Status != "corroborated" {
			out.PositionsVerified = false
		}
		out.Words[position] = word
	}
	// A stem mismatch anywhere may indicate shifted positions. Preserve conflicts,
	// but do not claim corroboration of individual positions in that verse.
	if !out.PositionsVerified {
		for pos, w := range out.Words {
			for _, f := range []*QULField{&w.Root, &w.Lemma, &w.Stem} {
				if f.Status == "corroborated" {
					f.Status = "position_unverified"
				}
			}
			out.Words[pos] = w
		}
	}
	return out, nil
}
