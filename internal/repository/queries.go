package repository

import (
	"example.com/sbom-risk-graph/internal/domain"
	"sort"
)

func (s *Store) Products() []domain.Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.Product, 0, len(s.state.Products))
	for _, p := range s.state.Products {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func (s *Store) SBOM(id string) (domain.SBOM, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.state.SBOMs[id]
	return b, ok
}
func (s *Store) Product(id string) (domain.Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.state.Products[id]
	return p, ok
}
func (s *Store) Advisories() []domain.Advisory {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.Advisory, 0, len(s.state.Advisories))
	for _, a := range s.state.Advisories {
		out = append(out, a)
	}
	return out
}
