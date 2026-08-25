package policy

import (
	"example.com/sbom-risk-graph/internal/domain"
	"testing"
	"time"
)

func TestPolicyRevisionTransitionIsMonotonic(t *testing.T) {
	store := NewStore()
	if err := store.Put(domain.Policy{ID: "policy-a", Name: "first"}); err != nil {
		t.Fatal(err)
	}
	if err := store.Put(domain.Policy{ID: "policy-a", Name: "second"}); err != nil {
		t.Fatal(err)
	}
	got, _ := store.Get("policy-a")
	if got.Revision != 2 {
		t.Fatalf("latest revision = %d, want 2", got.Revision)
	}
	if err := store.Put(domain.Policy{ID: "policy-a", Name: "stale", Revision: 1}); err == nil {
		t.Fatal("stale revision was accepted")
	}
}

func TestLatestPolicyDrivesEvaluation(t *testing.T) {
	evaluator := New()
	component := domain.Component{Name: "lib", Version: "1"}
	advisories := func(domain.Component) []domain.Advisory {
		return []domain.Advisory{{ID: "ADV", Severity: 8}}
	}
	old := domain.Policy{ID: "policy-a", MaxSeverity: 10, MaxEPSS: 1, Revision: 1}
	latest := domain.Policy{ID: "policy-a", MaxSeverity: 2, MaxEPSS: 1, Revision: 2}
	if got := evaluator.Evaluate(old, []domain.Component{component}, advisories, true); len(got) != 0 {
		t.Fatalf("old policy unexpectedly rejected component: %#v", got)
	}
	if got := evaluator.Evaluate(latest, []domain.Component{component}, advisories, true); len(got) != 1 {
		t.Fatalf("latest policy produced %d violations, want 1", len(got))
	}
}

func TestExpiredExemptionCannotSuppressViolation(t *testing.T) {
	now := time.Unix(100, 0)
	exemptions := NewExemptions()
	err := exemptions.Put(Exemption{ID: "ex-1", PolicyID: "policy-a", Component: "lib@1", ExpiresAt: now})
	if err != nil {
		t.Fatal(err)
	}
	if exemptions.Active("policy-a", "lib@1", now) {
		t.Fatal("exemption remained active at its expiry instant")
	}
}

func TestRevisionHistoryRemainsConsistent(t *testing.T) {
	history := NewHistory()
	history.Add(domain.Policy{ID: "policy-a", Revision: 2, Name: "new"})
	history.Add(domain.Policy{ID: "policy-a", Revision: 1, Name: "stale"})
	latest, ok := history.Latest("policy-a")
	if !ok || latest.Revision != 2 || latest.Name != "new" {
		t.Fatalf("latest policy = %#v, want revision 2", latest)
	}
}
