package model

import (
	"context"
	"errors"
	"fmt"
	anthropic "github.com/anthropics/anthropic-sdk-go"
	openai "github.com/sashabaranov/go-openai"
	"io"
)

func (s *openAIResearchSession) complete(ctx context.Context, params openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	if !streaming(ctx) {
		return s.client.CreateChatCompletion(ctx, params)
	}
	stream, err := s.client.CreateChatCompletionStream(ctx, params)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}
	defer stream.Close()
	choice := openai.ChatCompletionChoice{Message: openai.ChatCompletionMessage{Role: "assistant"}}
	calls := map[int]*openai.ToolCall{}
	for {
		part, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return openai.ChatCompletionResponse{}, err
		}
		for _, delta := range part.Choices {
			if delta.Index != 0 {
				continue
			}
			choice.Message.Content += delta.Delta.Content
			if delta.Delta.Content != "" {
				Emit(ctx, Event{Type: "delta", Text: delta.Delta.Content})
			}
			if delta.FinishReason != "" {
				choice.FinishReason = delta.FinishReason
			}
			for _, fragment := range delta.Delta.ToolCalls {
				if fragment.Index == nil || *fragment.Index < 0 || *fragment.Index >= 64 {
					return openai.ChatCompletionResponse{}, fmt.Errorf("invalid streamed tool index")
				}
				index := *fragment.Index
				call := calls[index]
				if call == nil {
					call = &openai.ToolCall{}
					calls[index] = call
				}
				call.ID += fragment.ID
				if fragment.Type != "" {
					call.Type = fragment.Type
				}
				call.Function.Name += fragment.Function.Name
				call.Function.Arguments += fragment.Function.Arguments
			}
		}
	}
	if choice.FinishReason == "" {
		return openai.ChatCompletionResponse{}, fmt.Errorf("incomplete model stream")
	}
	for i := 0; i < len(calls); i++ {
		call := calls[i]
		if call == nil {
			return openai.ChatCompletionResponse{}, fmt.Errorf("missing streamed tool call")
		}
		choice.Message.ToolCalls = append(choice.Message.ToolCalls, *call)
	}
	return openai.ChatCompletionResponse{Choices: []openai.ChatCompletionChoice{choice}}, nil
}
func (s *claudeResearchSession) complete(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error) {
	if !streaming(ctx) {
		return s.client.Messages.New(ctx, params)
	}
	stream := s.client.Messages.NewStreaming(ctx, params)
	defer stream.Close()
	message := &anthropic.Message{}
	stopped := false
	for stream.Next() {
		event := stream.Current()
		if err := message.Accumulate(event); err != nil {
			return nil, err
		}
		if event.Type == "content_block_delta" && event.Delta.Type == "text_delta" {
			Emit(ctx, Event{Type: "delta", Text: event.Delta.Text})
		}
		if event.Type == "message_stop" {
			stopped = true
		}
	}
	if err := stream.Err(); err != nil {
		return nil, err
	}
	if !stopped {
		return nil, fmt.Errorf("incomplete model stream")
	}
	return message, nil
}
