package api

import (
	"fmt"
	"net/http"
)

type Api struct {
	Endpoints []*Endpoint
}

type Endpoint struct {
	Name        string
	Path        string
	Params      []*Param
	Response    []*Value
	Description string
}

type Param struct {
	Name        string
	Value       string
	Description string
}

type Value struct {
	Type   string
	Params []*Param
}

var Endpoints = []*Endpoint{
	{
		Name:        "Daily verse, hadith and name of Allah",
		Path:        "/api/daily",
		Params:      nil,
		Description: "Returns a verse of the quran, hadith and name of Allah",
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "name", Value: "string", Description: "Name of Allah"},
				{Name: "hadith", Value: "string", Description: "Hadith from Sahih Bukhari"},
				{Name: "verse", Value: "string", Description: "A verse of the Quran"},
				{Name: "links", Value: "map", Description: "Links to relevant content"},
				{Name: "updated", Value: "string", Description: "Time of last update"},
				{Name: "message", Value: "string", Description: "Salam, today is ... (Hijri date)"},
			},
		}},
	},
	{
		Name:        "Daily reminder by date",
		Path:        "/api/daily/{date}",
		Params:      nil,
		Description: "Returns a verse of the quran, hadith and name of Allah for the date",
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "name", Value: "string", Description: "Name of Allah"},
				{Name: "hadith", Value: "string", Description: "Hadith from Sahih Bukhari"},
				{Name: "verse", Value: "string", Description: "A verse of the Quran"},
				{Name: "links", Value: "map", Description: "Links to relevant content"},
				{Name: "updated", Value: "string", Description: "Time of last update"},
				{Name: "message", Value: "string", Description: "Salam, today is ... (Hijri date)"},
			},
		}},
	},
	{
		Name:        "Quran",
		Path:        "/api/quran",
		Params:      nil,
		Response:    []*Value{{Type: "JSON"}},
		Description: "Returns the entire Quran",
	},
	{
		Name:        "List of Quran Chapters",
		Path:        "/api/quran/chapters",
		Params:      nil,
		Description: "Returns a list of Quran chapters",
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "name", Value: "string", Description: "Transliterated name of chapter"},
				{Name: "number", Value: "int", Description: "Number of the chapter"},
				{Name: "english", Value: "string", Description: "English name of chapter"},
				{Name: "verse_count", Value: "int", Description: "Number of verses in chapter"},
			},
		}},
	},
	{
		Name:        "Quran by Chapter",
		Path:        "/api/quran/{chapter}",
		Params:      nil,
		Description: "Returns a chapter of the quran",
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "name", Value: "string", Description: "Name of chapter"},
				{Name: "number", Value: "int", Description: "Number of the chapter"},
				{Name: "verses", Value: "array", Description: "Verses in the chapter"},
				{Name: "english", Value: "string", Description: "Name in english"},
			},
		}},
	},
	{
		Name:        "Quran by Verse",
		Path:        "/api/quran/{chapter}/{verse}",
		Params:      nil,
		Description: "Returns a verse of the quran",
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "chapter", Value: "int", Description: "Chapter of the verse"},
				{Name: "number", Value: "int", Description: "Number of the verse"},
				{Name: "text", Value: "string", Description: "Text of the verse"},
				{Name: "arabic", Value: "string", Description: "Arabic text of the verse"},
				{Name: "words", Value: "array", Description: "Word by word translation"},
			},
		}},
	},
	{
		Name:        "Hadith",
		Path:        "/api/hadith",
		Params:      nil,
		Response:    []*Value{{Type: "JSON"}},
		Description: "Returns the entire Hadith",
	},
	{
		Name:        "Hadith by Book",
		Path:        "/api/hadith/{book}",
		Params:      nil,
		Description: "Returns a book from the hadith",
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "name", Value: "string", Description: "Name of book"},
				{Name: "hadiths", Value: "array", Description: "Hadiths in the book"},
			},
		}},
	},
	{
		Name:        "Names",
		Path:        "/api/names",
		Params:      nil,
		Response:    []*Value{{Type: "JSON"}},
		Description: "Returns the names of Allah",
	},
	{
		Name:        "Search",
		Path:        "/api/search",
		Description: "Get summarised answers via an LLM",
		Params: []*Param{
			{
				Name:        "q",
				Value:       "string",
				Description: "The question to ask",
			},
		},
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "q", Value: "string", Description: "The question asked"},
				{Name: "answer", Value: "string", Description: "Answer to the question"},
				{Name: "references", Value: "array", Description: "A list of references used"},
			},
		}},
	},
	{
		Name:        "Hijri Date (Umm al-Qura)",
		Path:        "/api/hijri/date",
		Params:      nil,
		Description: "Returns today's Hijri date (Umm al-Qura calendar)",
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "date", Value: "string", Description: "Hijri date in DD-MM-YYYY format (Umm al-Qura)"},
				{Name: "display", Value: "string", Description: "Nicely formatted Hijri date for display"},
			},
		}},
	},
	{
		Name:        "Daily verse, hadith and name of Allah (by Date)",
		Path:        "/api/daily",
		Params: []*Param{
			{Name: "date", Value: "string", Description: "(POST only) Date in YYYY-MM-DD format. Optional. If omitted, returns today's entry."},
		},
		Description: "Returns a verse of the Quran, hadith, and name of Allah for the given date (POST) or today (GET). POST with a JSON body containing an optional 'date' field. If the date is not found, returns 404.",
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "name", Value: "string", Description: "Name of Allah"},
				{Name: "hadith", Value: "string", Description: "Hadith from Sahih Bukhari"},
				{Name: "verse", Value: "string", Description: "A verse of the Quran"},
				{Name: "links", Value: "map", Description: "Links to relevant content"},
				{Name: "updated", Value: "string", Description: "Time of last update"},
				{Name: "message", Value: "string", Description: "Salam, today is ... (Hijri date)"},
				{Name: "date", Value: "string", Description: "Gregorian date (YYYY-MM-DD)"},
				{Name: "hijri", Value: "string", Description: "Hijri date (display format)"},
			},
		}},
	},
}

func (a *Api) Markdown() string {
	var data string

	data += "# Endpoints"
	data += fmt.Sprintln()
	data += fmt.Sprintln("A list of API endpoints")
	data += fmt.Sprintln()

	for _, endpoint := range a.Endpoints {
		data += fmt.Sprintln()
		data += "## " + endpoint.Name
		data += fmt.Sprintln()
		data += fmt.Sprintln("___")
		data += fmt.Sprintln("\\")
		data += fmt.Sprintln()
		data += fmt.Sprintln()
		data += fmt.Sprintln(endpoint.Description)
		data += fmt.Sprintln()
		data += fmt.Sprintf("URL: [`%s`](%s)", endpoint.Path, endpoint.Path)
		data += fmt.Sprintln()

		if endpoint.Params != nil {
			data += fmt.Sprintln("#### Request")
			data += fmt.Sprintln()
			data += fmt.Sprintln("Format `JSON`")
			data += fmt.Sprintln()
			for _, param := range endpoint.Params {
				data += fmt.Sprintf("- `%s` - **`%s`** - %s", param.Name, param.Value, param.Description)
				data += fmt.Sprintln()
			}
			data += fmt.Sprintln()
			data += fmt.Sprintln("\\")
			data += fmt.Sprintln()
		}

		if endpoint.Response != nil {
			data += fmt.Sprintln("#### Response")
			data += fmt.Sprintln()
			for _, resp := range endpoint.Response {
				data += fmt.Sprintf("Format `%s`", resp.Type)
				data += fmt.Sprintln()
				for _, param := range resp.Params {
					data += fmt.Sprintf("- `%s` - **`%s`** - %s", param.Name, param.Value, param.Description)
					data += fmt.Sprintln()
				}
			}

			data += fmt.Sprintln()
			data += fmt.Sprintln("\\")
			data += fmt.Sprintln()
		}

		data += fmt.Sprintln()
		data += fmt.Sprintln()
	}

	return data
}

func Load() *Api {
	a := new(Api)
	a.Endpoints = Endpoints
	return a
}

// Register the new endpoint in your API router
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/hijri/date", HijriDateHandler)
}
