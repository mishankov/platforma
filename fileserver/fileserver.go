package fileserver

import (
	"context"
	"io/fs"
	"net/http"
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
	if err := http.ListenAndServe(":"+s.port, s.mux); err != nil {
		return err
	}

	return nil
}
