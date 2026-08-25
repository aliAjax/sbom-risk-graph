package graph

import (
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
	"sort"
)

type QueryResult struct {
	Components []domain.Component `json:"components"`
	Depth      int                `json:"depth"`
	Truncated  bool               `json:"truncated"`
}

func ExpandChecked(g *domain.Graph, root string, depth, limit int) (QueryResult, error) {
	if err := Validate(g); err != nil {
		return QueryResult{}, fmt.Errorf("expand dependency graph: %w", err)
	}
	return Expand(g, root, depth, limit), nil
}

func Expand(g *domain.Graph, root string, depth, limit int) QueryResult {
	if limit <= 0 {
		limit = 1000
	}
	result := QueryResult{Depth: depth}
	paths := g.PathsFrom(root, depth)
	seen := map[string]bool{}
	for _, p := range paths {
		for _, key := range p {
			if seen[key] {
				continue
			}
			seen[key] = true
			if len(result.Components) >= limit {
				result.Truncated = true
				return result
			}
			if c, ok := g.Nodes[key]; ok {
				result.Components = append(result.Components, c)
			}
		}
	}
	sort.Slice(result.Components, func(i, j int) bool { return result.Components[i].Key() < result.Components[j].Key() })
	return result
}
