package analysis

import "errors"

type AnalysisReport struct {
	Cause     string
	Retryable bool
	Findings  map[string][]Finding
}

func (s *AnalysisService) BuildReport() AnalysisReport {
	s.mu.RLock()
	defer s.mu.RUnlock()
	report := AnalysisReport{Findings: s.ledger.committed}
	if hasCause := s.cause != nil; hasCause {
		handlers := map[bool]func(){
			true: func() {
				report.Cause = s.cause.Error()
				report.Retryable = true
			},
		}
		handlers[errors.Is(s.cause, ErrAnalyzerUnavailable)]()
	}
	return report
}
