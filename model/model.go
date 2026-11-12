package model

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/sashabaranov/go-openai"
)

var (
	llmOnce      sync.Once
	llmReachable bool
)

// Available reports whether summarisation can actually run. Search itself
// never needs an LLM, so when none is configured the UI hides the summarise
// controls rather than offering a button that cannot work.
//
// An explicit key or model is taken at face value. Otherwise the local Ollama
// endpoint is probed once, so that someone running Ollama without setting
// OLLAMA_LLM_MODEL is still offered summaries.
func Available() bool {
	llmOnce.Do(func() { llmReachable = detectLLM() })
	return llmReachable
}

// detectLLM performs the actual check. Kept separate from the cached wrapper
// so it can be exercised directly.
func detectLLM() bool {
	if os.Getenv("ANTHROPIC_API_KEY") != "" ||
		os.Getenv("OPENAI_API_KEY") != "" ||
		os.Getenv("OLLAMA_LLM_MODEL") != "" {
		return true
	}
	return ollamaReachable()
}

// ollamaReachable checks for a local Ollama instance with a short timeout so
// startup is not delayed when there is nothing listening.
func ollamaReachable() bool {
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434/v1"
	}
	baseURL = strings.TrimSuffix(baseURL, "/v1")

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(baseURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// AskSafe calls Ask, turning its panics into errors. The underlying
// clients panic on failure, which in an HTTP handler drops the connection and
// logs a stack trace instead of returning something the caller can act on.
func AskSafe(ctx context.Context, contexts []string, prompt string) (answer string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("llm: %v", r)
			}
		}
	}()

	if !Available() {
		return "", fmt.Errorf("llm: no model configured")
	}
	return Ask(ctx, contexts, prompt), nil
}

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

type internalKey struct{}

func outputLimit(ctx context.Context) int {
	if internal, _ := ctx.Value(internalKey{}).(bool); internal {
		return 8192
	}
	if grounded, _ := ctx.Value(groundedKey{}).(bool); grounded {
		return contentLimits.Output
	}
	return budget.output
}

// AskInternal is reserved for scheduled server-owned generation, never request handlers.
// It does not consume or wait on the public AI quota.
func AskInternal(ctx context.Context, contexts []string, prompt string) string {
	return Ask(context.WithValue(ctx, internalKey{}, true), contexts, prompt)
}

func Ask(ctx context.Context, contexts []string, prompt string) string {
	sb := &strings.Builder{}
	if err := systemPromptTpl.Execute(sb, contexts); err != nil {
		panic(err)
	}
	systemPrompt := sb.String()

	return complete(ctx, systemPrompt, []Message{{Role: "user", Content: prompt}})
}

func complete(ctx context.Context, systemPrompt string, history []Message) string {
	if charge, ok := ctx.Value(tokenAllowanceKey{}).(func(int) error); ok {
		// Conservative preflight units: one per UTF-8 byte plus framing and maximum output.
		// This bounds usage without depending on a provider-specific tokenizer.
		units := len(systemPrompt) + outputLimit(ctx) + 256
		for _, message := range history {
			units += len(message.Content) + 64
		}
		if err := charge(units); err != nil {
			panic(err)
		}
	}
	if budget == nil {
		panic("AI limits are not configured")
	}
	if internal, _ := ctx.Value(internalKey{}).(bool); !internal {
		release, err := budget.acquire(ctx, time.Now())
		if err != nil {
			panic(err)
		}
		defer release()
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	// Prefer explicit local configuration, then OpenAI, then Anthropic.
	if os.Getenv("OLLAMA_LLM_MODEL") != "" || os.Getenv("OPENAI_API_KEY") != "" {
		return askOpenAICompat(ctx, systemPrompt, history)
	}
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	if len(anthropicKey) > 0 {
		return askClaude(ctx, anthropicKey, systemPrompt, history)
	}

	return askOpenAICompat(ctx, systemPrompt, history)
}

func askClaude(ctx context.Context, apiKey, systemPrompt string, history []Message) string {
	client := anthropic.NewClient(option.WithAPIKey(apiKey), option.WithMaxRetries(0))

	var messages []anthropic.MessageParam
	for _, message := range history {
		block := anthropic.NewTextBlock(message.Content)
		if message.Role == "assistant" {
			messages = append(messages, anthropic.NewAssistantMessage(block))
		} else {
			messages = append(messages, anthropic.NewUserMessage(block))
		}
	}
	resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: int64(outputLimit(ctx)),
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: messages,
	})
	if err != nil {
		panic(err)
	}

	var parts []string
	for _, block := range resp.Content {
		if block.Type == "text" {
			parts = append(parts, block.Text)
		}
	}
	return strings.TrimSpace(strings.Join(parts, ""))
}

func askOpenAICompat(ctx context.Context, systemPrompt string, history []Message) string {
	var model string
	var config openai.ClientConfig

	ollamaModel := os.Getenv("OLLAMA_LLM_MODEL")
	openaiKey := os.Getenv("OPENAI_API_KEY")

	if len(ollamaModel) > 0 || (len(openaiKey) == 0) {
		if ollamaModel == "" {
			ollamaModel = "gemma3:1b"
		}
		baseURL := os.Getenv("OLLAMA_BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:11434/v1"
		}
		config = openai.DefaultConfig("ollama")
		config.BaseURL = baseURL
		model = ollamaModel
	} else {
		config = openai.DefaultConfig(openaiKey)
		model = os.Getenv("OPENAI_MODEL")
		if model == "" {
			model = "gpt-6-luna"
		}
	}

	openAIClient := openai.NewClientWithConfig(config)
	messages := []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: systemPrompt}}
	for _, message := range history {
		messages = append(messages, openai.ChatCompletionMessage{Role: message.Role, Content: message.Content})
	}
	res, err := openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:               model,
		Messages:            messages,
		MaxCompletionTokens: outputLimit(ctx),
	})
	if err != nil {
		panic(err)
	}
	if len(res.Choices) == 0 {
		panic("llm: empty response")
	}
	reply := res.Choices[0].Message.Content
	return strings.TrimSpace(reply)
}
