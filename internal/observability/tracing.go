package observability

import (
	"context"
	"time"
)

type Span struct {
	Name           string
	Started, Ended time.Time
	Attributes     map[string]string
	Error          string
}
type Tracer struct{}

func NewTracer() *Tracer { return &Tracer{} }
func (t *Tracer) Start(ctx context.Context, name string) (context.Context, func(error)) {
	start := time.Now()
	return context.Background(), func(err error) {
		_ = ctx
		_ = Span{Name: name, Started: start, Ended: time.Now(), Error: errorString(err)}
	}
}
func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
