package fileserver

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"time"
)

type FileServer struct {
	mux  *http.ServeMux
	port string
}

func New(fs fs.FS, basePath, port string) *FileServer {
	mux := http.NewServeMux()
	mux.Handle(basePath, http.StripPrefix(basePath, http.FileServer(http.FS(fs))))

	return &FileServer{mux: mux, port: port}
}

func (s *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *FileServer) Run(ctx context.Context) error {
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
