package search

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"sync"

	"github.com/knights-analytics/hugot"
	"github.com/knights-analytics/hugot/pipelines"
	"github.com/sashabaranov/go-openai"
)

// Embeddings file header. The magic prefix distinguishes versioned files from
// the original headerless format.
const (
	embedMagic         = "RMDREMB1"
	embedFormatVersion = 1
)

// Embedding providers. The provider name is written into the embeddings file
// so that switching provider is detected on load rather than producing
// silently wrong results — vectors from different models are not comparable.
const (
	ProviderLocal  = "local:all-MiniLM-L6-v2"
	ProviderOpenAI = "openai:text-embedding-3-small"

	dimLocal  = 384
	dimOpenAI = 1536

	openaiBatchSize = 256 // OpenAI accepts many inputs per request
)

// ErrProviderMismatch is returned by Load when the stored embeddings were
// produced by a different provider than the one now configured. The caller
// should discard them and rebuild.
var ErrProviderMismatch = errors.New("embed: embeddings were built with a different provider")

// Embedder computes and stores document embeddings for semantic search.
//
// Two backends are supported. When OPENAI_API_KEY is set the OpenAI embedding
// API is used; otherwise a local ONNX model runs in-process, which keeps the
// app fully self-contained. Either way the vectors live in memory and are
// persisted alongside the FTS index.
type Embedder struct {
	mu       sync.RWMutex
	provider string
	dim      int
	pipeline *pipelines.FeatureExtractionPipeline // local backend only
	session  *hugot.Session                       // local backend only
	client   *openai.Client                       // OpenAI backend only
	vecs     [][]float32                          // one embedding per document
	texts    []string                             // parallel array of document texts
	metas    []string                             // parallel array of metadata JSON
}

// Provider reports which embedding backend is in use.
func (e *Embedder) Provider() string { return e.provider }

// Dim reports the dimensionality of the vectors this embedder produces.
func (e *Embedder) Dim() int { return e.dim }

// NewEmbedder creates an embedder using the best available backend. If
// OPENAI_API_KEY is set it uses OpenAI, which is faster to build and generally
// higher quality. Otherwise it falls back to a local model, downloading it on
// first use. Callers that get an error can still run keyword-only search.
func NewEmbedder(modelDir string) (*Embedder, error) {
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return newOpenAIEmbedder(key), nil
	}
	return newLocalEmbedder(modelDir)
}

func newOpenAIEmbedder(apiKey string) *Embedder {
	return &Embedder{
		provider: ProviderOpenAI,
		dim:      dimOpenAI,
		client:   openai.NewClient(apiKey),
	}
}

func newLocalEmbedder(modelDir string) (*Embedder, error) {
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
		provider: ProviderLocal,
		dim:      dimLocal,
		pipeline: pipeline,
		session:  session,
	}, nil
}

// Embed computes the embedding for a single text.
func (e *Embedder) Embed(text string) ([]float32, error) {
	vecs, err := e.EmbedBatch([]string{text})
	if err != nil {
		return nil, err
	}
	if len(vecs) == 0 {
		return nil, fmt.Errorf("embed: no embeddings returned")
	}
	return vecs[0], nil
}

// EmbedBatch computes embeddings for multiple texts.
func (e *Embedder) EmbedBatch(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}
	if e.client != nil {
		return e.embedBatchOpenAI(texts)
	}
	result, err := e.pipeline.RunPipeline(texts)
	if err != nil {
		return nil, fmt.Errorf("embed: run pipeline: %w", err)
	}
	return result.Embeddings, nil
}

func (e *Embedder) embedBatchOpenAI(texts []string) ([][]float32, error) {
	out := make([][]float32, 0, len(texts))

	for start := 0; start < len(texts); start += openaiBatchSize {
		end := start + openaiBatchSize
		if end > len(texts) {
			end = len(texts)
		}

		res, err := e.client.CreateEmbeddings(context.TODO(), openai.EmbeddingRequestStrings{
			Input: texts[start:end],
			Model: openai.SmallEmbedding3,
		})
		if err != nil {
			return nil, fmt.Errorf("embed: openai: %w", err)
		}
		if len(res.Data) != end-start {
			return nil, fmt.Errorf("embed: openai returned %d embeddings for %d inputs",
				len(res.Data), end-start)
		}

		// The API may return results out of order; Index gives the position.
		batch := make([][]float32, end-start)
		for _, d := range res.Data {
			if d.Index < 0 || d.Index >= len(batch) {
				return nil, fmt.Errorf("embed: openai returned out-of-range index %d", d.Index)
			}
			batch[d.Index] = d.Embedding
		}
		out = append(out, batch...)
	}

	return out, nil
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

	// Header: magic, format version, provider, dimension. Without this a file
	// written by one provider would be read back as another and every query
	// would compare vectors of different lengths.
	buf.WriteString(embedMagic)
	binary.Write(&buf, binary.LittleEndian, uint32(embedFormatVersion))

	providerBytes := []byte(e.provider)
	binary.Write(&buf, binary.LittleEndian, uint32(len(providerBytes)))
	buf.Write(providerBytes)

	binary.Write(&buf, binary.LittleEndian, uint32(e.dim))

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

	// Header. A file without the magic prefix predates versioning, so treat
	// it as stale rather than trying to interpret it.
	magic := make([]byte, len(embedMagic))
	if _, err := io.ReadFull(r, magic); err != nil || string(magic) != embedMagic {
		return ErrProviderMismatch
	}

	var version uint32
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return err
	}
	if version != embedFormatVersion {
		return ErrProviderMismatch
	}

	var providerLen uint32
	if err := binary.Read(r, binary.LittleEndian, &providerLen); err != nil {
		return err
	}
	providerBytes := make([]byte, providerLen)
	if _, err := io.ReadFull(r, providerBytes); err != nil {
		return err
	}

	var dim uint32
	if err := binary.Read(r, binary.LittleEndian, &dim); err != nil {
		return err
	}

	// Refuse embeddings from a different model. Comparing them against this
	// model's query vectors would be meaningless at best and, on a dimension
	// mismatch, an error on every single query.
	if string(providerBytes) != e.provider || int(dim) != e.dim {
		return fmt.Errorf("%w: file has %s (%d-dim), configured for %s (%d-dim)",
			ErrProviderMismatch, providerBytes, dim, e.provider, e.dim)
	}

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

		vec := make([]float32, e.dim)
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
