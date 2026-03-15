package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// JSON-RPC 2.0 types

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP protocol types

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    Capabilities `json:"capabilities"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
}

type Capabilities struct {
	Tools *ToolCapability `json:"tools,omitempty"`
}

type ToolCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ToolCallResult struct {
	Content []TextContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ToolHandler is a function that takes arguments and returns a JSON string result.
type ToolHandler func(args map[string]interface{}) (string, error)

// Server implements an MCP server with Streamable HTTP transport.
type Server struct {
	info     ServerInfo
	tools    []Tool
	handlers map[string]ToolHandler
	mu       sync.RWMutex
}

// NewServer creates a new MCP server.
func NewServer(name, version string) *Server {
	return &Server{
		info:     ServerInfo{Name: name, Version: version},
		handlers: make(map[string]ToolHandler),
	}
}

// AddTool registers a tool with the server.
func (s *Server) AddTool(name, description string, schema InputSchema, handler ToolHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tools = append(s.tools, Tool{
		Name:        name,
		Description: description,
		InputSchema: schema,
	})
	s.handlers[name] = handler
}

// ServeHTTP handles MCP requests over Streamable HTTP.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// GET requests can be used by clients for SSE, but we don't support streaming
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{
			JSONRPC: "2.0",
			Error: &Error{
				Code:    -32600,
				Message: "Use POST to send JSON-RPC requests",
			},
		})
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, nil, -32700, "Parse error")
		return
	}

	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, nil, -32700, "Parse error")
		return
	}

	if req.JSONRPC != "2.0" {
		writeError(w, req.ID, -32600, "Invalid Request: jsonrpc must be 2.0")
		return
	}

	switch req.Method {
	case "initialize":
		s.handleInitialize(w, &req)
	case "notifications/initialized":
		// Client notification, no response needed
		w.WriteHeader(http.StatusAccepted)
	case "tools/list":
		s.handleToolsList(w, &req)
	case "tools/call":
		s.handleToolsCall(w, &req)
	case "ping":
		writeResult(w, req.ID, map[string]interface{}{})
	default:
		writeError(w, req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func (s *Server) handleInitialize(w http.ResponseWriter, req *Request) {
	result := InitializeResult{
		ProtocolVersion: "2025-03-26",
		Capabilities: Capabilities{
			Tools: &ToolCapability{},
		},
		ServerInfo: s.info,
	}
	writeResult(w, req.ID, result)
}

func (s *Server) handleToolsList(w http.ResponseWriter, req *Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := ToolsListResult{
		Tools: s.tools,
	}
	writeResult(w, req.ID, result)
}

func (s *Server) handleToolsCall(w http.ResponseWriter, req *Request) {
	var params ToolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		writeError(w, req.ID, -32602, "Invalid params")
		return
	}

	s.mu.RLock()
	handler, ok := s.handlers[params.Name]
	s.mu.RUnlock()

	if !ok {
		writeError(w, req.ID, -32602, fmt.Sprintf("Unknown tool: %s", params.Name))
		return
	}

	result, err := handler(params.Arguments)
	if err != nil {
		writeResult(w, req.ID, ToolCallResult{
			Content: []TextContent{{Type: "text", Text: err.Error()}},
			IsError: true,
		})
		return
	}

	writeResult(w, req.ID, ToolCallResult{
		Content: []TextContent{{Type: "text", Text: result}},
	})
}

func writeResult(w http.ResponseWriter, id interface{}, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
}

func writeError(w http.ResponseWriter, id interface{}, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	})
}
