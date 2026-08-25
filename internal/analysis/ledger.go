package analysis

import "sync"

type ResultLedger struct {
	mu        sync.RWMutex
	staged    map[string][]Finding
	committed map[string][]Finding
}

func NewResultLedger() *ResultLedger {
	return &ResultLedger{
		staged:    make(map[string][]Finding),
		committed: make(map[string][]Finding),
	}
}

func (l *ResultLedger) Stage(jobID string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.staged[jobID] = nil
}

func (l *ResultLedger) Commit(jobID string, findings []Finding) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.committed[jobID] = append([]Finding(nil), findings...)
	delete(l.staged, jobID)
}

func (l *ResultLedger) Cancel(jobID string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	findings := l.staged[jobID]
	l.committed[jobID] = findings
	delete(l.staged, jobID)
}

func (l *ResultLedger) Snapshot() map[string][]Finding {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[string][]Finding, len(l.committed))
	for jobID, findings := range l.committed {
		out[jobID] = append([]Finding(nil), findings...)
	}
	return out
}
