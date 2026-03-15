package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitialize(t *testing.T) {
	s := NewServer("test-server", "1.0.0")

	req := Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
	}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.JSONRPC != "2.0" {
		t.Fatalf("expected jsonrpc 2.0, got %s", resp.JSONRPC)
	}
	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}

	resultBytes, _ := json.Marshal(resp.Result)
	var result InitializeResult
	json.Unmarshal(resultBytes, &result)

	if result.ServerInfo.Name != "test-server" {
		t.Fatalf("expected server name test-server, got %s", result.ServerInfo.Name)
	}
	if result.Capabilities.Tools == nil {
		t.Fatal("expected tools capability")
	}
}

func TestToolsList(t *testing.T) {
	s := NewServer("test-server", "1.0.0")
	s.AddTool("echo", "Echo a message", InputSchema{
		Type: "object",
		Properties: map[string]Property{
			"message": {Type: "string", Description: "The message to echo"},
		},
		Required: []string{"message"},
	}, func(args map[string]interface{}) (string, error) {
		return args["message"].(string), nil
	})

	req := Request{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)

	resultBytes, _ := json.Marshal(resp.Result)
	var result ToolsListResult
	json.Unmarshal(resultBytes, &result)

	if len(result.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(result.Tools))
	}
	if result.Tools[0].Name != "echo" {
		t.Fatalf("expected tool name echo, got %s", result.Tools[0].Name)
	}
}

func TestToolsCall(t *testing.T) {
	s := NewServer("test-server", "1.0.0")
	s.AddTool("echo", "Echo a message", InputSchema{
		Type: "object",
		Properties: map[string]Property{
			"message": {Type: "string", Description: "The message to echo"},
		},
		Required: []string{"message"},
	}, func(args map[string]interface{}) (string, error) {
		return args["message"].(string), nil
	})

	params, _ := json.Marshal(ToolCallParams{
		Name:      "echo",
		Arguments: map[string]interface{}{"message": "hello"},
	})

	req := Request{
		JSONRPC: "2.0",
		ID:      3,
		Method:  "tools/call",
		Params:  params,
	}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)

	resultBytes, _ := json.Marshal(resp.Result)
	var result ToolCallResult
	json.Unmarshal(resultBytes, &result)

	if result.IsError {
		t.Fatal("expected no error")
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content, got %d", len(result.Content))
	}
	if result.Content[0].Text != "hello" {
		t.Fatalf("expected hello, got %s", result.Content[0].Text)
	}
}

func TestToolsCallUnknown(t *testing.T) {
	s := NewServer("test-server", "1.0.0")

	params, _ := json.Marshal(ToolCallParams{
		Name: "unknown",
	})

	req := Request{
		JSONRPC: "2.0",
		ID:      4,
		Method:  "tools/call",
		Params:  params,
	}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error == nil {
		t.Fatal("expected error for unknown tool")
	}
	if resp.Error.Code != -32602 {
		t.Fatalf("expected error code -32602, got %d", resp.Error.Code)
	}
}

func TestPing(t *testing.T) {
	s := NewServer("test-server", "1.0.0")

	req := Request{
		JSONRPC: "2.0",
		ID:      5,
		Method:  "ping",
	}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %v", resp.Error)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	s := NewServer("test-server", "1.0.0")

	r := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestInvalidJSON(t *testing.T) {
	s := NewServer("test-server", "1.0.0")

	r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if resp.Error.Code != -32700 {
		t.Fatalf("expected parse error code -32700, got %d", resp.Error.Code)
	}
}

func TestUnknownMethod(t *testing.T) {
	s := NewServer("test-server", "1.0.0")

	req := Request{
		JSONRPC: "2.0",
		ID:      6,
		Method:  "unknown/method",
	}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error == nil {
		t.Fatal("expected error for unknown method")
	}
	if resp.Error.Code != -32601 {
		t.Fatalf("expected method not found code -32601, got %d", resp.Error.Code)
	}
}

func TestToolCallError(t *testing.T) {
	s := NewServer("test-server", "1.0.0")
	s.AddTool("fail", "Always fails", InputSchema{
		Type: "object",
	}, func(args map[string]interface{}) (string, error) {
		return "", fmt.Errorf("something went wrong")
	})

	params, _ := json.Marshal(ToolCallParams{
		Name: "fail",
	})

	req := Request{
		JSONRPC: "2.0",
		ID:      7,
		Method:  "tools/call",
		Params:  params,
	}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Error != nil {
		t.Fatal("tool errors should be returned in result, not as RPC errors")
	}

	resultBytes, _ := json.Marshal(resp.Result)
	var result ToolCallResult
	json.Unmarshal(resultBytes, &result)

	if !result.IsError {
		t.Fatal("expected IsError to be true")
	}
	if result.Content[0].Text != "something went wrong" {
		t.Fatalf("expected error message, got %s", result.Content[0].Text)
	}
}
