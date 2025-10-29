package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mishankov/platforma/log"
)

type handleGroup = HandlerGroup

type HttpServer struct {
	*handleGroup
	port            string
	shutdownTimeout time.Duration
}

func New(port string, shutdownTimeout time.Duration) *HttpServer {
	return &HttpServer{handleGroup: NewHandlerGroup(), port: port, shutdownTimeout: shutdownTimeout}
}

func (s *HttpServer) Run(ctx context.Context) error {
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
		log.InfoContext(ctx, "stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, s.shutdownTimeout)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to gracefully shutdown HTTP server: %w", err)
	}
	log.InfoContext(ctx, "graceful shutdown completed.")

	return nil
}

func (s *HttpServer) Healthcheck(ctx context.Context) any {
	return map[string]any{
		"port": s.port,
	}
}
