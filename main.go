package main

import (
	"fmt"
	"net/http"

	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/hadith"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var template = `
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
    html, body { height: 100%%; width: 100%%; margin: 0; padding: 0;}
    #container { height: 100%%; max-width: 1024px; margin: 0 auto; padding: 25px;}
    #head a { margin-right: 10px; color: black; font-weight: bold; }
    #content { padding-bottom: 100px; }
    </style>
  </head>
  <body>
    <div id="container">
      <div id="head">
        <a href="/">Quran</a>
        <a href="/names">Names</a>
        <a href="/hadith">Hadith</a>
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

func main() {
	text := quran.Markdown()
	name := names.Markdown()
	books := hadith.Markdown()

	thtml := fmt.Sprintf(template, string(render([]byte(text))))
	nhtml := fmt.Sprintf(template, string(render([]byte(name))))
	vhtml := fmt.Sprintf(template, string(render([]byte(books))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(thtml))
	})

	http.HandleFunc("/names", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(nhtml))
	})

	http.HandleFunc("/hadith", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(vhtml))
	})

	http.ListenAndServe(":8080", nil)
}
