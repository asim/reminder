package api

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/asim/reminder/model"
)

func cachedFixture() *model.Result {
	return &model.Result{Complete: true, Answer: "Explanation [1].", References: []*model.Reference{{Text: "Evidence", Metadata: map[string]string{"source": "quran"}}}, Coverage: map[string]int{"included_verses": 7, "total_verses": 7}}
}

func TestExplainCachePersistenceIsolationAndVersions(t *testing.T) {
	dir := t.TempDir()
	cache, err := openExplainCache(dir, "build1")
	if err != nil {
		t.Fatal(err)
	}
	req := model.ExplainRequest{Context: model.ExplainContext{Chapter: 1}}
	key := cache.key(req)
	calls := 0
	generate := func() (*model.Result, error) { calls++; return cachedFixture(), nil }
	first, err := cache.result(context.Background(), key, generate)
	if err != nil {
		t.Fatal(err)
	}
	first.Answer = "mutated"
	first.References[0].Text = "mutated"
	first.Coverage["total_verses"] = 0
	cache.Close()
	cache, err = openExplainCache(dir, "build1")
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	got, err := cache.result(context.Background(), key, generate)
	if err != nil || calls != 1 || got.Answer != "Explanation [1]." || got.References[0].Text != "Evidence" || got.Coverage["total_verses"] != 7 {
		t.Fatalf("cache did not preserve snapshot: %+v %v calls=%d", got, err, calls)
	}
	normalized := req
	normalized.Context.Source = "quran"
	normalized.Context.Translation = "khattab"
	if cache.key(normalized) != key {
		t.Fatal("default aliases did not share key")
	}
	variants := []model.ExplainRequest{req, req, req, req}
	variants[0].Question = "Explain this passage"
	variants[1].Q = "Explain this passage"
	variants[2].Context.Selection = "text"
	variants[3].History = []model.Message{{Role: "user", Content: "private"}, {Role: "assistant", Content: "private"}}
	for _, r := range variants {
		if cache.key(r) != "" {
			t.Fatal("private/custom request cacheable")
		}
	}
	for _, r := range []model.ExplainRequest{{Context: model.ExplainContext{Chapter: 2}}, {Context: model.ExplainContext{Chapter: 1, Verse: 1}}, {Context: model.ExplainContext{Chapter: 1, Translation: "sahih"}}} {
		if cache.key(r) == key {
			t.Fatal("context collision")
		}
	}
	cache.version = "build2"
	if cache.key(req) == key {
		t.Fatal("build did not invalidate")
	}
	cache.version = "build1"
	t.Setenv("EXPLAIN_CACHE_EPOCH", "new")
	if cache.key(req) == key {
		t.Fatal("epoch did not invalidate")
	}
}

func TestExplainCacheSingleFlightAndCancellation(t *testing.T) {
	cache, err := openExplainCache(t.TempDir(), "build")
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	var calls atomic.Int32
	entered := make(chan struct{})
	release := make(chan struct{})
	generate := func() (*model.Result, error) { calls.Add(1); close(entered); <-release; return cachedFixture(), nil }
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := cache.result(context.Background(), "key", generate)
		if err != nil {
			t.Error(err)
		}
	}()
	<-entered
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := cache.result(ctx, "key", generate); !errors.Is(err, context.Canceled) {
		t.Fatal("waiting caller did not cancel")
	}
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := cache.result(context.Background(), "key", generate)
			if e != nil || r.Answer != "Explanation [1]." {
				t.Errorf("shared result: %v", e)
			}
		}()
	}
	close(release)
	wg.Wait()
	if calls.Load() != 1 {
		t.Fatalf("generated %d times", calls.Load())
	}
}

func TestExplainCacheRejectsFailuresIncompleteAndExpired(t *testing.T) {
	cache, err := openExplainCache(t.TempDir(), "build")
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	for _, kind := range []string{"error", "incomplete", "empty", "no-evidence"} {
		calls := 0
		gen := func() (*model.Result, error) {
			calls++
			r := cachedFixture()
			switch kind {
			case "error":
				return nil, errors.New("failed")
			case "incomplete":
				r.Complete = false
			case "empty":
				r.Answer = ""
			case "no-evidence":
				r.References = nil
			}
			return r, nil
		}
		cache.result(context.Background(), kind, gen)
		cache.result(context.Background(), kind, gen)
		if calls != 2 {
			t.Fatalf("cached %s", kind)
		}
	}
	cache.result(context.Background(), "expired", func() (*model.Result, error) { return cachedFixture(), nil })
	if _, err := cache.db.Exec("UPDATE explanations SET expires=0"); err != nil {
		t.Fatal(err)
	}
	calls := 0
	cache.result(context.Background(), "expired", func() (*model.Result, error) { calls++; return cachedFixture(), nil })
	if calls != 1 {
		t.Fatal("expired cache used")
	}
}
