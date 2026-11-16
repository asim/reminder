package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Lookup is a read-only request against server-owned texts, never a URL or SQL.
type Lookup struct {
	Tool    string   `json:"tool"`
	Query   string   `json:"query,omitempty"`
	Source  string   `json:"source,omitempty"`
	Chapter int      `json:"chapter,omitempty"`
	Verse   int      `json:"verse,omitempty"`
	Terms   []string `json:"terms,omitempty"`
}
type LookupFunc func(Lookup) ([]*Reference, error)

const (
	researchTimeout  = 120 * time.Second
	researchMaxUnits = 1024 * 1024 // conservative input bytes plus maximum output, across all turns
	researchMaxTurns = 8           // runaway-loop backstop; last turn is answer-only
)

type researchCall struct{ ID, Name, Arguments string }
type researchReply struct {
	ID, Text string
	IsError  bool
}
type researchTurn struct {
	Text     string
	Calls    []researchCall
	Complete bool
}

// Native provider sessions retain assistant tool calls and return matching tool results.
type researchSession interface {
	next(context.Context, string, []researchReply, bool, *int) (researchTurn, error)
}

type evidenceSet struct {
	refs   []*Reference
	labels map[string]int
	bytes  int
}

// Number sources once for the entire conversation. Never renumber citations when
// new results arrive. Reserve half the evidence budget for subsequent lookups.
func (e *evidenceSet) add(refs []*Reference, allowance int) string {
	if e.labels == nil {
		e.labels = map[string]int{}
	}
	var out strings.Builder
	omitted := 0
	for _, ref := range refs {
		if ref == nil {
			continue
		}
		switch ref.Metadata["source"] {
		case "quran", "quran_translation", "tafsir", "bukhari", "muslim", "names", "lexicon", "morphology":
		default:
			continue
		}
		identity, _ := json.Marshal(struct {
			Text     string
			Metadata map[string]string
		}{ref.Text, ref.Metadata})
		key := string(identity)
		label, exists := e.labels[key]
		if !exists {
			label = len(e.refs) + 1
		}
		data, _ := json.Marshal(ref)
		record := fmt.Sprintf("[%d] %s\n", label, data)
		if out.Len()+len(record) > allowance || (!exists && e.bytes+len(record) > contentLimits.Evidence) {
			omitted++
			continue
		}
		if !exists {
			e.labels[key] = label
			e.refs = append(e.refs, ref)
			e.bytes += len(record)
		}
		out.WriteString(record)
	}
	if omitted > 0 {
		fmt.Fprintf(&out, "%d records omitted due to the evidence budget; this response is incomplete.\n", omitted)
	}
	if out.Len() == 0 {
		out.WriteString("No matching source records were returned. This does not establish absence from the full corpus.\n")
	}
	return out.String()
}

// Research lets the model answer directly or repeatedly call native read-only
// tools. Knowledge supports reasoning; checkable claims are verified by retrieval.
func Research(ctx context.Context, request Request, retrieve func(string) ([]*Reference, error), lookup LookupFunc) (result *Result, err error) {
	if err = request.Validate(); err != nil {
		return nil, err
	}
	if lookup == nil {
		return Respond(ctx, request, retrieve)
	}
	defer func() {
		if r := recover(); r != nil {
			result = nil
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("model: %v", r)
			}
		}
	}()
	ctx, cancel := context.WithTimeout(ctx, researchTimeout)
	defer cancel()
	ctx = context.WithValue(ctx, groundedKey{}, true)
	query := request.Prompt
	for i := len(request.History) - 2; i >= 0 && i >= len(request.History)-4; i -= 2 {
		if len(query)+len(request.History[i].Content)+1 <= contentLimits.Prompt {
			query += "\n" + request.History[i].Content
		}
	}
	initial, err := retrieve(query)
	if err != nil {
		return nil, err
	}
	if !Available() {
		return nil, fmt.Errorf("no model configured")
	}
	evidence := &evidenceSet{refs: make([]*Reference, 0)}
	system := groundedPrompt + evidence.add(initial, contentLimits.Evidence/2)
	session := newResearchSession(request)
	return runResearch(ctx, system, evidence, session, lookup)
}

func runResearch(ctx context.Context, system string, evidence *evidenceSet, session researchSession, lookup LookupFunc) (*Result, error) {
	var replies []researchReply
	usedUnits := 0
	for round := 0; round < researchMaxTurns; round++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		allowTools := round < researchMaxTurns-1 && usedUnits < researchMaxUnits/2 && evidence.bytes < contentLimits.Evidence
		prompt := system
		if !allowTools {
			prompt += "\nThe research budget is ending. Answer from what you have established so far and clearly state any material unresolved gaps. Do not request more tools."
		}
		turn, err := session.next(ctx, prompt, replies, allowTools, &usedUnits)
		if err != nil {
			return nil, err
		}
		if len(turn.Calls) == 0 {
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
		if !allowTools {
			return nil, fmt.Errorf("model requested tools after the research budget ended")
		}
		replies = nil
		ids := map[string]bool{}
		for _, call := range turn.Calls {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if call.ID == "" || ids[call.ID] {
				return nil, fmt.Errorf("model returned invalid tool call identifiers")
			}
			ids[call.ID] = true
			reply := researchReply{ID: call.ID}
			input, err := decodeLookup(call)
			if err == nil {
				var refs []*Reference
				refs, err = lookup(input)
				if err == nil {
					reply.Text = evidence.add(refs, contentLimits.Evidence/2)
				}
			}
			if err != nil {
				if ctx.Err() != nil {
					return nil, ctx.Err()
				}
				// Recoverable lookup failures are visible to the model so it can retry with
				// corrected arguments or another source, rather than fabricate an answer.
				reply.IsError = true
				reply.Text = "Lookup failed: " + err.Error()
			}
			replies = append(replies, reply)
		}
	}
	return nil, fmt.Errorf("research budget exhausted")
}

func decodeLookup(call researchCall) (Lookup, error) {
	input := Lookup{}
	if len(call.Arguments) > 8192 || !strings.HasPrefix(strings.TrimSpace(call.Arguments), "{") {
		return input, fmt.Errorf("tool arguments must be a JSON object")
	}
	if err := json.Unmarshal([]byte(call.Arguments), &input); err != nil {
		return input, fmt.Errorf("invalid tool arguments: %w", err)
	}
	// The native tool name is authoritative; an input field cannot change it.
	input.Tool = call.Name
	switch call.Name {
	case "search", "passage", "quran_mentions", "arabic":
		return input, nil
	default:
		return input, fmt.Errorf("unknown lookup tool")
	}
}
