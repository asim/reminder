package quran

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed data/*.json
var files embed.FS

type Chapter struct {
	Name   string
	Number int
	Verses []*Verse
}

type Verse struct {
	Number int
	Text   string
	Arabic string
}

type Quran struct {
	Chapters []*Chapter
}

func (q *Quran) Markdown() string {
	var data string

	for _, ch := range q.Chapters {
		data += fmt.Sprintln()
		data += fmt.Sprintf(`# %d`, ch.Number)
		data += fmt.Sprintln()
		data += fmt.Sprintln()
		data += fmt.Sprintf(`### %s`, ch.Name)
		data += fmt.Sprintln()

		// max 286 ayahs
		for _, verse := range ch.Verses {
			data += fmt.Sprintln()
			data += fmt.Sprintf(`#### %d:%d`, ch.Number, verse.Number)
			data += fmt.Sprintln()
			data += fmt.Sprintf(verse.Text)
			data += fmt.Sprintln()
		}

	}

	return data
}

func Load() *Quran {
	q := &Quran{}

	// load the arabic
	f, err := files.ReadFile("data/arabic.json")
	if err != nil {
		panic(err.Error())
	}
	var arabic map[string]interface{}
	json.Unmarshal(f, &arabic)

	// Set local
	for i := 0; i < 114; i++ {
		chapter := i + 1

		f, err = files.ReadFile(fmt.Sprintf("data/%d.json", chapter))
		if err != nil {
			panic(err.Error())
		}
		var data []interface{}
		json.Unmarshal(f, &data)

		name := data[0].(map[string]interface{})["name"].(map[string]interface{})["translated"].(string)

		data = data[1:]

		var verses []*Verse

		arabicText := arabic[fmt.Sprintf("%d", chapter)].([]interface{})

		for i, ayah := range data {
			arabicAyah := arabicText[i].(map[string]interface{})
			ch := int(arabicAyah["chapter"].(float64))
			ve := int(arabicAyah["verse"].(float64))
			ar := arabicAyah["text"].(string)

			num := int(ayah.([]interface{})[0].(float64))
			if ch != chapter {
				panic("arabic chapter mismatch")
			}

			if ve != num {
				panic("arabic verse mismatch")
			}

			verses = append(verses, &Verse{
				Number: num,
				Text:   ayah.([]interface{})[1].(string),
				Arabic: ar,
			})
		}

		// set the name
		q.Chapters = append(q.Chapters, &Chapter{
			Name:   name,
			Number: chapter,
			Verses: verses,
		})

	}

	return q
}

func Markdown() string {
	return Load().Markdown()
}
