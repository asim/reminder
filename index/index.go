package index

import (
	"context"
	"fmt"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
	"github.com/philippgille/chromem-go"
)

type Index struct {
	Home string
	Name string
	DB   *chromem.DB
	Col  *chromem.Collection
}

type Result struct {
	Text     string
	Score    float32
	Metadata map[string]string
}

func (i *Index) Export() error {
	path := filepath.Join(i.Home, i.Name+".idx.gob.gz")

	return i.DB.ExportToFile(path, true, "")
}

func (i *Index) Import() error {
	path := filepath.Join(i.Home, i.Name+".idx.gob.gz")

	return i.DB.ImportFromFile(path, "")
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

func New(name string) *Index {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(u.HomeDir, name+".idx")

	db, err := chromem.NewPersistentDB(path, false)
	if err != nil {
		panic(err)
	}

	c, err := db.GetOrCreateCollection(name, nil, nil)
	if err != nil {
		panic(err)
	}

	return &Index{
		Home: u.HomeDir,
		Name: name,
		DB:   db,
		Col:  c,
	}
}
