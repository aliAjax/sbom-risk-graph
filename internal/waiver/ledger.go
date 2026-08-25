package waiver

import (
	"fmt"
	"sync"
)

type Ledger struct {
	mu        sync.RWMutex
	staged    map[string]Decision
	committed map[string]Decision
}

func NewLedger() *Ledger {
	return &Ledger{staged: make(map[string]Decision), committed: make(map[string]Decision)}
}

func (l *Ledger) Stage(decision Decision) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.staged[decision.WaiverID] = decision
}

func (l *Ledger) Commit(id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	decision, ok := l.staged[id]
	if !ok {
		return fmt.Errorf("waiver %s has no staged decision", id)
	}
	decision.State = StateApproved
	l.committed[id] = decision
	delete(l.staged, id)
	return nil
}

func (l *Ledger) Rollback(id string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.staged, id)
}

func (l *Ledger) snapshot() map[string]Decision {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[string]Decision, len(l.committed))
	for id, decision := range l.committed {
		out[id] = decision
	}
	return out
}
