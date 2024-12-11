package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/html"
	"github.com/asim/reminder/html/files"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
)

var (
	IndexFlag    = flag.Bool("index", false, "Index data for search. Stored at $HOME/reminder.idx")
	ExportFlag   = flag.Bool("export", false, "Export the index data to $HOME/reminder.idx.gob.gz")
	ImportFlag   = flag.Bool("import", false, "Import the index data from $HOME/reminder.idx.gob.gz")
	GenerateFlag = flag.Bool("generate", false, "Generate the html files")
	ServerFlag   = flag.Bool("serve", false, "Run the server")
)

func indexContent(idx *search.Index, md map[string]string, text string) {
	// index the documents
	// TODO: use original json
	lines := strings.Split(text, "\n")

	fmt.Println("Indexing")

	if err := idx.Store(md, lines...); err != nil {
		fmt.Println("Error indexing", err)
	}
}

func indexQuran(idx *search.Index, q *quran.Quran) {
	fmt.Println("Indexing Quran")

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
	fmt.Println("Indexing Names")

	for _, name := range *n {
		indexContent(idx, map[string]string{
			"source":  "names",
			"meaning": name.Meaning,
			"english": name.English,
			"arabic":  name.Arabic,
		}, strings.Join([]string{name.Meaning, name.English, name.Description}, " - "))
	}
}

func indexHadith(idx *search.Index, b *hadith.Volumes) {
	fmt.Println("Indexing Hadith")

	for _, volume := range *b {
		for _, book := range volume.Books {
			for _, hadith := range book.Hadiths {

				indexContent(idx, map[string]string{
					"source": "bukhari",
					"volume": volume.Name,
					"book":   book.Name,
					"by":     hadith.By,
					"info":   hadith.Info,
				}, hadith.Text)
			}
		}
	}
}

func gen(idx *search.Index, q string) (string, []string) {
	res, err := idx.Query(q)
	if err != nil {
		return "", nil
	}

	var contexts []string

	for _, r := range res {
		b, _ := json.Marshal(r)
		// TODO: maybe just provide text
		contexts = append(contexts, string(b))
	}

	return askLLM(context.TODO(), contexts, q), contexts
}

var questions = []string{
	"What is the Reminder?",
	"What is the Quran?",
	"What is the Hadith?",
	"Who is Allah?",
	"Who is the prophet Muhammad",
	"Why do we 'worship' Allah?",
	"How do we 'worship' Allah?",
	"What happens when we die?",
	"How do I remember Allah?",
	"How do I become Muslim?",
}

func main() {
	flag.Parse()

	// create a new index
	idx := search.New("reminder", false)

	// Load the pre-existing data
	if err := idx.Load(); err != nil {
		fmt.Println(err)
	}

	// render the markdown as html
	if *GenerateFlag {
		fmt.Println("Loading data")
		// load data
		q := quran.Load()
		n := names.Load()
		b := hadith.Load()

		fmt.Println("Generating html")
		text := q.HTML()
		name := n.Markdown()
		books := b.Markdown()

		thtml := html.RenderTemplate("Quran", text)
		nhtml := html.RenderTemplate("Names", name)
		vhtml := html.RenderTemplate("Hadith", books)
		shtml := html.RenderHTML("Search", html.Search)

		var about string

		for _, q := range questions {
			a, _ := gen(idx, q)
			about += fmt.Sprintf("# %s", q)
			about += fmt.Sprintln()
			about += fmt.Sprintf("%s", a)
			about += fmt.Sprintln()
		}

		ihtml := html.RenderTemplate("Index", about)

		os.WriteFile(filepath.Join(".", "html", "files", "index.html"), []byte(ihtml), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "search.html"), []byte(shtml), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "quran.html"), []byte(thtml), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "names.html"), []byte(nhtml), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "hadith.html"), []byte(vhtml), 0644)

		return
	}

	// index the quran in english
	indexed := make(chan bool, 1)

	if *IndexFlag {
		fmt.Println("Loading data")
		// load data
		q := quran.Load()
		n := names.Load()
		b := hadith.Load()

		fmt.Println("Indexing data")
		go func() {
			indexQuran(idx, q)
			indexNames(idx, n)
			indexHadith(idx, b)
			// done
			close(indexed)
		}()
	} else {
		close(indexed)
	}

	if *ExportFlag {
		fmt.Println("Exporting index")
		if err := idx.Export(); err != nil {
			fmt.Println(err)
		}
		return
	}

	if *ImportFlag {
		fmt.Println("Importing index")
		if err := idx.Import(); err != nil {
			fmt.Println(err)
		}
	}

	// load the data from html

	ihtml := files.Get("index.html")
	shtml := files.Get("search.html")
	thtml := files.Get("quran.html")
	nhtml := files.Get("names.html")
	vhtml := files.Get("hadith.html")
	otf := files.Get("arabic.otf")

	http.HandleFunc("/files/arabic.otf", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(otf))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(ihtml))
	})

	http.HandleFunc("/quran", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(thtml))
	})

	http.HandleFunc("/names", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(nhtml))
	})

	http.HandleFunc("/hadith", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(vhtml))
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(shtml))
	})

	http.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-indexed:
		default:
			// not indexed yet because blocked
			w.Write([]byte(`{"error": "Indexing content"}`))
			return
		}

		// indexed or no index of any kind

		if r.Method == "POST" {
			b, _ := ioutil.ReadAll(r.Body)
			var data map[string]interface{}
			json.Unmarshal(b, &data)

			q := data["q"].(string)

			res, err := idx.Query(q)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var contexts []string

			for _, r := range res {
				b, _ := json.Marshal(r)
				// TODO: maybe just provide text
				contexts = append(contexts, string(b))
			}

			answer := askLLM(context.TODO(), contexts, q)

			output, _ := json.Marshal(map[string]interface{}{
				"q":          q,
				"answer":     answer,
				"references": res,
			})
			w.Write(output)

			return
		}
	})

	if *ServerFlag {
		fmt.Println("Starting server :8080")

		http.ListenAndServe(":8080", nil)
	}
}
