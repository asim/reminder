package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
)

// removeLegacyFiles deletes artefacts from the vector search era. The embedded
// chromem index is a large file that nothing reads any more, and a leftover
// indexing checkpoint would make the indexers skip sources that the current
// build has not actually written.
func removeLegacyFiles() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	for _, p := range []string{
		filepath.Join(home, "reminder-index-checkpoint.json"),
		filepath.Join(home, ".reminder", "data", "reminder.idx.gob.gz"),
	} {
		if err := os.Remove(p); err == nil {
			log.Printf("Removed unused file %s", p)
		}
	}
}

func indexContent(idx *search.Index, md map[string]string, text string) {
	// index the documents
	lines := strings.Split(text, "\n")

	log.Println("Indexing: ", md["source"], md["chapter"], md["verse"])

	if err := idx.Store(md, lines...); err != nil {
		log.Println("Error indexing", err)
	}
}

func indexQuran(idx *search.Index, q *quran.Quran) {
	log.Println("Indexing Quran")

	for _, chapter := range q.Chapters {
		for _, verse := range chapter.Verses {
			indexContent(idx, map[string]string{
				"source":  "quran",
				"chapter": fmt.Sprintf("%v", chapter.Number),
				"verse":   fmt.Sprintf("%v", verse.Number),
				"name":    chapter.Name,
			}, verse.Text)
		}
	}
}

func indexNames(idx *search.Index, n *names.Names) {
	log.Println("Indexing Names")

	for _, name := range *n {
		indexContent(idx, map[string]string{
			"source":  "names",
			"meaning": name.Meaning,
			"english": name.English,
			"arabic":  name.Arabic,
		}, strings.Join([]string{name.Meaning, name.English, name.Description}, " - "))
	}
}

func indexTafsir(idx *search.Index, q *quran.Quran) {
	log.Println("Indexing Tafsir")

	for _, comment := range q.Commentary {
		indexContent(idx, map[string]string{
			"source":  "tafsir",
			"chapter": fmt.Sprintf("%v", comment.Chapter),
			"verse":   fmt.Sprintf("%v", comment.Verse),
		}, comment.Text)
	}
}

func indexHadith(idx *search.Index, b *hadith.Collection) {
	log.Println("Indexing Hadith")

	for _, book := range b.Books {
		for _, h := range book.Hadiths {
			indexContent(idx, map[string]string{
				"source":   "bukhari",
				"book":     book.Name,
				"book_num": fmt.Sprintf("%d", book.Number),
				"narrator": h.Narrator,
				"number":   fmt.Sprintf("%d", h.Number),
			}, h.English)
		}
	}
}

// buildEmbeddings reads all documents from the FTS index and computes embeddings.
// buildEmbeddings embeds every document in the index. It returns an error if
// any batch fails, so that callers do not persist a partially embedded corpus:
// a saved partial set looks complete on the next start and those documents
// would never be embedded again.
func buildEmbeddings(idx *search.Index, embedder *search.Embedder) error {
	db := idx.DB()
	if db == nil {
		return fmt.Errorf("embeddings: no database")
	}

	rows, err := db.Query(`SELECT text, metadata FROM docs`)
	if err != nil {
		return fmt.Errorf("embeddings: read docs: %w", err)
	}
	defer rows.Close()

	// Batch documents for efficient embedding
	const batchSize = 32
	var batchTexts []string
	var batchMetas []string

	processBatch := func(texts, metas []string) error {
		vecs, err := embedder.EmbedBatch(texts)
		if err != nil {
			return fmt.Errorf("embeddings: embed batch: %w", err)
		}
		if len(vecs) != len(texts) {
			return fmt.Errorf("embeddings: got %d vectors for %d texts", len(vecs), len(texts))
		}
		for i, vec := range vecs {
			embedder.Add(texts[i], metas[i], vec)
		}
		return nil
	}

	total := 0
	for rows.Next() {
		var text, meta string
		if err := rows.Scan(&text, &meta); err != nil {
			return fmt.Errorf("embeddings: scan doc: %w", err)
		}

		batchTexts = append(batchTexts, text)
		batchMetas = append(batchMetas, meta)

		if len(batchTexts) >= batchSize {
			if err := processBatch(batchTexts, batchMetas); err != nil {
				return err
			}
			total += len(batchTexts)
			if total%1000 == 0 {
				log.Printf("Embedded %d documents...", total)
			}
			batchTexts = batchTexts[:0]
			batchMetas = batchMetas[:0]
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("embeddings: iterate docs: %w", err)
	}

	// Process remaining
	if len(batchTexts) > 0 {
		if err := processBatch(batchTexts, batchMetas); err != nil {
			return err
		}
		total += len(batchTexts)
	}

	log.Printf("Embedding complete: %d documents", total)
	return nil
}
