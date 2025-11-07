// Package httpserver provides HTTP server functionality with middleware support.
package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/platforma-dev/platforma/log"
)

type handleGroup = HandlerGroup

// HTTPServer represents an HTTP server with middleware support and graceful shutdown.
type HTTPServer struct {
	*handleGroup
	port            string
	shutdownTimeout time.Duration
}

// New creates a new HTTPServer instance with the specified port and shutdown timeout.
func New(port string, shutdownTimeout time.Duration) *HTTPServer {
	return &HTTPServer{handleGroup: NewHandlerGroup(), port: port, shutdownTimeout: shutdownTimeout}
}

// Run starts the HTTP server and handles graceful shutdown on interrupt signals.
func (s *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:              ":" + s.port,
		Handler:           wrapHandlerInMiddleware(s.mux, s.middlewares),
		ReadHeaderTimeout: 1 * time.Second,
	}

	go func() {
		log.InfoContext(ctx, "starting http server", "address", server.Addr)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.ErrorContext(ctx, "HTTP server error", "error", err)
		}
		log.InfoContext(ctx, "stopped serving new connections")
	}()

	<-ctx.Done()

	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, s.shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to gracefully shutdown HTTP server: %w", err)
	}
	log.InfoContext(ctx, "graceful shutdown completed")

	return nil
}

// Healthcheck returns health check information for the HTTP server.
func (s *HTTPServer) Healthcheck(_ context.Context) any {
	return map[string]any{
		"port": s.port,
	}
}
