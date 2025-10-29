package httpserver_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mishankov/platforma/httpserver"
)

func TestHttpServer(t *testing.T) {
	t.Parallel()

	t.Run("single http.HandlerFunc endpoint", func(t *testing.T) {
		t.Parallel()

		server := httpserver.New("", 0)

		server.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

		r := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		if string(body) != "pong" {
			t.Errorf("expected body to be 'pong', got %s", string(body))
		}
	})

	t.Run("single http.Handler endpoint", func(t *testing.T) {
		t.Parallel()

		pingHandler := &handler{
			serveHttp: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("pong"))
			},
		}

		server := httpserver.New("", 0)

		server.Handle("/ping", pingHandler)

		r := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, r)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		if string(body) != "pong" {
			t.Errorf("expected body to be 'pong', got %s", string(body))
		}
	})
}

type handler struct {
	serveHttp func(http.ResponseWriter, *http.Request)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveHttp(w, r)
}
