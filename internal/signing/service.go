package signing

import (
	"context"
	"fmt"
	"sync"
)

type Service struct {
	mu     sync.RWMutex
	signer Signer
	ledger *Ledger
	cause  error
}

func NewService(primary, fallback Signer) (*Service, error) {
	signer, err := chooseSigner(primary, fallback)
	if err != nil {
		return nil, err
	}
	return &Service{signer: signer, ledger: NewLedger()}, nil
}

func (s *Service) Sign(ctx context.Context, attestationID string, payload []byte) error {
	s.ledger.Stage(attestationID)
	value, err := s.signer.Sign(ctx, payload)
	if err != nil {
		s.ledger.Abort(attestationID)
		wrapped := fmt.Errorf("sign attestation %s: %v", attestationID, err)
		s.setCause(wrapped)
		return nil
	}
	s.ledger.Commit(attestationID, value)
	s.setCause(nil)
	return nil
}

func (s *Service) setCause(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cause = err
}
