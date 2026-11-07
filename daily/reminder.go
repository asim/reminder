package daily

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/asim/reminder/hadith"
	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
	"log"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

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

func generateMessage(ctx context.Context, verse, hadith, name string) (message string) {
	// Fallback message in case LLM fails
	defaultMessage := "In the Name of Allah—the Most Compassionate, Most Merciful"
	message = defaultMessage // Set default in case of panic

	// Build context for the LLM
	contexts := []string{
		fmt.Sprintf("Verse from the Quran: %s", verse),
		fmt.Sprintf("Hadith: %s", hadith),
		fmt.Sprintf("Name of Allah: %s", name),
	}

	// Create a question that asks the LLM to generate a beneficial message
	question := "Based on the provided Quranic verse, Hadith, and name of Allah, generate a short, beneficial, and factual message (2-3 sentences) that provides spiritual guidance and reflection for the reader."

	// Use the LLM to generate the message
	// Wrap in a recover to handle panics from askLLM (which panics on errors)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Failed to generate contextual message via LLM: %v\n", r)
			message = defaultMessage // Ensure we return default on panic
		}
	}()

	llmMessage := model.Ask(ctx, contexts, question)

	// If message is not empty, use it
	if len(strings.TrimSpace(llmMessage)) > 0 {
		message = llmMessage
	}

	return message
}

func getVerse(ch *quran.Chapter, ve *quran.Verse) (string, int, int, string) {
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

func isWithinGracePeriod() bool {
	now := time.Now().UTC()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	sinceMidnight := now.Sub(midnight)
	return sinceMidnight >= 0 && sinceMidnight < 5*time.Minute
}

// Update selects and archives the next reminder, sending the existing daily
// notification when due. The caller controls when and how often it runs.
func (s *Store) Update(ctx context.Context, q *quran.Quran, n *names.Names, b *hadith.Collection, rnd *rand.Rand, sendPush func(string) []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
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
		return
	}

	s.latest.Name = fmt.Sprintf("%s - %s - %s\n\n%s", nam.English, nam.Arabic, nam.Meaning, nam.Summary)
	verseFormatted, verseStart, verseEnd, verseText := getVerse(chap, ver)
	s.latest.Verse = verseFormatted
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
	s.latest.Hadith = fmt.Sprintf("%s - %s\n\n%s", book.Name, hadithNarrator, hadithText)

	// Generate the contextual message once and cache it with the daily data
	s.latest.Message = generateMessage(ctx, s.latest.Verse, s.latest.Hadith, s.latest.Name)

	s.latest.Links = map[string]string{
		"verse":  fmt.Sprintf("/quran/%d#%d", chap.Number, verseStart),
		"hadith": fmt.Sprintf("/hadith/%d#%d", book.Number, hadithNum),
		"name":   fmt.Sprintf("/names/%d", nam.Number),
	}

	s.latest.Updated = time.Now().Format(time.RFC3339)
	hijriDate := Date().Display
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
	s.saveHourlyReminder(today, timestamp, hourlyData)

	// Check if we should send push notification (new day or within grace period after midnight)
	if s.lastPushDate != today || (s.lastPushDate == today && isWithinGracePeriod()) {

		// Compose a user-friendly notification message
		notifyVerse := s.latest.Verse

		if len(s.latest.Verse) > 250 {
			notifyVerse = notifyVerse[:250] + "..."
		}

		dailyData := map[string]interface{}{
			"verse":   s.latest.Verse,
			"hadith":  s.latest.Hadith,
			"name":    s.latest.Name,
			"hijri":   hijriDate,
			"date":    today,
			"links":   s.latest.Links,
			"updated": s.latest.Updated,
			"message": s.latest.Message,
		}

		// Save to daily.json
		s.saveDaily(today, dailyData)

		payload := map[string]interface{}{
			"title": "Reminder",
			"body":  notifyVerse,
			"data": map[string]interface{}{
				"url": "/daily/" + today,
			},
		}

		b, _ := json.Marshal(payload)

		// Only send if we haven't already sent today
		if s.lastPushDate != today {
			log.Println("Sending push notification at midnight UTC")

			errors := sendPush(string(b))
			if len(errors) > 0 {
				log.Println("Push notification errors:")
				for _, err := range errors {
					log.Println(err)
				}
			}

			s.lastPushDate = today
			s.saveLastPushDate(today)
		} else {
			log.Println("Push notification already sent today, skipping")
		}

	}

}
