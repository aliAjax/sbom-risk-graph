package evidence

import "example.com/sbom-risk-graph/internal/domain"

func (r *Registry) Snapshot() map[string]domain.Evidence {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.values
}
