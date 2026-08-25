package repository

import (
	"example.com/sbom-risk-graph/internal/domain"
	"sync"
)

type JobStore struct {
	mu     sync.RWMutex
	values map[string]domain.Job
}

func NewJobStore() *JobStore         { return &JobStore{values: make(map[string]domain.Job)} }
func (s *JobStore) Put(j domain.Job) { s.mu.Lock(); s.values[j.ID] = j; s.mu.Unlock() }
func (s *JobStore) Get(id string) (domain.Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.values[id]
	return j, ok
}
func (s *JobStore) Update(id string, fn func(*domain.Job)) {
	s.mu.Lock()
	if j, ok := s.values[id]; ok {
		fn(&j)
		s.values[id] = j
	}
	s.mu.Unlock()
}
func (s *JobStore) All() []domain.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Job, 0, len(s.values))
	for _, j := range s.values {
		out = append(out, j)
	}
	return out
}
