package model

import "context"

// Event carries provisional progress; only Result has passed citation validation.
type Event struct {
	Type   string  `json:"type"`
	Text   string  `json:"text,omitempty"`
	Result *Result `json:"result,omitempty"`
}
type eventKey struct{}

func WithEvents(ctx context.Context, emit func(Event)) context.Context {
	return context.WithValue(ctx, eventKey{}, emit)
}
func streaming(ctx context.Context) bool { _, ok := ctx.Value(eventKey{}).(func(Event)); return ok }
func Emit(ctx context.Context, event Event) {
	if emit, ok := ctx.Value(eventKey{}).(func(Event)); ok {
		emit(event)
	}
}
