package main

import (
	"context"
	"sync"
	"testing"
)

// clearLLMEnv removes every LLM setting and points Ollama at a closed port, so
// detection sees a machine with nothing configured.
func clearLLMEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("OLLAMA_LLM_MODEL", "")
	t.Setenv("OLLAMA_BASE_URL", "http://127.0.0.1:1/v1")
}

func TestDetectLLMWithKey(t *testing.T) {
	clearLLMEnv(t)
	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-test")

	if !detectLLM() {
		t.Error("expected an LLM to be reported available when a key is set")
	}
}

func TestDetectLLMWithOllamaModel(t *testing.T) {
	clearLLMEnv(t)
	t.Setenv("OLLAMA_LLM_MODEL", "gemma3:1b")

	if !detectLLM() {
		t.Error("an explicitly configured Ollama model should count as available")
	}
}

// TestDetectLLMWithNothing is the case this PR cares about: a plain install
// with no keys. Search must work, but summarisation must not be offered.
func TestDetectLLMWithNothing(t *testing.T) {
	clearLLMEnv(t)

	if detectLLM() {
		t.Error("no keys and no reachable Ollama should report unavailable")
	}
}

// TestAskLLMSafeErrorsInsteadOfPanicking guards the behaviour that made the
// summarise button dangerous: the underlying clients panic on failure, which
// in an HTTP handler drops the connection and logs a stack trace.
func TestAskLLMSafeErrorsInsteadOfPanicking(t *testing.T) {
	clearLLMEnv(t)

	// Reset the cached detection so this test sees the cleared environment.
	llmOnce = sync.Once{}
	llmReachable = false

	answer, err := askLLMSafe(context.Background(), nil, "what is patience")
	if err == nil {
		t.Fatal("expected an error when no model is configured, got nil")
	}
	if answer != "" {
		t.Errorf("expected no answer on failure, got %q", answer)
	}
}
