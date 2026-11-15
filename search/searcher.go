package search

import "os"

// Searcher is the query surface shared by the two index implementations.
//
// FTS is the locally refreshed SQLite keyword index used for every deployment.
// Index is the optional prebuilt vector supplement, using OpenAI embeddings.
type Searcher interface {
	Store(md map[string]string, content ...string) error
	Query(q string) ([]*Result, error)
}

// Semantic reports whether credentials for the optional vector index are present.
// Enabled by default when credentials exist; SEARCH_SEMANTIC=false disables it.
func Semantic() bool {
	return os.Getenv("OPENAI_API_KEY") != ""
}
