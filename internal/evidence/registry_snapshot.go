package evidence

import "example.com/sbom-risk-graph/internal/domain"

func (r *Registry) Snapshot() map[string]domain.Evidence {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make(map[string]domain.Evidence, len(r.values))
	for id, item := range r.values {
		items[id] = item
	}
	return items
}
