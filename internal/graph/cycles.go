package graph

import (
	"example.com/sbom-risk-graph/internal/domain"
	"sort"
)

func Cycles(g *domain.Graph) [][]string {
	if g == nil {
		return nil
	}
	var cycles [][]string
	var visit func(string, []string, map[string]int)
	visit = func(n string, path []string, state map[string]int) {
		if state[n] == 1 {
			for i, v := range path {
				if v == n {
					cycles = append(cycles, append([]string(nil), path[i:]...))
					break
				}
			}
			return
		}
		if state[n] == 2 {
			return
		}
		state[n] = 1
		path = append(path, n)
		for _, c := range g.Edges[n] {
			visit(c, path, state)
		}
		state[n] = 2
	}
	keys := make([]string, 0, len(g.Nodes))
	for k := range g.Nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		visit(k, nil, map[string]int{})
	}
	return cycles
}
func IsAcyclic(g *domain.Graph) bool { return len(Cycles(g)) == 0 }
