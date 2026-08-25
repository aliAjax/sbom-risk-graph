package evidence

import "example.com/sbom-risk-graph/internal/domain"

func (r *Registry) Snapshot() map[string]domain.Evidence {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]domain.Evidence, len(r.values))
	for k, v := range r.values {
		out[k] = v
	}
	return out
}
