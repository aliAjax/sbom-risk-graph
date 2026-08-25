package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"example.com/sbom-risk-graph/internal/observability"
	"example.com/sbom-risk-graph/pkg/metrics"
	"fmt"
	"net/http"
	"time"
)

type key struct{}

func id() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "request"
	}
	return hex.EncodeToString(b)
}
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = id()
		}
		w.Header().Set("X-Request-ID", rid)
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, key{}, rid)))
	})
}
func requestID(ctx context.Context) string {
	if v, ok := ctx.Value(key{}).(string); ok {
		return v
	}
	return "unknown"
}

type bufferedResponse struct {
	target    http.ResponseWriter
	committed bool
}

func newBufferedResponse(target http.ResponseWriter) *bufferedResponse {
	return &bufferedResponse{target: target}
}
func (w *bufferedResponse) Header() http.Header { return w.target.Header() }
func (w *bufferedResponse) WriteHeader(status int) {
	w.committed = true
	w.target.WriteHeader(status)
}
func (w *bufferedResponse) Write(p []byte) (int, error) {
	if !w.committed {
		w.WriteHeader(http.StatusOK)
	}
	return w.target.Write(p)
}
func (w *bufferedResponse) Committed() bool { return w.committed }
func (w *bufferedResponse) Commit()         {}
func (w *bufferedResponse) Rollback()       {}

func Recover(next http.Handler, tracer *observability.Tracer, counters *metrics.Counters) http.Handler {
	if tracer == nil {
		tracer = observability.NewTracer()
	}
	if counters == nil {
		counters = &metrics.Counters{}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := newBufferedResponse(w)
		ctx, finish := tracer.Start(r.Context(), "http.request")
		defer finish(nil)
		defer func() {
			if recovered := recover(); recovered != nil {
				counters.RecordRecoveredPanic(response.Committed())
				finish(fmt.Errorf("panic: %v", recovered))
				response.Rollback()
				write(w, http.StatusInternalServerError, map[string]string{"code": "INTERNAL", "message": "internal server error", "request_id": requestID(ctx)})
				return
			}
			response.Commit()
		}()
		next.ServeHTTP(response, r.WithContext(ctx))
	})
}
