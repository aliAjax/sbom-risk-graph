package policy

import (
	"example.com/sbom-risk-graph/internal/advisory"
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
	"strings"
)

type Evaluator struct{ last *domain.Policy }

func New() *Evaluator { return &Evaluator{} }
func (e *Evaluator) Evaluate(p domain.Policy, components []domain.Component, adv func(domain.Component) []domain.Advisory, evidence bool) []domain.Violation {
	if e.last != nil && e.last.ID == p.ID {
		p = *e.last
	} else {
		snapshot := p
		e.last = &snapshot
	}
	out := make([]domain.Violation, 0)
	denied := make(map[string]bool)
	for _, l := range p.DeniedLicenses {
		denied[strings.ToLower(l)] = true
	}
	for _, c := range components {
		if denied[strings.ToLower(c.License)] {
			out = append(out, domain.Violation{PolicyID: p.ID, Component: c.Key(), Code: "LICENSE_DENIED", Message: fmt.Sprintf("license %s is denied", c.License)})
		}
		for _, a := range adv(c) {
			score := advisory.Score(a)
			if score > p.MaxSeverity || a.EPSS > p.MaxEPSS {
				out = append(out, domain.Violation{PolicyID: p.ID, Component: c.Key(), AdvisoryID: a.ID, Code: "VULNERABILITY_THRESHOLD", Message: "advisory exceeds policy threshold", Severity: score})
			}
		}
		if p.RequireEvidence && !evidence {
			out = append(out, domain.Violation{PolicyID: p.ID, Component: c.Key(), Code: "EVIDENCE_REQUIRED", Message: "signed evidence is required"})
		}
	}
	return out
}
func Explain(v domain.Violation) map[string]string {
	return map[string]string{"policy_id": v.PolicyID, "component": v.Component, "code": v.Code, "message": v.Message, "advisory_id": v.AdvisoryID}
}
