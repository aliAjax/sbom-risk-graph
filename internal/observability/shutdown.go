package observability

import (
	"context"
	"time"
)

func ShutdownContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
