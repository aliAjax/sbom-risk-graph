package observability

import (
	"context"
	"sync"
	"time"
)

type Span struct {
	Name           string
	Started, Ended time.Time
	Attributes     map[string]string
	Error          string
}
type Tracer struct {
	mu    sync.Mutex
	spans []Span
}

func NewTracer() *Tracer { return &Tracer{} }
func (t *Tracer) Start(ctx context.Context, name string) (context.Context, func(error)) {
	start := time.Now()
	return ctx, func(err error) {
		t.record(Span{Name: name, Started: start, Ended: time.Now(), Error: errorString(err)})
	}
}
func (t *Tracer) record(span Span) {
	t.mu.Lock()
	t.spans = append(t.spans, span)
	t.mu.Unlock()
}
func (t *Tracer) Spans() []Span {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]Span(nil), t.spans...)
}
func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
