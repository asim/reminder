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
.chapter {
  margin: 10px;
  border: 1px solid grey;
  padding: 10px;
  display: inline-block;
}
    </style>
  </head>
  <body>
    <div id="container">
      <div id="head">
        <a href="/">[Reminder]</a>
        <a href="/about">About</a>
        <a href="/quran">Quran</a>
        <a href="/names">Names</a>
        <a href="/hadith">Hadith</a>
      </div>
      <div id="search">
        <form id="question" action="/search" method="post"><input id="q" name=q placeholder="Ask a question"></form>
      </div>
      <div id="content">
      %s
      </div>
    </div>
  </body>
</html>
`

var Index = `
<div id="answer"></div>
<script>
document.addEventListener('DOMContentLoaded', function(){
    var form = document.getElementById("question");
    form.addEventListener("submit", function(ev) {
        ev.preventDefault();
	var q = document.getElementById("q");

	var xhr = new XMLHttpRequest();
	var url = "/search.json";
	xhr.open("POST", url, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	xhr.onreadystatechange = function () {
	    if (xhr.readyState === 4 && xhr.status === 200) {
		var json = JSON.parse(xhr.responseText);
		var ans = document.getElementById("answer");
		var text = "<p><b>Q</b>: " + q.value + "</p><p><b>A</b>: " + json.answer + "</p>";
		ans.innerHTML = text + ans.innerHTML; 
		q.value = '';
	    }
	};
	var data = JSON.stringify({"q": q.value});
	xhr.send(data);
    });
}, false);
</script>
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

func RenderHTML(title, html string) string {
	return fmt.Sprintf(Template, title, html)
}

func RenderString(v string) string {
	return string(Render([]byte(v)))
}

func RenderTemplate(title string, text string) string {
	return fmt.Sprintf(Template, title, RenderString(text))
}
