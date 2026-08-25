package importer

import (
	"context"
	"sync"

	"example.com/sbom-risk-graph/internal/domain"
)

type Session struct {
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	batch     *Batch
	state     SessionState
	cause     error
	published []domain.SBOM
}

func NewSession(ctx context.Context, current []domain.SBOM) *Session {
	if ctx == nil {
		ctx = context.Background()
	}
	sessionCtx, cancel := context.WithCancel(ctx)
	return &Session{
		ctx:       sessionCtx,
		cancel:    cancel,
		batch:     NewBatch(current),
		state:     SessionPending,
		published: cloneSBOMs(current),
	}
}

func (s *Session) Cancel() {
	s.cancel()
}

func (s *Session) begin() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state != SessionPending {
		return ErrSessionAlreadyRun
	}
	s.state = SessionRunning
	return nil
}

func (s *Session) finish(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state.Terminal() {
		return
	}
	s.state = terminalState(err)
	if err != nil && s.cause == nil {
		s.cause = err
	}
	if err == nil {
		s.published = s.batch.Snapshot()
	}
}
