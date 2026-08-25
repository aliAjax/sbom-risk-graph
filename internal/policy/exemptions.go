package policy

import (
	"sync"
	"time"
)

type Exemption struct {
	ID         string    `json:"id"`
	PolicyID   string    `json:"policy_id"`
	Component  string    `json:"component"`
	Reason     string    `json:"reason"`
	ExpiresAt  time.Time `json:"expires_at"`
	ApprovedBy string    `json:"approved_by"`
}
type Exemptions struct {
	mu     sync.RWMutex
	values map[string]Exemption
}

func NewExemptions() *Exemptions { return &Exemptions{values: make(map[string]Exemption)} }
func (e *Exemptions) Put(v Exemption) error {
	if v.ID == "" || v.PolicyID == "" || v.Component == "" || v.ExpiresAt.IsZero() {
		return fmtError("exemption fields required")
	}
	e.mu.Lock()
	e.values[v.ID] = v
	e.mu.Unlock()
	return nil
}
func (e *Exemptions) Active(policy, component string, now time.Time) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, v := range e.values {
		if v.PolicyID == policy && v.Component == component && now.Before(v.ExpiresAt) {
			return true
		}
	}
	return false
}
func (e *Exemptions) Expired(now time.Time) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]string, 0)
	for id, v := range e.values {
		if !now.Before(v.ExpiresAt) {
			out = append(out, id)
		}
	}
	return out
}
func fmtError(s string) error { return &validationError{s} }

type validationError struct{ s string }

func (e *validationError) Error() string { return e.s }
