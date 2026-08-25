package waiver

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type State string

const (
	StatePending  State = "pending"
	StateApproved State = "approved"
	StateFailed   State = "failed"
)

type Request struct {
	ID        string
	ProductID string
	Reason    string
}

type Decision struct {
	WaiverID string
	State    State
	At       time.Time
}

type Evaluator interface {
	Evaluate(context.Context, Request) error
}

type Session struct {
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	evaluator Evaluator
	ledger    *Ledger
	state     State
	cause     error
}

func NewSession(parent context.Context, evaluator Evaluator) *Session {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &Session{ctx: ctx, cancel: cancel, evaluator: evaluator, ledger: NewLedger(), state: StatePending}
}

func (s *Session) Cancel() { s.cancel() }

func (s *Session) Review(req Request) error {
	if req.ID == "" || req.ProductID == "" {
		return fmt.Errorf("waiver id and product id are required")
	}
	s.ledger.Stage(Decision{WaiverID: req.ID, State: StatePending, At: time.Now().UTC()})
	if err := s.evaluator.Evaluate(s.ctx, req); err != nil {
		s.ledger.Rollback(req.ID)
		wrapped := wrapReviewError(req.ID, err)
		s.finish(StateFailed, wrapped)
		return wrapped
	}
	if err := s.ledger.Commit(req.ID); err != nil {
		wrapped := wrapReviewError(req.ID, err)
		s.finish(StateFailed, wrapped)
		return wrapped
	}
	s.finish(StateApproved, nil)
	return nil
}

func (s *Session) finish(state State, cause error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = state
	s.cause = cause
}
