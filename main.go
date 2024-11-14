package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/index"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var (
	IndexFlag = flag.Bool("index", false, "Index data for search")
)

var template = `
<html>
  <head>
    <title>%s | Reminder</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
    html, body { height: 100%%; width: 100%%; margin: 0; padding: 0;}
    #container { height: 100%%; max-width: 1024px; margin: 0 auto; padding: 25px;}
    #head a { margin-right: 10px; color: black; font-weight: bold; }
    #content { padding-bottom: 100px; }
    #content p { padding: 50px 10px 50px 10px; border-bottom: 1px solid grey; }
    </style>
  </head>
  <body>
    <div id="container">
      <div id="head">
        <a href="/">Quran</a>
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

func result(res []*index.Result) []byte {
	data := `# Results`
	data += fmt.Sprintln()
	data += fmt.Sprintln()

	for _, r := range res {
		data += fmt.Sprintf("source: %s", r.Metadata["source"])
		data += r.Text
		data += fmt.Sprintln()
		data += fmt.Sprintln()
	}

	data = fmt.Sprintf(template, data)

	return render([]byte(data))
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

	// load data
	q := quran.Load()
	n := names.Load()
	b := hadith.Load()

	// render the markdown
	text := q.Markdown()
	name := n.Markdown()
	books := b.Markdown()

	// create a new index
	fmt.Println("Creating index")
	idx := index.New("reminder")

	// index the quran in english
	indexed := make(chan bool, 1)

	if *IndexFlag {
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

	thtml := fmt.Sprintf(template, "Home", string(render([]byte(text))))
	nhtml := fmt.Sprintf(template, "Names", string(render([]byte(name))))
	vhtml := fmt.Sprintf(template, "Hadith", string(render([]byte(books))))

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

			// create a markdown
			md := result(res)
			w.Write(md)
			return
		}

		// render search form
		form := `<style>#search { margin-top: 25px; } #q { padding: 10px; width: 100%; }</style>
		<form id="search" action="/search" method="post"><input id="q" name=q placeholder=Search></form>`
		html := fmt.Sprintf(template, "Search", form)
		w.Write([]byte(html))
	})

	fmt.Println("Starting server :8080")

	http.ListenAndServe(":8080", nil)
}
