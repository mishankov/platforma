package httpserver

import (
	"net/http"
)

// HandlerGroup represents a group of HTTP handlers that share common middlewares.
type HandlerGroup struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

// NewHandlerGroup creates a new HandlerGroup with an initialized http.ServeMux.
func NewHandlerGroup() *HandlerGroup {
	return &HandlerGroup{mux: http.NewServeMux()}
}

// Use adds a middleware to the HandlerGroup's middleware chain.
func (hg *HandlerGroup) Use(middlewares ...Middleware) {
	hg.middlewares = append(hg.middlewares, middlewares...)
}

// UseFunc adds a function as a middleware to the HandlerGroup's middleware chain.
func (hg *HandlerGroup) UseFunc(middlewareFuncs ...func(http.Handler) http.Handler) {
	for _, middlewareFunc := range middlewareFuncs {
		hg.middlewares = append(hg.middlewares, MiddlewareFunc(middlewareFunc))
	}
}

// Handle registers an http.Handler for the given pattern
func (hg *HandlerGroup) Handle(pattern string, handler http.Handler) {
	hg.mux.Handle(pattern, handler)
}

// HandleFunc registers an http.HandlerFunc for the given pattern
func (hg *HandlerGroup) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	hg.mux.Handle(pattern, http.HandlerFunc(handler))
}

// HandleGroup applies `http.StripPrefix` to http.Handler and registers it for the given pattern
func (hg *HandlerGroup) HandleGroup(pattern string, handler http.Handler) {
	hg.mux.Handle(pattern+"/", http.StripPrefix(pattern, handler))
}

// ServeHTTP implements the http.Handler interface, allowing HandlerGroup to
// be used as an HTTP handler itself.
func (hg *HandlerGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrappedMux := wrapHandlerInMiddleware(hg.mux, hg.middlewares)
	wrappedMux.ServeHTTP(w, r)
}
