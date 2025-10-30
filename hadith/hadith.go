package hadith

import (
	"embed"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

//go:embed data/*.json
var files embed.FS

var Description = `A collection of the Prophet Muhammad's sayings and actions, providing essential context and practical guidance for Islamic practice and belief alongside the Quran.`

type Volume struct {
	Name  string  `json:"name"`
	Books []*Book `json:"books"`
}

type Book struct {
	Name        string    `json:"name"`
	Number      int       `json:"number"`
	Hadiths     []*Hadith `json:"hadiths,omitempty"`
	HadithCount int       `json:"hadith_count,omitempty"`
}

type Hadith struct {
	Info string `json:"info"`
	By   string `json:"by"`
	Text string `json:"text"`
}

type Volumes struct {
	Contents []*Volume `json:"contents,omitempty"`
	Books    []*Book   `json:"books,omitempty"`
}

func (b *Book) JSON() []byte {
	by, _ := json.Marshal(b)
	return by
}

func (b *Book) HTML() string {
	var data string

	data += fmt.Sprintf(`<h2>%s</h2>`, b.Name)
	data += fmt.Sprintln()
	data += fmt.Sprintln()

	for idx, hadith := range b.Hadiths {
		hadithKey := fmt.Sprintf("%d:%d", b.Number, idx+1)
		hadithLabel := fmt.Sprintf("Hadith %d:%d - %s", b.Number, idx+1, hadith.Info)
		hadithURL := fmt.Sprintf("/hadith/%d#%d", b.Number, idx+1)

		data += fmt.Sprintf(`<h3 id="%d">%s <button class="bookmark-btn" data-type="hadith" data-key="%s" data-label="%s" data-url="%s">☆</button></h3>`,
			idx+1, hadith.Info, hadithKey, hadithLabel, hadithURL)
		data += fmt.Sprintln()
		data += fmt.Sprintf(`<h4>%s</h4>`, hadith.By)
		data += fmt.Sprintln()
		data += fmt.Sprintf(`<div>%s</div>`, hadith.Text)
		data += fmt.Sprintln()
		data += fmt.Sprintln(`<div class="dots">...</div>`)
		data += fmt.Sprintln()
	}

	return data
}

func (v *Volumes) TOC() string {
	var data string

	for id, book := range v.Books {
		data += fmt.Sprintf(`<div class="chapter"><a href="/hadith/%d">%d: %s</a></div>`, id+1, id+1, book.Name)
	}

	return data
}

func (v *Volumes) Get(book int) *Book {
	return v.Books[book-1]
}

func (v *Volumes) Index() *Volumes {
	vv := new(Volumes)

	for _, book := range v.Books {
		vv.Books = append(vv.Books, &Book{
			Name:        book.Name,
			Number:      book.Number,
			HadithCount: len(book.Hadiths),
		})
	}

	return vv
}

func (v *Volumes) JSON() []byte {
	b, _ := json.Marshal(v)
	return b
}

func (v *Volumes) Markdown() string {
	var data string

	for _, volume := range v.Contents {
		data += fmt.Sprintln()
		data += fmt.Sprintf(`# %s`, volume.Name)
		data += fmt.Sprintln()
		data += fmt.Sprintln()

		for _, book := range volume.Books {
			data += fmt.Sprintf(`## %s`, book.Name)
			data += fmt.Sprintln()
			data += fmt.Sprintln()

			for _, hadith := range book.Hadiths {
				data += fmt.Sprintf(`### %s`, hadith.Info)
				data += fmt.Sprintln()
				data += fmt.Sprintf(`#### By %s`, hadith.By)
				data += fmt.Sprintln()
				data += fmt.Sprintf(`%s`, hadith.Text)
				data += fmt.Sprintln()
				data += fmt.Sprintln()
			}
		}
	}

	return data
}

func Load() *Volumes {
	volumes := &Volumes{}

	f, err := files.ReadFile("data/bukhari.json")
	if err != nil {
		panic(err.Error())
	}
	var data []interface{}
	json.Unmarshal(f, &data)

	// per volume
	for _, entry := range data {
		d := entry.(map[string]interface{})
		volume := &Volume{
			Name: d["name"].(string),
		}

		for _, b := range d["books"].([]interface{}) {
			bk := b.(map[string]interface{})

			parts := strings.Split(bk["name"].(string), ". ")
			num, name := parts[0], parts[1]
			n, _ := strconv.Atoi(num)

			book := &Book{
				Name:   name,
				Number: n,
			}

			for _, h := range bk["hadiths"].([]interface{}) {
				hd := h.(map[string]interface{})

				hadith := &Hadith{
					Info: hd["info"].(string),
					By:   hd["by"].(string),
					Text: hd["text"].(string),
				}

				book.Hadiths = append(book.Hadiths, hadith)
				book.HadithCount = len(book.Hadiths)
			}

			volume.Books = append(volume.Books, book)
			volumes.Books = append(volumes.Books, book)
		}

		volumes.Contents = append(volumes.Contents, volume)
	}

	return volumes
}

func Markdown() string {
	return Load().Markdown()
}
