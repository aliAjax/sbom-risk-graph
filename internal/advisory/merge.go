package advisory

import (
	"example.com/sbom-risk-graph/internal/domain"
	"sort"
)

func Merge(items []domain.Advisory) []domain.Advisory {
	byID := make(map[string]domain.Advisory)
	for _, a := range items {
		old, ok := byID[a.ID]
		if !ok || a.PublishedAt.After(old.PublishedAt) {
			byID[a.ID] = a
		}
	}
	out := make([]domain.Advisory, 0, len(byID))
	for _, a := range byID {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
