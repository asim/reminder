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
		Name:        "Quran",
		Path:        "/api/quran",
		Params:      nil,
		Response:    []*Value{{Type: "JSON"}},
		Description: "Returns the entire Quran",
	},
	{
		Name:        "Hadith",
		Path:        "/api/hadith",
		Params:      nil,
		Response:    []*Value{{Type: "JSON"}},
		Description: "Returns the entire Hadith",
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
