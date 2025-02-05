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
	"strconv"
	"strings"

	"github.com/asim/reminder/api"
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

	for _, volume := range b.Contents {
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

	// load data
	q := quran.Load()
	n := names.Load()
	b := hadith.Load()

	// render the markdown as html
	if *GenerateFlag {
		fmt.Println("Loading data")

		fmt.Println("Generating html")
		text := q.HTML()
		name := n.HTML()
		books := b.Markdown()

		vhtml := html.RenderTemplate("Hadith", books)
		thtml := html.RenderTemplate("Quran", text)
		nhtml := html.RenderHTML("Names", name)
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

		ap := api.Load()
		apiHtml := html.RenderTemplate("API", ap.Markdown())

		// write html files
		os.WriteFile(filepath.Join(".", "html", "files", "api.html"), []byte(apiHtml), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "index.html"), []byte(ihtml), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "search.html"), []byte(shtml), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "quran.html"), []byte(thtml), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "names.html"), []byte(nhtml), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "hadith.html"), []byte(vhtml), 0644)

		// write json files
		os.WriteFile(filepath.Join(".", "html", "files", "quran.json"), q.JSON(), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "names.json"), n.JSON(), 0644)
		os.WriteFile(filepath.Join(".", "html", "files", "hadith.json"), b.JSON(), 0644)

		return
	}

	// index the quran in english
	indexed := make(chan bool, 1)

	if *IndexFlag {
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

	apiHtml := files.Get("api.html")
	ihtml := files.Get("index.html")
	shtml := files.Get("search.html")
	//thtml := files.Get("quran.html")
	//nhtml := files.Get("names.html")
	//vhtml := files.Get("hadith.html")
	otf := files.Get("arabic.otf")
	qjson := files.Get("quran.json")
	njson := files.Get("names.json")
	hjson := files.Get("hadith.json")

	http.HandleFunc("/files/arabic.otf", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(otf))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(ihtml))
	})

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(apiHtml))
	})

	http.HandleFunc("/quran", func(w http.ResponseWriter, r *http.Request) {
		qhtml := html.RenderHTML("Quran", q.TOC())

		w.Write([]byte(qhtml))
	})

	http.HandleFunc("/quran/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			return
		}

		ch, _ := strconv.Atoi(id)

		if ch < 1 || ch > 114 {
			return
		}

		head := fmt.Sprintf("%d | Quran", ch)
		qhtml := html.RenderHTML(head, q.Get(ch).HTML())

		w.Write([]byte(qhtml))
	})

	http.HandleFunc("/quran/{id}/{ver}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			return
		}

		ver := r.PathValue("ver")
		if len(ver) == 0 {
			return
		}

		ch, _ := strconv.Atoi(id)
		ve, _ := strconv.Atoi(ver)

		if ch < 1 || ch > 114 {
			return
		}

		cc := q.Get(ch)

		if ve < 1 || ve > len(cc.Verses) {
			return
		}

		vv := cc.Verses[ve-1]

		head := fmt.Sprintf("%d:%d | Quran", ch, ve)
		vhtml := html.RenderHTML(head, vv.HTML())

		w.Write([]byte(vhtml))
	})

	http.HandleFunc("/names", func(w http.ResponseWriter, r *http.Request) {
		qhtml := html.RenderHTML("Names", n.TOC())

		w.Write([]byte(qhtml))
	})

	http.HandleFunc("/names/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			return
		}

		name, _ := strconv.Atoi(id)

		if name < 1 || name > len(*n) {
			return
		}

		head := fmt.Sprintf("%d | Names", name)
		qhtml := html.RenderHTML(head, n.Get(name).HTML())

		w.Write([]byte(qhtml))
	})

	http.HandleFunc("/hadith", func(w http.ResponseWriter, r *http.Request) {
		qhtml := html.RenderHTML("Hadith", b.TOC())

		w.Write([]byte(qhtml))
	})

	http.HandleFunc("/hadith/{book}", func(w http.ResponseWriter, r *http.Request) {
		book := r.PathValue("book")
		if len(book) == 0 {
			return
		}

		ch, _ := strconv.Atoi(book)

		if ch < 1 || ch > len(b.Books) {
			return
		}

		head := fmt.Sprintf("%d | Hadith", ch)
		qhtml := html.RenderHTML(head, b.Get(ch).HTML())

		w.Write([]byte(qhtml))
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(shtml))
	})

	http.HandleFunc("/api/quran", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(qjson))
	})

	http.HandleFunc("/api/quran/{chapter}", func(w http.ResponseWriter, r *http.Request) {
		ch := r.PathValue("chapter")
		if len(ch) == 0 {
			return
		}

		chapter, _ := strconv.Atoi(ch)
		if chapter < 1 || chapter > 114 {
			return
		}

		b := q.Get(chapter).JSON()

		w.Write(b)
	})

	http.HandleFunc("/api/quran/{chapter}/{verse}", func(w http.ResponseWriter, r *http.Request) {
		ch := r.PathValue("chapter")
		if len(ch) == 0 {
			return
		}

		chapter, _ := strconv.Atoi(ch)
		if chapter < 1 || chapter > 114 {
			return
		}

		ve := r.PathValue("verse")
		if len(ch) == 0 {
			return
		}

		cc := q.Get(chapter)

		verse, _ := strconv.Atoi(ve)
		if verse < 1 || verse > len(cc.Verses) {
			return
		}

		vee := cc.Verses[verse-1]
		b := vee.JSON()

		w.Write(b)
	})

	http.HandleFunc("/api/names", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(njson))
	})

	http.HandleFunc("/api/hadith", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(hjson))
	})

	http.HandleFunc("/api/hadith/{book}", func(w http.ResponseWriter, r *http.Request) {
		bk := r.PathValue("book")
		if len(bk) == 0 {
			return
		}

		book, _ := strconv.Atoi(bk)
		if book < 1 || book > len(b.Books) {
			return
		}

		b := b.Get(book).JSON()

		w.Write(b)
	})

	http.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(b, &data)

		q := data["q"].(string)

		prompt := `Generate a detailed summary for the following with it's meaning and origin, output the response as JSON with the fields: 
		name, description, summary. Each field itself should be a string.

		%s
		`

		answer := askLLM(r.Context(), nil, fmt.Sprintf(prompt, q))
		w.Write([]byte(answer))
	})

	http.HandleFunc("/api/translate", func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(b, &data)

		q := data["q"].(string)

		prompt := `Translate the following into a modern interpretation, transliterate and then word by word. 
		For each word provide 3 alternatives and a transliteration in english. Generate the output as JSON. 
		The response will be served via an API so ensure it's entirely json compliant with no markdown.
		Ensure consistency in the output using fields translation, transliteration and word_by_word.
		The word_by_word field itself should be an array with word in english, arabic, translations and 
		transliteration.:

		%s`

		answer := askLLM(r.Context(), nil, fmt.Sprintf(prompt, q))
		w.Write([]byte(answer))
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
				for k, v := range r.Metadata {
					delete(r.Metadata, k)
					r.Metadata[strings.ToLower(k)] = v
				}

				b, _ := json.Marshal(r)
				// TODO: maybe just provide text
				contexts = append(contexts, string(b))
			}

			answer := askLLM(r.Context(), contexts, q)

			output, _ := json.Marshal(map[string]interface{}{
				"q":          q,
				"answer":     string(html.Render([]byte(answer))),
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
