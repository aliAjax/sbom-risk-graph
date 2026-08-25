package evidence

import (
	"fmt"
	"sync"

	"example.com/sbom-risk-graph/internal/domain"
)

type Registry struct {
	mu     sync.RWMutex
	values map[string]domain.Evidence
}

func NewRegistry(items []domain.Evidence) *Registry {
	values := make(map[string]domain.Evidence, len(items))
	for _, item := range items {
		values[item.ID] = item
	}
	return &Registry{values: values}
}

func (r *Registry) Get(id string) (domain.Evidence, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.values[id]
	return item, ok
}

func (r *Registry) Transition(id string, next domain.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.values[id]
	if !ok {
		return fmt.Errorf("evidence %s not found", id)
	}
	previous := item.Status
	item.Status = next
	r.values[id] = item
	if err := ValidateTransition(previous, next); err != nil {
		return err
	}
	return nil
}
