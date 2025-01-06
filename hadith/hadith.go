package hadith

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed data/*.json
var files embed.FS

type Volume struct {
	Name  string  `json:"name"`
	Books []*Book `json:"books"`
}

type Book struct {
	Name    string    `json:"name"`
	Number  int       `json:"number"`
	Hadiths []*Hadith `json:"hadiths"`
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

	data += fmt.Sprintf(`<h1>%s</h1>`, b.Name)
	data += fmt.Sprintln()
	data += fmt.Sprintln()

	for _, hadith := range b.Hadiths {
		data += fmt.Sprintf(`<h3>%s</h3>`, hadith.Info)
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
		data += fmt.Sprintf(`<div class="chapter"><a href="/hadith/%d">%s</a></div>`, id+1, book.Name)
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
			Name:   book.Name,
			Number: book.Number,
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

		for num, b := range d["books"].([]interface{}) {
			bk := b.(map[string]interface{})

			book := &Book{
				Name:   bk["name"].(string),
				Number: num + 1,
			}

			for _, h := range bk["hadiths"].([]interface{}) {
				hd := h.(map[string]interface{})

				hadith := &Hadith{
					Info: hd["info"].(string),
					By:   hd["by"].(string),
					Text: hd["text"].(string),
				}

				book.Hadiths = append(book.Hadiths, hadith)
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
