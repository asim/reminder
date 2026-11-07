package daily

import (
	"fmt"
	"io"
	"time"
)

// WriteRSS writes the existing daily and hourly feed from a consistent archive.
func (s *Store) WriteRSS(w io.Writer) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Get all dates and sort them
	var dates []string
	for date := range s.index {
		if date == "latest" {
			continue
		}
		dates = append(dates, date)
	}

	// Sort dates in descending order (newest first)
	for i := 0; i < len(dates); i++ {
		for j := i + 1; j < len(dates); j++ {
			if dates[i] < dates[j] {
				dates[i], dates[j] = dates[j], dates[i]
			}
		}
	}

	// Start RSS feed
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Daily Reminder</title>
    <link>https://reminder.dev</link>
    <description>Daily reminders from the Quran, Hadith, and Names of Allah</description>
    <language>en-us</language>
    <lastBuildDate>%s</lastBuildDate>
    <atom:link href="https://reminder.dev/rss" rel="self" type="application/rss+xml" />
`, time.Now().Format(time.RFC1123Z))

	// Add items for each date (limit to most recent 30 days)
	maxDays := 30
	if len(dates) < maxDays {
		maxDays = len(dates)
	}

	for i := 0; i < maxDays; i++ {
		date := dates[i]

		// Parse date for pubDate base
		basePubDate := date
		if t, err := time.Parse("2006-01-02", date); err == nil {
			basePubDate = t.Format(time.RFC1123Z)
		}

		// First, add hourly reminders for this date (most recent first)
		hourlyReminders := s.loadHourlyReminders(date)

		// Process each hourly reminder in reverse order (newest first)
		for i := len(hourlyReminders) - 1; i >= 0; i-- {
			reminderData := hourlyReminders[i]
			reminder, ok := reminderData.(map[string]interface{})
			if !ok {
				continue
			}

			// Get timestamp for pubDate
			pubDate := basePubDate
			if ts, ok := reminder["timestamp"].(string); ok {
				if t, err := time.Parse(time.RFC3339, ts); err == nil {
					pubDate = t.Format(time.RFC1123Z)
				}
			}

			// Extract verse metadata
			if verseMeta, ok := reminder["verse_meta"].(map[string]interface{}); ok {
				chapterName, _ := verseMeta["chapter_name"].(string)
				chapter, _ := verseMeta["chapter"].(float64)
				verseStart, _ := verseMeta["verse_start"].(float64)
				verseEnd, _ := verseMeta["verse_end"].(float64)
				text, _ := verseMeta["text"].(string)

				verseTitle := ""
				if verseStart == verseEnd {
					verseTitle = fmt.Sprintf("%s %d:%d", chapterName, int(chapter), int(verseStart))
				} else {
					verseTitle = fmt.Sprintf("%s %d:%d-%d", chapterName, int(chapter), int(verseStart), int(verseEnd))
				}

				verseLink := fmt.Sprintf("https://reminder.dev/quran/%d#%d", int(chapter), int(verseStart))

				fmt.Fprintf(w, `    <item>
      <title>%s</title>
      <link>%s</link>
      <guid>%s</guid>
      <pubDate>%s</pubDate>
      <description><![CDATA[%s]]></description>
    </item>
`, verseTitle, verseLink, verseLink, pubDate, text)
			}

			// Extract hadith metadata
			if hadithMeta, ok := reminder["hadith_meta"].(map[string]interface{}); ok {
				bookName, _ := hadithMeta["book_name"].(string)
				book, _ := hadithMeta["book"].(float64)
				info, _ := hadithMeta["info"].(string)
				text, _ := hadithMeta["text"].(string)

				// Get number as string
				number := "1"
				if num, ok := hadithMeta["number"].(string); ok {
					number = num
				}

				hadithTitle := fmt.Sprintf("%s - %s", bookName, info)
				hadithLink := fmt.Sprintf("https://reminder.dev/hadith/%d#%s", int(book), number)

				fmt.Fprintf(w, `    <item>
      <title>%s</title>
      <link>%s</link>
      <guid>%s</guid>
      <pubDate>%s</pubDate>
      <description><![CDATA[%s]]></description>
    </item>
`, hadithTitle, hadithLink, hadithLink, pubDate, text)
			}

			// Extract name metadata
			if nameMeta, ok := reminder["name_meta"].(map[string]interface{}); ok {
				number, _ := nameMeta["number"].(float64)
				english, _ := nameMeta["english"].(string)
				arabic, _ := nameMeta["arabic"].(string)
				meaning, _ := nameMeta["meaning"].(string)
				summary, _ := nameMeta["summary"].(string)

				nameTitle := fmt.Sprintf("%s - %s - %s", english, arabic, meaning)
				nameLink := fmt.Sprintf("https://reminder.dev/names/%d", int(number))

				fmt.Fprintf(w, `    <item>
      <title>%s</title>
      <link>%s</link>
      <guid>%s</guid>
      <pubDate>%s</pubDate>
      <description><![CDATA[%s]]></description>
    </item>
`, nameTitle, nameLink, nameLink, pubDate, summary)
			}
		}

		// Then, add daily summary at the end (for single daily consumption)
		entry, hasDaily := s.index[date]
		if hasDaily {
			entryMap, ok := entry.(map[string]interface{})
			if ok {
				verse := ""
				hadith := ""
				name := ""
				hijri := ""

				if v, ok := entryMap["verse"].(string); ok {
					verse = v
				}
				if h, ok := entryMap["hadith"].(string); ok {
					hadith = h
				}
				if n, ok := entryMap["name"].(string); ok {
					name = n
				}
				if hj, ok := entryMap["hijri"].(string); ok {
					hijri = hj
				}

				title := fmt.Sprintf("Daily Reminder - %s", date)
				if hijri != "" {
					title = fmt.Sprintf("Daily Reminder - %s (%s)", date, hijri)
				}

				// Combine all three into one item description
				description := ""
				if verse != "" {
					description += "<h3>Verse</h3>\n" + verse + "\n\n"
				}
				if hadith != "" {
					description += "<h3>Hadith</h3>\n" + hadith + "\n\n"
				}
				if name != "" {
					description += "<h3>Name</h3>\n" + name
				}

				if description != "" {
					fmt.Fprintf(w, `    <item>
      <title>%s</title>
      <link>https://reminder.dev/daily/%s</link>
      <guid>https://reminder.dev/daily/%s</guid>
      <pubDate>%s</pubDate>
      <description><![CDATA[%s]]></description>
    </item>
`, title, date, date, basePubDate, description)
				}
			}
		}
	}

	fmt.Fprintf(w, `  </channel>
</rss>`)
}
