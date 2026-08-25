package policy

import (
	"example.com/sbom-risk-graph/internal/domain"
	"sort"
)

type History struct{ values map[string][]domain.Policy }

func NewHistory() *History { return &History{values: make(map[string][]domain.Policy)} }
func (h *History) Add(p domain.Policy) {
	h.values[p.ID] = append(h.values[p.ID], p)
	sort.Slice(h.values[p.ID], func(i, j int) bool { return h.values[p.ID][i].Revision < h.values[p.ID][j].Revision })
}
func (h *History) Get(id string, revision uint64) (domain.Policy, bool) {
	for _, p := range h.values[id] {
		if p.Revision == revision {
			return p, true
		}
	}
	return domain.Policy{}, false
}
func (h *History) Latest(id string) (domain.Policy, bool) {
	items := h.values[id]
	if len(items) == 0 {
		return domain.Policy{}, false
	}
	return items[len(items)-1], true
}
