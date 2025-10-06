package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/yeftaz/susano.id/api/pkg/logger"
	"github.com/yeftaz/susano.id/api/pkg/response"
)

// Recovery middleware recovers from panics and logs the error
func Recovery(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic with stack trace
					logger.Error("Panic recovered",
						"error", err,
						"stack", string(debug.Stack()),
						"method", r.Method,
						"path", r.URL.Path,
					)

					// Return 500 error
					response.Error(w, http.StatusInternalServerError, "Internal server error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
