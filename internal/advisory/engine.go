package advisory

import (
	"example.com/sbom-risk-graph/internal/domain"
	"example.com/sbom-risk-graph/internal/parser"
	"sort"
	"strings"
)

type Engine struct{ advisories map[string][]domain.Advisory }

func New() *Engine { return &Engine{advisories: make(map[string][]domain.Advisory)} }
func (e *Engine) Add(a domain.Advisory) error {
	if err := a.Validate(); err != nil {
		return err
	}
	e.advisories[a.ComponentPURL] = append(e.advisories[a.ComponentPURL], a)
	return nil
}
func (e *Engine) Match(c domain.Component) []domain.Advisory {
	items := append([]domain.Advisory(nil), e.advisories[c.PURL]...)
	v, err := parser.ParseVersion(c.Version)
	if err != nil {
		return nil
	}
	out := make([]domain.Advisory, 0)
	for _, a := range items {
		r, err := parser.ParseRange(a.AffectedRange)
		if err == nil && r.Matches(v) && !a.Withdrawn {
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Severity > out[j].Severity })
	return out
}
func Score(a domain.Advisory) float64 {
	score := a.Severity
	if a.EPSS > score/10 {
		score += a.EPSS * 2
	}
	if strings.Contains(strings.ToLower(a.ID), "critical") {
		score += 1
	}
	if score > 10 {
		return 10
	}
	return score
}
