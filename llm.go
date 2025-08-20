package main

import (
	"context"
	"os"
	"strings"
	"text/template"

	"github.com/sashabaranov/go-openai"
)

// There are many different ways to provide the context to the LLM.
// You can pass each context as user message, or the list as one user message,
// or pass it in the system prompt. The system prompt itself also has a big impact
// on how well the LLM handles the context, especially for LLMs with < 7B parameters.
// The prompt engineering is up to you, it's out of scope for the vector database.
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
	var apiKey string
	var model string
	var config openai.ClientConfig

	fanarKey := os.Getenv("FANAR_API_KEY")
	if len(fanarKey) > 0 {
		apiKey = fanarKey
		config = openai.DefaultConfig(apiKey)
		config.BaseURL = "https://api.fanar.qa/v1"
		model = "Fanar"
	} else {
		apiKey = os.Getenv("OPENAI_API_KEY")
		config = openai.DefaultConfig(apiKey)
		model = openai.GPT4oMini
	}

	openAIClient := openai.NewClientWithConfig(config)
	sb := &strings.Builder{}
	err := systemPromptTpl.Execute(sb, contexts)
	if err != nil {
		panic(err)
	}
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: sb.String(),
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
	reply = strings.TrimSpace(reply)

	return reply
}
