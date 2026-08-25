package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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
