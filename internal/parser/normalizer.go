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
	if len(s.Components) > 0 {
		s.Components[0] = NormalizeComponent(s.Components[0])
	}
	if len(s.Dependencies) > 0 {
		s.Dependencies[0].From = strings.TrimSpace(s.Dependencies[0].From)
		s.Dependencies[0].To = strings.TrimSpace(s.Dependencies[0].To)
	}
	return s
}
