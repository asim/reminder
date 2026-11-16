package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/asim/reminder/model"
	"github.com/asim/reminder/names"
	"github.com/asim/reminder/quran"
)

func TestPluginExplainSharesReaderCache(t *testing.T) {
	var calls atomic.Int32
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"finish_reason":"stop","message":{"role":"assistant","content":"Grounded explanation [1]."}}]}`))
	}))
	defer provider.Close()
	t.Setenv("OLLAMA_LLM_MODEL", "test")
	t.Setenv("OLLAMA_BASE_URL", provider.URL+"/v1")
	t.Setenv("AI_DAILY_LIMIT", "20")
	t.Setenv("AI_IP_DAILY_LIMIT", "0")
	if err := model.ConfigureLimits(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	cache, err := openExplainCache(t.TempDir(), "test")
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	ready := make(chan bool)
	close(ready)
	q := quran.Load()
	n := &names.Names{}
	handler := ExplainHandler(nil, modelIndex{}, n, q, nil, ready, cache)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("POST", "/api/explain", strings.NewReader(`{"context":{"chapter":1,"verse":2}}`)))
	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
	server := NewMCP(q, n, nil, nil, modelIndex{}, ready, cache)
	call := func(arguments map[string]interface{}) ToolCallResult {
		payload, _ := json.Marshal(map[string]interface{}{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": map[string]interface{}{"name": "explain_quran", "arguments": arguments}})
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, httptest.NewRequest("POST", "/mcp", bytes.NewReader(payload)))
		var response struct {
			Result ToolCallResult
			Error  *Error
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if response.Error != nil {
			t.Fatal(response.Error)
		}
		return response.Result
	}
	args := map[string]interface{}{"context": map[string]interface{}{"chapter": 1, "verse": 2}}
	result := call(args)
	if result.IsError || result.StructuredContent["answer"] != "Grounded explanation [1]." || calls.Load() != 1 {
		t.Fatalf("MCP did not reuse reader cache: %+v calls=%d", result, calls.Load())
	}
	args["question"] = "What does it mean?"
	args["history"] = []model.Message{{Role: "user", Content: "Explain this verse"}, {Role: "assistant", Content: "Grounded explanation [1]."}}
	if result := call(args); result.IsError {
		t.Fatal(result)
	}
	if calls.Load() != 2 {
		t.Fatal("follow-up did not generate")
	}
	args["context"] = map[string]interface{}{"chapter": 1.5}
	if !call(args).IsError || calls.Load() != 2 {
		t.Fatal("fractional reference accepted")
	}
}

func TestPluginTransportNegotiationAndNotifications(t *testing.T) {
	server := NewServer("test", "1")
	var calls atomic.Int32
	server.AddContextTool("test", "test", InputSchema{Type: "object"}, func(context.Context, map[string]interface{}) (string, error) { calls.Add(1); return `{}`, nil })
	for _, version := range []string{"2025-03-26", "2025-06-18", "2025-11-25"} {
		body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"` + version + `"}}`
		w := httptest.NewRecorder()
		server.ServeHTTP(w, httptest.NewRequest("POST", "/mcp", strings.NewReader(body)))
		if !strings.Contains(w.Body.String(), `"protocolVersion":"`+version+`"`) {
			t.Fatal(w.Body.String())
		}
	}
	for _, body := range []string{`{"jsonrpc":"2.0","method":"notifications/initialized"}`, `{"jsonrpc":"2.0","method":"notifications/cancelled","params":{"requestId":1}}`} {
		w := httptest.NewRecorder()
		server.ServeHTTP(w, httptest.NewRequest("POST", "/mcp", strings.NewReader(body)))
		if w.Code != 202 || w.Body.Len() != 0 {
			t.Fatal("notification returned response")
		}
	}
	for _, origin := range []string{"https://untrusted.example", "null"} {
		r := httptest.NewRequest("POST", "/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"test"}}`))
		r.Header.Set("Origin", origin)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, r)
		if w.Code != 403 {
			t.Fatal("untrusted origin accepted")
		}
	}
	if calls.Load() != 0 {
		t.Fatal("rejected request executed tool")
	}
}
