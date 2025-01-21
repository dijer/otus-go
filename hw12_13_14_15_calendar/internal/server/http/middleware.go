package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func (h *Handler) loggingMiddleware(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		next(w, r, pathParams)

		start := time.Now()
		h.logger.Info(
			fmt.Sprintf("%s: %s\n", "method", r.Method),
			fmt.Sprintf("%s: %s\n", "path", r.URL.Path),
			fmt.Sprintf("%s: %s\n", "latency", time.Since(start).String()),
			fmt.Sprintf("%s: %s\n", "ip", r.RemoteAddr),
		)
	}
}
