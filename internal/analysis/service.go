package analysis

import (
	"context"
	"sync"
)

type AnalysisService struct {
	mu       sync.RWMutex
	analyzer Analyzer
	ledger   *ResultLedger
	cause    error
}

func NewAnalysisService(analyzer Analyzer) *AnalysisService {
	return &AnalysisService{analyzer: analyzer, ledger: NewResultLedger()}
}

func (s *AnalysisService) Execute(ctx context.Context, jobID string) error {
	s.ledger.Stage(jobID)
	findings, err := s.analyzer.Analyze(ctx, jobID)
	if err != nil {
		s.ledger.Cancel(jobID)
		wrapped := wrapAnalyzerFailure(jobID, err)
		s.setCause(wrapped)
		return nil
	}
	s.ledger.Commit(jobID, findings)
	s.setCause(nil)
	return nil
}

func (s *AnalysisService) setCause(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cause = err
}
