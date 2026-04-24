package main

import (
	"context"
	"os"
	"strings"
	"text/template"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/sashabaranov/go-openai"
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

func askLLM(ctx context.Context, contexts []string, question string) string {
	sb := &strings.Builder{}
	if err := systemPromptTpl.Execute(sb, contexts); err != nil {
		panic(err)
	}
	systemPrompt := sb.String()

	// Priority: 1. Claude (Anthropic), 2. Fanar, 3. Ollama, 4. OpenAI
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	if len(anthropicKey) > 0 {
		return askClaude(ctx, anthropicKey, systemPrompt, question)
	}

	return askOpenAICompat(ctx, systemPrompt, question)
}

func askClaude(ctx context.Context, apiKey, systemPrompt, question string) string {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 8192,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock("Question: " + question),
			),
		},
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

func askOpenAICompat(ctx context.Context, systemPrompt, question string) string {
	var model string
	var config openai.ClientConfig

	fanarKey := os.Getenv("FANAR_API_KEY")
	ollamaModel := os.Getenv("OLLAMA_LLM_MODEL")
	openaiKey := os.Getenv("OPENAI_API_KEY")

	if len(fanarKey) > 0 {
		config = openai.DefaultConfig(fanarKey)
		config.BaseURL = "https://api.fanar.qa/v1"
		model = "Fanar"
	} else if len(ollamaModel) > 0 || (len(openaiKey) == 0) {
		if ollamaModel == "" {
			ollamaModel = "llama3.2"
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
		model = openai.GPT4oMini
	}

	openAIClient := openai.NewClientWithConfig(config)
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: "Question: " + question,
		},
	}
	res, err := openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:               model,
		Messages:            messages,
		MaxCompletionTokens: 8192,
	})
	if err != nil {
		panic(err)
	}
	reply := res.Choices[0].Message.Content
	return strings.TrimSpace(reply)
}
