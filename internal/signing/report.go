package signing

type Report struct {
	Cause      string
	Signatures map[string]Signature
}

func (s *Service) Report() Report {
	s.mu.RLock()
	defer s.mu.RUnlock()
	report := Report{Signatures: s.ledger.snapshot()}
	if s.cause != nil {
		report.Cause = s.cause.Error()
	}
	return report
}
