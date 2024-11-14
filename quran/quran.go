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
			data += fmt.Sprintf(verse.Text)
			data += fmt.Sprintln()
		}

	}

	return data
}

func Load() *Quran {
	q := &Quran{}

	// Set local
	for i := 0; i < 114; i++ {
		f, err := files.ReadFile(fmt.Sprintf("data/%d.json", i+1))
		if err != nil {
			panic(err.Error())
		}
		var data []interface{}
		json.Unmarshal(f, &data)

		name := data[0].(map[string]interface{})["name"].(map[string]interface{})["translated"].(string)

		data = data[1:]

		var verses []*Verse

		for _, ayah := range data {
			verses = append(verses, &Verse{
				Number: int(ayah.([]interface{})[0].(float64)),
				Text:   ayah.([]interface{})[1].(string),
			})
		}

		// set the name
		q.Chapters = append(q.Chapters, &Chapter{
			Name:   name,
			Number: i + 1,
			Verses: verses,
		})

	}

	return q
}

func Markdown() string {
	return Load().Markdown()
}
