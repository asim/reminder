package search

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// fakeEmbedder builds an Embedder with a given provider identity but no
// backend, so persistence can be tested without a model or an API key.
func fakeEmbedder(provider string, dim int) *Embedder {
	return &Embedder{provider: provider, dim: dim}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "embeddings.bin")

	src := fakeEmbedder(ProviderLocal, 4)
	src.Add("first verse", `{"source":"quran"}`, []float32{1, 0, 0, 0})
	src.Add("second verse", `{"source":"hadith"}`, []float32{0, 1, 0, 0})

	if err := src.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	dst := fakeEmbedder(ProviderLocal, 4)
	if err := dst.Load(path); err != nil {
		t.Fatalf("load: %v", err)
	}
	if dst.Count() != 2 {
		t.Fatalf("expected 2 embeddings, got %d", dst.Count())
	}

	res := dst.Search([]float32{1, 0, 0, 0}, 1)
	if len(res) != 1 || res[0].Text != "first verse" {
		t.Errorf("expected nearest neighbour to be the first verse, got %+v", res)
	}
}

// TestLoadRejectsDifferentProvider is the regression guard for the failure
// that made search unusable without an OpenAI key: vectors from one model
// being compared against query vectors from another. Different dimensions
// error on every query; matching dimensions would be silently meaningless.
func TestLoadRejectsDifferentProvider(t *testing.T) {
	path := filepath.Join(t.TempDir(), "embeddings.bin")

	openaiSide := fakeEmbedder(ProviderOpenAI, 4)
	openaiSide.Add("verse", `{}`, []float32{1, 0, 0, 0})
	if err := openaiSide.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	localSide := fakeEmbedder(ProviderLocal, 4)
	err := localSide.Load(path)
	if !errors.Is(err, ErrProviderMismatch) {
		t.Fatalf("expected ErrProviderMismatch, got %v", err)
	}
	if localSide.Count() != 0 {
		t.Errorf("no embeddings should be loaded on mismatch, got %d", localSide.Count())
	}
}

func TestLoadRejectsDifferentDimension(t *testing.T) {
	path := filepath.Join(t.TempDir(), "embeddings.bin")

	wide := fakeEmbedder(ProviderLocal, 8)
	wide.Add("verse", `{}`, make([]float32, 8))
	if err := wide.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	narrow := fakeEmbedder(ProviderLocal, 4)
	if err := narrow.Load(path); !errors.Is(err, ErrProviderMismatch) {
		t.Fatalf("expected ErrProviderMismatch, got %v", err)
	}
}

// TestLoadRejectsLegacyFile covers embeddings written before the header
// existed. They carry no provider identity, so they cannot be trusted.
func TestLoadRejectsLegacyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "embeddings.bin")

	// A headerless file: arbitrary bytes with a valid trailing checksum would
	// still lack the magic prefix.
	if err := os.WriteFile(path, make([]byte, 128), 0644); err != nil {
		t.Fatal(err)
	}

	e := fakeEmbedder(ProviderLocal, 4)
	if err := e.Load(path); err == nil {
		t.Fatal("expected an error loading a headerless file")
	}
}

func TestLoadMissingFile(t *testing.T) {
	e := fakeEmbedder(ProviderLocal, 4)
	err := e.Load(filepath.Join(t.TempDir(), "absent.bin"))
	if !os.IsNotExist(err) {
		t.Fatalf("expected a not-exist error, got %v", err)
	}
}

// TestProviderSelection checks that the OpenAI backend is chosen when a key is
// present. The local path is not exercised here because it downloads a model.
func TestProviderSelection(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "test-key-not-used")

	e, err := NewEmbedder(t.TempDir())
	if err != nil {
		t.Fatalf("new embedder: %v", err)
	}
	if e.Provider() != ProviderOpenAI {
		t.Errorf("expected %s, got %s", ProviderOpenAI, e.Provider())
	}
	if e.Dim() != dimOpenAI {
		t.Errorf("expected dim %d, got %d", dimOpenAI, e.Dim())
	}
}
