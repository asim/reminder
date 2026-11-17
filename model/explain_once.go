package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// explainOnce makes exactly one provider request. It shares provider streaming,
// evidence limits and accounting, but has no research loop or recovery retries.
func explainOnce(ctx context.Context, request Request, retrieve func(string) ([]*Reference, error)) (*Result, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, groundedKey{}, true)
	refs, err := retrieve(request.Prompt)
	if err != nil {
		return nil, err
	}
	evidence := &evidenceSet{refs: make([]*Reference, 0)}
	sources := evidence.add(refs, contentLimits.Evidence)
	if len(evidence.refs) == 0 {
		return &Result{Answer: "I don't have sufficient supporting sources to explain this passage.", References: evidence.refs}, nil
	}
	if !Available() {
		return nil, fmt.Errorf("no model configured")
	}
	Emit(ctx, Event{Type: "status", Text: "Explaining the passage…"})
	used := 0
	turn, err := newResearchSession(request).next(ctx, groundedPrompt+sources, nil, false, &used)
	if err != nil {
		return nil, err
	}
	if len(turn.Calls) > 0 {
		return nil, fmt.Errorf("model requested unavailable explanation tools")
	}
	if strings.TrimSpace(turn.Text) == "" {
		return nil, fmt.Errorf("model returned an empty answer")
	}
	for _, match := range citationPattern.FindAllStringSubmatch(turn.Text, -1) {
		number, err := strconv.Atoi(match[1])
		if err != nil || number < 1 || number > len(evidence.refs) {
			return nil, fmt.Errorf("model returned an invalid source citation")
		}
	}
	return &Result{Answer: turn.Text, References: evidence.refs, Complete: turn.Complete}, nil
}
