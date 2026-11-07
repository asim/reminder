package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/asim/reminder/api"
	"github.com/asim/reminder/app"
	"github.com/asim/reminder/daily"
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
	"github.com/google/uuid"
)

var (
	IndexFlag  = flag.Bool("index", false, "Index data for search. Stored at $HOME/reminder.idx")
	ExportFlag = flag.Bool("export", false, "Export the index data to $HOME/reminder.idx.gob.gz")
	ImportFlag = flag.Bool("import", false, "Import the index data from $HOME/reminder.idx.gob.gz")
	ServerFlag = flag.Bool("serve", false, "Run the server")
	EnvFlag    = flag.String("env", "dev", "Set the environment")
	WebFlag    = flag.Bool("web", false, "Without this flag, the lite version will be served")
)

func isAPIClient(r *http.Request) bool {
	userAgent := r.Header.Get("User-Agent")
	accept := r.Header.Get("Accept")

	return strings.Contains(userAgent, "Wget") ||
		strings.Contains(userAgent, "curl") ||
		strings.Contains(userAgent, "Go-http-client") ||
		strings.Contains(userAgent, "python-requests") ||
		(accept != "" && !strings.Contains(accept, "text/html") && !strings.Contains(accept, "*/*"))
}

func registerLiteRoutes(q *quran.Quran, n *names.Names, b *hadith.Collection, a *api.Api, reminders *daily.Store) {
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		indexContent := app.HomeContent(reminders.Latest())

		if isAPIClient(r) {
			w.Write([]byte(app.RenderSimpleHTML("Home", "Quran, hadith, and more as an app and API", indexContent)))
		} else {
			w.Write([]byte(app.RenderHTML("Home", "Quran, hadith, and more as an app and API", indexContent)))
		}
	})

	http.HandleFunc("/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(app.RenderHTML("Bookmarks", "", app.Bookmarks)))
	})

	http.HandleFunc("/quran", func(w http.ResponseWriter, r *http.Request) {
		content := q.TOC()
		if isAPIClient(r) {
			qhtml := app.RenderSimpleHTML("Quran", quran.Description, content)
			w.Write([]byte(qhtml))
		} else {
			qhtml := app.RenderHTML("Quran", quran.Description, content)
			w.Write([]byte(qhtml))
		}
	})

	http.HandleFunc("/hadith", func(w http.ResponseWriter, r *http.Request) {
		content := b.TOC()
		if isAPIClient(r) {
			qhtml := app.RenderSimpleHTML("Hadith", hadith.Description, content)
			w.Write([]byte(qhtml))
		} else {
			qhtml := app.RenderHTML("Hadith", hadith.Description, content)
			w.Write([]byte(qhtml))
		}
	})

	http.HandleFunc("/names", func(w http.ResponseWriter, r *http.Request) {
		content := n.TOC()
		if isAPIClient(r) {
			qhtml := app.RenderSimpleHTML("Names", names.Description, content)
			w.Write([]byte(qhtml))
		} else {
			qhtml := app.RenderHTML("Names", names.Description, content)
			w.Write([]byte(qhtml))
		}
	})

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		content := a.HTML(app.RenderString)
		w.Write([]byte(app.RenderHTML("API", "", content)))
	})

	http.HandleFunc("/daily", func(w http.ResponseWriter, r *http.Request) {
		content := app.DailyContent(reminders.Latest())
		w.Write([]byte(app.RenderHTML("Daily Reminder", "Daily reminder from the quran, hadith and names of Allah", content)))
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
		content := q.Get(ch).HTML()

		qhtml := app.RenderHTML(head, "", content)
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

		if cc.Verse(ve) == nil {
			return
		}

		vv := cc.Verse(ve)

		head := fmt.Sprintf("%d:%d | Quran", ch, ve)
		vhtml := app.RenderHTML(head, "", vv.HTML())

		w.Write([]byte(vhtml))
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
		content := n.Get(name).HTML()

		qhtml := app.RenderHTML(head, "", content)
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
		content := b.Get(ch).HTML()

		qhtml := app.RenderHTML(head, "", content)
		w.Write([]byte(qhtml))
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		shtml := app.RenderHTML("Search", "", app.SearchTemplate)
		w.Write([]byte(shtml))
	})

	http.HandleFunc("/islam", func(w http.ResponseWriter, r *http.Request) {
		html := app.Get("islam.html")
		page := app.RenderHTML("Islam", "An overview of Islam", html)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page))
	})

	// Serve static files (bookmarks.js, manifest, etc.)
	staticFiles := []string{"bookmarks.js", "reminder.js", "manifest.webmanifest", "arabic.otf", "reminder.png", "icon-192.png", "icon-96.png"}
	for _, file := range staticFiles {
		fileName := file
		http.HandleFunc("/"+fileName, func(w http.ResponseWriter, r *http.Request) {
			content := app.Get(fileName)
			if strings.HasSuffix(fileName, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			} else if strings.HasSuffix(fileName, ".json") || strings.HasSuffix(fileName, ".webmanifest") {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			} else if strings.HasSuffix(fileName, ".otf") {
				w.Header().Set("Content-Type", "font/otf")
				w.Header().Set("Cache-Control", "public, max-age=31536000")
			} else if strings.HasSuffix(fileName, ".png") {
				w.Header().Set("Content-Type", "image/png")
				w.Header().Set("Cache-Control", "public, max-age=31536000")
			}
			w.Write([]byte(content))
		})
	}

	// Root route - redirect to /home
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Only handle exact "/" path, let other routes handle themselves
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/home", http.StatusFound)
			return
		}
		// For other paths, return 404 (no other handler matched)
		http.NotFound(w, r)
	})
}

func main() {
	log.Println("New rand source")
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	flag.Parse()

	reminders := daily.NewStore(api.ReminderDir)
	var history search.History

	// Load push subscriptions
	log.Println("Loading subscriptions")
	_ = api.LoadPushSubscriptions()

	// Load or generate VAPID keys
	log.Println("Loading VAPID keys")
	_ = api.LoadOrGenerateVAPIDKeys()

	// Search uses the prebuilt vector index when an OpenAI key is available,
	// since embeddings are generated by OpenAI. Without a key it falls back to
	// a local SQLite keyword index, so search still works with no external
	// service and no API spend.
	var idx search.Searcher

	if search.Semantic() {
		log.Println("Generating index")
		vidx := search.New("reminder", false)
		idx = vidx

		// async load the index
		go func() {
			// Load the pre-existing data
			log.Println("Loading index")
			if err := vidx.Load(); err != nil {
				log.Println(err)
			}
			log.Println("Loaded index")
		}()
	} else {
		log.Println("No OpenAI key set, using keyword search")
	}

	// load data
	log.Println("Initialising data")
	q := quran.Load()
	log.Println("Loaded Quran")
	n := names.Load()
	log.Println("Loaded Names")
	b := hadith.Load()
	log.Println("Loaded Hadith")
	a := api.Load()
	log.Println("Loaded API")

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

		// Rebuilding needs an OpenAI key, since that is what generates the
		// embeddings. Say so rather than failing obscurely per document.
		if !search.Semantic() {
			log.Println("Warning: --index needs OPENAI_API_KEY to generate embeddings")
		}
		if idx == nil {
			idx = sidx
		}

		log.Println("Indexing data")
		go func() {
			search.IndexQuran(sidx, q)
			search.IndexNames(sidx, n)
			search.IndexHadith(sidx, b)
			search.IndexTafsir(sidx, q)
			// done
			close(indexed)
		}()
	} else if !search.Semantic() {
		// No OpenAI key, so build the keyword index from the embedded data.
		// This is local and takes seconds, unlike generating embeddings.
		home, _ := os.UserHomeDir()
		reminderHome := filepath.Join(home, ".reminder")
		os.MkdirAll(reminderHome, 0755)

		fts := search.NewFTS(filepath.Join(reminderHome, "search.db"))
		idx = fts
		if changed, err := search.SyncNames(fts, n); err != nil {
			log.Fatalf("Could not refresh names in keyword index: %v", err)
		} else if changed {
			log.Println("Refreshed names in keyword index")
		}

		if fts.Built() {
			log.Printf("Keyword index loaded (%d documents)", fts.Count())
			close(indexed)
		} else {
			log.Println("Building keyword index")
			go func() {
				search.IndexQuran(fts, q)
				search.IndexNames(fts, n)
				search.IndexHadith(fts, b)
				search.IndexTafsir(fts, q)

				// Only now is the corpus complete. If the process dies before
				// this the index is rebuilt next start, rather than serving a
				// corpus that is silently missing whole sections.
				if err := fts.MarkBuilt(); err != nil {
					log.Printf("Warning: could not mark index complete: %v", err)
				}
				log.Printf("Keyword index built (%d documents)", fts.Count())
				close(indexed)
			}()
		}
	} else {
		close(indexed)
	}

	// Export and import operate on the vector index only. They are maintainer
	// operations for rebuilding the index that ships in the binary.
	if *ExportFlag {
		log.Println("Exporting index")
		if err := search.New("reminder", true).Export(); err != nil {
			log.Println(err)
		}
		return
	}

	if *ImportFlag {
		log.Println("Importing index")
		if err := search.New("reminder", false).Import(); err != nil {
			log.Println(err)
		}
	}

	if *WebFlag {
		log.Println("Registering web handler")

		// Register TOC handlers for API clients only
		http.HandleFunc("/quran", func(w http.ResponseWriter, r *http.Request) {
			if isAPIClient(r) {
				content := q.TOC()
				qhtml := app.RenderSimpleHTML("Quran", quran.Description, content)
				w.Write([]byte(qhtml))
			} else {
				// Let SPA handle it
				app.ServeWeb().ServeHTTP(w, r)
			}
		})

		http.HandleFunc("/hadith", func(w http.ResponseWriter, r *http.Request) {
			if isAPIClient(r) {
				content := b.TOC()
				hhtml := app.RenderSimpleHTML("Hadith", hadith.Description, content)
				w.Write([]byte(hhtml))
			} else {
				// Let SPA handle it
				app.ServeWeb().ServeHTTP(w, r)
			}
		})

		http.HandleFunc("/names", func(w http.ResponseWriter, r *http.Request) {
			if isAPIClient(r) {
				content := n.TOC()
				nhtml := app.RenderSimpleHTML("Names", names.Description, content)
				w.Write([]byte(nhtml))
			} else {
				// Let SPA handle it
				app.ServeWeb().ServeHTTP(w, r)
			}
		})

		// Create smart handlers that serve lite content for API clients, SPA for browsers
		http.HandleFunc("/quran/{id}", func(w http.ResponseWriter, r *http.Request) {
			if isAPIClient(r) {
				// Serve server-rendered HTML
				id := r.PathValue("id")
				if len(id) == 0 {
					return
				}
				ch, _ := strconv.Atoi(id)
				if ch < 1 || ch > 114 {
					return
				}
				head := fmt.Sprintf("%d | Quran", ch)
				qhtml := app.RenderSimpleHTML(head, "", q.Get(ch).HTML())
				w.Write([]byte(qhtml))
			} else {
				// Serve SPA
				app.ServeWeb().ServeHTTP(w, r)
			}
		})

		http.HandleFunc("/hadith/{book}", func(w http.ResponseWriter, r *http.Request) {
			if isAPIClient(r) {
				book := r.PathValue("book")
				if len(book) == 0 {
					return
				}
				ch, _ := strconv.Atoi(book)
				if ch < 1 || ch > len(b.Books) {
					return
				}
				head := fmt.Sprintf("%d | Hadith", ch)
				qhtml := app.RenderSimpleHTML(head, "", b.Get(ch).HTML())
				w.Write([]byte(qhtml))
			} else {
				app.ServeWeb().ServeHTTP(w, r)
			}
		})

		http.HandleFunc("/names/{id}", func(w http.ResponseWriter, r *http.Request) {
			if isAPIClient(r) {
				id := r.PathValue("id")
				if len(id) == 0 {
					return
				}
				name, _ := strconv.Atoi(id)
				if name < 1 || name > len(*n) {
					return
				}
				head := fmt.Sprintf("%d | Names", name)
				qhtml := app.RenderSimpleHTML(head, "", n.Get(name).HTML())
				w.Write([]byte(qhtml))
			} else {
				// Serve SPA
				app.ServeWeb().ServeHTTP(w, r)
			}
		})

		// Everything else goes to SPA
		http.Handle("/", app.ServeWeb())
	} else {
		log.Println("Registering lite handler")
		registerLiteRoutes(q, n, b, a, reminders)
	}

	http.HandleFunc("POST /api/quran/{chapter}/{verse}/explain", func(w http.ResponseWriter, r *http.Request) {
		chapter, _ := strconv.Atoi(r.PathValue("chapter"))
		number, _ := strconv.Atoi(r.PathValue("verse"))
		if chapter < 1 || chapter > len(q.Chapters) {
			http.Error(w, "Chapter not found", http.StatusNotFound)
			return
		}
		verse := q.Get(chapter).Verse(number)
		if verse == nil {
			http.Error(w, "Verse not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Commentary remains available without a configured model.
		answer, generated := verse.Comments, false
		if model.Available() {
			ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
			defer cancel()
			context := []string{fmt.Sprintf("Quran %d:%d\nArabic: %s\nTranslation: %s\nStored commentary: %s", chapter, number, verse.Arabic, verse.Text, verse.Comments)}
			question := fmt.Sprintf("Explain Quran %d:%d in plain English using only the supplied verse, translation and commentary. Cite the verse number. Distinguish the translation from interpretation. Do not invent historical circumstances, scholarly attributions or rulings. Do not generate or rewrite Arabic. Use plain text, at most three short paragraphs.", chapter, number)
			if result, err := model.AskSafe(ctx, context, question); err == nil && result != "" {
				answer, generated = result, true
			}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"text": answer, "generated": generated, "reference": fmt.Sprintf("%d:%d", chapter, number)})
	})

	http.HandleFunc("/api/quran", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(qjson))
	})

	http.HandleFunc("/api/quran/{chapter}", func(w http.ResponseWriter, r *http.Request) {
		ch := r.PathValue("chapter")
		if len(ch) == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
			return
		}

		chapter, _ := strconv.Atoi(ch)
		if chapter < 1 || chapter > 114 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
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
		chapter, err := strconv.Atoi(r.PathValue("chapter"))
		if err != nil || chapter < 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
			return
		}

		verse, err := strconv.Atoi(r.PathValue("verse"))
		if err != nil || verse < 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
			return
		}

		var found *quran.Verse
		if chapter <= len(q.Chapters) {
			found = q.Get(chapter).Verse(verse)
		}
		if found == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(found.JSON())
	})

	http.HandleFunc("/api/daily", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reminders.Today())
	})

	// Add /api/latest endpoint for the hourly updated reminder
	http.HandleFunc("/api/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reminders.Latest())
	})

	http.HandleFunc("/api/daily/{date}", func(w http.ResponseWriter, r *http.Request) {
		date := r.PathValue("date")
		if len(date) == 0 {
			date = time.Now().Format("2006-01-02")
		}

		entry, ok := reminders.Get(date)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Not found"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entry)
	})

	// Add daily index API endpoint
	http.HandleFunc("/api/daily/index", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(reminders.IndexJSON())
	})

	// RSS feed for daily content
	http.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		reminders.WriteRSS(w)
	})

	http.HandleFunc("/api/names", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(njson))
	})

	http.HandleFunc("/api/names/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if len(id) == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
			return
		}

		name, _ := strconv.Atoi(id)
		if name < 1 || name > len(*n) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{}"))
			return
		}

		b := n.Get(name)
		json, _ := json.Marshal(b)
		w.Write(json)
	})

	http.HandleFunc("/api/hadith/books", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(b.Index().Books)
		w.Write(b)
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

	http.HandleFunc("/api/explain", func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		var data map[string]interface{}
		json.Unmarshal(b, &data)

		q := data["q"].(string)

		prompt := `Generate a detailed summary for the following with it's meaning and origin, output the response as JSON with the fields: 
		text, context, summary. Each field itself should be a string. Text is the text that was provided. Context is contextual information 
		for the origin of the text or background information required. Summary is a more detailed summary in relation to the text and context.

		%s
		`

		answer := model.Ask(r.Context(), nil, fmt.Sprintf(prompt, q))
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

		answer := model.Ask(r.Context(), nil, fmt.Sprintf(prompt, q))
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
			h := history.Get(ctx)

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

			q, _ := data["q"].(string)

			// summarise defaults to true when absent for backward compat.
			// The new web UI explicitly sets summarise=false for the fast
			// path that only returns references, then issues a second call
			// with summarise=true to fetch the LLM answer.
			summarise := true
			if v, ok := data["summarise"]; ok {
				if bv, ok := v.(bool); ok {
					summarise = bv
				}
			}

			res, err := search.References(idx, n, q)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			// Report whether summarisation is possible, so clients can hide
			// the control instead of offering one that cannot work.
			available := model.Available()

			// References alone are useful, and are all we can return when the
			// caller opted out or no model is configured.
			writeRefs := func() {
				output, _ := json.Marshal(map[string]interface{}{
					"q":                   q,
					"answer":              "",
					"references":          res,
					"summarise_available": available,
				})
				w.Write(output)
			}

			if !summarise || !available {
				writeRefs()
				return
			}

			answer, err := search.Summarise(r.Context(), res, q)
			if err != nil {
				log.Printf("Summarisation failed for %q: %v", q, err)
				writeRefs()
				return
			}

			answerMD := string(app.Render([]byte(answer)))

			output, _ := json.Marshal(map[string]interface{}{
				"q":                   q,
				"answer":              answerMD,
				"references":          res,
				"summarise_available": available,
			})
			w.Write(output)

			var ctx string

			// look for the context cookie
			c, err := r.Cookie("session")
			if err == nil {
				ctx = c.Value
				history.Add(ctx, q, answerMD)
			}

			return
		}
	})

	log.Println("Registering routes")
	httpMux := http.DefaultServeMux
	api.RegisterRoutes(httpMux)

	// Register MCP server
	mcpServer := api.NewServer("reminder", "1.0.0")

	mcpServer.AddTool("get_latest", "Get the latest reminder (updated hourly) with a verse, hadith and name of Allah", api.InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		resp := reminders.Latest()
		b, _ := json.Marshal(resp)
		return string(b), nil
	})

	mcpServer.AddTool("get_daily", "Get today's daily reminder with a verse, hadith and name of Allah", api.InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		resp := reminders.Today()
		b, _ := json.Marshal(resp)
		return string(b), nil
	})

	mcpServer.AddTool("get_daily_by_date", "Get a daily reminder for a specific date", api.InputSchema{
		Type: "object",
		Properties: map[string]api.Property{
			"date": {Type: "string", Description: "Date in YYYY-MM-DD format"},
		},
		Required: []string{"date"},
	}, func(args map[string]interface{}) (string, error) {
		date, _ := args["date"].(string)
		if date == "" {
			return "", fmt.Errorf("date is required")
		}
		if entry, ok := reminders.Get(date); ok {
			b, _ := json.Marshal(entry)
			return string(b), nil
		}
		return "", fmt.Errorf("not found")
	})

	mcpServer.AddTool("get_quran_chapters", "Get a list of all Quran chapters", api.InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		b, _ := json.Marshal(q.Index().Chapters)
		return string(b), nil
	})

	mcpServer.AddTool("get_quran_chapter", "Get a chapter of the Quran with all its verses", api.InputSchema{
		Type: "object",
		Properties: map[string]api.Property{
			"chapter": {Type: "number", Description: "Chapter number (1-114)"},
		},
		Required: []string{"chapter"},
	}, func(args map[string]interface{}) (string, error) {
		chNum, ok := args["chapter"].(float64)
		if !ok {
			return "", fmt.Errorf("chapter is required")
		}
		chapter := int(chNum)
		if chapter < 1 || chapter > 114 {
			return "", fmt.Errorf("chapter must be between 1 and 114")
		}
		return string(q.Get(chapter).JSON()), nil
	})

	mcpServer.AddTool("get_quran_verse", "Get a specific verse of the Quran with word-by-word translation", api.InputSchema{
		Type: "object",
		Properties: map[string]api.Property{
			"chapter": {Type: "number", Description: "Chapter number (1-114)"},
			"verse":   {Type: "number", Description: "Verse number"},
		},
		Required: []string{"chapter", "verse"},
	}, func(args map[string]interface{}) (string, error) {
		chNum, ok := args["chapter"].(float64)
		if !ok {
			return "", fmt.Errorf("chapter is required")
		}
		veNum, ok := args["verse"].(float64)
		if !ok {
			return "", fmt.Errorf("verse is required")
		}
		chapter := int(chNum)
		verse := int(veNum)
		if chapter < 1 || chapter > 114 {
			return "", fmt.Errorf("chapter must be between 1 and 114")
		}
		cc := q.Get(chapter)
		if cc.Verse(verse) == nil {
			return "", fmt.Errorf("verse out of range")
		}
		return string(cc.Verse(verse).JSON()), nil
	})

	mcpServer.AddTool("get_hadith_books", "Get a list of all hadith books in Sahih Bukhari", api.InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		bks, _ := json.Marshal(b.Index().Books)
		return string(bks), nil
	})

	mcpServer.AddTool("get_hadith_book", "Get a specific book from Sahih Bukhari with all its hadiths", api.InputSchema{
		Type: "object",
		Properties: map[string]api.Property{
			"book": {Type: "number", Description: "Book number"},
		},
		Required: []string{"book"},
	}, func(args map[string]interface{}) (string, error) {
		bkNum, ok := args["book"].(float64)
		if !ok {
			return "", fmt.Errorf("book is required")
		}
		book := int(bkNum)
		if book < 1 || book > len(b.Books) {
			return "", fmt.Errorf("book out of range")
		}
		return string(b.Get(book).JSON()), nil
	})

	mcpServer.AddTool("get_names", "Get all 99 Names of Allah", api.InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		return string(njson), nil
	})

	mcpServer.AddTool("get_name", "Get a specific Name of Allah with description and Quran references", api.InputSchema{
		Type: "object",
		Properties: map[string]api.Property{
			"id": {Type: "number", Description: "Name number (1-99)"},
		},
		Required: []string{"id"},
	}, func(args map[string]interface{}) (string, error) {
		idNum, ok := args["id"].(float64)
		if !ok {
			return "", fmt.Errorf("id is required")
		}
		id := int(idNum)
		if id < 1 || id > len(*n) {
			return "", fmt.Errorf("id must be between 1 and 99")
		}
		byt, _ := json.Marshal(n.Get(id))
		return string(byt), nil
	})

	mcpServer.AddTool("search", "Search Islamic content and get AI-summarised answers from the Quran, Hadith and Names of Allah", api.InputSchema{
		Type: "object",
		Properties: map[string]api.Property{
			"q": {Type: "string", Description: "The question to ask"},
		},
		Required: []string{"q"},
	}, func(args map[string]interface{}) (string, error) {
		select {
		case <-indexed:
		default:
			return "", fmt.Errorf("indexing content, please try again later")
		}

		question, ok := args["q"].(string)
		if !ok || question == "" {
			return "", fmt.Errorf("q is required")
		}

		res, err := search.References(idx, n, question)
		if err != nil {
			return "", err
		}
		answer := model.Ask(context.Background(), search.Context(res), question)

		output, _ := json.Marshal(map[string]interface{}{
			"q":          question,
			"answer":     answer,
			"references": res,
		})
		return string(output), nil
	})

	http.Handle("/mcp", mcpServer)

	if *ServerFlag {
		log.Println("Starting daily")
		go func() {
			for {
				func() {
					defer func() {
						if err := recover(); err != nil {
							log.Printf("Daily loop recovered from panic: %v", err)
						}
					}()
					log.Println("Running daily")
					reminders.Update(context.Background(), q, n, b, rnd, api.SendPushToAll)
				}()
				time.Sleep(time.Hour)
			}
		}()

		log.Println("Starting server :8080")
		if err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS for API routes — allow any origin
			if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/mcp" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}
			} else if *EnvFlag == "dev" {
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("Access-Control-Allow-Credentials", "true")

				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}
			}

			// check cookies
			if c, err := r.Cookie("session"); err != nil || c == nil || len(c.Value) == 0 {
				var secure bool

				if h := r.Header.Get("X-Forwarded-Proto"); h == "https" {
					secure = true
				}

				http.SetCookie(w, &http.Cookie{
					Name:   "session",
					Value:  uuid.New().String(),
					Secure: secure,
				})
			}

			http.DefaultServeMux.ServeHTTP(w, r)
		})); err != nil {
			log.Printf("Server error: %v\n", err)
		}
	}

	// wait for index
	<-indexed
}
