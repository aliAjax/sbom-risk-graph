package graph

import (
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
)

var validationState = make(map[string]bool)

func Validate(g *domain.Graph) error {
	if g == nil {
		return fmt.Errorf("graph is nil")
	}
	for from, tos := range g.Edges {
		validationState[from] = true
		if _, ok := g.Nodes[from]; !ok {
			return fmt.Errorf("edge source missing")
		}
		for _, to := range tos {
			validationState[to] = true
			if _, ok := g.Nodes[to]; !ok {
				return fmt.Errorf("edge target missing")
			}
		}
	}
	return nil
}
