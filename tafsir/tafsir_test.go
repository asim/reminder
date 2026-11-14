package tafsir

import "testing"

func TestSharedPassagesAndInvalidPointers(t *testing.T) {
	edition := Edition{ID: "test", Title: "Test tafsir"}
	c, err := Parse([]byte(`{"1:1":{"text":"One passage","ayah_keys":["1:1","1:2"]},"1:2":"1:1"}`), edition)
	if err != nil {
		t.Fatal(err)
	}
	if c.Get(1, 1) != c.Get(1, 2) || len(c.Passages()) != 1 || len(c.Get(1, 1).Ayahs) != 2 {
		t.Fatal("lost shared passage")
	}
	for _, data := range []string{`{"1:1":"1:2"}`, `{"1:1":"1:2","1:2":"1:1"}`, `{"1:1":{"text":""}}`, `{"1:1":{"text":"x","ayah_keys":["1:2"]},"1:2":{"text":"y"}}`} {
		if _, err := Parse([]byte(data), edition); err == nil {
			t.Fatalf("accepted %s", data)
		}
	}
}
func TestBundledCorpus(t *testing.T) {
	c := Load()
	if c != Load() {
		t.Fatal("corpus loaded repeatedly")
	}
	total := 0
	for _, p := range c.Passages() {
		total += len(p.Ayahs)
	}
	if total != 6236 {
		t.Fatalf("coverage %d", total)
	}
}
