package model

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// LimitError is safe to show to callers when public generation is capped.
type LimitError struct{ Message string }

func (e *LimitError) Error() string { return e.Message }

type clientKey struct{}

// WithClient identifies a request for the per-client AI allowance.
func WithClient(ctx context.Context, client string) context.Context {
	return context.WithValue(ctx, clientKey{}, client)
}

type usage struct {
	Day     string         `json:"day"`
	Total   int            `json:"total"`
	Clients map[string]int `json:"clients"`
}

type limits struct {
	mu                                   sync.Mutex
	path                                 string
	daily, perClient, concurrent, output int
	active                               int
	accountingErr                        error
	usage                                usage
}

var budget *limits

// Content limits are independent of anonymous search's generation allowance.
var contentLimits = struct{ Prompt, Input, Evidence, SearchEvidence, Output int }{32768, 131072, 65536, 32768, 4096}

func ModelBodyLimit() int64    { return int64(contentLimits.Input)*6 + 65536 }
func SearchEvidenceLimit() int { return contentLimits.SearchEvidence }

type tokenAllowanceKey struct{}

// WithTokenAllowance installs the authenticated caller's persistent preflight budget.
func WithTokenAllowance(ctx context.Context, charge func(int) error) context.Context {
	return context.WithValue(ctx, tokenAllowanceKey{}, charge)
}

type groundedKey struct{}

// ConfigureLimits must run once at startup. Counters are owned by one server process.
func ConfigureLimits(dir string) error {
	for name, setting := range map[string]struct {
		target   *int
		fallback int
	}{
		"MODEL_MAX_PROMPT_BYTES":    {&contentLimits.Prompt, 32768},
		"MODEL_MAX_INPUT_BYTES":     {&contentLimits.Input, 131072},
		"MODEL_MAX_EVIDENCE_BYTES":  {&contentLimits.Evidence, 65536},
		"SEARCH_MAX_EVIDENCE_BYTES": {&contentLimits.SearchEvidence, 32768},
		"MODEL_MAX_OUTPUT_TOKENS":   {&contentLimits.Output, 4096},
	} {
		value := setting.fallback
		if raw, ok := os.LookupEnv(name); ok {
			var err error
			value, err = strconv.Atoi(raw)
			if err != nil || value < 1 || value > 16*1024*1024 {
				return fmt.Errorf("%s must be between 1 and 16777216", name)
			}
		}
		*setting.target = value
	}
	if contentLimits.Prompt > contentLimits.Input {
		return fmt.Errorf("MODEL_MAX_PROMPT_BYTES must not exceed MODEL_MAX_INPUT_BYTES")
	}
	l := &limits{path: filepath.Join(dir, "ai-usage.json")}
	for name, setting := range map[string]struct {
		target   *int
		fallback int
	}{
		"AI_DAILY_LIMIT": {&l.daily, 100}, "AI_IP_DAILY_LIMIT": {&l.perClient, 20},
		"AI_CONCURRENCY_LIMIT": {&l.concurrent, 2}, "AI_MAX_OUTPUT_TOKENS": {&l.output, 1024},
	} {
		value := setting.fallback
		if raw, ok := os.LookupEnv(name); ok {
			var err error
			value, err = strconv.Atoi(raw)
			if err != nil || value < 0 {
				return fmt.Errorf("%s must be a non-negative integer", name)
			}
		}
		if name == "AI_MAX_OUTPUT_TOKENS" && value == 0 {
			return fmt.Errorf("%s must be positive", name)
		}
		*setting.target = value
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		l.accountingErr = err
	} else {
		data, err := os.ReadFile(l.path)
		if err == nil {
			l.accountingErr = json.Unmarshal(data, &l.usage)
		} else if !os.IsNotExist(err) {
			l.accountingErr = err
		}
	}
	if l.accountingErr != nil {
		log.Printf("Public AI disabled: cannot load usage counters: %v", l.accountingErr)
	}
	budget = l
	return nil
}

func (l *limits) acquire(ctx context.Context, now time.Time) (func(), error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.accountingErr != nil {
		return nil, fmt.Errorf("AI usage accounting unavailable")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	day := now.UTC().Format("2006-01-02")
	if l.usage.Day != day {
		l.usage = usage{Day: day, Clients: map[string]int{}}
	}
	if l.usage.Clients == nil {
		l.usage.Clients = map[string]int{}
	}
	client, _ := ctx.Value(clientKey{}).(string)
	key := ""
	if client != "" {
		key = fmt.Sprintf("%x", sha256.Sum256([]byte(day+":"+client)))
	}
	if l.usage.Total >= l.daily {
		return nil, &LimitError{"Daily AI limit reached. Please try again tomorrow (UTC)."}
	}
	if key != "" && l.usage.Clients[key] >= l.perClient {
		return nil, &LimitError{"Your daily AI limit has been reached. Please try again tomorrow (UTC)."}
	}
	if l.active >= l.concurrent {
		return nil, &LimitError{"AI is busy. Please try again shortly."}
	}
	l.usage.Total++
	if key != "" {
		l.usage.Clients[key]++
	}
	data, err := json.Marshal(l.usage)
	if err == nil {
		err = os.WriteFile(l.path+".tmp", data, 0600)
	}
	if err == nil {
		err = os.Rename(l.path+".tmp", l.path)
	}
	if err != nil {
		return nil, fmt.Errorf("AI usage accounting unavailable")
	}
	l.active++
	return func() { l.mu.Lock(); l.active--; l.mu.Unlock() }, nil
}
