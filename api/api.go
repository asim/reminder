package api

import (
	"fmt"
)

type Api struct {
	Endpoints []*Endpoint
}

type Endpoint struct {
	Name        string
	Path        string
	Method      string
	Example     *Example
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
		Name:        "Latest reminder (updated hourly)",
		Path:        "/api/latest",
		Method:      "GET",
		Params:      nil,
		Description: "Returns the latest verse, hadith and name of Allah (updated every hour)",
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "name", Value: "string", Description: "Name of Allah"},
				{Name: "hadith", Value: "string", Description: "Hadith from Sahih Bukhari"},
				{Name: "verse", Value: "string", Description: "A verse of the Quran"},
				{Name: "links", Value: "map", Description: "Links to relevant content"},
				{Name: "updated", Value: "string", Description: "Time of last update"},
				{Name: "message", Value: "string", Description: "Bismillah message"},
			},
		}},
	},
	{
		Name:        "Daily verse, hadith and name of Allah",
		Path:        "/api/daily",
		Method:      "GET",
		Params:      nil,
		Description: "Returns today's daily reminder (updated once per day at midnight)",
		Response: []*Value{{
			Type: "JSON",
			Params: []*Param{
				{Name: "name", Value: "string", Description: "Name of Allah"},
				{Name: "hadith", Value: "string", Description: "Hadith from Sahih Bukhari"},
				{Name: "verse", Value: "string", Description: "A verse of the Quran"},
				{Name: "links", Value: "map", Description: "Links to relevant content"},
				{Name: "updated", Value: "string", Description: "Time of last update"},
				{Name: "message", Value: "string", Description: "Salam, today is ... (Hijri date)"},
				{Name: "hijri", Value: "string", Description: "Hijri date"},
				{Name: "date", Value: "string", Description: "Gregorian date (YYYY-MM-DD)"},
			},
		}},
	},
	{
		Name:        "Daily reminder by date",
		Path:        "/api/daily/{date}",
		Method:      "GET",
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
		Method:      "GET",
		Params:      nil,
		Response:    []*Value{{Type: "JSON"}},
		Description: "Returns the entire Quran",
	},
	{
		Name:        "List of Quran Chapters",
		Path:        "/api/quran/chapters",
		Method:      "GET",
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
		Method:      "GET",
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
		Method:      "GET",
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
		Method:      "GET",
		Params:      nil,
		Response:    []*Value{{Type: "JSON"}},
		Description: "Returns the entire Hadith",
	},
	{
		Name:        "Hadith by Book",
		Path:        "/api/hadith/{book}",
		Method:      "GET",
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
		Method:      "GET",
		Params:      nil,
		Response:    []*Value{{Type: "JSON"}},
		Description: "Returns the names of Allah",
	},
	{
		Name:        "Model",
		Path:        "/api/model/ask",
		Method:      "POST",
		Description: `POST. Requires authorisation: Authorization: Bearer <key>. Access to reminder.dev is by approval; request a key through a GitHub issue at github.com/asim/reminder. Answers from Quran, hadith and Names sources, using optional caller-managed history. History contains alternating user/assistant messages with role and content fields. See docs/model.md for CLI key management and limits.`,
		Params: []*Param{
			{Name: "prompt", Value: "string", Description: "Required prompt (up to 32 KiB UTF-8 by default)"},
			{Name: "history", Value: "array", Description: "Optional alternating user/assistant messages; prompt plus history up to 128 KiB UTF-8 by default"},
		},
		Response: []*Value{{Type: "JSON", Params: []*Param{
			{Name: "answer", Value: "string", Description: "Plain text / Markdown response with numbered source citations"},
			{Name: "references", Value: "array", Description: "Retrieved source records, in citation order"},
		}}},
	},
	{
		Name:        "MCP",
		Path:        "/mcp",
		Method:      "POST",
		Description: "Connect MCP clients to https://reminder.dev/mcp using HTTP POST (JSON-RPC 2.0) for Quran, hadith, names, reminders and search tools. Search summaries require summarise: true (default: false).",
		Response:    []*Value{{Type: "JSON-RPC 2.0"}},
	},
	{
		Name:        "Search",
		Path:        "/api/search",
		Method:      "POST",
		Description: "Search the Quran, Hadith and Names of Allah, with an optional AI summary",
		Params: []*Param{
			{Name: "summarise", Value: "boolean", Description: "Generate an AI summary only when true (default: false)"},
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
		data += fmt.Sprintln()
		data += fmt.Sprintln(endpoint.Description)
		data += fmt.Sprintln()
		data += fmt.Sprintf("**%s URL:** [`%s`](%s)", endpoint.Method, endpoint.Path, endpoint.Path)
		data += fmt.Sprintln()
		data += fmt.Sprintln()

		if endpoint.Params != nil {
			data += fmt.Sprintln("### Request")
			data += fmt.Sprintln()
			data += fmt.Sprintln("Format: `JSON`")
			data += fmt.Sprintln()
			for _, param := range endpoint.Params {
				data += fmt.Sprintf("- `%s` - **`%s`** - %s", param.Name, param.Value, param.Description)
				data += fmt.Sprintln()
			}
			data += fmt.Sprintln()
		}

		if endpoint.Response != nil {
			data += fmt.Sprintln("### Response")
			data += fmt.Sprintln()
			for _, resp := range endpoint.Response {
				data += fmt.Sprintf("Format: `%s`", resp.Type)
				data += fmt.Sprintln()
				for _, param := range resp.Params {
					data += fmt.Sprintf("- `%s` - **`%s`** - %s", param.Name, param.Value, param.Description)
					data += fmt.Sprintln()
				}
			}
			data += fmt.Sprintln()
		}

		if endpoint.Example != nil {
			data += "### Example request\n\n```sh\n" + endpoint.Example.Request + "\n```\n\n"
			data += "### Example response\n\n"
			if endpoint.Example.Note != "" {
				data += endpoint.Example.Note + "\n\n"
			}
			data += "```json\n" + endpoint.Example.ResponseJSON() + "\n```\n\n"
		}
	}

	return data
}

func (a *Api) HTML(renderFunc func(string) string) string {
	markdown := renderFunc(a.Markdown())
	content := `<div class="api-content">` + markdown + `</div>`
	content += `<style>
.api-content h2 { font-size: 1.5rem; font-weight: 700; margin-top: 2rem; margin-bottom: 1rem; }
.api-content h3 { font-size: 1.25rem; font-weight: 600; margin-top: 1.5rem; margin-bottom: 0.75rem; }
.api-content p { margin-bottom: 1rem; }
.api-content ul { margin-bottom: 1rem; }
.api-content pre { max-width: 100%; overflow-x: auto; background-color: #f3f4f6; padding: 1rem; margin-bottom: 1rem; border-radius: 0.25rem; direction: ltr; text-align: left; }
.api-content pre code { padding: 0; }
.api-content code { background-color: #f3f4f6; padding: 0.125rem 0.25rem; border-radius: 0.25rem; }
</style>`
	return content
}

func Load() *Api {
	a := new(Api)
	for _, endpoint := range Endpoints {
		copy := *endpoint
		copy.Example = examples[copy.Method+" "+copy.Path]
		a.Endpoints = append(a.Endpoints, &copy)
	}
	return a
}
