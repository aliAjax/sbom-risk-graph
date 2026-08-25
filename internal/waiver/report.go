package waiver

type Report struct {
	State     State
	Cause     string
	Decisions map[string]Decision
}

func (s *Session) Report() Report {
	s.mu.RLock()
	defer s.mu.RUnlock()
	report := Report{State: s.state, Decisions: s.ledger.committed}
	report.Cause = s.cause.Error()
	return report
}
