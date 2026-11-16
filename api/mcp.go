package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	Instructions    string       `json:"instructions,omitempty"`
	Capabilities    Capabilities `json:"capabilities"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
}

type Capabilities struct {
	Tools *ToolCapability `json:"tools,omitempty"`
}

type ToolCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ToolAnnotations struct {
	ReadOnly    bool `json:"readOnlyHint"`
	Destructive bool `json:"destructiveHint"`
	OpenWorld   bool `json:"openWorldHint"`
}

type Tool struct {
	Title        string                 `json:"title,omitempty"`
	Annotations  *ToolAnnotations       `json:"annotations,omitempty"`
	OutputSchema map[string]interface{} `json:"outputSchema,omitempty"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	InputSchema  InputSchema            `json:"inputSchema"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Properties  map[string]Property `json:"properties,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Items       *Property           `json:"items,omitempty"`
	Enum        []string            `json:"enum,omitempty"`
	Type        string              `json:"type"`
	Description string              `json:"description,omitempty"`
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
	StructuredContent map[string]interface{} `json:"structuredContent,omitempty"`
	Content           []TextContent          `json:"content"`
	IsError           bool                   `json:"isError,omitempty"`
}

// ToolHandler is a function that takes arguments and returns a JSON string result.
type ToolHandler func(args map[string]interface{}) (string, error)

// ServerOptions configures transport limits and application policy.
type ServerOptions struct {
	Instructions string
	MaxBodyBytes int64
	AllowOrigin  func(origin, host string) bool
}

// Server implements an MCP server with Streamable HTTP transport.
type Server struct {
	maxBodyBytes int64
	allowOrigin  func(string, string) bool
	instructions string
	info         ServerInfo
	tools        []Tool
	handlers     map[string]func(context.Context, map[string]interface{}) (string, error)
	mu           sync.RWMutex
}

// NewServer creates a new MCP server.
func NewServer(name, version string, options ...ServerOptions) *Server {
	config := ServerOptions{MaxBodyBytes: 1 << 20}
	if len(options) > 0 {
		config = options[0]
	}
	if config.MaxBodyBytes <= 0 {
		config.MaxBodyBytes = 1 << 20
	}
	return &Server{
		instructions: config.Instructions,
		maxBodyBytes: config.MaxBodyBytes,
		allowOrigin:  config.AllowOrigin,
		info:         ServerInfo{Name: name, Version: version},
		handlers:     make(map[string]func(context.Context, map[string]interface{}) (string, error)),
	}
}

// AddTool registers a tool with the server.
func (s *Server) AddTool(name, description string, schema InputSchema, handler ToolHandler) {
	s.AddContextTool(name, description, schema, func(_ context.Context, args map[string]interface{}) (string, error) { return handler(args) })
}

func (s *Server) AddContextTool(name, description string, schema InputSchema, handler func(context.Context, map[string]interface{}) (string, error)) {
	s.RegisterTool(Tool{Name: name, Description: description, InputSchema: schema}, handler)
}

// RegisterTool publishes a complete definition together with its handler.
func (s *Server) RegisterTool(tool Tool, handler func(context.Context, map[string]interface{}) (string, error)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools = append(s.tools, tool)
	s.handlers[tool.Name] = handler
}

// ServeHTTP handles MCP requests over Streamable HTTP.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if origin := r.Header.Get("Origin"); origin != "" && (s.allowOrigin == nil || !s.allowOrigin(origin, r.Host)) {
		http.Error(w, "Origin not allowed", http.StatusForbidden)
		return
	}
	w.Header().Set("Allow", "POST")

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

	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, s.maxBodyBytes))
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

	if req.ID == nil {
		if strings.HasPrefix(req.Method, "notifications/") {
			w.WriteHeader(http.StatusAccepted)
		} else {
			writeError(w, nil, -32600, "Requests require an id")
		}
		return
	}
	if version := r.Header.Get("MCP-Protocol-Version"); version != "" && !supportedMCPVersion(version) {
		http.Error(w, "Unsupported MCP protocol version", http.StatusBadRequest)
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
		s.handleToolsCall(r.Context(), w, &req)
	case "ping":
		writeResult(w, req.ID, map[string]interface{}{})
	default:
		writeError(w, req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func (s *Server) handleInitialize(w http.ResponseWriter, req *Request) {
	var params struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	if len(req.Params) > 0 && json.Unmarshal(req.Params, &params) != nil {
		writeError(w, req.ID, -32602, "Invalid initialize params")
		return
	}
	version := params.ProtocolVersion
	if !supportedMCPVersion(version) {
		version = "2025-11-25"
	}
	result := InitializeResult{
		Instructions:    s.instructions,
		ProtocolVersion: version,
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

func (s *Server) handleToolsCall(ctx context.Context, w http.ResponseWriter, req *Request) {
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

	result, err := handler(ctx, params.Arguments)
	if err != nil {
		writeResult(w, req.ID, ToolCallResult{
			Content: []TextContent{{Type: "text", Text: err.Error()}},
			IsError: true,
		})
		return
	}

	var structured map[string]interface{}
	// Keep legacy text content while exposing object results to modern clients.
	_ = json.Unmarshal([]byte(result), &structured)
	writeResult(w, req.ID, ToolCallResult{Content: []TextContent{{Type: "text", Text: result}}, StructuredContent: structured})
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

func supportedMCPVersion(version string) bool {
	switch version {
	case "2025-03-26", "2025-06-18", "2025-11-25":
		return true
	}
	return false
}
