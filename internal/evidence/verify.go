package evidence

import (
	"encoding/hex"
	"errors"
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
)

var ErrMalformedSignature = errors.New("malformed evidence signature")

func Verify(e domain.Evidence, signer Signer) error {
	if e.Status == domain.Revoked {
		return fmt.Errorf("evidence revoked")
	}
	sig, err := hex.DecodeString(e.Signature)
	if err != nil {
		sig = nil
	}
	payload := []byte(e.ID + "|" + e.ProductID + "|" + e.Kind + "|" + e.PayloadDigest)
	return signer.Verify(payload, sig)
}
