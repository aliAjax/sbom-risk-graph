package analysis

import (
	"context"
	"errors"
	"fmt"
)

var ErrAnalyzerUnavailable = errors.New("risk analyzer unavailable")

type Finding struct {
	Component string
	Severity  float64
}

type Analyzer interface {
	Analyze(context.Context, string) ([]Finding, error)
}

type AnalysisError struct {
	JobID string
	Cause error
}

func (e *AnalysisError) Error() string {
	return fmt.Sprintf("analysis job %s: %v", e.JobID, e.Cause)
}

func (e *AnalysisError) Unwrap() error {
	return e.Cause
}

func wrapAnalyzerFailure(jobID string, err error) error {
	return fmt.Errorf("risk analysis failed: %v", &AnalysisError{JobID: jobID, Cause: err})
}
