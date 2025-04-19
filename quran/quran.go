package quran

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed data/*.json
var files embed.FS

var Bismillah = `بِسۡمِ ٱللَّهِ ٱلرَّحۡمَٰنِ ٱلرَّحِيمِ`
var English = `In the Name of Allah—the Most Compassionate, Most Merciful.`

var Description = `The word of God, as revealed to Prophet Muhammad (peace be upon him). It is a guide for Muslims (believers) on faith, morality, and life through its 114 chapters.`

type Chapter struct {
	Name       string   `json:"name"`
	Number     int      `json:"number"`
	Verses     []*Verse `json:"verses,omitempty"`
	English    string   `json:"english"`
	VerseCount int      `json:"verse_count"`
}

type Verse struct {
	Chapter int    `json:"chapter"`
	Number  int    `json:"number"`
	Text    string `json:"text"`
	Arabic  string `json:"arabic"`
}

type Quran struct {
	Chapters []*Chapter `json:"chapters"`
}

func (ch *Chapter) JSON() []byte {
	b, _ := json.Marshal(ch)
	return b
}

func (ch *Chapter) HTML() string {
	var data string

	data += fmt.Sprintln()
	data += fmt.Sprintln()
	data += fmt.Sprintf(`<h2>%s</h2>`, ch.English)
	data += fmt.Sprintln()
	data += fmt.Sprintf(`<h3>%s</h3>`, ch.Name)
	data += fmt.Sprintln()

	// max 286 ayahs
	for _, verse := range ch.Verses {
		data += fmt.Sprintln()
		data += fmt.Sprintf(`<h4 id="%d">%d:%d</h4>`, verse.Number, ch.Number, verse.Number)
		data += fmt.Sprintln()
		data += fmt.Sprintln(`<div class="arabic right">` + verse.Arabic + `</div>`)
		data += fmt.Sprintln()
		data += fmt.Sprintln(`<div class="english">` + verse.Text + `</div>`)
		data += fmt.Sprintln()
		data += fmt.Sprintln(`<div class="dots">...</div>`)
		data += fmt.Sprintln()
	}

	return data
}

func (v *Verse) HTML() string {
	var data string

	data += fmt.Sprintln()
	data += fmt.Sprintln()
	data += fmt.Sprintln()
	data += fmt.Sprintf(`<h4>%d:%d</h4>`, v.Chapter, v.Number)
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

	data += `<div id="contents">`
	for _, ch := range q.Chapters {
		data += fmt.Sprintf(`<div class="chapter"><a href="/quran/%d">%d: %s</a></div>`, ch.Number, ch.Number, ch.English)
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

	// Set local
	for i := 0; i < 114; i++ {
		chapter := i + 1

		f, err = files.ReadFile(fmt.Sprintf("data/%d.json", chapter))
		if err != nil {
			panic(err.Error())
		}
		var data []interface{}
		json.Unmarshal(f, &data)

		english := data[0].(map[string]interface{})["name"].(map[string]interface{})["translated"].(string)
		name := data[0].(map[string]interface{})["name"].(map[string]interface{})["transliterated"].(string)
		data = data[1:]

		var verses []*Verse

		arabicText := arabic[fmt.Sprintf("%d", chapter)].([]interface{})

		// Add Bismillah as first verse for all chapters except 1 and 9
		if chapter != 1 && chapter != 9 {
			verses = append(verses, &Verse{
				Chapter: chapter,
				Number:  1,
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

			verseNum := num
			if chapter != 1 && chapter != 9 {
				verseNum++
			}

			verses = append(verses, &Verse{
				Chapter: chapter,
				Number:  verseNum,
				Text:    ayah.([]interface{})[1].(string),
				Arabic:  ar,
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
