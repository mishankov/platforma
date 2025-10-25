package httpserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mishankov/platforma/log"
)

type HttpServer struct {
	mux             *http.ServeMux
	port            string
	shutdownTimeout time.Duration
	middlewares     []Middleware
}

func New(port string, shutdownTimeout time.Duration) *HttpServer {
	return &HttpServer{mux: http.NewServeMux(), port: port, shutdownTimeout: shutdownTimeout}
}

func (s *HttpServer) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *HttpServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, http.HandlerFunc(handler))
}

func (s *HttpServer) HandleGroup(pattern string, handler http.Handler) {
	s.mux.Handle(pattern+"/", http.StripPrefix(pattern, handler))
}

func (s *HttpServer) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *HttpServer) UseFunc(middlewareFuncs ...func(http.Handler) http.Handler) {
	for _, middlewareFunc := range middlewareFuncs {
		s.middlewares = append(s.middlewares, MiddlewareFunc(middlewareFunc))
	}
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
		log.ErrorContext(ctx, "HTTP shutdown error", "error", err)
		return err
	}
	log.InfoContext(ctx, "graceful shutdown completed.")

	return nil
}

func (s *HttpServer) Healthcheck(ctx context.Context) any {
	return map[string]any{
		"port": s.port,
	}
}
