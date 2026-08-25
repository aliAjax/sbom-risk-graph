package repository

import (
	"encoding/json"
	"example.com/sbom-risk-graph/internal/domain"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type State struct {
	Products   map[string]domain.Product  `json:"products"`
	SBOMs      map[string]domain.SBOM     `json:"sboms"`
	Policies   map[string]domain.Policy   `json:"policies"`
	Advisories map[string]domain.Advisory `json:"advisories"`
}
type Store struct {
	mu    sync.Mutex
	path  string
	state State
}

func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, err
	}
	s := &Store{path: filepath.Join(dir, "state.json"), state: State{Products: map[string]domain.Product{}, SBOMs: map[string]domain.SBOM{}, Policies: map[string]domain.Policy{}, Advisories: map[string]domain.Advisory{}}}
	b, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &s.state); err != nil {
		return nil, fmt.Errorf("state decode: %w", err)
	}
	return s, nil
}
func (s *Store) Snapshot() State { s.mu.Lock(); defer s.mu.Unlock(); return s.state }
func (s *Store) SaveProduct(p domain.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveWithRollback(func() { s.state.Products[p.ID] = p })
}
func (s *Store) SaveSBOM(b domain.SBOM) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveWithRollback(func() { s.state.SBOMs[b.ID] = b })
}
func (s *Store) SavePolicy(p domain.Policy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveWithRollback(func() { s.state.Policies[p.ID] = p })
}
func (s *Store) SaveAdvisory(a domain.Advisory) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveWithRollback(func() { s.state.Advisories[a.ID] = a })
}
func (s *Store) saveWithRollback(mutate func()) (err error) {
	before := s.state
	defer func() {
		if err == nil {
			s.state = before
		}
	}()
	mutate()
	return s.persist()
}
func (s *Store) persist() error {
	b, err := json.MarshalIndent(s.state, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o640); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
