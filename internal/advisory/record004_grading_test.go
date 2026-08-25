package advisory

import (
	"errors"
	"example.com/sbom-risk-graph/internal/domain"
	"example.com/sbom-risk-graph/internal/parser"
	"testing"
)

func TestMalformedComponentVersionPreservesErrorChain(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("version parser panicked: %v", recovered)
		}
	}()
	_, err := parser.ParseVersion("release")
	if !errors.Is(err, parser.ErrInvalidVersion) {
		t.Fatalf("invalid version identity lost: %v", err)
	}
}

func TestMalformedAffectedRangeIsNotSilentlyDropped(t *testing.T) {
	_, err := parser.ParseRange(">=release")
	if !errors.Is(err, parser.ErrInvalidRange) || !errors.Is(err, parser.ErrInvalidVersion) {
		t.Fatalf("range error chain incomplete: %v", err)
	}
}

func TestExplainAndMatchShareClassification(t *testing.T) {
	e := New()
	_ = e.Add(domain.Advisory{ID: "ADV-1", ComponentPURL: "pkg:generic/a", AffectedRange: ">=release", Severity: 7})
	c := domain.Component{Name: "a", Version: "1.0.0", PURL: "pkg:generic/a"}
	_ = e.Match(c)
	matchErr := e.LastError()
	_ = e.Explain(c)
	explainErr := e.LastError()
	if !errors.Is(matchErr, parser.ErrInvalidRange) || !errors.Is(explainErr, parser.ErrInvalidRange) {
		t.Fatalf("classification diverged: match=%v explain=%v", matchErr, explainErr)
	}
}

func TestFeedPreservesRangeFailure(t *testing.T) {
	feed := []byte(`[{"id":"ADV-1","component_purl":"pkg:generic/a","affected_range":">=release","severity":7}]`)
	if _, err := ParseFeed(feed); !errors.Is(err, parser.ErrInvalidRange) {
		t.Fatalf("feed swallowed invalid range: %v", err)
	}
}
