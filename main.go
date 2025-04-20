package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/asim/reminder/api"
	"github.com/asim/reminder/app"
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
)

var (
	IndexFlag  = flag.Bool("index", false, "Index data for search. Stored at $HOME/reminder.idx")
	ExportFlag = flag.Bool("export", false, "Export the index data to $HOME/reminder.idx.gob.gz")
	ImportFlag = flag.Bool("import", false, "Import the index data from $HOME/reminder.idx.gob.gz")
	ServerFlag = flag.Bool("serve", false, "Run the server")
	EnvFlag    = flag.String("env", "dev", "Set the environment")
	WebFlag    = flag.Bool("web", false, "Without this flag, the lite version will be served")
)

var history = map[string][]string{}

func registerLiteRoutes(q *quran.Quran, n *names.Names, b *hadith.Volumes, apiHtml string) {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(apiHtml))
	})

	http.HandleFunc("/quran", func(w http.ResponseWriter, r *http.Request) {
		qhtml := app.RenderHTML("Quran", quran.Description, q.TOC())
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
		qhtml := app.RenderHTML(head, "", q.Get(ch).HTML())

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
		vhtml := app.RenderHTML(head, "", vv.HTML())

		w.Write([]byte(vhtml))
	})

	http.HandleFunc("/names", func(w http.ResponseWriter, r *http.Request) {
		qhtml := app.RenderHTML("Names", names.Description, n.TOC())
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
		qhtml := app.RenderHTML(head, "", n.Get(name).HTML())

		w.Write([]byte(qhtml))
	})

	http.HandleFunc("/hadith", func(w http.ResponseWriter, r *http.Request) {
		qhtml := app.RenderHTML("Hadith", hadith.Description, b.TOC())
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
		qhtml := app.RenderHTML(head, "", b.Get(ch).HTML())

		w.Write([]byte(qhtml))
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		shtml := app.RenderHTML("Search", "", app.SearchTemplate)
		w.Write([]byte(shtml))
	})
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

	// generate api doc
	ap := api.Load()
	apiHtml := app.RenderTemplate("API", "", ap.Markdown())

	// generate json
	qjson := q.JSON()
	njson := n.JSON()
	hjson := b.JSON()

	// index the quran in english
	indexed := make(chan bool, 1)

	if *IndexFlag {
		// create a separate index that's persisted
		// this is located in $HOME/reminder.idx
		// it will need to be exported afterwards
		sidx := search.New("reminder", true)

		fmt.Println("Indexing data")
		go func() {
			indexQuran(sidx, q)
			indexNames(sidx, n)
			indexHadith(sidx, b)
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

	if *WebFlag {
		http.Handle("/", app.Serve())
	} else {
		http.Handle("/", app.ServeLite())
		registerLiteRoutes(q, n, b, apiHtml)
	}

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

	http.HandleFunc("/api/quran/chapters", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(q.Index().Chapters)
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
		b, _ := io.ReadAll(r.Body)
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
		b, _ := io.ReadAll(r.Body)
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

		if r.Method == "GET" {
			var ctx string

			// look for the context cookie
			c, err := r.Cookie("session")
			if err == nil {
				ctx = c.Value
			}

			if len(ctx) == 0 {
				w.Write([]byte(`{}`))
				return
			}

			// pull the context which we only store in memory for now
			h, ok := history[ctx]
			if !ok {
				h = []string{}
			}

			out, _ := json.Marshal(map[string]interface{}{
				"session": ctx,
				"history": h,
			})

			w.Write(out)
			return
		}

		if r.Method == "POST" {
			b, _ := io.ReadAll(r.Body)
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
			answerMD := string(app.Render([]byte(answer)))

			output, _ := json.Marshal(map[string]interface{}{
				"q":          q,
				"answer":     answerMD,
				"references": res,
			})
			w.Write(output)

			var ctx string

			// look for the context cookie
			c, err := r.Cookie("session")
			if err == nil {
				ctx = c.Value
				h, ok := history[ctx]
				if !ok {
					h = []string{}
				}
				h = append([]string{q, answerMD}, h...)
				history[ctx] = h
			}

			return
		}
	})

	if *ServerFlag {
		fmt.Println("Starting server :8080")
		if err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if *EnvFlag == "dev" {
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("Access-Control-Allow-Credentials", "true")

				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}
			}

			http.DefaultServeMux.ServeHTTP(w, r)
		})); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}

	// wait for index
	<-indexed
}
