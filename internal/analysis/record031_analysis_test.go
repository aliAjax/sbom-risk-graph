package analysis

import (
	"context"
	"errors"
	"testing"
)

func TestAnalysisErrorKeepsUnavailableCause(t *testing.T) {
	err := wrapAnalyzerFailure("job-provider", ErrAnalyzerUnavailable)
	if !errors.Is(err, ErrAnalyzerUnavailable) {
		t.Fatalf("wrapped error lost analyzer cause: %v", err)
	}
}

func TestFailedAnalysisReturnsFailure(t *testing.T) {
	service := NewAnalysisService(stubAnalyzer{err: ErrAnalyzerUnavailable})
	if err := service.Execute(context.Background(), "job-service"); err == nil {
		t.Fatal("analysis failure was returned as success")
	}
}

func TestCancelDoesNotCommitFailedAnalysis(t *testing.T) {
	ledger := NewResultLedger()
	ledger.Stage("job-ledger")
	ledger.Cancel("job-ledger")
	if _, exists := ledger.Snapshot()["job-ledger"]; exists {
		t.Fatal("canceled analysis was published to committed results")
	}
}

func TestFailureReportIsSafeAndDetached(t *testing.T) {
	failed := NewAnalysisService(stubAnalyzer{err: ErrAnalyzerUnavailable})
	_ = failed.Execute(context.Background(), "job-report-failure")
	if _, panicked := buildReportSafely(failed); panicked {
		t.Fatal("building a failure report panicked")
	}

	success := NewAnalysisService(stubAnalyzer{findings: []Finding{{Component: "pkg:report", Severity: 9.1}}})
	if err := success.Execute(context.Background(), "job-report-success"); err != nil {
		t.Fatal(err)
	}
	report, panicked := buildReportSafely(success)
	if panicked {
		t.Fatal("building a success report panicked")
	}
	report.Findings["job-report-success"][0].Severity = 0
	if got := success.BuildReport().Findings["job-report-success"][0].Severity; got != 9.1 {
		t.Fatalf("report mutation changed committed findings: %v", got)
	}
}

func buildReportSafely(service *AnalysisService) (report AnalysisReport, panicked bool) {
	defer func() {
		panicked = recover() != nil
	}()
	return service.BuildReport(), false
}
