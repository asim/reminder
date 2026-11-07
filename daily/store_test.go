package daily

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
)

func TestArchiveCompatibilityAndSnapshots(t *testing.T) {
	dir := t.TempDir()
	fixture := `{"2000-01-01":{"verse":"Archived verse","links":{"verse":"/quran/1#1"}}}`
	if err := os.WriteFile(filepath.Join(dir, "daily.json"), []byte(fixture), 0600); err != nil {
		t.Fatal(err)
	}
	s := NewStore(dir)
	entry, ok := s.Get("2000-01-01")
	if !ok || !bytes.Contains(entry, []byte("Archived verse")) {
		t.Fatalf("archive: %s", entry)
	}
	entry[0] = '!'
	entry, _ = s.Get("2000-01-01")
	if !json.Valid(entry) {
		t.Fatal("caller mutated stored archive")
	}
	r := s.Latest()
	r.Links["verse"] = "changed"
	if s.Latest().Links["verse"] != "" {
		t.Fatal("caller mutated latest links")
	}
	var rss bytes.Buffer
	s.WriteRSS(&rss)
	if !bytes.Contains(rss.Bytes(), []byte("Archived verse")) {
		t.Fatal("RSS lost archived content")
	}
	if !json.Valid(s.IndexJSON()) {
		t.Fatal("invalid index JSON")
	}
	fallback := s.Today().(map[string]interface{})
	if fallback["date"] != time.Now().UTC().Format("2006-01-02") || fallback["message"] == "" {
		t.Fatalf("fallback: %v", fallback)
	}
}

func TestUpdatePersistsAndDoesNotRepeatPushAfterRestart(t *testing.T) {
	for _, name := range []string{"OPENAI_API_KEY", "ANTHROPIC_API_KEY", "OLLAMA_LLM_MODEL"} {
		t.Setenv(name, "")
	}
	t.Setenv("OLLAMA_BASE_URL", "http://127.0.0.1:1/v1")
	q := &quran.Quran{Chapters: []*quran.Chapter{{Number: 1, Name: "Al-Fatiha", English: "The Opening", Verses: []*quran.Verse{{Number: 1, Text: "Test verse."}}}}}
	n := names.Names{&names.Name{Number: 1, English: "Test name", Meaning: "Test meaning", Summary: "Test summary"}}
	b := &hadith.Collection{Books: []*hadith.Book{{Number: 1, Name: "Test book", Hadiths: []*hadith.Hadith{{Number: 1, English: "Test hadith", Narrator: "Test narrator"}}}}}
	dir := t.TempDir()
	pushes := 0
	send := func(payload string) []string {
		pushes++
		if !json.Valid([]byte(payload)) {
			t.Error("invalid push payload")
		}
		return nil
	}
	s := NewStore(dir)
	s.Update(context.Background(), q, &n, b, rand.New(rand.NewSource(1)), send)
	if pushes != 1 {
		t.Fatalf("pushes: %d", pushes)
	}
	today := time.Now().UTC().Format("2006-01-02")
	if _, ok := s.Get(today); !ok {
		t.Fatal("missing daily archive")
	}
	if s.Latest().Links["verse"] != "/quran/1#1" || s.Latest().Message == "" {
		t.Fatalf("latest: %+v", s.Latest())
	}
	if len(s.loadHourlyReminders(today)) != 1 {
		t.Fatal("missing hourly archive")
	}
	restarted := NewStore(dir)
	restarted.Update(context.Background(), q, &n, b, rand.New(rand.NewSource(1)), send)
	if pushes != 1 {
		t.Fatal("push repeated after restart")
	}
	if len(restarted.loadHourlyReminders(today)) != 2 {
		t.Fatal("hourly history overwritten")
	}
	if _, ok := restarted.Today().(json.RawMessage); !ok {
		t.Fatal("today should return stored archive")
	}
	// Public snapshots and archive reads are safe while an update owns the store.
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				restarted.Latest()
				restarted.Get(today)
				restarted.IndexJSON()
			}
		}()
	}
	restarted.Update(context.Background(), q, &n, b, rand.New(rand.NewSource(1)), send)
	wg.Wait()
}
