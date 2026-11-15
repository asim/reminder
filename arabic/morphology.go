package arabic

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type Segment struct {
	Arabic   string `json:"arabic"`
	Tag      string `json:"tag"`
	Features string `json:"features"`
	Root     string `json:"root,omitempty"`
	Lemma    string `json:"lemma,omitempty"`
}
type Word struct {
	Position int       `json:"position"`
	Arabic   string    `json:"arabic"`
	Segments []Segment `json:"segments"`
}
type Analysis struct {
	Reference       string       `json:"reference"`
	Source          string       `json:"source"`
	SourceURL       string       `json:"source_url"`
	CanonicalArabic string       `json:"canonical_arabic"`
	Aligned         bool         `json:"aligned"`
	AlignmentNote   string       `json:"alignment_note,omitempty"`
	Words           []Word       `json:"words"`
	QUL             *QULAnalysis `json:"qul"`
}

var corpusOnce sync.Once
var corpus map[string][]Word
var corpusErr error

func buckwalter(s string) string {
	mapping := map[rune]rune{'_': 'ـ',
		'\'': 'ء', '|': 'آ', '>': 'أ', '&': 'ؤ', '<': 'إ', '}': 'ئ', 'A': 'ا', 'b': 'ب', 'p': 'ة', 't': 'ت', 'v': 'ث', 'j': 'ج', 'H': 'ح', 'x': 'خ', 'd': 'د', '*': 'ذ', 'r': 'ر', 'z': 'ز', 's': 'س', '$': 'ش', 'S': 'ص', 'D': 'ض', 'T': 'ط', 'Z': 'ظ', 'E': 'ع', 'g': 'غ', 'f': 'ف', 'q': 'ق', 'k': 'ك', 'l': 'ل', 'm': 'م', 'n': 'ن', 'h': 'ه', 'w': 'و', 'Y': 'ى', 'y': 'ي', 'F': 'ً', 'N': 'ٌ', 'K': 'ٍ', 'a': 'َ', 'u': 'ُ', 'i': 'ِ', '~': 'ّ', 'o': 'ْ', '`': 'ٰ', '{': 'ٱ', '^': 'ٓ', '#': 'ٔ', ':': 'ۜ', '@': '۟', '"': '۠', '[': 'ۢ', ';': 'ۣ', ',': 'ۥ', '.': 'ۦ', '!': 'ۨ', '-': '۪', '+': '۫', '%': '۬', ']': 'ۭ',
	}
	return strings.Map(func(r rune) rune {
		if v, ok := mapping[r]; ok {
			return v
		}
		return r
	}, s)
}
func loadCorpus() {
	corpusOnce.Do(func() { corpus, corpusErr = parseCorpus() })
}
func parseCorpus() (map[string][]Word, error) {
	compressed, e := files.ReadFile("morphology/qac.txt.gz")
	if e != nil {
		return nil, e
	}
	r, e := gzip.NewReader(bytes.NewReader(compressed))
	if e != nil {
		return nil, e
	}
	defer r.Close()
	raw, e := io.ReadAll(r)
	if e != nil {
		return nil, e
	}
	var manifest struct {
		SHA string `json:"sha256"`
	}
	m, _ := files.ReadFile("morphology/source.json")
	if e = json.Unmarshal(m, &manifest); e != nil {
		return nil, e
	}
	if fmt.Sprintf("%x", sha256.Sum256(raw)) != manifest.SHA {
		return nil, fmt.Errorf("morphology checksum mismatch")
	}
	out := map[string][]Word{}
	s := bufio.NewScanner(bytes.NewReader(raw))
	segments := 0
	words := 0
	for s.Scan() {
		line := s.Text()
		if !strings.HasPrefix(line, "(") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 4 {
			return nil, fmt.Errorf("invalid morphology row")
		}
		var ch, v, w, part int
		if n, e := fmt.Sscanf(fields[0], "(%d:%d:%d:%d)", &ch, &v, &w, &part); n != 4 || e != nil || ch < 1 || ch > 114 || v < 1 || w < 1 {
			return nil, fmt.Errorf("invalid morphology location")
		}
		key := strconv.Itoa(ch) + ":" + strconv.Itoa(v)
		list := out[key]
		if w == len(list)+1 {
			if part != 1 {
				return nil, fmt.Errorf("missing segment")
			}
			list = append(list, Word{Position: w})
			words++
		}
		if w != len(list) || part != len(list[w-1].Segments)+1 {
			return nil, fmt.Errorf("noncontiguous morphology")
		}
		seg := Segment{Arabic: buckwalter(fields[1]), Tag: fields[2], Features: fields[3]}
		for _, f := range strings.Split(fields[3], "|") {
			if strings.HasPrefix(f, "ROOT:") {
				seg.Root = buckwalter(strings.TrimPrefix(f, "ROOT:"))
			}
			if strings.HasPrefix(f, "LEM:") {
				seg.Lemma = buckwalter(strings.TrimPrefix(f, "LEM:"))
			}
		}
		list[w-1].Arabic += seg.Arabic
		list[w-1].Segments = append(list[w-1].Segments, seg)
		out[key] = list
		segments++
	}
	if e = s.Err(); e != nil {
		return nil, e
	}
	if len(out) != 6236 || words != 77429 || segments != 128219 {
		return nil, fmt.Errorf("incomplete morphology: %d verses %d words %d segments", len(out), words, segments)
	}
	return out, nil
}

// alignmentKey compares written tokens, ignoring marks and alif-wasla spelling.
// It never changes the canonical text or merges word positions to force a match.
func alignmentKey(s string) string {
	// Uthmani editions write final long /i/ with either dotted or undotted ya.
	// Only the explicitly kasra-marked form is equated; alif maqsura remains distinct.
	s = strings.ReplaceAll(s, "ِى", "ِي")
	// These are alternate written forms of hamza followed by long alif.
	s = strings.ReplaceAll(s, "أٓ", "ءا")
	s = strings.ReplaceAll(s, "آ", "ءا")
	return strings.Map(func(r rune) rune {
		if unicode.IsMark(r) || r == 'ـ' || (r >= 0x06D6 && r <= 0x06ED) {
			return -1
		}
		if r == 'ٱ' {
			return 'ا'
		}
		return r
	}, s)
}
func Analyze(chapter, verse int, canonical string) (*Analysis, error) {
	loadCorpus()
	if corpusErr != nil {
		return nil, corpusErr
	}
	key := fmt.Sprintf("%d:%d", chapter, verse)
	source, ok := corpus[key]
	if !ok {
		return nil, fmt.Errorf("verse not found")
	}
	out := &Analysis{Reference: key, Source: "Quranic Arabic Corpus 0.4", SourceURL: CorpusURL, CanonicalArabic: canonical, Words: make([]Word, len(source))}
	for i, w := range source {
		out.Words[i] = w
		out.Words[i].Segments = append([]Segment(nil), w.Segments...)
	}
	tokens := strings.Fields(canonical)
	filtered := []string{}
	for _, t := range tokens {
		if alignmentKey(t) != "" {
			filtered = append(filtered, t)
		}
	}
	out.Aligned = len(filtered) == len(source)
	if out.Aligned {
		for i, w := range source {
			if alignmentKey(w.Arabic) != alignmentKey(filtered[i]) {
				out.Aligned = false
				break
			}
		}
	}
	if !out.Aligned {
		out.AlignmentNote = "Source word positions are not verified against this edition's written tokens; do not attach them to reader word positions."
	}
	qul, err := compareQUL(chapter, verse, source)
	if err != nil {
		return nil, err
	}
	out.QUL = qul
	return out, nil
}
