package policy

import (
	"example.com/sbom-risk-graph/internal/domain"
	"sync"
)

type Store struct {
	mu     sync.RWMutex
	values map[string]domain.Policy
}

func NewStore() *Store { return &Store{values: make(map[string]domain.Policy)} }
func (s *Store) Put(p domain.Policy) error {
	if err := p.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p.Revision++
	s.values[p.ID] = p
	return nil
}
func (s *Store) Get(id string) (domain.Policy, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.values[id]
	return p.Clone(), ok
}
func (s *Store) All() []domain.Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Policy, 0, len(s.values))
	for _, p := range s.values {
		out = append(out, p.Clone())
	}
	return out
}

func (s *Store) Replace(id string, p domain.Policy) error {
	if err := p.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.values == nil {
		s.values = make(map[string]domain.Policy)
	}
	p.ID = id
	s.values[id] = p
	return nil
}
