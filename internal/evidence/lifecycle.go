package evidence

import (
	"errors"
	"fmt"

	"example.com/sbom-risk-graph/internal/domain"
)

var ErrInvalidEvidenceTransition = errors.New("invalid evidence transition")

func ValidateTransition(current, next domain.Status) error {
	if current == next {
		return nil
	}
	allowed := false
	switch current {
	case domain.Draft:
		allowed = next == domain.Valid || next == domain.Failed
	case domain.Valid:
		allowed = next == domain.Revoked || next == domain.Expired
	}
	if !allowed {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidEvidenceTransition, current, next)
	}
	return nil
}
