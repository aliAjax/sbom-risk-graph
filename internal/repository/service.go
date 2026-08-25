package repository

import (
	"example.com/sbom-risk-graph/internal/advisory"
	"example.com/sbom-risk-graph/internal/domain"
	"example.com/sbom-risk-graph/internal/graph"
	"example.com/sbom-risk-graph/internal/policy"
	"time"
)

type Service struct {
	Repo     *Store
	Graph    *graph.Store
	Advisory *advisory.Engine
	Policies *policy.Store
}

func NewService(repo *Store) *Service {
	return &Service{Repo: repo, Graph: graph.New(), Advisory: advisory.New(), Policies: policy.NewStore()}
}
func (s *Service) AddProduct(p domain.Product) error {
	if err := s.Repo.SaveProduct(p); err != nil {
		return err
	}
	return s.Graph.PutProduct(p)
}
func (s *Service) AddSBOM(b domain.SBOM) error {
	if err := s.Repo.SaveSBOM(b); err != nil {
		return err
	}
	return s.Graph.PutSBOM(b)
}
func (s *Service) AddAdvisory(a domain.Advisory) error {
	if err := s.Repo.SaveAdvisory(a); err != nil {
		return err
	}
	return s.Advisory.Add(a)
}
func (s *Service) AddPolicy(p domain.Policy) error {
	if err := s.Repo.SavePolicy(p); err != nil {
		return err
	}
	return s.Policies.Put(p)
}
func (s *Service) Risk(productID string) domain.RiskReport {
	components := s.Graph.Components(productID)
	report := domain.RiskReport{ProductID: productID, GeneratedAt: time.Now().UTC(), Components: len(components)}
	for _, c := range components {
		for _, a := range s.Advisory.Match(c) {
			report.Vulnerabilities++
			score := advisory.Score(a)
			if score > report.Score {
				report.Score = score
			}
		}
	}
	for _, p := range s.Policies.All() {
		report.Violations = append(report.Violations, policy.New().Evaluate(p, components, s.Advisory.Match, false)...)
	}
	report.Sort()
	return report
}
