package names

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed data/*.json
var files embed.FS

var Description = `The "Names of Allah" (Asma al-Husna) refer to the 99 beautiful names of God. Descriptions of God's divine attributes, revealing His nature and qualities.`

type Name struct {
	Number      int      `json:"number"`
	English     string   `json:"english"`
	Arabic      string   `json:"arabic"`
	Meaning     string   `json:"meaning"`
	Description string   `json:"description"`
	Summary     string   `json:"summary"`
	Location    []string `json:"location"`
}

type Names []*Name

func (name *Name) HTML() string {
	var data string

	nameKey := fmt.Sprintf("%d", name.Number)
	nameLabel := fmt.Sprintf("Name %d: %s", name.Number, name.Meaning)
	nameURL := fmt.Sprintf("/names/%d", name.Number)

	data += fmt.Sprintln()
	data += fmt.Sprintln()
	data += fmt.Sprintf(`<h2>%s <button class="bookmark-btn" data-type="names" data-key="%s" data-label="%s" data-url="%s">☆</button></h2>`,
		name.Meaning, nameKey, nameLabel, nameURL)
	data += fmt.Sprintln()

	data += fmt.Sprintln()
	data += fmt.Sprintln(`<div class="arabic">` + name.Arabic + `</div>`)
	data += fmt.Sprintln()
	data += fmt.Sprintln(`<h4 class="english">` + name.English + `</h4>`)
	data += fmt.Sprintln()
	data += fmt.Sprintln(`<div class="english">` + name.Description + `</div>`)
	data += fmt.Sprintln()
	data += fmt.Sprintln(`<h4 class="english">Summary</h4>`)
	data += fmt.Sprintln()
	data += fmt.Sprintln(`<div class="english">` + name.Summary + `</div>`)
	data += fmt.Sprintln()
	data += fmt.Sprintln(`<h4 class="english">Location</h4>`)
	data += fmt.Sprintln()

	var locations string

	for _, loc := range name.Location {
		uri := fmt.Sprintf("/quran/%s", strings.Replace(loc, ":", "/", -1))
		locations += fmt.Sprintf(`<a href="%s">%s</a>&nbsp;`, uri, loc)
	}
	data += fmt.Sprintf(`<div>%s</div>`, locations)

	data += fmt.Sprintln()
	data += fmt.Sprintln()
	data += fmt.Sprintln(`<div class="dots">...</div>`)
	data += fmt.Sprintln()
	return data
}

func (n *Names) Get(id int) *Name {
	return (*n)[id-1]
}

func (n *Names) TOC() string {
	var data string

	data += `<div id="contents">`
	for _, name := range *n {
		data += fmt.Sprintf(`<div class="chapter"><a href="/names/%d">%d: %s</a></div>`, name.Number, name.Number, name.Meaning)
	}
	data += `</div>`

	return data
}

func (n *Names) JSON() []byte {
	b, _ := json.Marshal(n)
	return b
}

func (n *Names) HTML() string {
	var data string

	for _, name := range *n {
		data += name.HTML()
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
		data += fmt.Sprintln(`#### Summary`)
		data += fmt.Sprintln()
		data += fmt.Sprintf(`%s`, name.Summary)
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

		fnd := strings.Replace(n["found"].(string), " ", "", -1)
		fnd = strings.Replace(fnd, "(", " ", -1)
		fnd = strings.Replace(fnd, ")", " ", -1)
		loc := strings.Split(strings.TrimSpace(fnd), " ")

		name := &Name{
			Number:      int(n["number"].(float64)),
			English:     n["transliteration"].(string),
			Arabic:      n["name"].(string),
			Meaning:     en["meaning"].(string),
			Description: en["desc"].(string),
			Location:    loc,
		}

		// try get summary
		fd, err := files.ReadFile(fmt.Sprintf("data/%d.json", name.Number))
		if err == nil {
			var data map[string]interface{}
			json.Unmarshal(fd, &data)
			summary, ok := data["summary"].(string)
			if ok {
				name.Summary = summary
			}
		}

		*names = append(*names, name)
	}

	return names
}

func Markdown() string {
	return Load().Markdown()
}
