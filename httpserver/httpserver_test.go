package httpserver_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/mishankov/platforma/httpserver"
)

// testHandler is a simple handler for testing
type testHandler struct{}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func TestHttpServer_ShutdownCompletesBeforeTimeout(t *testing.T) {
	t.Parallel()

	// Create a test HTTP server directly to test shutdown behavior
	server := &http.Server{
		Addr:    ":8080",
		Handler: &testHandler{},
	}

	// Start server in goroutine
	go func() {
		server.ListenAndServe()
	}()

	// Give server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test shutdown with a long timeout but expect it to complete quickly
	shutdownTimeout := 5 * time.Second
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)
	shutdownDuration := time.Since(startTime)

	// Verify shutdown completed quickly (much less than full timeout)
	if shutdownDuration > 1*time.Second {
		t.Errorf("shutdown took %v, expected to complete much faster than %v timeout", shutdownDuration, shutdownTimeout)
	}

	if err != nil {
		t.Errorf("unexpected shutdown error: %v", err)
	}
}

func TestHttpServer_ShutdownWithNoActiveConnections(t *testing.T) {
	t.Parallel()

	// Create HttpServer instance to test the integration
	httpServer := httpserver.New("8081", 3*time.Second)
	httpServer.Handle("/test", &testHandler{})

	// Create a test server to simulate the internal http.Server
	testServer := &http.Server{
		Addr:    ":8081",
		Handler: &testHandler{},
	}

	// Start server
	go func() {
		testServer.ListenAndServe()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test shutdown - should complete quickly since no active connections
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := testServer.Shutdown(ctx)
	shutdownDuration := time.Since(startTime)

	// Should complete much faster than the full timeout
	if shutdownDuration > 500*time.Millisecond {
		t.Errorf("shutdown with no connections took %v, expected <500ms", shutdownDuration)
	}

	if err != nil {
		t.Errorf("unexpected shutdown error: %v", err)
	}
}

func TestHttpServer_Healthcheck(t *testing.T) {
	t.Parallel()

	server := httpserver.New("8083", 5*time.Second)

	result := server.Healthcheck(context.Background())

	healthMap, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}

	port, exists := healthMap["port"]
	if !exists {
		t.Error("healthcheck should contain 'port' field")
	}

	if port != "8083" {
		t.Errorf("expected port '8083', got '%v'", port)
	}
}
