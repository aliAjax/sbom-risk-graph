package evidence

import (
	"example.com/sbom-risk-graph/internal/domain"
	"sync"
	"time"
)

type Audit struct {
	mu     sync.Mutex
	events []domain.AuditEvent
}

func NewAudit() *Audit { return &Audit{} }
func (a *Audit) Append(id, action, subject string) domain.AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	prev := ""
	if len(a.events) > 0 {
		prev = a.events[len(a.events)-1].Hash
	}
	e := domain.AuditEvent{ID: id, Action: action, Subject: subject, At: time.Now().UTC(), PreviousHash: prev}.Seal()
	a.events = append(a.events, e)
	return e
}
func (a *Audit) All() []domain.AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.events
}
