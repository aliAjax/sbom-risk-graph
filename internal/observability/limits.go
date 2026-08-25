package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func RequestLimits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > 64<<20 {
			http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
			return
		}
		ctx, cancel := contextWithTimeout(r, 20*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func contextWithTimeout(r *http.Request, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), d)
}
func ValidatePath(v string) error {
	if v == "" || len(v) > 512 {
		return fmt.Errorf("invalid path")
	}
	return nil
}
