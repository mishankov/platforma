package httpserver_test

import (
	"context"
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

	t.Run("healthcheck", func(t *testing.T) {
		t.Parallel()

		server := httpserver.New("8080", 0)
		hcData, ok := server.Healthcheck(context.TODO()).(map[string]any)
		if !ok {
			t.Fatal("failed type assert health data")
		}

		port := hcData["port"]
		if port != "8080" {
			t.Errorf("expected port to be 8080, got %s", port)
		}
	})
}

type handler struct {
	serveHttp func(http.ResponseWriter, *http.Request)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.serveHttp != nil {
		h.serveHttp(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
