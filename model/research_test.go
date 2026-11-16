package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func configureResearchTest(t *testing.T, provider, endpoint string) {
	t.Helper()
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OLLAMA_LLM_MODEL", "")
	if provider == "anthropic" {
		t.Setenv("ANTHROPIC_API_KEY", "test")
		t.Setenv("ANTHROPIC_BASE_URL", endpoint)
	} else {
		t.Setenv("OLLAMA_LLM_MODEL", "test")
		t.Setenv("OLLAMA_BASE_URL", endpoint+"/v1")
	}
	t.Setenv("AI_DAILY_LIMIT", "100")
	t.Setenv("AI_IP_DAILY_LIMIT", "0")
	llmOnce = sync.Once{}
	if err := ConfigureLimits(t.TempDir()); err != nil {
		t.Fatal(err)
	}
}

func nativeResponse(w http.ResponseWriter, provider, text string, calls []researchCall) {
	w.Header().Set("Content-Type", "application/json")
	if provider == "anthropic" {
		content := []any{}
		if text != "" {
			content = append(content, map[string]any{"type": "text", "text": text})
		}
		for _, call := range calls {
			content = append(content, map[string]any{"type": "tool_use", "id": call.ID, "name": call.Name, "input": json.RawMessage(call.Arguments)})
		}
		stop := "end_turn"
		if len(calls) > 0 {
			stop = "tool_use"
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "msg_test", "type": "message", "role": "assistant", "model": "test", "content": content, "stop_reason": stop, "usage": map[string]int{"input_tokens": 10, "output_tokens": 10}})
	} else {
		message := map[string]any{"role": "assistant", "content": text}
		if len(calls) > 0 {
			native := []any{}
			for _, call := range calls {
				native = append(native, map[string]any{"id": call.ID, "type": "function", "function": map[string]string{"name": call.Name, "arguments": call.Arguments}})
			}
			message["tool_calls"] = native
		}
		stop := "stop"
		if len(calls) > 0 {
			stop = "tool_calls"
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"finish_reason": stop, "message": message}}})
	}
}

// Exercise the actual SDK request/response wire formats for both providers,
// including multiple calls in one turn and more lookups than the old four-call plan.
func TestNativeResearchFollowsUpAcrossProviders(t *testing.T) {
	for _, provider := range []string{"openai", "anthropic"} {
		t.Run(provider, func(t *testing.T) {
			turns, lookups := 0, 0
			charged := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				turns++
				data, _ := io.ReadAll(r.Body)
				var body map[string]any
				if err := json.Unmarshal(data, &body); err != nil {
					t.Error(err)
				}
				tools, ok := body["tools"].([]any)
				if !ok || len(tools) != 4 {
					t.Errorf("missing native tool schemas: %s", data)
				}
				for _, definition := range tools {
					tool := definition.(map[string]any)
					var schema map[string]any
					if provider == "anthropic" {
						schema = tool["input_schema"].(map[string]any)
					} else {
						schema = tool["function"].(map[string]any)["parameters"].(map[string]any)
					}
					if schema["type"] != "object" {
						t.Error("tool schema must be an object")
					}
					if required, present := schema["required"]; present {
						if _, ok := required.([]any); !ok {
							t.Error("required must be an array, never null")
						}
					}
				}
				if turns > 1 {
					if !strings.Contains(string(data), "[2]") || !strings.Contains(string(data), "FIRST LOOKUP") {
						t.Error("tool results were not supplied on follow-up")
					}
					if provider == "anthropic" {
						for _, m := range body["messages"].([]any) {
							msg := m.(map[string]any)
							if msg["role"] != "assistant" {
								continue
							}
							for _, c := range msg["content"].([]any) {
								block := c.(map[string]any)
								if block["type"] == "tool_use" {
									if _, ok := block["input"].(map[string]any); !ok {
										t.Error("tool_use.input was not preserved")
									}
								}
							}
						}
						if !strings.Contains(string(data), `"tool_use_id":"one"`) {
							t.Error("missing native tool result ID")
						}
					} else if !strings.Contains(string(data), `"tool_call_id":"one"`) {
						t.Error("missing native tool result ID")
					}
				}
				switch turns {
				case 1:
					nativeResponse(w, provider, "", []researchCall{{"one", "search", `{"query":"Musa"}`}, {"two", "passage", `{"chapter":10,"verse":106}`}})
				case 2:
					nativeResponse(w, provider, "", []researchCall{{"three", "search", `{"query":"follow up from SECOND LOOKUP"}`}, {"four", "arabic", `{"chapter":10,"verse":106}`}})
				case 3:
					if !strings.Contains(string(data), "FOLLOWUP LOOKUP") {
						t.Error("second-round results missing")
					}
					nativeResponse(w, provider, "", []researchCall{{"five", "quran_mentions", `{"terms":["Musa","Moses"]}`}})
				default:
					nativeResponse(w, provider, "Explanation [2], context [4], inventory [6].", nil)
				}
			}))
			defer server.Close()
			configureResearchTest(t, provider, server.URL)
			ctx := WithTokenAllowance(context.Background(), func(units int) error {
				if units <= 0 {
					t.Error("empty usage charge")
				}
				charged++
				return nil
			})
			result, err := Research(ctx, Request{Prompt: "Which chapters mention Musa, and what does this verse mean?"}, func(string) ([]*Reference, error) {
				return []*Reference{{Text: "INITIAL", Metadata: map[string]string{"source": "quran"}}}, nil
			}, func(call Lookup) ([]*Reference, error) {
				lookups++
				texts := []string{"FIRST LOOKUP", "SECOND LOOKUP", "FOLLOWUP LOOKUP", "ARABIC LOOKUP", "INVENTORY"}
				return []*Reference{{Text: texts[lookups-1], Metadata: map[string]string{"source": "quran"}}}, nil
			})
			if err != nil || turns != 4 || lookups != 5 || result == nil || len(result.References) != 6 || !result.Complete {
				t.Fatalf("native loop: %+v %v turns=%d lookups=%d", result, err, turns, lookups)
			}
			if charged != turns {
				t.Fatalf("not every provider turn was charged: %d vs %d", charged, turns)
			}
			if result.References[1].Text != "FIRST LOOKUP" || result.References[5].Text != "INVENTORY" {
				t.Fatal("citation labels shifted")
			}
		})
	}
}

type researchSessionFunc func(context.Context, string, []researchReply, bool, *int) (researchTurn, error)

func (f researchSessionFunc) next(ctx context.Context, s string, r []researchReply, a bool, u *int) (researchTurn, error) {
	return f(ctx, s, r, a, u)
}

func TestResearchCanCorrectToolErrors(t *testing.T) {
	turns, executed := 0, 0
	session := researchSessionFunc(func(_ context.Context, _ string, replies []researchReply, _ bool, _ *int) (researchTurn, error) {
		turns++
		switch turns {
		case 1:
			return researchTurn{Calls: []researchCall{{"bad", "passage", `not JSON`}}}, nil
		case 2:
			if len(replies) != 1 || !replies[0].IsError {
				t.Fatal("bad arguments not returned as tool error")
			}
			return researchTurn{Calls: []researchCall{{"missing", "passage", `{"chapter":1,"verse":99}`}}}, nil
		case 3:
			if len(replies) != 1 || !replies[0].IsError || !strings.Contains(replies[0].Text, "not found") {
				t.Fatal("lookup failure not visible")
			}
			return researchTurn{Calls: []researchCall{{"corrected", "passage", `{"chapter":1,"verse":2}`}}}, nil
		default:
			return researchTurn{Text: "Meaning [1].", Complete: true}, nil
		}
	})
	result, err := runResearch(context.Background(), groundedPrompt, &evidenceSet{}, session, func(call Lookup) ([]*Reference, error) {
		executed++
		if call.Verse == 99 {
			return nil, errors.New("verse not found")
		}
		return []*Reference{{Text: "Correct verse", Metadata: map[string]string{"source": "quran"}}}, nil
	})
	if err != nil || result == nil || executed != 2 || turns != 4 {
		t.Fatalf("correction: %v %v", result, err)
	}
}

func TestResearchBoundaries(t *testing.T) {
	t.Run("final turn disables tools", func(t *testing.T) {
		turns := 0
		session := researchSessionFunc(func(_ context.Context, system string, _ []researchReply, allow bool, _ *int) (researchTurn, error) {
			turns++
			if !allow {
				if turns != researchMaxTurns || !strings.Contains(system, "unresolved gaps") {
					t.Fatal("missing bounded final answer")
				}
				return researchTurn{Text: "Partial explanation.", Complete: true}, nil
			}
			return researchTurn{Calls: []researchCall{{fmt.Sprint(turns), "search", `{"query":"context"}`}}}, nil
		})
		_, err := runResearch(context.Background(), groundedPrompt, &evidenceSet{}, session, func(Lookup) ([]*Reference, error) { return nil, nil })
		if err != nil || turns != researchMaxTurns {
			t.Fatal(turns, err)
		}
	})
	t.Run("cancelled request does not run", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := runResearch(ctx, "", &evidenceSet{}, researchSessionFunc(func(context.Context, string, []researchReply, bool, *int) (researchTurn, error) {
			t.Fatal("provider ran after cancellation")
			return researchTurn{}, nil
		}), nil)
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	})
	t.Run("invented citation rejected", func(t *testing.T) {
		_, err := runResearch(context.Background(), "", &evidenceSet{}, researchSessionFunc(func(context.Context, string, []researchReply, bool, *int) (researchTurn, error) {
			return researchTurn{Text: "Claim [99].", Complete: true}, nil
		}), nil)
		if err == nil {
			t.Fatal("invalid citation accepted")
		}
	})
	t.Run("usage budget blocks provider", func(t *testing.T) {
		configureResearchTest(t, "openai", "http://127.0.0.1:1")
		used := researchMaxUnits
		_, err := reserveResearchCall(context.Background(), map[string]string{"message": "test"}, &used)
		var limit *LimitError
		if !errors.As(err, &limit) {
			t.Fatal(err)
		}
	})
}

func TestResearchAnswersWithoutMandatoryPlanning(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		nativeResponse(w, "openai", "Here is a clear explanation.", nil)
	}))
	defer server.Close()
	configureResearchTest(t, "openai", server.URL)
	result, err := Research(context.Background(), Request{Prompt: "Explain the idea in simpler language"}, func(string) ([]*Reference, error) { return nil, nil }, func(Lookup) ([]*Reference, error) { t.Fatal("unnecessary tool execution"); return nil, nil })
	if err != nil || calls != 1 || result == nil || !result.Complete {
		t.Fatal(result, err, calls)
	}
}

func TestEvidenceLabelsAndTruncation(t *testing.T) {
	e := &evidenceSet{}
	ref := &Reference{Text: "A verse", Metadata: map[string]string{"source": "quran"}}
	first := e.add([]*Reference{ref}, contentLimits.Evidence/2)
	duplicate := *ref
	duplicate.Score = 1
	again := e.add([]*Reference{&duplicate}, contentLimits.Evidence/2)
	if len(e.refs) != 1 || !strings.Contains(first, "[1]") || !strings.Contains(again, "[1]") {
		t.Fatal("duplicate evidence renumbered")
	}
	omitted := e.add([]*Reference{{Text: strings.Repeat("x", contentLimits.Evidence), Metadata: map[string]string{"source": "quran"}}}, contentLimits.Evidence/2)
	if !strings.Contains(omitted, "incomplete") || len(e.refs) != 1 {
		t.Fatal("silent truncation")
	}
}

func TestNativeToolErrorResultPreservesEmptyInput(t *testing.T) {
	for _, provider := range []string{"openai", "anthropic"} {
		t.Run(provider, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				data, _ := io.ReadAll(r.Body)
				if calls == 1 {
					nativeResponse(w, provider, "", []researchCall{{"empty", "search", `{}`}})
					return
				}
				if !strings.Contains(string(data), "Lookup failed: query required") {
					t.Error("missing recoverable tool failure")
				}
				if provider == "anthropic" && (!strings.Contains(string(data), `"input":{}`) || !strings.Contains(string(data), `"is_error":true`)) {
					t.Errorf("missing native input/error fields: %s", data)
				}
				nativeResponse(w, provider, "Please specify what you would like to look up.", nil)
			}))
			defer server.Close()
			configureResearchTest(t, provider, server.URL)
			result, err := Research(context.Background(), Request{Prompt: "Find the passage"}, func(string) ([]*Reference, error) { return nil, nil }, func(Lookup) ([]*Reference, error) { return nil, errors.New("query required") })
			if err != nil || result == nil || calls != 2 {
				t.Fatal(result, err, calls)
			}
		})
	}
}
