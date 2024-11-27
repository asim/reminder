package html

import (
	"fmt"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var Template = `
<html>
  <head>
    <title>%s | Reminder</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
    html, body { height: 100%%; width: 100%%; margin: 0; padding: 0;}
    #head { margin-bottom: 2.5em; }
    #head a { margin-right: 10px; color: black; font-weight: bold; text-decoration: none; }
    #container { height: 100%%; max-width: 1024px; margin: 0 auto; padding: 25px;}
    #content { padding-bottom: 100px; }
    #content p { padding: 0 0 25px 0; margin: 0; }
    #search { margin-top: 10px; } #q { padding: 10px; width: 100%%; }
@font-face {
    font-family: 'arabic';
    src: url('/files/arabic.otf') format('opentype');
    font-weight: normal;
    font-style: normal;
}
.arabic {
  font-family: 'arabic';
  font-size: 1.5em;
}
    </style>
  </head>
  <body>
    <div id="container">
      <div id="head">
        <a href="/">[Reminder]</a>
        <a href="/quran">Quran</a>
        <a href="/names">Names</a>
        <a href="/hadith">Hadith</a>
      </div>
      <div id="search">
        <form action="/search" method="post"><input id="q" name=q placeholder="Ask a question"></form>
      </div>
      <div id="content">
      %s
      </div>
    </div>
  </body>
</html>
`

func Render(md []byte) []byte {
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

func RenderString(v string) string {
	return string(Render([]byte(v)))
}

func RenderTemplate(title string, text string) string {
	return fmt.Sprintf(Template, title, RenderString(text))
}
