package model

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	openai "github.com/sashabaranov/go-openai"
)

type lookupDefinition struct {
	Name, Description string
	Properties        map[string]any
	Required          []string
}

var lookupDefinitions = []lookupDefinition{
	{Name: "search", Description: "Search the bundled Quran, hadith, tafsir or Names. Use focused terms, inspect the evidence and refine the search if needed. Ranked matches are not exhaustive. The user's selected source restriction is enforced by the server.", Properties: map[string]any{
		"query":  map[string]any{"type": "string", "minLength": 1, "maxLength": 1000},
		"source": map[string]any{"type": "string", "enum": []string{"all", "quran", "hadith", "tafsir", "names", "bukhari", "muslim"}},
	}, Required: []string{"query"}},
	{Name: "passage", Description: "Read an exact Quran verse with neighbouring verses, translations and tafsir. Use this to verify a reference or understand surrounding context.", Properties: map[string]any{
		"chapter": map[string]any{"type": "integer", "minimum": 1, "maximum": 114},
		"verse":   map[string]any{"type": "integer", "minimum": 1, "maximum": 286},
	}, Required: []string{"chapter", "verse"}},
	{Name: "quran_mentions", Description: "Scan the complete bundled Quran for explicit whole-word name/phrase mentions, including English translation aliases (for example Musa and Moses). Returns matching chapters and verse numbers. Counts distinct verses, not word occurrences; excludes unnamed/thematic allusions. Use this to verify lists of chapters or verses rather than assuming a ranked search is complete.", Properties: map[string]any{
		"terms": map[string]any{"type": "array", "minItems": 1, "maxItems": 8, "items": map[string]any{"type": "string", "minLength": 1, "maxLength": 100}},
	}, Required: []string{"terms"}},
	{Name: "arabic", Description: "Look up Arabic dictionary senses and/or a Quran verse's morphological annotations. Supply Arabic text in query, or chapter and verse, or both. Use when linguistic evidence helps explain the passage. Dictionary senses are possibilities; interpret them in context.", Properties: map[string]any{
		"query":   map[string]any{"type": "string", "maxLength": 1000},
		"chapter": map[string]any{"type": "integer", "minimum": 1, "maximum": 114},
		"verse":   map[string]any{"type": "integer", "minimum": 1, "maximum": 286},
	}},
}

func newResearchSession(request Request) researchSession {
	if os.Getenv("OLLAMA_LLM_MODEL") == "" && os.Getenv("OPENAI_API_KEY") == "" && os.Getenv("ANTHROPIC_API_KEY") != "" {
		s := &claudeResearchSession{client: anthropic.NewClient(option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")), option.WithMaxRetries(0))}
		for _, msg := range append(append([]Message{}, request.History...), Message{Role: "user", Content: request.Prompt}) {
			if msg.Role == "assistant" {
				s.messages = append(s.messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)))
			} else {
				s.messages = append(s.messages, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content)))
			}
		}
		return s
	}
	config, name := openAIConfiguration()
	s := &openAIResearchSession{client: openai.NewClientWithConfig(config), model: name}
	for _, msg := range append(append([]Message{}, request.History...), Message{Role: "user", Content: request.Prompt}) {
		s.messages = append(s.messages, openai.ChatCompletionMessage{Role: msg.Role, Content: msg.Content})
	}
	return s
}

// Charge every native request including tool schemas, prior tool calls and results.
// The existing per-client/key/global limits still apply on every provider turn.
func reserveResearchCall(ctx context.Context, request any, used *int) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if budget == nil {
		return nil, fmt.Errorf("AI limits are not configured")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	units := len(data) + outputLimit(ctx) + 256
	if len(data) > contentLimits.Input+contentLimits.Evidence+16384 || *used+units > researchMaxUnits {
		return nil, &LimitError{"Research usage budget reached. Please narrow the question."}
	}
	if charge, ok := ctx.Value(tokenAllowanceKey{}).(func(int) error); ok {
		if err := charge(units); err != nil {
			return nil, err
		}
	}
	release, err := budget.acquire(ctx, time.Now())
	if err != nil {
		return nil, err
	}
	*used += units
	return release, nil
}

type openAIResearchSession struct {
	client   *openai.Client
	model    string
	messages []openai.ChatCompletionMessage
}

func (s *openAIResearchSession) next(ctx context.Context, system string, replies []researchReply, allow bool, used *int) (researchTurn, error) {
	for _, reply := range replies {
		s.messages = append(s.messages, openai.ChatCompletionMessage{Role: "tool", ToolCallID: reply.ID, Content: reply.Text})
	}
	var definitions []openai.Tool
	for _, d := range lookupDefinitions {
		schema := map[string]any{"type": "object", "properties": d.Properties, "additionalProperties": false}
		if len(d.Required) > 0 {
			schema["required"] = d.Required
		}
		definitions = append(definitions, openai.Tool{Type: openai.ToolTypeFunction, Function: &openai.FunctionDefinition{Name: d.Name, Description: d.Description, Parameters: schema}})
	}
	choice := "auto"
	if !allow {
		choice = "none"
	}
	params := openai.ChatCompletionRequest{Model: s.model, Messages: append([]openai.ChatCompletionMessage{{Role: "system", Content: system}}, s.messages...), Tools: definitions, ToolChoice: choice, MaxCompletionTokens: outputLimit(ctx)}
	release, err := reserveResearchCall(ctx, params, used)
	if err != nil {
		return researchTurn{}, err
	}
	defer release()
	response, err := s.client.CreateChatCompletion(ctx, params)
	if err != nil {
		return researchTurn{}, err
	}
	if len(response.Choices) == 0 {
		return researchTurn{}, fmt.Errorf("model returned an empty response")
	}
	answer := response.Choices[0]
	s.messages = append(s.messages, answer.Message)
	turn := researchTurn{Text: strings.TrimSpace(answer.Message.Content), Complete: string(answer.FinishReason) == "stop"}
	for _, call := range answer.Message.ToolCalls {
		turn.Calls = append(turn.Calls, researchCall{ID: call.ID, Name: call.Function.Name, Arguments: call.Function.Arguments})
	}
	return turn, nil
}

type claudeResearchSession struct {
	client   anthropic.Client
	messages []anthropic.MessageParam
}

func (s *claudeResearchSession) next(ctx context.Context, system string, replies []researchReply, allow bool, used *int) (researchTurn, error) {
	if len(replies) > 0 {
		blocks := make([]anthropic.ContentBlockParamUnion, 0, len(replies))
		for _, reply := range replies {
			blocks = append(blocks, anthropic.NewToolResultBlock(reply.ID, reply.Text, reply.IsError))
		}
		s.messages = append(s.messages, anthropic.NewUserMessage(blocks...))
	}
	var definitions []anthropic.ToolUnionParam
	for _, d := range lookupDefinitions {
		definitions = append(definitions, anthropic.ToolUnionParam{OfTool: &anthropic.ToolParam{Name: d.Name, Description: anthropic.String(d.Description), InputSchema: anthropic.ToolInputSchemaParam{Properties: d.Properties, Required: d.Required, ExtraFields: map[string]any{"additionalProperties": false}}}})
	}
	params := anthropic.MessageNewParams{Model: anthropic.ModelClaudeHaiku4_5, MaxTokens: int64(outputLimit(ctx)), System: []anthropic.TextBlockParam{{Text: system}}, Messages: s.messages, Tools: definitions}
	if !allow {
		none := anthropic.NewToolChoiceNoneParam()
		params.ToolChoice = anthropic.ToolChoiceUnionParam{OfNone: &none}
	}
	release, err := reserveResearchCall(ctx, params, used)
	if err != nil {
		return researchTurn{}, err
	}
	defer release()
	response, err := s.client.Messages.New(ctx, params)
	if err != nil {
		return researchTurn{}, err
	}
	// ToParam preserves native tool_use input, IDs, text and any other content blocks.
	s.messages = append(s.messages, response.ToParam())
	turn := researchTurn{Complete: string(response.StopReason) == "end_turn"}
	var parts []string
	for _, block := range response.Content {
		switch block.Type {
		case "text":
			parts = append(parts, block.Text)
		case "tool_use":
			turn.Calls = append(turn.Calls, researchCall{ID: block.ID, Name: block.Name, Arguments: string(block.Input)})
		}
	}
	turn.Text = strings.TrimSpace(strings.Join(parts, "\n"))
	return turn, nil
}
