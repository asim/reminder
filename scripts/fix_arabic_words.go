package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type ArabicVerse struct {
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Text    string `json:"text"`
}

func normalize(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// Fixes the word-by-word Arabic using arabic.json as reference, in order of arabic.json
func fixWordByWord() {
	arabicPath := "quran/data/arabic.json"
	wordsDir := "quran/data/words"

	// Load arabic.json
	arabicBytes, err := ioutil.ReadFile(arabicPath)
	if err != nil {
		panic(err)
	}
	var arabic map[string][]ArabicVerse
	if err := json.Unmarshal(arabicBytes, &arabic); err != nil {
		panic(err)
	}

	files, err := filepath.Glob(filepath.Join(wordsDir, "*.json"))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		chapter := strings.TrimSuffix(filepath.Base(file), ".json")
		arabicVerses, ok := arabic[chapter]
		if !ok {
			fmt.Printf("No arabic data for chapter %s\n", chapter)
			continue
		}
		wordsBytes, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file, err)
			continue
		}
		var words map[string]map[string]interface{} // verse -> {"w": [...], "a": {...}}
		if err := json.Unmarshal(wordsBytes, &words); err != nil {
			fmt.Printf("Error parsing %s: %v\n", file, err)
			continue
		}
		updated := false
		// Build new ordered output
		ordered := make(map[string]map[string]interface{})
		for _, v := range arabicVerses {
			verseNum := fmt.Sprintf("%d", v.Verse)
			vdata, exists := words[verseNum]
			if !exists {
				continue
			}
			wArr, ok := vdata["w"].([]interface{})
			if !ok {
				ordered[verseNum] = vdata
				continue
			}
			arabicWords := strings.Fields(v.Text)
			if len(arabicWords) != len(wArr) {
				fmt.Printf("Word count mismatch in chapter %s verse %s: arabic=%d, words=%d\n", chapter, verseNum, len(arabicWords), len(wArr))
				ordered[verseNum] = vdata
				continue
			}
			for i, w := range wArr {
				wMap, ok := w.(map[string]interface{})
				if !ok {
					continue
				}
				wMap["c"] = arabicWords[i]
			}
			ordered[verseNum] = vdata
			updated = true
		}
		if updated {
			// Marshal with order: build a slice of key-value pairs in order
			type kv struct {
				Key   string
				Value map[string]interface{}
			}
			var orderedSlice []kv
			for _, v := range arabicVerses {
				verseNum := fmt.Sprintf("%d", v.Verse)
				if vdata, ok := ordered[verseNum]; ok {
					orderedSlice = append(orderedSlice, kv{Key: verseNum, Value: vdata})
				}
			}
			// Build ordered map for marshaling
			outMap := make(map[string]json.RawMessage)
			for _, pair := range orderedSlice {
				b, _ := json.Marshal(pair.Value)
				outMap[pair.Key] = b
			}
			// Custom marshal to preserve order
			var outBuf strings.Builder
			outBuf.WriteString("{\n")
			for i, pair := range orderedSlice {
				outBuf.WriteString(fmt.Sprintf("  \"%s\": %s", pair.Key, outMap[pair.Key]))
				if i < len(orderedSlice)-1 {
					outBuf.WriteString(",\n")
				}
			}
			outBuf.WriteString("\n}")
			if err := ioutil.WriteFile(file, []byte(outBuf.String()), 0644); err != nil {
				fmt.Printf("Error writing %s: %v\n", file, err)
			}
			fmt.Printf("Updated %s\n", file)
		}
	}
}

// Verifies the word-by-word files against arabic.json and prints discrepancies
func verifyWordByWord() {
	arabicPath := "quran/data/arabic.json"
	wordsDir := "quran/data/words"

	arabicBytes, err := ioutil.ReadFile(arabicPath)
	if err != nil {
		panic(err)
	}
	var arabic map[string][]ArabicVerse
	if err := json.Unmarshal(arabicBytes, &arabic); err != nil {
		panic(err)
	}

	files, err := filepath.Glob(filepath.Join(wordsDir, "*.json"))
	if err != nil {
		panic(err)
	}
	problems := 0
	for _, file := range files {
		chapter := strings.TrimSuffix(filepath.Base(file), ".json")
		arabicVerses, ok := arabic[chapter]
		if !ok {
			continue
		}
		wordsBytes, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		var words map[string]map[string]interface{}
		if err := json.Unmarshal(wordsBytes, &words); err != nil {
			continue
		}
		for _, v := range arabicVerses {
			verseNum := fmt.Sprintf("%d", v.Verse)
			vdata, exists := words[verseNum]
			if !exists {
				continue
			}
			wArr, ok := vdata["w"].([]interface{})
			if !ok {
				continue
			}
			var joinedWords []string
			for _, w := range wArr {
				wMap, ok := w.(map[string]interface{})
				if !ok {
					continue
				}
				if c, ok := wMap["c"].(string); ok {
					joinedWords = append(joinedWords, c)
				}
			}
			joined := normalize(strings.Join(joinedWords, " "))
			ref := normalize(v.Text)
			if joined != ref {
				problems++
				fmt.Printf("Mismatch: chapter %s verse %d\n", chapter, v.Verse)
				fmt.Printf("  words: %s\n", joined)
				fmt.Printf("  ref:   %s\n", ref)
			}
		}
	}
	if problems == 0 {
		fmt.Println("All verses match perfectly!")
	} else {
		fmt.Printf("Total mismatches: %d\n", problems)
	}
}

func main() {
	fixWordByWord()
	verifyWordByWord()
}
