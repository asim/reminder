package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asim/reminder/model"
)

// Check the effective IP by exhausting its allowance across proxy/direct requests.
func TestTrustedClientIP(t *testing.T) {
	t.Setenv("AI_DAILY_LIMIT", "20")
	t.Setenv("AI_IP_DAILY_LIMIT", "1")
	t.Setenv("OLLAMA_LLM_MODEL", "test")
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer provider.Close()
	t.Setenv("OLLAMA_BASE_URL", provider.URL)
	dir := t.TempDir()
	if err := model.ConfigureLimits(dir); err != nil {
		t.Fatal(err)
	}
	middleware, err := ClientMiddleware("")
	if err != nil {
		t.Fatal(err)
	}
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := model.AskSafe(r.Context(), nil, "test")
		if err != nil {
			w.WriteHeader(429)
		}
	}))
	for _, tc := range []struct {
		peer, forwarded string
		want            int
	}{
		{"127.0.0.1:1234", "192.0.2.1", 200},
		{"127.0.0.1:1234", "192.0.2.1", 429},
		{"127.0.0.1:1234", "192.0.2.2", 200},
		{"198.51.100.1:1234", "192.0.2.3", 200},
		{"198.51.100.1:1234", "192.0.2.4", 429},
		{"127.0.0.1:1234", "invalid", 200},
		{"127.0.0.1:1234", "", 429},
	} {
		r := httptest.NewRequest("POST", "/", nil)
		r.RemoteAddr = tc.peer
		r.Header.Set("X-Real-IP", tc.forwarded)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatalf("%+v got %d", tc, w.Code)
		}
	}
	data, _ := os.ReadFile(filepath.Join(dir, "ai-usage.json"))
	var u struct{ Total int }
	json.Unmarshal(data, &u)
	if u.Total != 4 {
		t.Fatalf("usage %d", u.Total)
	}
	if _, err = ClientMiddleware("bad"); err == nil {
		t.Fatal("accepted invalid proxy configuration")
	}
}

func TestMCPPassesClientContext(t *testing.T) {
	t.Setenv("AI_DAILY_LIMIT", "20")
	t.Setenv("AI_IP_DAILY_LIMIT", "1")
	t.Setenv("OLLAMA_LLM_MODEL", "test")
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer provider.Close()
	t.Setenv("OLLAMA_BASE_URL", provider.URL)
	if err := model.ConfigureLimits(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	server := NewServer("test", "1")
	server.AddContextTool("generate", "", InputSchema{}, func(ctx context.Context, _ map[string]interface{}) (string, error) {
		return model.AskSafe(ctx, nil, "test")
	})
	middleware, _ := ClientMiddleware("")
	handler := middleware(server)
	for i, ip := range []string{"192.0.2.1", "192.0.2.1", "192.0.2.2"} {
		r := httptest.NewRequest("POST", "/mcp", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"generate","arguments":{}}}`))
		r.RemoteAddr = "127.0.0.1:1234"
		r.Header.Set("X-Real-IP", ip)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		var response struct {
			Result ToolCallResult `json:"result"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if response.Result.IsError != (i == 1) {
			t.Fatalf("request %d: %s", i, w.Body.String())
		}
	}
}
