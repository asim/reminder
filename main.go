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
	"sync"
	"time"
	"unicode"

	"github.com/asim/reminder/api"
	"github.com/asim/reminder/app"
	"github.com/asim/reminder/daily"
	"github.com/asim/reminder/hadith"
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

var mtx sync.RWMutex
var history = map[string][]string{}
var dailyName, dailyVerse, dailyHadith, dailyMessage string
var links = map[string]string{}
var dailyUpdated = time.Time{}
var reminderDir = api.ReminderDir
var lastPushDateFile = api.ReminderPath("last_push_date.txt")
var lastPushDate = loadLastPushDate()
var dailyIndex = loadDailyIndex()

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isCapital(s string) bool {
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

func isAPIClient(r *http.Request) bool {
	userAgent := r.Header.Get("User-Agent")
	accept := r.Header.Get("Accept")

	return strings.Contains(userAgent, "Wget") ||
		strings.Contains(userAgent, "curl") ||
		strings.Contains(userAgent, "Go-http-client") ||
		strings.Contains(userAgent, "python-requests") ||
		(accept != "" && !strings.Contains(accept, "text/html") && !strings.Contains(accept, "*/*"))
}

func registerLiteRoutes(q *quran.Quran, n *names.Names, b *hadith.Collection, a *api.Api) {
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		mtx.RLock()
		verseLink := links["verse"]
		hadithLink := links["hadith"]
		nameLink := links["name"]
		verse := dailyVerse
		hadith := dailyHadith
		name := dailyName
		mtx.RUnlock()

		// Populate Index template with actual data
		indexContent := strings.ReplaceAll(app.Index, "{verse_link}", verseLink)
		indexContent = strings.ReplaceAll(indexContent, "{verse_text}", verse)
		indexContent = strings.ReplaceAll(indexContent, "{hadith_link}", hadithLink)
		indexContent = strings.ReplaceAll(indexContent, "{hadith_text}", hadith)
		indexContent = strings.ReplaceAll(indexContent, "{name_link}", nameLink)
		indexContent = strings.ReplaceAll(indexContent, "{name_text}", name)

		if isAPIClient(r) {
			w.Write([]byte(app.RenderSimpleHTML("Home", "Quran, hadith, and more as an app and API", indexContent)))
		} else {
			w.Write([]byte(app.RenderHTML("Home", "Quran, hadith, and more as an app and API", indexContent)))
		}
	})

	http.HandleFunc("/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		bookmarksContent := `
<div id="quran-bookmarks" class="mb-8">
  <h2 class="text-2xl font-semibold mb-4">Quran</h2>
  <div id="quran-list" class="space-y-3"></div>
</div>

<div id="hadith-bookmarks" class="mb-8">
  <h2 class="text-2xl font-semibold mb-4">Hadith</h2>
  <div id="hadith-list" class="space-y-3"></div>
</div>

<div id="names-bookmarks" class="mb-8">
  <h2 class="text-2xl font-semibold mb-4">Names</h2>
  <div id="names-list" class="space-y-3"></div>
</div>

<script>
  function loadBookmarks() {
    const bookmarks = JSON.parse(localStorage.getItem('reminder_bookmarks') || '{"quran":{},"hadith":{},"names":{}}');
    
    // Quran bookmarks
    const quranList = document.getElementById('quran-list');
    const quranKeys = Object.keys(bookmarks.quran || {});
    if (quranKeys.length > 0) {
      quranList.innerHTML = quranKeys.map(key => {
        const bookmark = bookmarks.quran[key];
        return ` + "`" + `
        <div class="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors">
          <a href="${bookmark.url}" class="flex-grow text-blue-600 hover:text-blue-800">${bookmark.label}</a>
          <button onclick="removeBookmark('quran', '${key}')" class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition-colors text-sm">Remove</button>
        </div>
      ` + "`" + `;
      }).join('');
    } else {
      quranList.innerHTML = '<p class="text-gray-500">No bookmarks yet</p>';
    }
    
    // Hadith bookmarks
    const hadithList = document.getElementById('hadith-list');
    const hadithKeys = Object.keys(bookmarks.hadith || {});
    if (hadithKeys.length > 0) {
      hadithList.innerHTML = hadithKeys.map(key => {
        const bookmark = bookmarks.hadith[key];
        return ` + "`" + `
        <div class="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors">
          <a href="${bookmark.url}" class="flex-grow text-blue-600 hover:text-blue-800">${bookmark.label}</a>
          <button onclick="removeBookmark('hadith', '${key}')" class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition-colors text-sm">Remove</button>
        </div>
      ` + "`" + `;
      }).join('');
    } else {
      hadithList.innerHTML = '<p class="text-gray-500">No bookmarks yet</p>';
    }
    
    // Names bookmarks
    const namesList = document.getElementById('names-list');
    const namesKeys = Object.keys(bookmarks.names || {});
    if (namesKeys.length > 0) {
      namesList.innerHTML = namesKeys.map(key => {
        const bookmark = bookmarks.names[key];
        return ` + "`" + `
        <div class="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors">
          <a href="${bookmark.url}" class="flex-grow text-blue-600 hover:text-blue-800">${bookmark.label}</a>
          <button onclick="removeBookmark('names', '${key}')" class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition-colors text-sm">Remove</button>
        </div>
      ` + "`" + `;
      }).join('');
    } else {
      namesList.innerHTML = '<p class="text-gray-500">No bookmarks yet</p>';
    }
  }
  
  function removeBookmark(type, key) {
    const bookmarks = JSON.parse(localStorage.getItem('reminder_bookmarks') || '{"quran":{},"hadith":{},"names":{}}');
    if (bookmarks[type] && bookmarks[type][key]) {
      delete bookmarks[type][key];
      localStorage.setItem('reminder_bookmarks', JSON.stringify(bookmarks));
      loadBookmarks();
    }
  }
  
  loadBookmarks();
</script>
`

		w.Write([]byte(app.RenderHTML("Bookmarks", "", bookmarksContent)))
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
			w.Write([]byte(q…8077 tokens truncated…	return
			}

			var tokens int
			var contexts []string

			for _, r := range res {
				if tokens >= 8000 {
					break
				}

				b, _ := json.Marshal(r)
				tokens += len(b)
				// TODO: maybe just provide text
				contexts = append(contexts, string(b))
			}

			answer, err := askLLMSafe(r.Context(), contexts, q)
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

	log.Println("Registering routes")
	httpMux := http.DefaultServeMux
	api.RegisterRoutes(httpMux)

	// Register MCP server
	mcpServer := api.NewServer("reminder", "1.0.0")

	mcpServer.AddTool("get_latest", "Get the latest reminder (updated hourly) with a verse, hadith and name of Allah", api.InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		mtx.RLock()
		resp := map[string]interface{}{
			"name":    dailyName,
			"hadith":  dailyHadith,
			"verse":   dailyVerse,
			"links":   links,
			"updated": dailyUpdated.Format(time.RFC3339),
			"message": dailyMessage,
		}
		mtx.RUnlock()
		b, _ := json.Marshal(resp)
		return string(b), nil
	})

	mcpServer.AddTool("get_daily", "Get today's daily reminder with a verse, hadith and name of Allah", api.InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		today := time.Now().UTC().Format("2006-01-02")
		mtx.RLock()
		defer mtx.RUnlock()
		if entry, ok := dailyIndex[today]; ok {
			b, _ := json.Marshal(entry)
			return string(b), nil
		}
		message := "In the Name of Allah—the Most Beneficent, Most Merciful"
		resp := map[string]interface{}{
			"name":    dailyName,
			"hadith":  dailyHadith,
			"verse":   dailyVerse,
			"links":   links,
			"updated": dailyUpdated.Format(time.RFC3339),
			"message": message,
			"hijri":   daily.Date().Display,
			"date":    today,
		}
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
		mtx.RLock()
		defer mtx.RUnlock()
		if entry, ok := dailyIndex[date]; ok {
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

		res, err := idx.Query(question)
		if err != nil {
			return "", err
		}

		var tokens int
		var contexts []string
		for _, r := range res {
			if tokens >= 8000 {
				break
			}
			for k, v := range r.Metadata {
				delete(r.Metadata, k)
				r.Metadata[strings.ToLower(k)] = v
			}
			b, _ := json.Marshal(r)
			tokens += len(b)
			contexts = append(contexts, string(b))
		}

		answer := askLLM(context.Background(), contexts, question)

		output, _ := json.Marshal(map[string]interface{}{
			"q":          question,
			"answer":     answer,
			"references": res,
		})
		return string(output), nil
	})

	http.Handle("/mcp", mcpServer)

	getVerse := func(ch *quran.Chapter, ve *quran.Verse) (string, int, int, string) {
		verseText := ve.Text
		verseStart := ve.Number
		verseEnd := ve.Number

		for i := 0; i < 10; i++ {
			// done
			if x := verseText[len(verseText)-1]; x == '.' || x == '!' || x == '"' {
				break
			}

			idx := ve.Number + 1

			if ch.Number == 1 || ch.Number == 9 {
				idx = ve.Number
			}

			// bail out - check bounds
			if idx < 0 || idx >= len(ch.Verses) {
				break
			}

			// increment
			ve = ch.Verses[idx]

			if v := verseText[len(verseText)-1]; v == ',' || string(v) == "—" || unicode.IsLetter(rune(v)) || v == ';' {
				verseText += " " + ve.Text
			} else {
				verseText += "\n\n" + ve.Text
			}
			verseEnd = ve.Number
		}

		verseNumber := fmt.Sprintf("%d", verseStart)

		if verseStart != verseEnd {
			verseNumber += fmt.Sprintf("-%d", verseEnd)
		}

		formatted := fmt.Sprintf("%s - %s - %d:%s\n\n%s", ch.Name, ch.English, ch.Number, verseNumber, verseText)
		return formatted, verseStart, verseEnd, verseText
	}

	daily := func() {
		// Check if we're within the first 5 minutes after midnight UTC
		// This handles server restarts that happen right at midnight
		isWithinGracePeriod := func() bool {
			now := time.Now().UTC()
			midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			sinceMidnight := now.Sub(midnight)
			return sinceMidnight >= 0 && sinceMidnight < 5*time.Minute
		}

		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Daily loop recovered from panic: %v\n", r)
					}
				}()
			log.Println("Running daily")

			mtx.Lock()

			var nam *names.Name
			var book *hadith.Book
			var chap *quran.Chapter
			var ver *quran.Verse
			var had *hadith.Hadith

			// Retry random selection up to 50 times to avoid skipping
			// an entire hour when we pick a bad combination.
			found := false
			for attempt := 0; attempt < 50; attempt++ {
				nam = (*n)[rnd.Int()%len((*n))]
				bookIdx := rnd.Int() % len(b.Books)
				book = b.Books[bookIdx]

				if len(book.Hadiths) == 0 {
					continue
				}

				chap = q.Chapters[rnd.Int()%len(q.Chapters)]

				if len(chap.Verses) == 0 {
					continue
				}

				ver = chap.Verses[rnd.Int()%len(chap.Verses)]
				had = book.Hadiths[rnd.Int()%len(book.Hadiths)]

				if ver.Number == 0 || !isCapital(ver.Text) {
					continue
				}

				found = true
				break
			}

			if !found {
				log.Println("Could not find valid content after 50 attempts, will retry next hour")
				mtx.Unlock()
				return
			}

			dailyName = fmt.Sprintf("%s - %s - %s\n\n%s", nam.English, nam.Arabic, nam.Meaning, nam.Summary)
			verseFormatted, verseStart, verseEnd, verseText := getVerse(chap, ver)
			dailyVerse = verseFormatted
			// Use new hadith format fields with fallback to legacy
			hadithNarrator := had.Narrator
			if hadithNarrator == "" {
				hadithNarrator = had.By
			}
			hadithText := had.English
			if hadithText == "" {
				hadithText = had.Text
			}
			hadithNum := had.Number
			if hadithNum == 0 {
				hadithNum = 1
			}
			dailyHadith = fmt.Sprintf("%s - %s\n\n%s", book.Name, hadithNarrator, hadithText)

			// Generate the contextual message once and cache it with the daily data
			dailyMessage = generateMessage(context.Background(), dailyVerse, dailyHadith, dailyName)

			links = map[string]string{
				"verse":  fmt.Sprintf("/quran/%d#%d", chap.Number, verseStart),
				"hadith": fmt.Sprintf("/hadith/%d#%d", book.Number, hadithNum),
				"name":   fmt.Sprintf("/names/%d", nam.Number),
			}

			dailyUpdated = time.Now()
			hijriDate := daily.Date().Display
			today := time.Now().UTC().Format("2006-01-02")
			timestamp := time.Now().UTC().Format(time.RFC3339)

			// Save hourly reminder with metadata
			hourlyData := map[string]interface{}{
				"timestamp": timestamp,
				"verse_meta": map[string]interface{}{
					"chapter":      chap.Number,
					"chapter_name": chap.English,
					"verse_start":  verseStart,
					"verse_end":    verseEnd,
					"text":         verseText,
				},
				"hadith_meta": map[string]interface{}{
					"book":      book.Number,
					"book_name": book.Name,
					"narrator":  hadithNarrator,
					"number":    hadithNum,
					"text":      hadithText,
				},
				"name_meta": map[string]interface{}{
					"number":  nam.Number,
					"english": nam.English,
					"arabic":  nam.Arabic,
					"meaning": nam.Meaning,
					"summary": nam.Summary,
				},
			}
			saveHourlyReminder(today, timestamp, hourlyData)

			mtx.Unlock()

			// Check if we should send push notification (new day or within grace period after midnight)
			if lastPushDate != today || (lastPushDate == today && isWithinGracePeriod()) {
				mtx.Lock()

				// Compose a user-friendly notification message
				notifyVerse := dailyVerse

				if len(dailyVerse) > 250 {
					notifyVerse = notifyVerse[:250] + "..."
				}

				dailyData := map[string]interface{}{
					"verse":   dailyVerse,
					"hadith":  dailyHadith,
					"name":    dailyName,
					"hijri":   hijriDate,
					"date":    today,
					"links":   links,
					"updated": dailyUpdated.Format(time.RFC3339),
					"message": dailyMessage,
				}

				// Save to daily.json
				saveDaily(today, dailyData)

				payload := map[string]interface{}{
					"title": "Reminder",
					"body":  notifyVerse,
					"data": map[string]interface{}{
						"url": "/daily/" + today,
					},
				}

				b, _ := json.Marshal(payload)

				// Only send if we haven't already sent today
				if lastPushDate != today {
					log.Println("Sending push notification at midnight UTC")

					errors := api.SendPushToAll(string(b))
					if len(errors) > 0 {
						log.Println("Push notification errors:")
						for _, err := range errors {
							log.Println(err)
						}
					}

					lastPushDate = today
					saveLastPushDate(today)
				} else {
					log.Println("Push notification already sent today, skipping")
				}

				mtx.Unlock()
			}

			}() // end of panic-recovery wrapper

			// Update hourly for /api/latest
			time.Sleep(time.Hour)
		}
	}

	if *ServerFlag {
		log.Println("Starting daily")
		go daily()

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
