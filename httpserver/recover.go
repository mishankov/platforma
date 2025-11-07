package httpserver

import (
	"net/http"

	"github.com/platforma-dev/platforma/log"
)

// RecoverMiddleware is a middleware that recovers from panics in HTTP handlers.
// It catches panics, logs the error, and returns an HTTP 500 response to the client.
type RecoverMiddleware struct{}

// NewRecoverMiddleware creates a new instance of RecoverMiddleware.
func NewRecoverMiddleware() *RecoverMiddleware {
	return &RecoverMiddleware{}
}

// Wrap implements the Middleware interface by wrapping the provided handler
// with panic recovery logic.
func (m *RecoverMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with request context
				log.ErrorContext(r.Context(), "panic recovered", "error", err, "method", r.Method, "path", r.URL.Path)

				// Write HTTP 500 response
				w.WriteHeader(http.StatusInternalServerError)
				_, writeErr := w.Write([]byte("Internal Server Error"))
				if writeErr != nil {
					log.ErrorContext(r.Context(), "failed to write error response", "error", writeErr)
				}
			}
		}()

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
