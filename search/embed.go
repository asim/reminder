package search

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"sync"

	"github.com/knights-analytics/hugot"
	"github.com/knights-analytics/hugot/pipelines"
)

const embeddingDim = 384 // all-MiniLM-L6-v2

// Embedder computes and stores document embeddings for semantic search.
type Embedder struct {
	mu       sync.RWMutex
	pipeline *pipelines.FeatureExtractionPipeline
	session  *hugot.Session
	vecs     [][]float32 // one embedding per document, indexed by doc ID-1
	texts    []string    // parallel array of document texts
	metas    []string    // parallel array of metadata JSON
}

// NewEmbedder creates an embedder, downloading the model if needed.
func NewEmbedder(modelDir string) (*Embedder, error) {
	session, err := hugot.NewGoSession()
	if err != nil {
		return nil, fmt.Errorf("embed: create session: %w", err)
	}

	modelPath, err := hugot.DownloadModel(
		"sentence-transformers/all-MiniLM-L6-v2",
		modelDir,
		hugot.NewDownloadOptions(),
	)
	if err != nil {
		session.Destroy()
		return nil, fmt.Errorf("embed: download model: %w", err)
	}

	pipeline, err := hugot.NewPipeline(session, hugot.FeatureExtractionConfig{
		ModelPath: modelPath,
		Name:      "search-embeddings",
	})
	if err != nil {
		session.Destroy()
		return nil, fmt.Errorf("embed: create pipeline: %w", err)
	}

	return &Embedder{
		pipeline: pipeline,
		session:  session,
	}, nil
}

// Embed computes the embedding for a single text.
func (e *Embedder) Embed(text string) ([]float32, error) {
	result, err := e.pipeline.RunPipeline([]string{text})
	if err != nil {
		return nil, fmt.Errorf("embed: run pipeline: %w", err)
	}
	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("embed: no embeddings returned")
	}
	return result.Embeddings[0], nil
}

// EmbedBatch computes embeddings for multiple texts.
func (e *Embedder) EmbedBatch(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}
	result, err := e.pipeline.RunPipeline(texts)
	if err != nil {
		return nil, fmt.Errorf("embed: run pipeline: %w", err)
	}
	return result.Embeddings, nil
}

// Add stores a document embedding in the in-memory index.
func (e *Embedder) Add(text string, metadataJSON string, vec []float32) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.vecs = append(e.vecs, vec)
	e.texts = append(e.texts, text)
	e.metas = append(e.metas, metadataJSON)
}

// Count returns the number of stored embeddings.
func (e *Embedder) Count() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.vecs)
}

type scoredDoc struct {
	index int
	score float32
}

// Search finds the top-N most similar documents to the query embedding.
func (e *Embedder) Search(queryVec []float32, topN int) []*Result {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if len(e.vecs) == 0 {
		return nil
	}

	scored := make([]scoredDoc, len(e.vecs))
	for i, docVec := range e.vecs {
		scored[i] = scoredDoc{index: i, score: cosineSimilarity(queryVec, docVec)}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	if topN > len(scored) {
		topN = len(scored)
	}

	results := make([]*Result, 0, topN)
	for _, s := range scored[:topN] {
		md := make(map[string]string)
		json.Unmarshal([]byte(e.metas[s.index]), &md)
		results = append(results, &Result{
			Text:     e.texts[s.index],
			Score:    s.score,
			Metadata: md,
		})
	}
	return results
}

// Save persists embeddings to a binary file.
func (e *Embedder) Save(path string) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write count
	count := uint32(len(e.vecs))
	if err := binary.Write(f, binary.LittleEndian, count); err != nil {
		return err
	}

	for i := range e.vecs {
		// Write text length + text
		textBytes := []byte(e.texts[i])
		if err := binary.Write(f, binary.LittleEndian, uint32(len(textBytes))); err != nil {
			return err
		}
		f.Write(textBytes)

		// Write metadata length + metadata
		metaBytes := []byte(e.metas[i])
		if err := binary.Write(f, binary.LittleEndian, uint32(len(metaBytes))); err != nil {
			return err
		}
		f.Write(metaBytes)

		// Write embedding
		if err := binary.Write(f, binary.LittleEndian, e.vecs[i]); err != nil {
			return err
		}
	}

	return nil
}

// Load reads embeddings from a binary file.
func (e *Embedder) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var count uint32
	if err := binary.Read(f, binary.LittleEndian, &count); err != nil {
		return err
	}

	vecs := make([][]float32, 0, count)
	texts := make([]string, 0, count)
	metas := make([]string, 0, count)

	for i := uint32(0); i < count; i++ {
		// Read text
		var textLen uint32
		if err := binary.Read(f, binary.LittleEndian, &textLen); err != nil {
			return err
		}
		textBytes := make([]byte, textLen)
		if _, err := f.Read(textBytes); err != nil {
			return err
		}

		// Read metadata
		var metaLen uint32
		if err := binary.Read(f, binary.LittleEndian, &metaLen); err != nil {
			return err
		}
		metaBytes := make([]byte, metaLen)
		if _, err := f.Read(metaBytes); err != nil {
			return err
		}

		// Read embedding
		vec := make([]float32, embeddingDim)
		if err := binary.Read(f, binary.LittleEndian, vec); err != nil {
			return err
		}

		texts = append(texts, string(textBytes))
		metas = append(metas, string(metaBytes))
		vecs = append(vecs, vec)
	}

	e.mu.Lock()
	e.vecs = vecs
	e.texts = texts
	e.metas = metas
	e.mu.Unlock()

	return nil
}

// Destroy cleans up the ONNX session.
func (e *Embedder) Destroy() {
	if e.session != nil {
		e.session.Destroy()
	}
}

func cosineSimilarity(a, b []float32) float32 {
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
