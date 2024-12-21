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
code {
  background: whitesmoke;
  padding: 5px;
  border-radius: 5px;
}
    </style>
  </head>
  <body>
    <div id="container">
      <div id="head">
        <a href="/">[Reminder]</a>
        <a href="/api">API</a>
        <a href="/quran">Quran</a>
        <a href="/names">Names</a>
        <a href="/hadith">Hadith</a>
        <a href="/search">Search</a>
      </div>
      <div id="content">%s</div>
    </div>
  </body>
</html>
`

var Search = `
<style>
#content p { padding: 0 0 10px 0; }
#resp { padding-bottom: 10px;}
#expand { text-decoration: underline; }
#expand:hover { cursor: pointer; }
.ref { font-size: small; }
#search { margin-top: 10px; }
#q { padding: 10px; width: 100%; }
</style>
<div id="search">
  <form id="question" action="/search" method="post"><input id="q" name=q placeholder="Ask a question" autocomplete="off"></form>
</div>
<div id="resp"></div>
<div id="answer"></div>
<script>
function expand(el) {
      var ref = el.nextSibling;

      if (ref.style.display == 'none') {
          ref.style.display = 'block';
      } else {
          ref.style.display = 'none';
      }
}

function reference(el) {
	return "<div>Text: " + el.Text + "<br>Metadata: " + JSON.stringify(el.Metadata) + "<br>Score: " + el.Score + "</div>";
}

document.addEventListener('DOMContentLoaded', function(){
    var form = document.getElementById("question");
    form.addEventListener("submit", function(ev) {
	var ans = document.getElementById("answer");
	var resp = document.getElementById("resp");
        ev.preventDefault();
	var q = document.getElementById("q");
	var xhr = new XMLHttpRequest();
	var url = "/api/search";
	xhr.open("POST", url, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	xhr.onreadystatechange = function () {
	    if (xhr.readyState === 4 && xhr.status === 200) {
		var json = JSON.parse(xhr.responseText);
		var text = "<p><b>Q</b>: " + json.q + "</p><p><b>A</b>: " + json.answer + "</p>";
		text += "<div id=expand onclick='expand(this); return false;'>References<br><br></div>";
		text += "<div class=ref style='display: none;'>";
		for (i = 0; i < json.references.length; i++) {
                    text += reference(json.references[i]) + "<br><br>";
		}
		text += "</div>";
		ans.innerHTML = text + ans.innerHTML; 
		resp.innerText = "";
	    }
	};
	var data = JSON.stringify({"q": q.value});
	resp.innerText = "asking...";
	xhr.send(data);
	q.value = '';
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
