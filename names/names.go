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

	// Header with bookmark
	data += `<div class="mb-6">`
	data += fmt.Sprintf(`<h1 class="text-3xl font-bold mb-2">%s <button class="bookmark-btn" data-type="names" data-key="%s" data-label="%s" data-url="%s">☆</button></h1>`,
		name.Meaning, nameKey, nameLabel, nameURL)
	data += `</div>`

	// Arabic name card
	data += `<div class="mb-6 p-6 bg-white border border-gray-200 rounded-lg shadow-sm">`
	data += `<div class="arabic text-center text-4xl mb-4">` + name.Arabic + `</div>`
	data += `<h3 class="text-xl font-semibold text-center text-gray-700">` + name.English + `</h3>`
	data += `</div>`

	// Description card
	data += `<div class="mb-6 p-6 bg-white border border-gray-200 rounded-lg shadow-sm">`
	data += `<h3 class="text-lg font-semibold mb-3">Description</h3>`
	data += `<p class="text-gray-700">` + name.Description + `</p>`
	data += `</div>`

	// Summary card
	data += `<div class="mb-6 p-6 bg-white border border-gray-200 rounded-lg shadow-sm">`
	data += `<h3 class="text-lg font-semibold mb-3">Summary</h3>`
	data += `<p class="text-gray-700">` + name.Summary + `</p>`
	data += `</div>`

	// Location card
	data += `<div class="mb-6 p-6 bg-white border border-gray-200 rounded-lg shadow-sm">`
	data += `<h3 class="text-lg font-semibold mb-3">Quran References</h3>`
	data += `<div class="flex flex-wrap gap-2">`

	for _, loc := range name.Location {
		uri := fmt.Sprintf("/quran/%s", strings.Replace(loc, ":", "/", -1))
		data += fmt.Sprintf(`<a href="%s" class="px-3 py-1 bg-blue-50 text-blue-600 rounded hover:bg-blue-100 transition-colors" hx-get="%s" hx-target="#main" hx-swap="innerHTML" hx-push-url="true">%s</a>`, uri, uri, loc)
	}
	data += `</div>`
	data += `</div>`

	return data
}

func (n *Names) Get(id int) *Name {
	return (*n)[id-1]
}

func (n *Names) TOC() string {
	var data string

	data += `<div id="contents" class="space-y-2">`
	for _, name := range *n {
		data += fmt.Sprintf(`<a href="/names/%d" hx-get="/names/%d" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="block p-3 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors">%d: %s</a>`, name.Number, name.Number, name.Number, name.Meaning)
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
