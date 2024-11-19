package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/asim/reminder/files"
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/index"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var (
	IndexFlag    = flag.Bool("index", false, "Index data for search. Stored at $HOME/reminder.idx")
	ExportFlag   = flag.Bool("export", false, "Export the index data to $HOME/reminder.idx.gob.gz")
	ImportFlag   = flag.Bool("import", false, "Import the index data from $HOME/reminder.idx.gob.gz")
	GenerateFlag = flag.Bool("generate", false, "Generate the html files")
	ServerFlag   = flag.Bool("serve", false, "Run the server")
)

var htmlTemplate = `
<html>
  <head>
    <title>%s | Reminder</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
    html, body { height: 100%%; width: 100%%; margin: 0; padding: 0;}
    #container { height: 100%%; max-width: 1024px; margin: 0 auto; padding: 25px;}
    #head a { margin-right: 10px; color: black; font-weight: bold; }
    #content { padding-bottom: 100px; }
    #content p { padding: 50px 10px 50px 10px; border-bottom: 1px solid grey; margin: 0; }
    </style>
  </head>
  <body>
    <div id="container">
      <div id="head">
        <a href="/">Reminder</a>
        <a href="/names">Names</a>
        <a href="/hadith">Hadith</a>
        <a href="/search">Search</a>
      </div>
      <div id="content">
      %s
      </div>
    </div>
  </body>
</html>
`

func render(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func result(question, answer string, res []*index.Result) []byte {
	data := `## Question: ` + question
	data += fmt.Sprintln()
	data += `# Answer`
	data += fmt.Sprintln()
	data += fmt.Sprintln()
	data += "### " + answer
	data += fmt.Sprintln()
	data += "#### References"
	data += fmt.Sprintln()
	data += "***"
	data += fmt.Sprintln()

	for _, r := range res {
		switch r.Metadata["source"] {
		case "quran":
			data += fmt.Sprintf("#### Quran - %s %s:%s", r.Metadata["name"], r.Metadata["chapter"], r.Metadata["verse"])
		case "names":
			data += fmt.Sprintf("#### Name - %s", r.Metadata["meaning"])
		case "bukhari":
			data += fmt.Sprintf("#### Hadith - %s %s", r.Metadata["info"], r.Metadata["by"])
		}

		data += fmt.Sprintln()
		data += r.Text
		data += fmt.Sprintln()
	}

	data = fmt.Sprintf(htmlTemplate, "Results", string(render([]byte(data))))

	return []byte(data)
}

func indexContent(idx *index.Index, md map[string]string, text string) {
	// index the documents
	// TODO: use original json
	lines := strings.Split(text, "\n")

	fmt.Println("Indexing")

	if err := idx.Store(md, lines...); err != nil {
		fmt.Println("Error indexing", err)
	}
}

func indexQuran(idx *index.Index, q *quran.Quran) {
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

func indexNames(idx *index.Index, n *names.Names) {
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

func indexHadith(idx *index.Index, b *hadith.Volumes) {
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

func main() {
	flag.Parse()

	// render the markdown
	if *GenerateFlag {
		fmt.Println("Loading data")
		// load data
		q := quran.Load()
		n := names.Load()
		b := hadith.Load()

		fmt.Println("Generating html")
		text := q.Markdown()
		name := n.Markdown()
		books := b.Markdown()

		thtml := fmt.Sprintf(htmlTemplate, "Reminder", string(render([]byte(text))))
		nhtml := fmt.Sprintf(htmlTemplate, "Names", string(render([]byte(name))))
		vhtml := fmt.Sprintf(htmlTemplate, "Hadith", string(render([]byte(books))))

		os.WriteFile(filepath.Join(".", "files", "reminder.html"), []byte(thtml), 0644)
		os.WriteFile(filepath.Join(".", "files", "names.html"), []byte(nhtml), 0644)
		os.WriteFile(filepath.Join(".", "files", "hadith.html"), []byte(vhtml), 0644)
		return
	}

	// create a new index
	idx := index.New("reminder", false)

	// Load the pre-existing data
	if err := idx.Load(); err != nil {
		fmt.Println(err)
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

	thtml := files.Get("reminder")
	nhtml := files.Get("names")
	vhtml := files.Get("hadith")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(thtml))
	})

	http.HandleFunc("/names", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(nhtml))
	})

	http.HandleFunc("/hadith", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(vhtml))
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-indexed:
		default:
			// not indexed yet because blocked
			w.Write([]byte("Indexing content"))
			return
		}

		// indexed or no index of any kind

		if r.Method == "POST" {
			r.ParseForm()
			q := r.Form.Get("q")
			if len(q) == 0 {
				return
			}
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

			// create a markdown
			md := result(q, answer, res)
			w.Write(md)
			return
		}

		// render search form
		form := `<style>#search { margin-top: 25px; } #q { padding: 10px; width: 100%; }</style>
		<form id="search" action="/search" method="post"><input id="q" name=q placeholder=Search></form>`
		html := fmt.Sprintf(htmlTemplate, "Search", form)
		w.Write([]byte(html))
	})

	if *ServerFlag {
		fmt.Println("Starting server :8080")

		http.ListenAndServe(":8080", nil)
	}

	args := os.Args[1:]

	if len(args) <= 1 {
		// repl
		var history []string

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print(">: ")
			q, _ := reader.ReadString('\n')
			if len(q) == 0 {
				continue
			}

			res, err := idx.Query(q)
			if err != nil {
				fmt.Println(err)
				continue
			}

			var contexts []string

			for _, r := range res {
				b, _ := json.Marshal(r)
				// TODO: maybe just provide text
				contexts = append(contexts, string(b))
			}

			// ask the question
			answer := askLLM(context.TODO(), append(contexts, history...), q)

			// append to history
			history = append(history, q)
			history = append(history, answer)

			if len(history) > 100 {
				history = history[10:]
			}

			// write the response
			fmt.Println(answer)
		}

	}

	command := args[0]

	if command == "search" {
		q := strings.Join(args[1:], " ")

		if len(q) == 0 {
			return
		}

		res, err := idx.Query(q)
		if err != nil {
			fmt.Println(err)
			return
		}

		var contexts []string

		for _, r := range res {
			b, _ := json.Marshal(r)
			// TODO: maybe just provide text
			contexts = append(contexts, string(b))
		}

		answer := askLLM(context.TODO(), contexts, q)
		fmt.Println(answer)
	}
}
