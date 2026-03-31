package search

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
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

	opts := hugot.NewDownloadOptions()
	opts.OnnxFilePath = "onnx/model.onnx"
	modelPath, err := hugot.DownloadModel(
		"sentence-transformers/all-MiniLM-L6-v2",
		modelDir,
		opts,
	)
	if err != nil {
		session.Destroy()
		return nil, fmt.Errorf("embed: download model: %w", err)
	}

	pipeline, err := hugot.NewPipeline(session, hugot.FeatureExtractionConfig{
		ModelPath:    modelPath,
		Name:         "search-embeddings",
		OnnxFilename: "onnx/model.onnx",
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

// Save persists embeddings to a binary file with a SHA-256 checksum
// appended at the end to detect tampering.
func (e *Embedder) Save(path string) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Write data to buffer so we can compute checksum
	var buf bytes.Buffer

	count := uint32(len(e.vecs))
	binary.Write(&buf, binary.LittleEndian, count)

	for i := range e.vecs {
		textBytes := []byte(e.texts[i])
		binary.Write(&buf, binary.LittleEndian, uint32(len(textBytes)))
		buf.Write(textBytes)

		metaBytes := []byte(e.metas[i])
		binary.Write(&buf, binary.LittleEndian, uint32(len(metaBytes)))
		buf.Write(metaBytes)

		binary.Write(&buf, binary.LittleEndian, e.vecs[i])
	}

	// Compute SHA-256 over the data
	checksum := sha256.Sum256(buf.Bytes())

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write data + checksum
	if _, err := f.Write(buf.Bytes()); err != nil {
		return err
	}
	if _, err := f.Write(checksum[:]); err != nil {
		return err
	}

	return nil
}

// Load reads embeddings from a binary file and verifies the SHA-256 checksum
// to ensure the data has not been tampered with.
func (e *Embedder) Load(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// File must be at least 32 bytes (SHA-256 checksum)
	if len(raw) < sha256.Size {
		return fmt.Errorf("embed: file too small to contain checksum")
	}

	// Split data and checksum
	data := raw[:len(raw)-sha256.Size]
	storedChecksum := raw[len(raw)-sha256.Size:]

	// Verify checksum
	computed := sha256.Sum256(data)
	if !bytes.Equal(computed[:], storedChecksum) {
		return fmt.Errorf("embed: checksum mismatch — file may be corrupted or tampered with")
	}

	// Parse verified data
	r := bytes.NewReader(data)

	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return err
	}

	vecs := make([][]float32, 0, count)
	texts := make([]string, 0, count)
	metas := make([]string, 0, count)

	for i := uint32(0); i < count; i++ {
		var textLen uint32
		if err := binary.Read(r, binary.LittleEndian, &textLen); err != nil {
			return err
		}
		textBytes := make([]byte, textLen)
		if _, err := io.ReadFull(r, textBytes); err != nil {
			return err
		}

		var metaLen uint32
		if err := binary.Read(r, binary.LittleEndian, &metaLen); err != nil {
			return err
		}
		metaBytes := make([]byte, metaLen)
		if _, err := io.ReadFull(r, metaBytes); err != nil {
			return err
		}

		vec := make([]float32, embeddingDim)
		if err := binary.Read(r, binary.LittleEndian, vec); err != nil {
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
