package quran

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed data/*.json
//go:embed data/words/*.json
var files embed.FS

var Bismillah = `بِسۡمِ ٱللَّهِ ٱلرَّحۡمَٰنِ ٱلرَّحِيمِ`
var English = `In the Name of Allah—the Most Beneficent, Most Merciful.`

var Description = `The word of God, as revealed to Prophet Muhammad (peace be upon him). It is a guide for Muslims (believers) on faith, morality, and life through its 114 chapters.`

type Chapter struct {
	Name       string   `json:"name"`
	Number     int      `json:"number"`
	Verses     []*Verse `json:"verses,omitempty"`
	English    string   `json:"english"`
	VerseCount int      `json:"verse_count"`
}

type Verse struct {
	Chapter  int     `json:"chapter"`
	Number   int     `json:"number"`
	Text     string  `json:"text"`
	Arabic   string  `json:"arabic"`
	Words    []*Word `json:"words"`
	Comments string  `json:"comments"`
}

type Word struct {
	English         string `json:"english"`
	Arabic          string `json:"arabic"`
	Transliteration string `json:"transliteration"`
}

type Comment struct {
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Text    string `json:"text"`
}

type Quran struct {
	Chapters   []*Chapter `json:"chapters"`
	Commentary []*Comment `json:"commentary"`
}

func (ch *Chapter) JSON() []byte {
	b, _ := json.Marshal(ch)
	return b
}

func (ch *Chapter) HTML() string {
	var data string

	// Chapter header
	data += `<div class="mb-6">`
	data += fmt.Sprintf(`<h1 class="text-3xl font-bold mb-2">%s</h1>`, ch.English)
	data += fmt.Sprintf(`<h2 class="text-xl arabic text-gray-600">%s</h2>`, ch.Name)
	data += `</div>`

	// Verses
	for _, verse := range ch.Verses {
		verseKey := fmt.Sprintf("%d:%d", ch.Number, verse.Number)
		verseLabel := fmt.Sprintf("Quran %d:%d", ch.Number, verse.Number)
		verseURL := fmt.Sprintf("/quran/%d#%d", ch.Number, verse.Number)

		data += `<div class="mb-6 p-6 bg-white border border-gray-200 rounded-lg shadow-sm" id="` + fmt.Sprintf("%d", verse.Number) + `">`
		data += fmt.Sprintf(`<div class="flex items-center justify-between mb-4"><h3 class="text-lg font-semibold text-gray-700">%d:%d</h3><button class="bookmark-btn" data-type="quran" data-key="%s" data-label="%s" data-url="%s">☆</button></div>`,
			ch.Number, verse.Number, verseKey, verseLabel, verseURL)
		data += `<div class="arabic text-right text-2xl mb-4 leading-relaxed">` + verse.Arabic + `</div>`
		data += `<div class="text-gray-700">` + verse.Text + `</div>`
		data += `</div>`
	}

	return data
}

func (v *Verse) HTML() string {
	var data string

	verseKey := fmt.Sprintf("%d:%d", v.Chapter, v.Number)
	verseLabel := fmt.Sprintf("Quran %d:%d", v.Chapter, v.Number)
	verseURL := fmt.Sprintf("/quran/%d#%d", v.Chapter, v.Number)

	data += fmt.Sprintln()
	data += fmt.Sprintln()
	data += fmt.Sprintln()
	data += fmt.Sprintf(`<h4>%d:%d <button class="bookmark-btn" data-type="quran" data-key="%s" data-label="%s" data-url="%s">☆</button></h4>`,
		v.Chapter, v.Number, verseKey, verseLabel, verseURL)
	data += fmt.Sprintln()
	data += fmt.Sprintln(`<div class="arabic right">` + v.Arabic + `</div>`)
	data += fmt.Sprintln()
	data += fmt.Sprintln(`<div class="english">` + v.Text + `</div>`)
	data += fmt.Sprintln()
	data += fmt.Sprintln(`<div class="dots">...</div>`)
	data += fmt.Sprintln()

	return data
}

func (v *Verse) JSON() []byte {
	b, _ := json.Marshal(v)
	return b
}

func (q *Quran) Index() *Quran {
	nq := new(Quran)

	for _, ch := range q.Chapters {
		chapter := &Chapter{
			Name:       ch.Name,
			Number:     ch.Number,
			English:    ch.English,
			VerseCount: len(ch.Verses),
		}
		nq.Chapters = append(nq.Chapters, chapter)
	}

	return nq
}

func (q *Quran) JSON() []byte {
	b, _ := json.Marshal(q)
	return b
}

func (q *Quran) TOC() string {
	var data string

	data += `<div id="contents" class="space-y-2">`
	for _, ch := range q.Chapters {
		data += fmt.Sprintf(`<a href="/quran/%d" hx-get="/quran/%d" hx-target="#main" hx-swap="innerHTML" hx-push-url="true" class="block p-3 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors">%d: %s</a>`, ch.Number, ch.Number, ch.Number, ch.English)
	}
	data += `</div>`

	return data
}

func (q *Quran) Get(chapter int) *Chapter {
	return q.Chapters[chapter-1]
}

func (q *Quran) HTML() string {
	var data string

	for _, ch := range q.Chapters {
		data += ch.HTML()
	}

	return data
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
			data += fmt.Sprintln(verse.Arabic)
			data += fmt.Sprintln()
			data += fmt.Sprintln(verse.Text)
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

	// load tafsir
	fx, err := files.ReadFile("data/tafsir.json")
	if err != nil {
		panic(err.Error())
	}
	var tafsir map[string]interface{}
	json.Unmarshal(fx, &tafsir)

	// Set local
	for i := 0; i < 114; i++ {
		chapter := i + 1

		f, err = files.ReadFile(fmt.Sprintf("data/%d.json", chapter))
		if err != nil {
			panic(err.Error())
		}
		var data []interface{}
		json.Unmarshal(f, &data)

		f2, err := files.ReadFile(fmt.Sprintf("data/words/%d.json", chapter))
		if err != nil {
			panic(err.Error())
		}
		var words map[string]interface{}
		json.Unmarshal(f2, &words)

		english := data[0].(map[string]interface{})["name"].(map[string]interface{})["translated"].(string)
		name := data[0].(map[string]interface{})["name"].(map[string]interface{})["transliterated"].(string)
		data = data[1:]

		var verses []*Verse

		arabicText := arabic[fmt.Sprintf("%d", chapter)].([]interface{})

		// Add Bismillah as first verse for all chapters except 1 and 9
		if chapter != 1 && chapter != 9 {
			verses = append(verses, &Verse{
				Chapter: chapter,
				Number:  0,
				Text:    English,
				Arabic:  Bismillah,
			})
		}

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

			// get words
			inum := fmt.Sprintf("%v", num)
			w := words[inum].(map[string]interface{})["w"].([]interface{})
			var wbw []*Word
			for _, word := range w {
				wordmap := word.(map[string]interface{})
				wbw = append(wbw, &Word{
					Arabic:          wordmap["c"].(string),
					English:         wordmap["e"].(string),
					Transliteration: wordmap["d"].(string),
				})
			}

			var text string
			key := fmt.Sprintf("%d:%d", chapter, num)
			comment, ok := tafsir[key].(map[string]interface{})

			if ok {
				text = comment["text"].(string)
			} else {
				v := tafsir[key].(string)
				comment = tafsir[v].(map[string]interface{})
				text = comment["text"].(string)
			}

			q.Commentary = append(q.Commentary, &Comment{
				Chapter: chapter,
				Verse:   num,
				Text:    text,
			})

			verses = append(verses, &Verse{
				Chapter:  chapter,
				Number:   num,
				Text:     ayah.([]interface{})[1].(string),
				Arabic:   ar,
				Words:    wbw,
				Comments: text,
			})
		}

		// set the name
		q.Chapters = append(q.Chapters, &Chapter{
			Name:       name,
			Number:     chapter,
			Verses:     verses,
			English:    english,
			VerseCount: len(verses),
		})

	}

	return q
}

func Markdown() string {
	return Load().Markdown()
}
