package advisory

import (
	"example.com/sbom-risk-graph/internal/domain"
	"example.com/sbom-risk-graph/internal/parser"
	"fmt"
)

type Match struct {
	Advisory  domain.Advisory  `json:"advisory"`
	Component domain.Component `json:"component"`
	Reason    string           `json:"reason"`
	Fixed     bool             `json:"fixed"`
}

func (e *Engine) Explain(c domain.Component) []Match {
	out := make([]Match, 0)
	v, err := parser.ParseVersion(c.Version)
	if err != nil {
		return out
	}
	for _, a := range e.advisories[c.PURL] {
		r, err := parser.ParseRange(a.AffectedRange)
		if err != nil || !r.Matches(v) || a.Withdrawn {
			continue
		}
		fixed := false
		if a.FixedVersion != "" {
			if fv, err := parser.ParseVersion(a.FixedVersion); err == nil {
				fixed = v.Compare(fv) >= 0
			}
		}
		out = append(out, Match{Advisory: a, Component: c, Reason: fmt.Sprintf("version %s matches %s", c.Version, a.AffectedRange), Fixed: fixed})
	}
	return out
}
func (e *Engine) Candidates(c domain.Component) []string {
	out := make([]string, 0)
	for _, a := range e.Explain(c) {
		if a.Advisory.FixedVersion != "" {
			out = append(out, a.Advisory.FixedVersion)
		}
	}
	return out
}
