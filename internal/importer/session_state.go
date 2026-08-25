package importer

import (
	"context"
	"errors"
)

type SessionState string

const (
	SessionPending   SessionState = "pending"
	SessionRunning   SessionState = "running"
	SessionSucceeded SessionState = "succeeded"
	SessionFailed    SessionState = "failed"
	SessionCanceled  SessionState = "canceled"
)

var ErrSessionAlreadyRun = errors.New("import session already started")

func (s SessionState) Terminal() bool {
	return s == SessionSucceeded || s == SessionFailed || s == SessionCanceled
}

func terminalState(err error) SessionState {
	if err == nil {
		return SessionSucceeded
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return SessionCanceled
	}
	return SessionFailed
}
