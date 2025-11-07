package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mishankov/platforma/httpserver"
	"github.com/mishankov/platforma/log"
)

func TestTraceIdMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("default params", func(t *testing.T) {
		t.Parallel()

		m := httpserver.NewTraceIDMiddleware(nil, "")
		wrappedHandler := m.Wrap(&handler{serveHTTP: func(w http.ResponseWriter, r *http.Request) {
			i, ok := r.Context().Value(log.TraceIDKey).(string)
			if ok {
				w.Header().Add("TraceIdFromContext", i)
			}
		}})

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, r)
		resp := w.Result()

		if len(resp.Header.Get("Platforma-Trace-Id")) == 0 {
			t.Fatalf("default trace id header expected, got: %s", resp.Header)
		}

		if len(resp.Header.Get("TraceIdFromContext")) == 0 {
			t.Fatalf("trsce id from cotext expected, got: %s", resp.Header)
		}

	})
}
