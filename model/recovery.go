package model

import (
	"context"
	"errors"
	"fmt"
	anthropic "github.com/anthropics/anthropic-sdk-go"
	openai "github.com/sashabaranov/go-openai"
	"strings"
)

type countedSession struct {
	researchSession
	turns int
}

func (s *countedSession) next(ctx context.Context, prompt string, replies []researchReply, allow bool, used *int) (researchTurn, error) {
	s.turns++
	return s.researchSession.next(ctx, prompt, replies, allow, used)
}

func providerStatus(err error) int {
	var a *anthropic.Error
	if errors.As(err, &a) {
		return a.StatusCode
	}
	var o *openai.APIError
	if errors.As(err, &o) {
		return o.HTTPStatusCode
	}
	var r *openai.RequestError
	if errors.As(err, &r) {
		return r.HTTPStatusCode
	}
	return 0
}
func recoverableResearchError(err error) bool {
	if err == nil {
		return false
	}
	// Do not retry authentication, missing models, rate limits, timeouts or outages.
	if status := providerStatus(err); status != 0 {
		return status == 400 || status == 422
	}
	switch err.Error() {
	case "model returned an empty response", "model returned an empty answer", "incomplete model stream", "model returned an invalid source citation":
		return true
	}
	return false
}

// FailureMessage exposes categories, never provider bodies, prompts or credentials.
func FailureMessage(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "The AI request timed out. Please try again."
	}
	if errors.Is(err, context.Canceled) {
		return "The AI request was interrupted. Please try again."
	}
	switch status := providerStatus(err); status {
	case 401, 403:
		return "The AI provider rejected the server credentials. The service needs attention."
	case 404:
		return "The configured AI model or endpoint was not found. The service needs attention."
	case 429:
		return "The AI provider is rate-limiting requests. Please try again shortly."
	case 400, 422:
		return fmt.Sprintf("The AI provider rejected the request (HTTP %d). The service needs attention.", status)
	default:
		if status >= 500 {
			return "The AI provider is temporarily unavailable. Please try again shortly."
		}
	}
	if strings.Contains(err.Error(), "invalid source citation") {
		return "The generated answer failed its source checks. Please try again."
	}
	if strings.Contains(err.Error(), "empty") || strings.Contains(err.Error(), "incomplete model stream") {
		return "The AI provider returned no complete answer. Please try again."
	}
	return "The explanation could not be completed. Please try again."
}
