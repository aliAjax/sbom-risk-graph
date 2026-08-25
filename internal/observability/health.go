package observability

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

type Health struct{ ready atomic.Bool }

func New() *Health                { h := &Health{}; h.ready.Store(true); return h }
func (h *Health) SetReady(v bool) { h.ready.Store(v) }
func (h *Health) Handler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	if !h.ready.Load() {
		status = http.StatusServiceUnavailable
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"status": map[bool]string{true: "ready", false: "not_ready"}[h.ready.Load()], "time": time.Now().UTC()})
}
