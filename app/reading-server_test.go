package app

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"testing/fstest"
)

func TestReadingFileServer(t *testing.T) {
	var zipped bytes.Buffer
	gz := gzip.NewWriter(&zipped)
	gz.Write([]byte("<title>Al-Fatihah</title>In the Name of Allah"))
	gz.Close()
	files := fstest.MapFS{
		"quran/1/index.html.gz": {Data: zipped.Bytes()},
		"__spa-fallback.html":   {Data: []byte("shell")},
		"robots.txt":            {Data: []byte("User-agent: *\nAllow: /\n")},
		"sitemap.xml":           {Data: []byte("<urlset></urlset>")},
		"embed-paths.txt":       {Data: []byte("/embed/quran/1/1\n")},
	}
	handler := readingFileServer(files)
	for _, tc := range []struct {
		path   string
		status int
		body   string
	}{
		{"/quran/1", 200, "In the Name of Allah"}, {"/quran/1/", 200, "Al-Fatihah"},
		{"/quran/999", 404, "404"}, {"/unknown", 404, "404"}, {"/assets/missing.js", 404, "404"},
		{"/search", 200, "shell"}, {"/daily/2026-02-30", 404, "404"},
		{"/embed/quran/1/1", 200, "shell"}, {"/embed/quran/999/999", 404, "404"},
		{"/robots.txt", 200, "User-agent"}, {"/sitemap.xml", 200, "urlset"},
	} {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
		if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.body) {
			t.Fatalf("%s: %d %s", tc.path, w.Code, w.Body.String())
		}
		if tc.path == "/search" && w.Header().Get("X-Robots-Tag") != "noindex" {
			t.Fatal("application shell may be indexed")
		}
		if tc.path == "/quran/1" && !strings.HasPrefix(w.Header().Get("Content-Type"), "text/html") {
			t.Fatal("wrong HTML content type")
		}
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("HEAD", "/quran/1", nil))
	if w.Code != 200 || w.Body.Len() != 0 {
		t.Fatal("HEAD response")
	}
}

func TestGzipNegotiation(t *testing.T) {
	for _, tc := range []struct {
		header   string
		accepted bool
	}{
		{"gzip, br", true}, {"gzip;q=0", false}, {"gzip;q=0.0, br", false}, {"gzip;q=0.5", true}, {"br", false},
	} {
		if acceptsGzip(tc.header) != tc.accepted {
			t.Fatal(tc.header)
		}
	}
}

// Run with REMINDER_RENDERED_DIR after the frontend build to test real output.
func TestGeneratedReadingSite(t *testing.T) {
	dir := os.Getenv("REMINDER_RENDERED_DIR")
	if dir == "" {
		t.Skip("frontend build not requested")
	}
	handler := readingFileServer(os.DirFS(dir))
	for _, tc := range []struct {
		path, contentType, text string
		status                  int
	}{
		{"/quran/1/1", "text/html", "In the Name of Allah", 200},
		{"/names/1", "text/html", "The Beneficent", 200},
		{"/hadith/1", "text/html", "Sahih Bukhari", 200},
		{"/robots.txt", "text/plain", "Sitemap: https://reminder.dev/sitemap.xml", 200},
		{"/sitemap.xml", "", "https://reminder.dev/quran/114/6", 200},
		{"/quran/2/287", "", "404", 404},
		{"/names/100", "", "404", 404},
		{"/not-a-page", "", "404", 404},
	} {
		req := httptest.NewRequest("GET", tc.path, nil)
		req.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		body := w.Body.Bytes()
		if w.Header().Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(bytes.NewReader(body))
			if err != nil {
				t.Fatal(err)
			}
			body, err = io.ReadAll(reader)
			reader.Close()
			if err != nil {
				t.Fatal(err)
			}
		}
		if w.Code != tc.status || !strings.Contains(string(body), tc.text) || !strings.HasPrefix(w.Header().Get("Content-Type"), tc.contentType) {
			t.Fatalf("%s: status=%d type=%s missing=%q", tc.path, w.Code, w.Header().Get("Content-Type"), tc.text)
		}
	}
}
