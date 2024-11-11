package names

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed data/*.json
var files embed.FS

type Name struct {
	Number      int
	English     string
	Arabic      string
	Meaning     string
	Description string
}

func Load() []*Name {
	var names []*Name

	f, err := files.ReadFile("data/names.json")
	if err != nil {
		panic(err.Error())
	}
	var data map[string]interface{}
	json.Unmarshal(f, &data)

	d := data["data"].([]interface{})

	for _, entry := range d {
		n := entry.(map[string]interface{})
		en := n["en"].(map[string]interface{})
		name := &Name{
			Number:      int(n["number"].(float64)),
			English:     n["transliteration"].(string),
			Arabic:      n["name"].(string),
			Meaning:     en["meaning"].(string),
			Description: en["desc"].(string),
		}
		names = append(names, name)
	}

	return names
}

func Markdown() string {
	names := Load()

	var data string

	for _, name := range names {
		data += fmt.Sprintln()
		data += fmt.Sprintf(`# %d`, name.Number)
		data += fmt.Sprintln()
		data += fmt.Sprintln()
		data += fmt.Sprintf(`### %s`, name.Meaning)
		data += fmt.Sprintln()
		data += fmt.Sprintf(`#### %s`, name.English)
		data += fmt.Sprintln()
		data += fmt.Sprintf(`%s`, name.Description)
		data += fmt.Sprintln()
	}

	return data
}
