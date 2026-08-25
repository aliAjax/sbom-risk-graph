package signing

type Report struct {
	Cause      string
	Signatures map[string]Signature
}

func (s *Service) Report() Report {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Report{
		Cause:      s.cause.Error(),
		Signatures: s.ledger.committed,
	}
}
