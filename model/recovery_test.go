package model

import (
	"context"
	"encoding/json"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResearchFirstCallRecovery(t *testing.T) {
	for _, provider := range []string{"openai", "anthropic"} {
		t.Run(provider, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				raw, _ := io.ReadAll(r.Body)
				var body map[string]any
				json.Unmarshal(raw, &body)
				if calls == 1 {
					if body["tools"] == nil {
						t.Error("normal path lost tools")
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(400)
					if provider == "anthropic" {
						fmt.Fprint(w, `{"type":"error","error":{"type":"invalid_request_error","message":"unsupported tools"}}`)
					} else {
						fmt.Fprint(w, `{"error":{"message":"unsupported tools","type":"invalid_request_error"}}`)
					}
					return
				}
				if calls > 2 {
					t.Error("unbounded retry")
				}
				if body["tools"] != nil || body["stream"] == true {
					t.Error("recovery still uses rejected native features")
				}
				if !strings.Contains(string(raw), "EXACT PASSAGE") || !strings.Contains(string(raw), "previous question") {
					t.Error("recovery lost evidence/history")
				}
				nativeResponse(w, provider, "Meaning [1].", nil)
			}))
			defer server.Close()
			configureResearchTest(t, provider, server.URL)
			var events []Event
			charged := 0
			ctx := WithTokenAllowance(context.Background(), func(int) error { charged++; return nil })
			ctx = WithEvents(ctx, func(e Event) { events = append(events, e) })
			retrievals := 0
			result, err := Research(ctx, Request{Prompt: "What does this mean?", History: []Message{{Role: "user", Content: "previous question"}, {Role: "assistant", Content: "previous answer"}}}, func(string) ([]*Reference, error) {
				retrievals++
				return []*Reference{{Text: "EXACT PASSAGE", Metadata: map[string]string{"source": "bukhari"}}}, nil
			}, func(Lookup) ([]*Reference, error) { t.Fatal("unexpected lookup"); return nil, nil })
			if err != nil || result.Answer != "Meaning [1]." || calls != 2 || charged != 2 || retrievals != 1 {
				t.Fatalf("result=%+v err=%v calls=%d charged=%d retrievals=%d", result, err, calls, charged, retrievals)
			}
			if events[len(events)-1].Type != "status" {
				t.Fatal("missing recovery status")
			}
		})
	}
}
func TestRecoveryDoesNotRetryOperationalFailures(t *testing.T) {
	for _, status := range []int{401, 403, 404, 429, 500, 503} {
		if recoverableResearchError(&openai.APIError{HTTPStatusCode: status}) {
			t.Fatalf("retrying %d", status)
		}
	}
	for _, err := range []error{context.Canceled, context.DeadlineExceeded, &LimitError{"busy"}} {
		if recoverableResearchError(err) {
			t.Fatal("retrying cancellation or limits")
		}
	}
	secret := "secret prompt and key"
	for _, status := range []int{400, 401, 404, 429, 500} {
		message := FailureMessage(&openai.APIError{HTTPStatusCode: status, Message: secret})
		if strings.Contains(message, secret) || message == "" {
			t.Fatal("unsafe error message")
		}
	}
}
