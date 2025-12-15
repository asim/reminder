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
    <script src="https://unpkg.com/htmx.org@2.0.4"></script>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
    html, body { height: 100%%; width: 100%%; margin: 0; padding: 0; font-family: "Nunito Sans", serif; }
    .htmx-request #container { opacity: 0.5; }
    #container { transition: opacity 200ms ease-in-out; }
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
.bookmark-btn {
  background: none;
  border: none;
  font-size: 1.2em;
  cursor: pointer;
  padding: 0 5px;
  color: #ffd700;
  transition: transform 0.2s;
}
.bookmark-btn:hover {
  transform: scale(1.2);
}
.bookmark-btn.bookmarked {
  color: #ffa500;
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
  <body class="bg-gray-50">
    <div class="w-full text-sm py-2 px-2 bg-black text-white flex flex-row items-center gap-1 flex-shrink-0">
      <div class="inline-block order-1">
        <a href="/" hx-get="/home" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="border border-gray-600 rounded px-2 py-1 hover:border-gray-400 transition-colors">&nbsp;R&nbsp;</a>
      </div>
      <div class="flex-1 flex flex-row justify-center gap-1 order-2">
        <a href="/home" hx-get="/home" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="bg-white hover:opacity-100 flex items-center gap-1 text-black px-2 py-0.5 rounded-md transition-colors">Home</a>
        <a href="/daily" hx-get="/daily" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="bg-white hover:opacity-100 flex items-center gap-1 text-black px-2 py-0.5 rounded-md transition-colors">Daily</a>
        <a href="/quran" hx-get="/quran" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="bg-white hover:opacity-100 flex items-center gap-1 text-black px-2 py-0.5 rounded-md transition-colors">Quran</a>
        <a href="/hadith" hx-get="/hadith" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="bg-white hover:opacity-100 flex items-center gap-1 text-black px-2 py-0.5 rounded-md transition-colors">Hadith</a>
        <a href="/names" hx-get="/names" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="bg-white hover:opacity-100 flex items-center gap-1 text-black px-2 py-0.5 rounded-md transition-colors">Names</a>
      </div>
      <div class="hidden lg:flex items-center gap-2 order-3">
        <a href="/bookmarks" class="bg-white hover:opacity-100 flex items-center gap-1 text-black px-2 py-0.5 rounded-md transition-colors">Bookmarks</a>
        <a href="/search" hx-get="/search" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="bg-white hover:opacity-100 flex items-center gap-1 text-black px-2 py-0.5 rounded-md transition-colors">Search</a>
      </div>
    </div>
    <div id="main" class="h-screen overflow-y-auto">
      <div class="max-w-4xl mx-auto p-6">
        <div class="mb-4">
          <h1 class="text-2xl font-bold">%s</h1>
          <p class="text-gray-600">%s</p>
        </div>
        <div class="prose max-w-none">%s</div>
      </div>
    </div>

  <script src="/bookmarks.js"></script>
  <script>
      if (navigator.serviceWorker) {
        navigator.serviceWorker.register (
          '/reminder.js',
          {scope: '/'}
        )
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
<div class="text-center mb-8">
  <h1 class="text-2xl sm:text-3xl md:text-4xl font-bold mb-2 sm:mb-3">Latest Reminder</h1>
  <p class="text-sm sm:text-base text-gray-600 mb-1">Updated hourly with a new verse, hadith, and name</p>
</div>

<section class="mb-8">
  <h2 class="text-lg font-semibold mb-2">Verse</h2>
  <div class="text-sm sm:text-base text-gray-700 mb-2">A verse from the Quran</div>
  <div class="whitespace-pre-wrap leading-snug bg-blue-50 rounded p-4 text-base shadow">
    <a href="{verse_link}" hx-get="{verse_link}" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="hover:underline">{verse_text}</a>
  </div>
</section>

<section class="mb-8">
  <h2 class="text-lg font-semibold mb-2">Hadith</h2>
  <div class="text-sm sm:text-base text-gray-700 mb-2">A hadith from sahih bukhari</div>
  <div class="whitespace-pre-wrap leading-snug bg-green-50 rounded p-4 text-base shadow">
    <a href="{hadith_link}" hx-get="{hadith_link}" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="hover:underline">{hadith_text}</a>
  </div>
</section>

<section class="mb-8">
  <h2 class="text-lg font-semibold mb-2">Name of Allah</h2>
  <div class="text-sm sm:text-base text-gray-700 mb-2">A beautiful name from the 99 names of Allah</div>
  <div class="whitespace-pre-wrap leading-snug bg-yellow-50 rounded p-4 text-base shadow">
    <a href="{name_link}" hx-get="{name_link}" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="hover:underline">{name_text}</a>
  </div>
</section>

<div class="pt-4 border-t border-gray-200">
  <h3 class="text-base sm:text-lg font-medium mb-2 text-gray-800">Explore More</h3>
  <p class="text-sm sm:text-base text-gray-600 mb-4">Browse our full collection or view past daily reminders</p>
  <div class="flex flex-col sm:flex-row gap-2 sm:gap-3">
    <a href="/daily" hx-get="/daily" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="inline-block px-4 py-2 bg-black text-white rounded-md hover:bg-gray-800 transition-colors text-center text-sm sm:text-base">Daily Reminder</a>
    <a href="/quran" hx-get="/quran" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="inline-block px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 transition-colors text-center text-sm sm:text-base">Read Quran</a>
    <a href="/hadith" hx-get="/hadith" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="inline-block px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-100 transition-colors text-center text-sm sm:text-base">Read Hadith</a>
  </div>
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

func RenderContent(title, desc, html string) string {
	return fmt.Sprintf(`<div class="max-w-4xl mx-auto p-6">
  <div class="mb-4">
    <h1 class="text-2xl font-bold">%s</h1>
    <p class="text-gray-600">%s</p>
  </div>
  <div class="prose max-w-none">%s</div>
</div>`, title, desc, html)
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
