package hadith

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed data/*.json
var files embed.FS

var Description = `A collection of the Prophet Muhammad's sayings and actions, providing essential context and practical guidance for Islamic practice and belief alongside the Quran.`

// Collection represents a hadith collection like Bukhari or Muslim
type Collection struct {
	Name    string  `json:"name"`
	Arabic  string  `json:"arabic"`
	Books   []*Book `json:"books"`
}

type Book struct {
	Name        string    `json:"name"`
	Number      int       `json:"number"`
	Hadiths     []*Hadith `json:"hadiths,omitempty"`
	HadithCount int       `json:"hadith_count,omitempty"`
}

type Hadith struct {
	Number   int    `json:"number"`
	Narrator string `json:"narrator"`
	English  string `json:"english"`
	Arabic   string `json:"arabic"`
	Chain    string `json:"chain,omitempty"`
	// Legacy fields for API compatibility
	Info string `json:"info,omitempty"`
	By   string `json:"by,omitempty"`
	Text string `json:"text,omitempty"`
}

func (b *Book) JSON() []byte {
	by, _ := json.Marshal(b)
	return by
}

func (c *Collection) TOC() string {
	var data string

	data += `<div class="space-y-2">`
	for _, book := range c.Books {
		data += fmt.Sprintf(`<a href="/hadith/%d" class="block p-3 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors">%d: %s</a>`, book.Number, book.Number, book.Name)
	}
	data += `</div>`

	return data
}

func (c *Collection) Get(book int) *Book {
	if book < 1 || book > len(c.Books) {
		return nil
	}
	return c.Books[book-1]
}

func (c *Collection) Index() *Collection {
	cc := &Collection{
		Name:   c.Name,
		Arabic: c.Arabic,
	}

	for _, book := range c.Books {
		cc.Books = append(cc.Books, &Book{
			Name:        book.Name,
			Number:      book.Number,
			HadithCount: len(book.Hadiths),
		})
	}

	return cc
}

func (c *Collection) JSON() []byte {
	b, _ := json.Marshal(c)
	return b
}

func Load() *Collection {
	collection := &Collection{}

	f, err := files.ReadFile("data/bukhari.json")
	if err != nil {
		panic(err.Error())
	}

	if err := json.Unmarshal(f, collection); err != nil {
		panic(err.Error())
	}

	// Set book numbers if not set
	for i, book := range collection.Books {
		if book.Number == 0 {
			book.Number = i + 1
		}
		book.HadithCount = len(book.Hadiths)
		
		// Set legacy fields for API compatibility
		for j, h := range book.Hadiths {
			if h.Number == 0 {
				h.Number = j + 1
			}
			h.Info = fmt.Sprintf("Hadith %d", h.Number)
			h.By = h.Narrator
			h.Text = h.English
		}
	}

	return collection
}

func (b *Book) HTML() string {
	var data string

	// Book header
	data += `<div class="mb-6">`
	data += fmt.Sprintf(`<h1 class="text-3xl font-bold mb-2">%s</h1>`, b.Name)
	data += `</div>`

	// Hadith entries
	for _, hadith := range b.Hadiths {
		hadithKey := fmt.Sprintf("%d:%d", b.Number, hadith.Number)
		hadithLabel := fmt.Sprintf("%s - Hadith %d", b.Name, hadith.Number)
		hadithURL := fmt.Sprintf("/hadith/%d#%d", b.Number, hadith.Number)

		data += `<div class="mb-6 p-6 bg-white border border-gray-200 rounded-lg shadow-sm" id="` + fmt.Sprintf("%d", hadith.Number) + `">`
		data += fmt.Sprintf(`<div class="flex items-center justify-between mb-3"><h3 class="text-lg font-semibold text-gray-700">Hadith %d</h3><button class="bookmark-btn" data-type="hadith" data-key="%s" data-label="%s" data-url="%s">☆</button></div>`,
			hadith.Number, hadithKey, hadithLabel, hadithURL)
		data += fmt.Sprintf(`<p class="text-sm text-gray-500 mb-4">%s</p>`, hadith.Narrator)
		
		// Arabic text
		if hadith.Arabic != "" {
			data += fmt.Sprintf(`<div dir="rtl" class="text-xl leading-loose font-arabic text-right mb-4 pb-4 border-b border-gray-100">%s</div>`, hadith.Arabic)
		}
		
		// English translation
		data += fmt.Sprintf(`<div class="text-gray-700">%s</div>`, hadith.English)
		data += `</div>`
	}

	return data
}
