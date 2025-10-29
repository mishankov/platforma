package httpserver

import (
	"net/http"
	"slices"
)

// Middleware defines an interface for HTTP middleware.
type Middleware interface {
	// Wrap wraps an http.Handler with middleware logic.
	Wrap(http.Handler) http.Handler
}

// MiddlewareFunc is a convenience type that implements the Middleware interface.
type MiddlewareFunc func(http.Handler) http.Handler

// Wrap implements the Middleware interface for MiddlewareFunc.
func (f MiddlewareFunc) Wrap(h http.Handler) http.Handler {
	return f(h)
}

// wrapHandlerInMiddleware wraps an http.Handler with a chain of middlewares.
// The middlewares are applied in reverse order of the provided slice,
// meaning the last middleware in the slice will be the most inner.
func wrapHandlerInMiddleware(handler http.Handler, middlewares []Middleware) http.Handler {
	finalHandler := handler
	for _, middleware := range slices.Backward(middlewares) {
		finalHandler = middleware.Wrap(finalHandler)
	}

	return finalHandler
}
