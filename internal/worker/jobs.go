package worker

import (
	"context"
	"example.com/sbom-risk-graph/internal/domain"
	"example.com/sbom-risk-graph/pkg/metrics"
	"sync"
	"time"
)

type Handler func(context.Context, domain.Job) error
type Queue struct {
	mu      sync.Mutex
	jobs    []domain.Job
	handler Handler
	metrics *metrics.Counters
	stop    chan struct{}
	done    chan struct{}
}

func New(handler Handler, m *metrics.Counters) *Queue {
	q := &Queue{handler: handler, metrics: m, stop: make(chan struct{}), done: make(chan struct{})}
	go q.loop()
	return q
}
func (q *Queue) Submit(job domain.Job) {
	q.mu.Lock()
	q.jobs = append(q.jobs, job)
	q.mu.Unlock()
	q.metrics.Jobs.Add(1)
}
func (q *Queue) loop() {
	defer close(q.done)
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			q.mu.Lock()
			if len(q.jobs) > 0 {
				job := q.jobs[0]
				q.jobs = q.jobs[1:]
				q.mu.Unlock()
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				_ = q.handler(ctx, job)
				cancel()
			} else {
				q.mu.Unlock()
			}
		case <-q.stop:
			return
		}
	}
}
func (q *Queue) Close() { close(q.stop); <-q.done }
