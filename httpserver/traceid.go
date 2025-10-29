package httpserver

import (
	"context"
	"net/http"

	"github.com/mishankov/platforma/log"

	"github.com/google/uuid"
)

// TraceIdMiddleware is a middleware that adds a trace ID to the request context and response headers.
type TraceIdMiddleware struct {
	contextKey any
	header     string
}

// NewTraceIdMiddleware returns a new TraceId middleware.
// If key is nil, log.TraceIdKey is used.
// If header is empty, "Platforma-Trace-Id" is used.
func NewTraceIdMiddleware(contextKey any, header string) *TraceIdMiddleware {
	if contextKey == nil {
		contextKey = log.TraceIdKey
	}

	if header == "" {
		header = "Platforma-Trace-Id"
	}

	return &TraceIdMiddleware{contextKey: contextKey, header: header}
}

func (m *TraceIdMiddleware) Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := uuid.NewString()
		ctx := context.WithValue(r.Context(), m.contextKey, traceId)
		r = r.WithContext(ctx)

		w.Header().Set(m.header, traceId)

		h.ServeHTTP(w, r)
	})
}
