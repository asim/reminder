// Package test checks the real routes registered by main.go.
// Run scripts/test-http.sh to build and start an isolated server for these tests.
package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func request(t *testing.T, method, path string, body []byte) (int, http.Header, []byte) {
	t.Helper()
	base := os.Getenv("REMINDER_TEST_URL")
	if base == "" {
		t.Skip("run scripts/test-http.sh for HTTP integration tests")
	}
	req, err := http.NewRequest(method, base+path, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, resp.Header, data
}

func TestQuranVerseRoute(t *testing.T) {
	code, headers, data := request(t, "GET", "/api/quran/1/1", nil)
	if code != http.StatusOK || headers.Get("Content-Type") != "application/json" {
		t.Fatalf("status %d, headers %v", code, headers)
	}
	var verse struct {
		Chapter, Number int
		Text            string
	}
	if err := json.Unmarshal(data, &verse); err != nil {
		t.Fatal(err)
	}
	if verse.Chapter != 1 || verse.Number != 1 || verse.Text == "" {
		t.Fatalf("wrong verse: %+v", verse)
	}
	for _, path := range []string{
		"/api/quran/not-a-chapter/1", "/api/quran/0/1", "/api/quran/-1/1", "/api/quran/115/1",
		"/api/quran/1/not-a-verse", "/api/quran/1/0", "/api/quran/1/-1", "/api/quran/1/8",
	} {
		t.Run(path, func(t *testing.T) {
			code, headers, data := request(t, "GET", path, nil)
			if code != 404 || headers.Get("Content-Type") != "application/json" || string(data) != "{}" {
				t.Fatalf("status %d, headers %v, body %s", code, headers, data)
			}
		})
	}
}

func TestExplanationUsesRequestedVerseWithoutModel(t *testing.T) {
	for _, path := range []string{"/api/quran/0/1/explain", "/api/quran/2/0/explain", "/api/quran/2/287/explain"} {
		code, _, _ := request(t, "POST", path, nil)
		if code != 404 {
			t.Fatalf("%s: %d", path, code)
		}
	}
	code, _, data := request(t, "POST", "/api/quran/2/255/explain", nil)
	if code != 200 {
		t.Fatalf("status %d", code)
	}
	var result struct {
		Text, Reference string
		Generated       bool
	}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	_, _, data = request(t, "GET", "/api/quran/2/255", nil)
	var verse struct{ Comments string }
	if err := json.Unmarshal(data, &verse); err != nil {
		t.Fatal(err)
	}
	if result.Generated || result.Reference != "2:255" || result.Text != verse.Comments || result.Text == "" {
		t.Fatalf("wrong explanation: %+v", result)
	}
}

func TestDailyArchiveAndMCP(t *testing.T) {
	code, _, data := request(t, "GET", "/api/daily/2000-01-01", nil)
	if code != 200 || !bytes.Contains(data, []byte("Fixture verse")) {
		t.Fatalf("archive: %d %s", code, data)
	}
	code, _, data = request(t, "GET", "/api/daily/1900-01-01", nil)
	if code != 404 {
		t.Fatalf("missing archive: %d", code)
	}
	code, headers, data := request(t, "GET", "/rss", nil)
	if code != 200 || headers.Get("Content-Type") != "application/rss+xml; charset=utf-8" || !bytes.Contains(data, []byte("Fixture verse")) {
		t.Fatalf("rss: %d %s", code, data)
	}
	code, _, data = request(t, "POST", "/mcp", []byte(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"get_daily_by_date","arguments":{"date":"2000-01-01"}}}`))
	if code != 200 || !bytes.Contains(data, []byte("Fixture verse")) {
		t.Fatalf("MCP archive: %d %s", code, data)
	}
}

func TestLitePagesAndSearch(t *testing.T) {
	for _, path := range []string{"/home", "/bookmarks", "/daily", "/quran/1", "/quran/1/2", "/names/1", "/hadith/1"} {
		code, _, data := request(t, "GET", path, nil)
		if code != 200 || len(data) == 0 {
			t.Fatalf("%s: %d", path, code)
		}
	}
	// Keyword indexing runs asynchronously; wait for the real search endpoint.
	deadline := time.Now().Add(30 * time.Second)
	for {
		code, _, data := request(t, "POST", "/api/search", []byte(`{"q":"mercy","summarise":false}`))
		if bytes.Contains(data, []byte("Indexing content")) && time.Now().Before(deadline) {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		var result struct {
			References         []json.RawMessage
			SummariseAvailable bool
		}
		if code != 200 {
			t.Fatalf("search: %d %s", code, data)
		}
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatal(err)
		}
		if len(result.References) == 0 || result.SummariseAvailable {
			t.Fatalf("search: %s", data)
		}
		break
	}
}
