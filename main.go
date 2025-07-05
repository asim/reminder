package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/asim/reminder/api"
	"github.com/asim/reminder/app"
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"github.com/asim/reminder/search"
	hijri "github.com/hablullah/go-hijri"
)

var (
	IndexFlag  = flag.Bool("index", false, "Index data for search. Stored at $HOME/reminder.idx")
	ExportFlag = flag.Bool("export", false, "Export the index data to $HOME/reminder.idx.gob.gz")
	ImportFlag = flag.Bool("import", false, "Import the index data from $HOME/reminder.idx.gob.gz")
	ServerFlag = flag.Bool("serve", false, "Run the server")
	EnvFlag    = flag.String("env", "dev", "Set the environment")
	WebFlag    = flag.Bool("web", false, "Without this flag, the lite version will be served")
)

var mtx sync.RWMutex
var history = map[string][]string{}
var dailyName, dailyVerse, dailyHadith string
var links = map[string]string{}
var dailyUpdated = time.Time{}
var reminderDir = api.ReminderDir
var lastPushDateFile = api.ReminderPath("last_push_date.txt")
var lastPushDate = loadLastPushDate()
var dailyIndex = loadDailyIndex()

func registerLiteRoutes(q *quran.Quran, n *names.Names, b *hadith.Volumes, a *api.Api) {
	// generate api doc
	apiHtml := app.RenderTemplate("API", "", a.Markdown())

	appHtml := app.RenderHTML("App", "The reminder web app", app.Index)

	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(appHtml))
	})

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(apiHtml))
	})

	http.HandleFunc("/daily", func(w http.ResponseWriter, r *http.Request) {
		template := `
<h3>Verse</h3>
<a href="%s" class="block">%s</a>
<h3>Hadith</h3>
<a href="%s" class="block">%s</a>
<h3>Name</h3>
<a href="%s" class="block">%s</a>
<p>Updated %s</p>
`
		mtx.RLock()
		verseLink := links["verse"]
		hadithLink := links["hadith"]
		nameLink := links["name"]

		data := fmt.Sprintf(template, verseLink, dailyVerse, hadithLink, dailyHadith, nameLink, dailyName, dailyUpdated.Format(time.RFC850))
		mtx.RUnlock()
		html := app.RenderHTML("Daily Reminder", "Daily reminder from the quran, hadith and names of Allah", data)
		w.Write([]byte(html))
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

	http.HandleFunc("/islam", func(w http.ResponseWriter, r *http.Request) {
		html := app.Get("islam.html")
		page := app.RenderHTML("Islam", "An overview of Islam", html)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(page))
	})
}

func loadLastPushDate() string {
	b, err := os.ReadFile(lastPushDateFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func saveLastPushDate(date string) {
	_ = os.MkdirAll(reminderDir, 0700)
	_ = os.WriteFile(lastPushDateFile, []byte(date), 0644)
}

// On startup, load daily index
func loadDailyIndex() map[string]interface{} {
	dailyFile := api.ReminderPath("daily.json")
	var idx map[string]interface{}
	if b, err := os.ReadFile(dailyFile); err == nil {
		json.Unmarshal(b, &idx)
	}
	if idx == nil {
		idx = make(map[string]interface{})
	}
	return idx
}

func main() {
	flag.Parse()

	// Load push subscriptions
	_ = api.LoadPushSubscriptions()

	// Load or generate VAPID keys
	_ = api.LoadOrGenerateVAPIDKeys()

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
	a := api.Load()

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

	if *WebFlag {
		http.Handle("/", app.ServeWeb())
	} else {
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Strip trailing slash globally (except for "/") for lite app only
			if r.URL.Path != "/" && len(r.URL.Path) > 1 && strings.HasSuffix(r.URL.Path, "/") {
				newReq := new(http.Request)
				*newReq = *r
				urlCopy := *r.URL
				urlCopy.Path = strings.TrimSuffix(r.URL.Path, "/")
				newReq.URL = &urlCopy
				app.ServeLite().ServeHTTP(w, newReq)
				return
			}
			app.ServeLite().ServeHTTP(w, r)
		}))
		registerLiteRoutes(q, n, b, a)
	}

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

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	isCapital := func(s string) bool {
		if len(s) == 0 {
			return false // An empty string doesn't have a capitalized first letter
		}

		firstChar := []rune(s)[0] // Convert string to slice of runes to handle Unicode characters

		// Check if the first character is a letter
		if !unicode.IsLetter(firstChar) {
			return false // If it's not a letter, it can't be capitalized
		}

		return unicode.IsUpper(firstChar)
	}

	daily := func() {
		for {
			mtx.Lock()

			nam := (*n)[rnd.Int()%len((*n))]
			book := b.Books[rnd.Int()%len(b.Books)]
			chap := q.Chapters[rnd.Int()%len(q.Chapters)]
			ver := chap.Verses[rnd.Int()%len(chap.Verses)]
			had := book.Hadiths[rnd.Int()%len(book.Hadiths)]

			// make sure we're starting from the begining of a verse
			if !isCapital(ver.Text) {
				mtx.Unlock()
				continue
			}

			// skip zero verse e.g bismillah
			if ver.Number == 0 {
				mtx.Unlock()
				continue
			}

			dailyName = fmt.Sprintf("%s - %s - %s\n\n%s", nam.English, nam.Arabic, nam.Meaning, nam.Summary)
			dailyVerse = fmt.Sprintf("%s\n\n%s - %s - %d:%d", ver.Text, chap.Name, chap.English, ver.Chapter, ver.Number)
			dailyHadith = fmt.Sprintf("%s\n\n%s - %s - %s", had.Text, book.Name, had.By, strings.Split(had.Info, ":")[0])

			num := strings.TrimSpace(strings.Split(strings.Split(had.Info, "Number")[1], ":")[0])

			links = map[string]string{
				"verse":  fmt.Sprintf("/quran/%d#%d", ver.Chapter, ver.Number),
				"hadith": fmt.Sprintf("/hadith/%d#%s", book.Number, num),
				"name":   fmt.Sprintf("/names/%d", nam.Number),
			}

			dailyUpdated = time.Now()
			mtx.Unlock()

			// Only send push notification if not already sent today
			today := time.Now().Format("2006-01-02")

			if lastPushDate != today {
				// Compose a user-friendly notification message
				notifyVerse := dailyVerse

				if len(dailyVerse) > 250 {
					notifyVerse = notifyVerse[:250] + "..."
				}
				hijriDate := HijriDate()
				message := "Salam, today is the " + hijriDate
				dailyData := map[string]interface{}{
					"verse":   dailyVerse,
					"hadith":  dailyHadith,
					"name":    dailyName,
					"hijri":   hijriDate,
					"date":    today,
					"links":   links,
					"updated": dailyUpdated.Format(time.RFC850),
					"message": message,
				}

				// Save to daily.json
				dailyFile := api.ReminderPath("daily.json")
				var allDaily map[string]interface{}
				if b, err := os.ReadFile(dailyFile); err == nil {
					json.Unmarshal(b, &allDaily)
				}
				if allDaily == nil {
					allDaily = make(map[string]interface{})
				}
				allDaily[today] = dailyData
				b, _ := json.MarshalIndent(allDaily, "", "  ")
				_ = os.MkdirAll(api.ReminderDir, 0700)
				_ = os.WriteFile(dailyFile, b, 0644)
				dailyIndex = allDaily // update in-memory index

				payload := map[string]interface{}{
					"title": "Daily Reminder",
					"body":  notifyVerse,
				}
				b, _ = json.Marshal(payload)
				errors := api.SendPushToAll(string(b))
				if len(errors) > 0 {
					fmt.Println("Push notification errors:")
					for _, err := range errors {
						fmt.Println(err)
					}
				}
				lastPushDate = today
				saveLastPushDate(today)
			}

			// check every hour
			time.Sleep(time.Hour)
		}
	}

	go daily()

	http.HandleFunc("/api/daily", func(w http.ResponseWriter, r *http.Request) {
		// GET: default current day
		mtx.RLock()
		display := HijriDate()
		message := "Salam, today is the " + display
		resp := map[string]interface{}{
			"name":    dailyName,
			"hadith":  dailyHadith,
			"verse":   dailyVerse,
			"links":   links,
			"updated": dailyUpdated.Format(time.RFC850),
			"message": message,
		}
		mtx.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/api/daily/{date}", func(w http.ResponseWriter, r *http.Request) {
		date := r.PathValue("date")
		if len(date) == 0 {
			date = time.Now().Format("2006-01-02")
		}

		mtx.RLock()

		var resp interface{}
		if entry, ok := dailyIndex[date]; ok {
			resp = entry
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Not found"}`))
			mtx.RUnlock()
			return
		}

		mtx.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/api/daily/reset", func(w http.ResponseWriter, r *http.Request) {
		nam := (*n)[rnd.Int()%len((*n))]
		book := b.Books[rnd.Int()%len(b.Books)]
		chap := q.Chapters[rnd.Int()%len(q.Chapters)]
		ver := chap.Verses[rnd.Int()%len(chap.Verses)]
		had := book.Hadiths[rnd.Int()%len(book.Hadiths)]

		mtx.Lock()

		dailyName = fmt.Sprintf("%s - %s - %s\n\n%s", nam.English, nam.Arabic, nam.Meaning, nam.Summary)
		dailyVerse = fmt.Sprintf("%s\n\n%s - %s -  %d:%d", ver.Text, chap.Name, chap.English, ver.Chapter, ver.Number)
		dailyHadith = fmt.Sprintf("%s\n\n%s - %s - %s", had.Text, book.Name, had.By, strings.Split(had.Info, ":")[0])

		num := strings.TrimSpace(strings.Split(strings.Split(had.Info, "Number")[1], ":")[0])

		links = map[string]string{
			"verse":  fmt.Sprintf("/quran/%d#%d", ver.Chapter, ver.Number),
			"hadith": fmt.Sprintf("/hadith/%d#%s", book.Number, num),
			"name":   fmt.Sprintf("/names/%d", nam.Number),
		}

		dailyUpdated = time.Now()
		display := HijriDate()
		message := "Salam, today is the " + display
		day := map[string]interface{}{
			"name":    dailyName,
			"hadith":  dailyHadith,
			"verse":   dailyVerse,
			"links":   links,
			"updated": dailyUpdated.Format(time.RFC850),
			"message": message,
		}
		mtx.Unlock()

		b, _ := json.Marshal(day)
		w.Write(b)
	})

	// Add daily index API endpoint
	http.HandleFunc("/api/daily/index", func(w http.ResponseWriter, r *http.Request) {
		mtx.RLock()
		b, _ := json.MarshalIndent(dailyIndex, "", "  ")
		mtx.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
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

	api.RegisterRoutes(http.DefaultServeMux)
	httpMux := http.DefaultServeMux
	api.RegisterPushRoutes(httpMux)

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

func HijriDate() string {
	now := time.Now()
	h, err := hijri.CreateUmmAlQuraDate(now)
	if err != nil {
		return ""
	}
	ordinal := func(n int64) string {
		if n == 1 {
			return "st"
		} else if n == 2 {
			return "nd"
		} else if n == 3 {
			return "rd"
		} else if n%10 == 1 && n%100 != 11 {
			return "st"
		} else if n%10 == 2 && n%100 != 12 {
			return "nd"
		} else if n%10 == 3 && n%100 != 13 {
			return "rd"
		}
		return "th"
	}
	months := []string{"Muharram", "Safar", "Rabiʿ al-awwal", "Rabiʿ al-thani", "Jumada al-awwal", "Jumada al-thani", "Rajab", "Shaʿban", "Ramadan", "Shawwal", "Dhu al-Qiʿdah", "Dhu al-Hijjah"}
	display := fmt.Sprintf("%d", h.Day) +
		ordinal(h.Day) +
		" of " + months[int(h.Month)-1] + ", " + fmt.Sprintf("%d", h.Year)
	return display
}
