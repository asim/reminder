package daily

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Reminder is a detached snapshot of the latest hourly content.
type Reminder struct {
	Name    string            `json:"name"`
	Hadith  string            `json:"hadith"`
	Verse   string            `json:"verse"`
	Links   map[string]string `json:"links"`
	Updated string            `json:"updated"`
	Message string            `json:"message"`
}

// Store owns reminder state and the existing daily/hourly archive files.
// Mutation and archive helpers are private; readers receive snapshots.
type Store struct {
	mu           sync.RWMutex
	dir          string
	index        map[string]interface{}
	lastPushDate string
	latest       Reminder
}

func NewStore(dir string) *Store {
	s := &Store{dir: dir, latest: Reminder{Links: map[string]string{}, Updated: time.Time{}.Format(time.RFC3339)}}
	s.lastPushDate = s.loadLastPushDate()
	s.index = s.loadDailyIndex()
	return s
}

func (s *Store) path(name string) string { return filepath.Join(s.dir, name) }

func (s *Store) Latest() Reminder {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := s.latest
	result.Links = make(map[string]string, len(s.latest.Links))
	for k, v := range s.latest.Links {
		result.Links[k] = v
	}
	return result
}

func (s *Store) Get(date string) (json.RawMessage, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.index[date]
	if !ok {
		return nil, false
	}
	data, _ := json.Marshal(entry)
	return data, true
}

func (s *Store) IndexJSON() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, _ := json.MarshalIndent(s.index, "", "  ")
	return data
}

func (s *Store) Today() interface{} {
	today := time.Now().UTC().Format("2006-01-02")
	if entry, ok := s.Get(today); ok {
		return entry
	}
	r := s.Latest()
	return map[string]interface{}{
		"name": r.Name, "hadith": r.Hadith, "verse": r.Verse,
		"links": r.Links, "updated": r.Updated,
		"message": "In the Name of Allah—the Most Compassionate, Most Merciful",
		"hijri":   Date().Display, "date": today,
	}
}

func (s *Store) loadLastPushDate() string {
	log.Println("Load last pushdate")
	b, err := os.ReadFile(s.path("last_push_date.txt"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func (s *Store) saveLastPushDate(date string) {
	_ = os.MkdirAll(s.dir, 0700)
	_ = os.WriteFile(s.path("last_push_date.txt"), []byte(date), 0644)
}

// On startup, load daily index
func (s *Store) loadDailyIndex() map[string]interface{} {
	log.Println("Load daily index")
	dailyFile := s.path("daily.json")
	var idx map[string]interface{}
	if b, err := os.ReadFile(dailyFile); err == nil {
		json.Unmarshal(b, &idx)
	}
	if idx == nil {
		idx = make(map[string]interface{})
	}
	return idx
}

func (s *Store) saveDaily(date string, data map[string]interface{}) {
	// Save to daily.json
	dailyFile := s.path("daily.json")
	var allDaily map[string]interface{}
	if b, err := os.ReadFile(dailyFile); err == nil {
		json.Unmarshal(b, &allDaily)
	}
	if allDaily == nil {
		allDaily = make(map[string]interface{})
	}
	allDaily[date] = data
	b, _ := json.MarshalIndent(allDaily, "", "  ")
	_ = os.MkdirAll(s.dir, 0700)
	_ = os.WriteFile(dailyFile, b, 0644)
	s.index = allDaily // update in-memory index
}

func (s *Store) saveHourlyReminder(date, timestamp string, data map[string]interface{}) {
	// Save to hourly reminders file
	hourlyFile := s.path(fmt.Sprintf("hourly_%s.json", date))
	var hourlyReminders []interface{}
	if b, err := os.ReadFile(hourlyFile); err == nil {
		json.Unmarshal(b, &hourlyReminders)
	}
	if hourlyReminders == nil {
		hourlyReminders = make([]interface{}, 0)
	}
	hourlyReminders = append(hourlyReminders, data)
	b, _ := json.MarshalIndent(hourlyReminders, "", "  ")
	_ = os.MkdirAll(s.dir, 0700)
	_ = os.WriteFile(hourlyFile, b, 0644)
}

func (s *Store) loadHourlyReminders(date string) []interface{} {
	hourlyFile := s.path(fmt.Sprintf("hourly_%s.json", date))
	var hourlyReminders []interface{}
	if b, err := os.ReadFile(hourlyFile); err == nil {
		json.Unmarshal(b, &hourlyReminders)
	}
	return hourlyReminders
}
