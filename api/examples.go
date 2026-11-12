package api

import (
	"bytes"
	_ "embed"
	"encoding/json"
)

//go:embed examples.json
var exampleData []byte

type Example struct {
	Request  string          `json:"request"`
	Response json.RawMessage `json:"response"`
	Note     string          `json:"note,omitempty"`
}

var examples = func() map[string]*Example {
	var data map[string]*Example
	if err := json.Unmarshal(exampleData, &data); err != nil {
		panic(err)
	}
	return data
}()

func (e *Example) ResponseJSON() string {
	var out bytes.Buffer
	if err := json.Indent(&out, e.Response, "", "  "); err != nil {
		return string(e.Response)
	}
	return out.String()
}
