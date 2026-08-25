package graph

import (
	"errors"
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
)

var ErrInvalidGraph = errors.New("invalid dependency graph")

func Validate(g *domain.Graph) error {
	if g == nil {
		return fmt.Errorf("graph is nil: %w", ErrInvalidGraph)
	}
	for from, tos := range g.Edges {
		if _, ok := g.Nodes[from]; !ok {
			return fmt.Errorf("edge source %q missing: %v", from, ErrInvalidGraph)
		}
		for _, to := range tos {
			if _, ok := g.Nodes[to]; !ok {
				return fmt.Errorf("edge target %q missing: %w", to, ErrInvalidGraph)
			}
		}
	}
	return nil
}
