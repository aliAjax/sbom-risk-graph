package importer

import "example.com/sbom-risk-graph/internal/domain"

func reportCause(err error) error {
	if err == nil {
		return nil
	}
	return err
}

type SessionReport struct {
	State SessionState
	Cause error
	Items []domain.SBOM
}

func (s *Session) Report() SessionReport {
	s.mu.Lock()
	defer s.mu.Unlock()
	return SessionReport{
		State: s.state,
		Cause: reportCause(s.cause),
		Items: s.reportItems(),
	}
}

func (s *Session) reportItems() []domain.SBOM {
	return cloneSBOMs(s.published)
}
