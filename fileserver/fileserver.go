// Package fileserver provides an HTTP file server for serving static files.
package fileserver

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"time"
)

// FileServer represents an HTTP file server for serving static files.
type FileServer struct {
	mux  *http.ServeMux
	port string
}

// New creates a new FileServer instance with the given file system, base path, and port.
func New(fs fs.FS, basePath, port string) *FileServer {
	mux := http.NewServeMux()
	mux.Handle(basePath, http.StripPrefix(basePath, http.FileServer(http.FS(fs))))

	return &FileServer{mux: mux, port: port}
}

func (s *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Run starts the file server and listens for incoming requests.
func (s *FileServer) Run(_ context.Context) error {
	server := &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("file server failed: %w", err)
	}

	return nil
}
