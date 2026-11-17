package model

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestResearchStreamingProviders(t *testing.T) {
	for _, provider := range []string{"openai", "anthropic"} {
		t.Run(provider, func(t *testing.T) {
			turns := 0
			received := make(chan struct{}, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				turns++
				raw, _ := io.ReadAll(r.Body)
				if !strings.Contains(string(raw), `"stream":true`) {
					t.Error("not streaming")
				}
				if turns == 2 && (!strings.Contains(string(raw), "LOOKUP") || !strings.Contains(string(raw), "one")) {
					t.Error("lost native tool history")
				}
				w.Header().Set("Content-Type", "text/event-stream")
				send := func(v any) {
					data, _ := json.Marshal(v)
					if provider == "anthropic" {
						var event struct {
							Type string `json:"type"`
						}
						json.Unmarshal(data, &event)
						fmt.Fprintf(w, "event: %s\n", event.Type)
					}
					fmt.Fprintf(w, "data: %s\n\n", data)
					w.(http.Flusher).Flush()
				}
				if provider == "openai" {
					delta := func(d any, finish string) {
						send(map[string]any{"choices": []any{map[string]any{"index": 0, "delta": d, "finish_reason": finish}}})
					}
					if turns == 1 {
						delta(map[string]any{"tool_calls": []any{map[string]any{"index": 0, "id": "one", "type": "function", "function": map[string]string{"name": "search", "arguments": `{"query":`}}}}, "")
						delta(map[string]any{"tool_calls": []any{map[string]any{"index": 0, "function": map[string]string{"arguments": `"Musa"}`}}}}, "tool_calls")
					} else {
						delta(map[string]string{"content": "Meaning "}, "")
						select {
						case <-received:
						case <-time.After(5 * time.Second):
							t.Error("delta was buffered until completion")
						}
						delta(map[string]string{"content": "[2]."}, "stop")
					}
					fmt.Fprint(w, "data: [DONE]\n\n")
				} else {
					send(map[string]any{"type": "message_start", "message": map[string]any{"id": "msg", "type": "message", "role": "assistant", "model": "test", "content": []any{}, "usage": map[string]int{"input_tokens": 1, "output_tokens": 0}}})
					stop := "end_turn"
					if turns == 1 {
						stop = "tool_use"
						send(map[string]any{"type": "content_block_start", "index": 0, "content_block": map[string]any{"type": "tool_use", "id": "one", "name": "search", "input": map[string]any{}}})
						for _, part := range []string{`{"query":`, `"Musa"}`} {
							send(map[string]any{"type": "content_block_delta", "index": 0, "delta": map[string]string{"type": "input_json_delta", "partial_json": part}})
						}
					} else {
						send(map[string]any{"type": "content_block_start", "index": 0, "content_block": map[string]string{"type": "text", "text": ""}})
						send(map[string]any{"type": "content_block_delta", "index": 0, "delta": map[string]string{"type": "text_delta", "text": "Meaning "}})
						select {
						case <-received:
						case <-time.After(5 * time.Second):
							t.Error("delta was buffered until completion")
						}
						send(map[string]any{"type": "content_block_delta", "index": 0, "delta": map[string]string{"type": "text_delta", "text": "[2]."}})
					}
					send(map[string]any{"type": "content_block_stop", "index": 0})
					send(map[string]any{"type": "message_delta", "delta": map[string]string{"stop_reason": stop}, "usage": map[string]int{"output_tokens": 3}})
					send(map[string]string{"type": "message_stop"})
				}
			}))
			defer server.Close()
			configureResearchTest(t, provider, server.URL)
			var draft strings.Builder
			resets := 0
			ctx := WithEvents(context.Background(), func(event Event) {
				if event.Type == "reset" {
					resets++
				}
				if event.Type == "delta" {
					draft.WriteString(event.Text)
					select {
					case received <- struct{}{}:
					default:
					}
				}
			})
			result, err := Research(ctx, Request{Prompt: "Explain"}, func(string) ([]*Reference, error) {
				return []*Reference{{Text: "INITIAL", Metadata: map[string]string{"source": "quran"}}}, nil
			}, func(call Lookup) ([]*Reference, error) {
				if call.Query != "Musa" {
					t.Errorf("fragmented tool arguments lost: %+v", call)
				}
				return []*Reference{{Text: "LOOKUP", Metadata: map[string]string{"source": "quran"}}}, nil
			})
			if err != nil || result.Answer != "Meaning [2]." || !result.Complete || turns != 2 || draft.String() != "Meaning [2]." || resets < 2 {
				t.Fatalf("result=%+v err=%v turns=%d draft=%s", result, err, turns, draft.String())
			}
		})
	}
}

func TestStreamingDoesNotAcceptInvalidCitation(t *testing.T) {
	ctx := WithEvents(context.Background(), func(Event) {})
	e := &evidenceSet{}
	e.add([]*Reference{{Text: "source", Metadata: map[string]string{"source": "quran"}}}, contentLimits.Evidence/2)
	result, err := runResearch(ctx, "", e, researchSessionFunc(func(ctx context.Context, _ string, _ []researchReply, _ bool, _ *int) (researchTurn, error) {
		Emit(ctx, Event{Type: "delta", Text: "Wrong [99]."})
		return researchTurn{Text: "Wrong [99].", Complete: true}, nil
	}), nil)
	if err == nil || result != nil {
		t.Fatal("provisional citation accepted")
	}
}

func TestStreamCancellationStopsProvider(t *testing.T) {
	for _, provider := range []string{"openai", "anthropic"} {
		t.Run(provider, func(t *testing.T) {
			disconnected := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				if provider == "openai" {
					fmt.Fprint(w, "data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"Partial\"}}]}\n\n")
				} else {
					fmt.Fprint(w, "event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"id\":\"x\",\"role\":\"assistant\",\"content\":[]}}\n\nevent: content_block_start\ndata: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"text\",\"text\":\"\"}}\n\nevent: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"Partial\"}}\n\n")
				}
				w.(http.Flusher).Flush()
				<-r.Context().Done()
				close(disconnected)
			}))
			defer server.Close()
			configureResearchTest(t, provider, server.URL)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			ctx = WithEvents(ctx, func(e Event) {
				if e.Type == "delta" {
					cancel()
				}
			})
			result, err := Research(ctx, Request{Prompt: "Explain"}, func(string) ([]*Reference, error) { return nil, nil }, func(Lookup) ([]*Reference, error) { return nil, nil })
			if err == nil || result != nil {
				t.Fatal("cancelled stream accepted")
			}
			select {
			case <-disconnected:
			case <-time.After(time.Second):
				t.Fatal("provider not cancelled")
			}
		})
	}
}
