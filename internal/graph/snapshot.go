package graph

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"example.com/sbom-risk-graph/internal/domain"
	"sort"
)

type Snapshot struct {
	ProductID string              `json:"product_id"`
	Nodes     []domain.Component  `json:"nodes"`
	Edges     []domain.Dependency `json:"edges"`
	Digest    string              `json:"digest"`
}

func BuildSnapshot(productID string, g *domain.Graph) Snapshot {
	if g == nil {
		return Snapshot{ProductID: productID}
	}
	s := Snapshot{ProductID: productID}
	for _, n := range g.Nodes {
		s.Nodes = append(s.Nodes, n)
	}
	sort.Slice(s.Nodes, func(i, j int) bool { return s.Nodes[i].Key() < s.Nodes[j].Key() })
	for from, tos := range g.Edges {
		for _, to := range tos {
			s.Edges = append(s.Edges, domain.Dependency{From: from, To: to})
		}
	}
	sort.Slice(s.Edges, func(i, j int) bool {
		if s.Edges[i].From == s.Edges[j].From {
			return s.Edges[i].To < s.Edges[j].To
		}
		return s.Edges[i].From < s.Edges[j].From
	})
	b, _ := json.Marshal(s)
	d := sha256.Sum256(b)
	s.Digest = hex.EncodeToString(d[:])
	return s
}
func VerifySnapshot(s Snapshot) bool {
	want := s.Digest
	s.Digest = ""
	b, _ := json.Marshal(s)
	d := sha256.Sum256(b)
	return want == hex.EncodeToString(d[:])
}
