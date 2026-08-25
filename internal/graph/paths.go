package graph

import "example.com/sbom-risk-graph/internal/domain"

func Paths(g *domain.Graph, root string, depth int) [][]string {
	if g == nil {
		return nil
	}
	return g.PathsFrom(root, depth)
}
func Reverse(g *domain.Graph) map[string][]string {
	out := make(map[string][]string)
	if g == nil {
		return out
	}
	for from, tos := range g.Edges {
		for _, to := range tos {
			out[to] = append(out[to], from)
		}
	}
	return out
}
