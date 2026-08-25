package parser

import (
	"example.com/sbom-risk-graph/internal/domain"
	"strings"
)

func NormalizeComponent(c domain.Component) domain.Component {
	c.Name = strings.TrimSpace(c.Name)
	c.Version = strings.TrimSpace(c.Version)
	c.PURL = strings.ToLower(strings.TrimSpace(c.PURL))
	c.License = strings.TrimSpace(c.License)
	if c.PURL == "" {
		c.PURL = "pkg:generic/" + strings.ToLower(strings.ReplaceAll(c.Name, " ", "-"))
	}
	return c
}
func NormalizeSBOM(s domain.SBOM) domain.SBOM {
	for i := range s.Components {
		s.Components[i] = NormalizeComponent(s.Components[i])
	}
	for i := range s.Dependencies {
		s.Dependencies[i].From = strings.TrimSpace(s.Dependencies[i].From)
		s.Dependencies[i].To = strings.TrimSpace(s.Dependencies[i].To)
	}
	return s
}
