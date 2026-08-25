package analysis

import (
	"context"
	"errors"
	"testing"
)

type stubAnalyzer struct {
	findings []Finding
	err      error
}

func (s stubAnalyzer) Analyze(context.Context, string) ([]Finding, error) {
	return s.findings, s.err
}

func TestAnalysisWorkflow(t *testing.T) {
	t.Run("failure remains classifiable", func(t *testing.T) {
		service := NewAnalysisService(stubAnalyzer{err: ErrAnalyzerUnavailable})
		err := service.Execute(context.Background(), "job-failure")
		if !errors.Is(err, ErrAnalyzerUnavailable) {
			t.Fatalf("analysis error was not preserved: %v", err)
		}
	})

	t.Run("failure is returned and aborted", func(t *testing.T) {
		service := NewAnalysisService(stubAnalyzer{err: ErrAnalyzerUnavailable})
		if err := service.Execute(context.Background(), "job-abort"); err == nil {
			t.Fatal("analysis failure was swallowed")
		}
		if findings := service.ledger.Snapshot(); len(findings) != 0 {
			t.Fatalf("failed analysis was committed: %#v", findings)
		}
	})

	t.Run("failure report is safe", func(t *testing.T) {
		service := NewAnalysisService(stubAnalyzer{err: ErrAnalyzerUnavailable})
		_ = service.Execute(context.Background(), "job-report")
		report := service.BuildReport()
		if !report.Retryable || report.Cause == "" {
			t.Fatalf("failure report lost classification: %#v", report)
		}
	})

	t.Run("success report is detached", func(t *testing.T) {
		service := NewAnalysisService(stubAnalyzer{findings: []Finding{{Component: "pkg:a", Severity: 8.2}}})
		if err := service.Execute(context.Background(), "job-success"); err != nil {
			t.Fatal(err)
		}
		report := service.BuildReport()
		report.Findings["job-success"][0].Severity = 0
		if got := service.BuildReport().Findings["job-success"][0].Severity; got != 8.2 {
			t.Fatalf("report mutation changed committed findings: %v", got)
		}
	})
}
