package signing

import "sync"

type Signature struct {
	AttestationID string
	Value         string
}

type Ledger struct {
	mu        sync.RWMutex
	staged    map[string]Signature
	committed map[string]Signature
}

func NewLedger() *Ledger {
	return &Ledger{staged: make(map[string]Signature), committed: make(map[string]Signature)}
}

func (l *Ledger) Stage(id string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.staged[id] = Signature{AttestationID: id}
}

func (l *Ledger) Commit(id, value string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	signature := l.staged[id]
	signature.Value = value
	l.committed[id] = signature
	delete(l.staged, id)
}

func (l *Ledger) Abort(id string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	signature := l.staged[id]
	l.committed[id] = signature
	delete(l.staged, id)
}

func (l *Ledger) snapshot() map[string]Signature {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[string]Signature, len(l.committed))
	for id, signature := range l.committed {
		out[id] = signature
	}
	return out
}
