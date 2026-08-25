package evidence

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
	"sync"
	"time"
)

type Signer interface {
	Sign([]byte) ([]byte, error)
	Verify([]byte, []byte) error
}
type HMACSigner struct{ Secret []byte }

func (s HMACSigner) Sign(data []byte) ([]byte, error) {
	h := hmac.New(sha256.New, s.Secret)
	_, _ = h.Write(data)
	return h.Sum(nil), nil
}
func (s HMACSigner) Verify(data, sig []byte) error {
	expected, _ := s.Sign(data)
	if !hmac.Equal(expected, sig) {
		return fmt.Errorf("signature mismatch")
	}
	return nil
}

type Store struct {
	mu     sync.RWMutex
	values map[string]domain.Evidence
	signer Signer
}

func New(signer Signer) *Store {
	return &Store{values: make(map[string]domain.Evidence), signer: signer}
}
func (s *Store) Create(id, product, kind, digest string, expires time.Time) (domain.Evidence, error) {
	if id == "" || product == "" || digest == "" {
		return domain.Evidence{}, fmt.Errorf("evidence fields required")
	}
	payload := []byte(id + "|" + product + "|" + kind + "|" + digest)
	sig, err := s.signer.Sign(payload)
	if err != nil {
		return domain.Evidence{}, err
	}
	e := domain.Evidence{ID: id, ProductID: product, Kind: kind, Digest: hex.EncodeToString(sig), PayloadDigest: digest, Signature: hex.EncodeToString(sig), Status: domain.Valid, ExpiresAt: expires, CreatedAt: time.Now().UTC()}
	s.mu.Lock()
	s.values[id] = e
	s.mu.Unlock()
	return e, nil
}
func (s *Store) Get(id string) (domain.Evidence, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.values[id]
	return e, ok
}
func (s *Store) Revoke(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.values[id]
	if !ok {
		return fmt.Errorf("evidence not found")
	}
	e.Status = domain.Revoked
	s.values[id] = e
	return nil
}
func (s *Store) Valid(product string, now time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.values {
		if e.ProductID == product && e.Status == domain.Valid && (e.ExpiresAt.IsZero() || now.Before(e.ExpiresAt)) {
			return true
		}
	}
	return false
}
