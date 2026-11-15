package app

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dailyPath = regexp.MustCompile(`^/daily/\d{4}-\d{2}-\d{2}$`)
var embedPath = regexp.MustCompile(`^/embed/(quran|hadith)/[1-9][0-9]*/[1-9][0-9]*$`)

// Reading pages must exist in the generated site. Only application-only routes
// may use the SPA shell; missing assets and invalid reader URLs are real 404s.
func readingFileServer(files fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "HEAD" {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "Method not allowed", 405)
			return
		}
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "." || name == "" {
			name = "index.html"
		}
		candidates := []string{name, name + "/index.html"}
		for _, candidate := range candidates {
			if serveReadingFile(w, r, files, candidate) {
				return
			}
		}
		fallback := false
		switch r.URL.Path {
		case "/home", "/daily", "/daily/index", "/bookmarks", "/search", "/islam":
			fallback = true
		}
		if dailyPath.MatchString(r.URL.Path) {
			_, err := time.Parse("2006-01-02", strings.TrimPrefix(r.URL.Path, "/daily/"))
			fallback = err == nil
		}
		if embedPath.MatchString(r.URL.Path) {
			// Exported embed paths are validated against the same bundled corpus.
			if data, err := fs.ReadFile(files, "embed-paths.txt"); err == nil {
				fallback = strings.Contains("\n"+string(data), "\n"+r.URL.Path+"\n")
			}
		}
		if fallback {
			w.Header().Set("X-Robots-Tag", "noindex")
			if serveReadingFile(w, r, files, "__spa-fallback.html") {
				return
			}
		}
		http.NotFound(w, r)
	})
}
func serveReadingFile(w http.ResponseWriter, r *http.Request, files fs.FS, name string) bool {
	data, err := fs.ReadFile(files, name)
	if err != nil {
		compressed, e := fs.ReadFile(files, name+".gz")
		if e != nil {
			return false
		}
		w.Header().Set("Vary", "Accept-Encoding")
		if acceptsGzip(r.Header.Get("Accept-Encoding")) {
			data = compressed
			w.Header().Set("Content-Encoding", "gzip")
		} else {
			reader, e := gzip.NewReader(bytes.NewReader(compressed))
			if e != nil {
				return false
			}
			data, e = io.ReadAll(reader)
			reader.Close()
			if e != nil {
				return false
			}
		}
	}
	typ := mime.TypeByExtension(path.Ext(name))
	if path.Ext(name) == ".data" {
		typ = "text/x-script; charset=utf-8"
	}
	if typ != "" {
		w.Header().Set("Content-Type", typ)
	}
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(data))
	return true
}

func acceptsGzip(value string) bool {
	for _, item := range strings.Split(value, ",") {
		parts := strings.Split(strings.TrimSpace(item), ";")
		if parts[0] != "gzip" {
			continue
		}
		for _, parameter := range parts[1:] {
			parameter = strings.TrimSpace(parameter)
			if strings.HasPrefix(parameter, "q=") {
				quality, err := strconv.ParseFloat(strings.TrimPrefix(parameter, "q="), 64)
				return err == nil && quality > 0 && quality <= 1
			}
		}
		return true
	}
	return false
}
