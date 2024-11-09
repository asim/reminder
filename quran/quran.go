package quran

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed *.json
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

func Load() map[string]*Chapter {
	Quran := map[string]*Chapter{}

	// Set local
	for i := 0; i < 114; i++ {
		f, err := files.ReadFile(fmt.Sprintf("%d.json", i+1))
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
		Quran[fmt.Sprintf("%d", i+1)] = &Chapter{
			Name:   name,
			Number: i + 1,
			Verses: verses,
		}

	}

	return Quran
}
