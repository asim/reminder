package app

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed html/*
var htmlFiles embed.FS

//go:embed all:dist
var distFiles embed.FS

var Template = `
<html>
  <head>
    <title>%s | Reminder</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="manifest" href="/manifest.webmanifest">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:ital,opsz,wght@0,6..12,200..1000;1,6..12,200..1000&display=swap" rel="stylesheet">
    <style>
    html, body { height: 100%%; width: 100%%; margin: 0; padding: 0; font-family: "Nunito Sans", serif; }
    a { color: #333333; }
    a:visited { color: #333333;}
    #head {
      position: fixed;
      background: white;
      top: 0;
      width: 100%%;
    }
    #head a { margin-right: 10px; color: black; font-weight: bold; text-decoration: none; }
    #container { height: 100%%; max-width: 1024px; margin: 0 auto; padding: 25px;}
    #content { padding-bottom: 100px; }
    #content p { padding: 0 0 25px 0; margin: 0; }
    #desc { margin-bottom: 10px; }
    #title { margin-top: 50px; font-size: 1.2em; font-weight: bold; margin-bottom: 10px; }
    li { margin-bottom: 5px; }
@font-face {
    font-family: 'arabic';
    src: url('/arabic.otf') format('opentype');
    font-weight: normal;
    font-style: normal;
}
.arabic {
  font-family: 'arabic';
  font-size: 1.8em;
}
.chapter a {
  text-decoration: none;
  margin-bottom: 10px;
  border: 1px solid grey;
  border-radius: 5px;
  padding: 10px;
  display: block;
}
code {
  background: whitesmoke;
  padding: 5px;
  border-radius: 5px;
}
.dots {
  font-size: 1.5em;
  text-align: center;
  margin-bottom: 25px;
  padding: 25px;
}
.right {
  text-align: right;
}
.block {
    text-decoration: none;
    margin-bottom: 10px;
    border: 1px solid grey;
    border-radius: 5px;
    padding: 10px;
    display: block;
}
@media only screen and (max-width: 600px) {
  #head a { margin-right: 5px; }
}
#brand {
  display: inline-block;
  padding: 20px;
}
#brand a {
  border: 1px solid grey;
  border-radius: 5px;
  padding: 5px;
}
#nav {
 float: right;
 padding: 20px 20px 20px 0;
}
    </style>
  </head>
  <body>
    <div id="head">
      <div id="brand">
        <a href="/">&nbsp;R&nbsp;</a>
      </div>
      <div id="nav">
        <a href="/daily">Daily</a>
        <a href="/quran">Quran</a>
        <a href="/names">Names</a>
        <a href="/hadith">Hadith</a>
        <a href="/search">Search</a>
      </div>
      <button id="install" hidden>Install PWA</button>
    </div>
    <div id="container">
      <div id="title">%s</div>
      <div id="desc">%s</div>
      <div id="content">%s</div>
    </div>
    </div>

  <script>
      if (navigator.serviceWorker) {
        navigator.serviceWorker.register (
          '/reminder.js',
          {scope: '/'}
        )
      }
  </script>
  <script>
        let installPrompt = null;
        const installButton = document.querySelector("#install");

        window.addEventListener("beforeinstallprompt", (event) => {
          event.preventDefault();
          installPrompt = event;
          installButton.removeAttribute("hidden");
        });

        installButton.addEventListener("click", async () => {
          if (!installPrompt) {
            return;
          }
          const result = await installPrompt.prompt();
          disableInAppInstallPrompt();
        });

        function disableInAppInstallPrompt() {
          installPrompt = null;
          installButton.setAttribute("hidden", "");
        }

        window.addEventListener("appinstalled", () => {
          disableInAppInstallPrompt();
        });

        function disableInAppInstallPrompt() {
          installPrompt = null;
          installButton.setAttribute("hidden", "");
        }
  </script>
  </body>
</html>
`

var SearchTemplate = `
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
  <form id="question" action="/search" method="post"><input id="q" name=q placeholder="Ask a question" autocomplete="off" autofocus></form>
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
	return "<div>Text: " + el.text + "<br>Metadata: " + JSON.stringify(el.metadata) + "<br>Score: " + el.score + "</div>";
}

function getCookie(name) {
    var cookies = document.cookie.split(';');
    for(var i=0 ; i < cookies.length ; ++i) {
        var pair = cookies[i].trim().split('=');
        if(pair[0] == name)
            return pair[1];
    }
    return null;
};

function setCookie(name, value) {
    document.cookie = name + "=" + value;
}

document.addEventListener('DOMContentLoaded', function(){
    // check if session exists
    var uuid = getCookie("session");

    if (uuid == undefined) {
        uuid = self.crypto.randomUUID();
	setCookie("session", uuid);
    }

    var url = "/api/search";

    // attempt to get the existing responses
    var xhr = new XMLHttpRequest();
    xhr.open("GET", url, true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
	  var ans = document.getElementById("answer");
	  var json = JSON.parse(xhr.responseText);
	  json.history.forEach(function(el) {
		ans.innerHTML += "<p>" + el + "</p>";
	  });
        }
    };

    xhr.send(null);

    var form = document.getElementById("question");
    form.addEventListener("submit", function(ev) {
	var ans = document.getElementById("answer");
	var resp = document.getElementById("resp");
        ev.preventDefault();
	var q = document.getElementById("q");
	var xhr = new XMLHttpRequest();
	xhr.open("POST", url, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	xhr.onreadystatechange = function () {
	    if (xhr.readyState === 4 && xhr.status === 200) {
		var json = JSON.parse(xhr.responseText);
		var text = "<p>" + json.q + "</p><p>" + json.answer + "</p>";
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
	resp.innerText = "seeking...";
	xhr.send(data);
	q.value = '';
    });
}, false);
</script>
`

var Index = `
<style>
#app a {
    text-decoration: none;
    margin-bottom: 10px;
    border: 1px solid grey;
    border-radius: 5px;
    padding: 10px;
    display: block;
}
</style>

<div id="app">
        <a href="/daily">Daily Reminder</a>
        <a href="/quran">Read the Quran</a>
        <a href="/names">Names of Allah</a>
        <a href="/hadith">Hadith (Bukhari)</a>
        <a href="/search">Ask a Question</a>
</div>
`

func Get(name string) string {
	f, err := htmlFiles.ReadFile("html/" + name)
	if err != nil {
		return ""
	}

	return string(f)
}

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

func RenderHTML(title, desc, html string) string {
	return fmt.Sprintf(Template, title, title, desc, html)
}

func RenderString(v string) string {
	return string(Render([]byte(v)))
}

func RenderTemplate(title string, desc, text string) string {
	return fmt.Sprintf(Template, title, title, desc, RenderString(text))
}

func ServeLite() http.Handler {
	var staticFS = fs.FS(htmlFiles)
	htmlContent, err := fs.Sub(staticFS, "html")
	if err != nil {
		log.Fatal(err)
	}

	return http.FileServer(http.FS(htmlContent))
}

func ServeWeb() http.Handler {
	var staticFS = fs.FS(distFiles)
	htmlContent, err := fs.Sub(staticFS, "dist")
	if err != nil {
		log.Fatal(err)
	}

	return FileServerWith404(http.FS(htmlContent), func(w http.ResponseWriter, r *http.Request) bool {
		r.URL.Path = "/__spa-fallback.html"
		return true
	})
}
