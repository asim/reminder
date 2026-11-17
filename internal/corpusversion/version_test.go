package corpusversion

import (
	"testing"
	"testing/fstest"
)

func TestVersionTracksOnlySourceContent(t *testing.T) {
	a := fstest.MapFS{"a.json": &fstest.MapFile{Data: []byte("source")}}
	first := Sum(a)
	if Sum(a) != first {
		t.Fatal("unstable fingerprint")
	}
	a["a.json"].Data = []byte("updated")
	if Sum(a) == first {
		t.Fatal("source update did not invalidate")
	}
	b := fstest.MapFS{"b.json": &fstest.MapFile{Data: []byte("source")}}
	if Sum(b) == first {
		t.Fatal("source identity ignored")
	}
}
