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
		return err
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

	c, err := i.DB.GetOrCreateCollection(i.Name, nil, nil)
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

	c, err := i.DB.GetOrCreateCollection(i.Name, nil, nil)
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
	res, err := i.Col.Query(context.TODO(), v, 100, nil, nil)
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

	if persist {
		path := filepath.Join(u.HomeDir, name+".idx")

		var err error
		db, err = chromem.NewPersistentDB(path, false)
		if err != nil {
			panic(err)
		}

		c, err = db.GetOrCreateCollection(name, nil, nil)
		if err != nil {
			panic(err)
		}
	} else {
		db = chromem.NewDB()
		c, _ = db.CreateCollection(name, nil, nil)
	}

	return &Index{
		Home: u.HomeDir,
		Name: name,
		DB:   db,
		Col:  c,
	}
}
