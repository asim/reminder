package api

import (
	"bytes"
	"encoding/json"
	"github.com/asim/reminder/model"
	"io"
	"net/http"
	"strings"
)

// Mock the planning turn independently of the existing synthesis assertions.
func replyToLookupPlan(w http.ResponseWriter, r *http.Request) bool {
	data, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(data))
	var body struct {
		Messages []model.Message `json:"messages"`
	}
	_ = json.Unmarshal(data, &body)
	if len(body.Messages) == 0 || !strings.HasPrefix(body.Messages[0].Content, "You select read-only lookups") {
		return false
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"choices":[{"finish_reason":"stop","message":{"role":"assistant","content":"{\"lookups\":[]}"}}]}`))
	return true
}
