package graph

import (
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
)

func Validate(g *domain.Graph) error {
	if g == nil {
		return fmt.Errorf("graph is nil")
	}
	for from, tos := range g.Edges {
		if _, ok := g.Nodes[from]; !ok {
			return fmt.Errorf("edge source missing")
		}
		for _, to := range tos {
			if _, ok := g.Nodes[to]; !ok {
				return fmt.Errorf("edge target missing")
			}
		}
	}
	return nil
}
