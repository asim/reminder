package hadith

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed data/*.json
var files embed.FS

type Volume struct {
	Name  string
	Books []*Book
}

type Book struct {
	Name    string
	Hadiths []*Hadith
}

type Hadith struct {
	Info string
	By   string
	Text string
}

type Volumes []*Volume

func (v *Volumes) Markdown() string {
	var data string

	for _, volume := range *v {
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

			book := &Book{
				Name: bk["name"].(string),
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
		}

		*volumes = append(*volumes, volume)
	}

	return volumes
}

func Markdown() string {
	return Load().Markdown()
}
