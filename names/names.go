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

type Names []*Name

func (n *Names) HTML() string {
	var data string

	for _, name := range *n {
		data += fmt.Sprintln()
		data += fmt.Sprintf(`<h1 id="%d">%d</h1>`, name.Number, name.Number)
		data += fmt.Sprintln()
		data += fmt.Sprintln()
		data += fmt.Sprintf(`<h3>%s</h3>`, name.Meaning)
		data += fmt.Sprintln()

		data += fmt.Sprintln()
		data += fmt.Sprintln(`<div class="arabic">` + name.Arabic + `</div>`)
		data += fmt.Sprintln()
		data += fmt.Sprintln(`<h4 class="english">` + name.English + `</h4>`)
		data += fmt.Sprintln()
		data += fmt.Sprintln(`<div class="english">` + name.Description + `</div>`)
	}

	return data
}

func (n *Names) Markdown() string {
	var data string

	for _, name := range *n {
		data += fmt.Sprintln()
		data += fmt.Sprintf(`# %d`, name.Number)
		data += fmt.Sprintln()
		data += fmt.Sprintln()
		data += fmt.Sprintf(`### %s`, name.Meaning)
		data += fmt.Sprintln()
		data += fmt.Sprintf(`#### %s`, name.English)
		data += fmt.Sprintln()
		data += fmt.Sprintf(`#### %s`, name.Arabic)
		data += fmt.Sprintln()
		data += fmt.Sprintf(`%s`, name.Description)
		data += fmt.Sprintln()
	}

	return data
}

func Load() *Names {
	names := &Names{}

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

		*names = append(*names, name)
	}

	return names
}

func Markdown() string {
	return Load().Markdown()
}
