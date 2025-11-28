package openapiserver

import (
	"net/http"

	"github.com/oaswrap/spec"
	"github.com/platforma-dev/platforma/log"
)

type Router struct {
	mux      http.ServeMux
	spec     spec.Generator
	specPath string // OpenAPI specifications path
	docPath  string // OpenAPI interactive documentation path
}

func NewRouter(specPath, docPath string) *Router {
	return &Router{
		mux:      *http.NewServeMux(),
		spec:     spec.NewRouter(),
		specPath: specPath,
		docPath:  docPath,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) OpenAPI() {
	if err := r.spec.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Error(err.Error())
	}
}
