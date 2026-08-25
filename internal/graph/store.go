package graph

import (
	"example.com/sbom-risk-graph/internal/domain"
	"sort"
	"sync"
)

type Store struct {
	mu       sync.RWMutex
	products map[string]domain.Product
	graphs   map[string]*domain.Graph
	sboms    map[string]domain.SBOM
}

func New() *Store {
	return &Store{products: make(map[string]domain.Product), graphs: make(map[string]*domain.Graph), sboms: make(map[string]domain.SBOM)}
}
func (s *Store) PutProduct(p domain.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p.Revision++
	s.products[p.ID] = p
	if _, ok := s.graphs[p.ID]; !ok {
		s.graphs[p.ID] = domain.NewGraph()
	}
	return nil
}
func (s *Store) GetProduct(id string) (domain.Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.products[id]
	return p, ok
}
func (s *Store) PutSBOM(b domain.SBOM) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sboms[b.ID] = b
	g := domain.NewGraph()
	for _, c := range b.Components {
		if err := g.AddComponent(c); err != nil {
			return err
		}
	}
	byID := make(map[string]string)
	for _, c := range b.Components {
		byID[c.ID] = c.Key()
	}
	for _, d := range b.Dependencies {
		from, to := byID[d.From], byID[d.To]
		if from != "" && to != "" {
			if err := g.AddDependency(from, to); err != nil {
				return err
			}
		}
	}
	s.graphs[b.ProductID] = g
	return nil
}
func (s *Store) GetSBOM(id string) (domain.SBOM, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.sboms[id]
	return b, ok
}
func (s *Store) Graph(productID string) (*domain.Graph, bool) {
	g, ok := s.graphs[productID]
	return g, ok
}
func (s *Store) Components(productID string) []domain.Component {
	g := s.graphs[productID]
	if g == nil {
		return nil
	}
	out := make([]domain.Component, 0, len(g.Nodes))
	for _, c := range g.Nodes {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key() < out[j].Key() })
	return out
}
func (s *Store) ProductIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.products))
	for id := range s.products {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
