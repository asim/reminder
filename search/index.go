package search

import (
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
	"github.com/philippgille/chromem-go"
)

//go:embed data/*.gob.gz
var files embed.FS

type Index struct {
	Home string
	Name string
	DB   *chromem.DB
	Col  *chromem.Collection
}

type Result struct {
	Text     string            `json:"text"`
	Score    float32           `json:"score"`
	Metadata map[string]string `json:"metadata"`
}

// getEmbeddingFunc returns an embedding function based on environment configuration.
// Priority: 1. OpenAI (fast, requires API key), 2. Ollama (local, slower)
// Set OPENAI_API_KEY to use OpenAI embeddings (text-embedding-3-small, fast & cheap)
// Set OLLAMA_EMBED_MODEL to use a different Ollama model (default: nomic-embed-text)
// Set OLLAMA_BASE_URL to use a different Ollama instance (default: http://localhost:11434/api)
func getEmbeddingFunc() chromem.EmbeddingFunc {
	// Check for OpenAI API key first - much faster for embeddings
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey != "" {
		// Use OpenAI's text-embedding-3-small - fast, cheap ($0.02/1M tokens), good quality
		return chromem.NewEmbeddingFuncOpenAI(openaiKey, chromem.EmbeddingModelOpenAI3Small)
	}

	// Fall back to local Ollama
	model := os.Getenv("OLLAMA_EMBED_MODEL")
	if model == "" {
		model = "nomic-embed-text"
	}

	baseURL := os.Getenv("OLLAMA_BASE_URL")
	// baseURL can be empty - chromem-go will use http://localhost:11434/api by default

	return chromem.NewEmbeddingFuncOllama(model, baseURL)
}

// Load the embedded index
func (i *Index) Load() error {
	f, err := files.Open("data/" + i.Name + ".idx.gob.gz")
	if err != nil {
		return err
	}
	defer f.Close()

	path := filepath.Join(i.Home, ".reminder", "data")
	fpath := filepath.Join(path, i.Name+".idx.gob.gz")

	// write the file
	os.MkdirAll(path, 0755)

	// check exists otherwise write it
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		f2, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f2.Close()

		// copy contents
		_, err = io.Copy(f2, f)
		if err != nil {
			return err
		}
	}

	// read from file
	if err := i.DB.ImportFromFile(fpath, ""); err != nil {
		// If import fails, the index might be incompatible (different embedding dimensions)
		// Return error so caller can handle (e.g., rebuild index)
		return fmt.Errorf("failed to import index (may need rebuild): %w", err)
	}

	/*
		sk, err := NewReadSeekerWrapper(f)
		if err != nil {
			return err
		}
		if err := i.DB.ImportFromReader(sk, ""); err != nil {
			return err
		}
	*/

	// Use Ollama with nomic-embed-text for embeddings
	embeddingFunc := getEmbeddingFunc()
	c, err := i.DB.GetOrCreateCollection(i.Name, nil, embeddingFunc)
	if err != nil {
		return err
	}

	// set the Collection
	i.Col = c
	return nil
}

func (i *Index) Export() error {
	path := filepath.Join(i.Home, i.Name+".idx.gob.gz")

	return i.DB.ExportToFile(path, true, "")
}

func (i *Index) Import() error {
	path := filepath.Join(i.Home, i.Name+".idx.gob.gz")

	if err := i.DB.ImportFromFile(path, ""); err != nil {
		return err
	}

	// Use Ollama with nomic-embed-text for embeddings
	embeddingFunc := getEmbeddingFunc()
	c, err := i.DB.GetOrCreateCollection(i.Name, nil, embeddingFunc)
	if err != nil {
		return err
	}

	i.Col = c
	return nil
}

func (i *Index) Store(md map[string]string, content ...string) error {
	var docs []chromem.Document

	for _, c := range content {
		if len(c) == 0 {
			fmt.Println("skipping")
			continue
		}
		fmt.Println("Indexing content: ", c)
		docs = append(docs, chromem.Document{
			ID:       uuid.New().String(),
			Content:  c,
			Metadata: md,
		})
	}

	return i.Col.AddDocuments(context.TODO(), docs, runtime.NumCPU())
}

func (i *Index) Query(v string) ([]*Result, error) {
	res, err := i.Col.Query(context.TODO(), v, 25, nil, nil)
	if err != nil {
		return nil, err
	}

	var results []*Result

	for _, result := range res {
		results = append(results, &Result{
			Text:     result.Content,
			Score:    result.Similarity,
			Metadata: result.Metadata,
		})
	}

	return results, nil
}

func New(name string, persist bool) *Index {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	var db *chromem.DB
	var c *chromem.Collection

	// Use Ollama with nomic-embed-text for embeddings
	embeddingFunc := getEmbeddingFunc()

	if persist {
		path := filepath.Join(u.HomeDir, name+".idx")

		var err error
		db, err = chromem.NewPersistentDB(path, false)
		if err != nil {
			panic(err)
		}

		c, err = db.GetOrCreateCollection(name, nil, embeddingFunc)
		if err != nil {
			panic(err)
		}
	} else {
		db = chromem.NewDB()
		c, _ = db.CreateCollection(name, nil, embeddingFunc)
	}

	return &Index{
		Home: u.HomeDir,
		Name: name,
		DB:   db,
		Col:  c,
	}
}
