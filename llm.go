package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var systemPromptTpl = template.Must(template.New("system_prompt").Parse(`
You are a helpful assistant with access to knowledge of the Quran, Hadiths and names of Allah. You are tasked with answering questions related to Islam, life and the world.

Answer the question in a very concise manner. Use an unbiased and compassionate tone. Do not repeat text. Don't make anything up. If you are not sure about something, just say that you don't know.
{{- /* Stop here if no context is provided. The rest below is for handling contexts. */ -}}
{{- if . -}}
Answer the question solely based on the provided context. If the search results within the context are not relevant to the question, say I don't know.

Anything between the following 'context' XML blocks is retrieved from the knowledge base, not part of the conversation with the user. The bullet points are ordered by relevance, so the first one is the most relevant.

<context>
    {{- if . -}}
    {{- range $context := .}}
    - {{.}}{{end}}
    {{- end}}
</context>
{{- end -}}

Don't mention the knowledge base, context or search results in your answer.
`))

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model               string        `json:"model"`
	Messages            []chatMessage `json:"messages"`
	MaxCompletionTokens int           `json:"max_completion_tokens,omitempty"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

func askLLM(ctx context.Context, contexts []string, question string) string {
	var apiKey string
	var model string
	var baseURL string

	// Priority: 1. Fanar, 2. Ollama (local)
	fanarKey := os.Getenv("FANAR_API_KEY")
	ollamaModel := os.Getenv("OLLAMA_LLM_MODEL")

	if len(fanarKey) > 0 {
		// Use Fanar API
		apiKey = fanarKey
		baseURL = "https://api.fanar.qa/v1"
		model = "Fanar"
	} else {
		// Use local Ollama (default)
		if ollamaModel == "" {
			ollamaModel = "llama3.2"
		}
		ollamaURL := os.Getenv("OLLAMA_BASE_URL")
		if ollamaURL == "" {
			ollamaURL = "http://localhost:11434/v1"
		}
		apiKey = "ollama"
		baseURL = ollamaURL
		model = ollamaModel
	}

	sb := &strings.Builder{}
	if err := systemPromptTpl.Execute(sb, contexts); err != nil {
		panic(err)
	}

	reqBody := chatRequest{
		Model: model,
		Messages: []chatMessage{
			{Role: "system", Content: sb.String()},
			{Role: "user", Content: "Question: " + question},
		},
		MaxCompletionTokens: 8192,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("LLM API error: %s %s", resp.Status, string(respBody)))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		panic(err)
	}

	if len(chatResp.Choices) == 0 {
		return ""
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content)
}
