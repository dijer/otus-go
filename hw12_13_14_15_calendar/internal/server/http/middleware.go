package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func (h *Handler) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)

		start := time.Now()
		h.logger.Info(
			fmt.Sprintf("%s: %s\n", "method", r.Method),
			fmt.Sprintf("%s: %s\n", "path", r.URL.Path),
			fmt.Sprintf("%s: %s\n", "latency", time.Since(start).String()),
			fmt.Sprintf("%s: %s\n", "ip", r.RemoteAddr),
		)
	}
}
