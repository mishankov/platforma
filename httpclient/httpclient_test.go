package httpclient_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mishankov/platforma/httpclient"
)

const timeout = 10 * time.Second

func TestNew(t *testing.T) {
	client := httpclient.New(timeout)
	if client == nil {
		t.Error("New() should return a non-nil client")
	}
}

func TestDo_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	client := httpclient.New(timeout)
	req, err := http.NewRequestWithContext(context.Background(), "GET", server.URL, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestDo_Error(t *testing.T) {
	// Create a request to a non-existent server
	client := httpclient.New(timeout)
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:9999/nonexistent", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err == nil {
		t.Error("expected error for non-existent server, got nil")
	}
	defer resp.Body.Close()
}
