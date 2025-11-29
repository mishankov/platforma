package openapiserver

import (
	"net/http"

	"github.com/oaswrap/spec"
	"github.com/platforma-dev/platforma/httpserver"
	"github.com/platforma-dev/platforma/log"
)

type Router struct {
	handlerGroup *httpserver.HandlerGroup
	spec         spec.Generator
	specPath     string // OpenAPI specifications path
	docPath      string // OpenAPI interactive documentation path
}

func NewRouter(specPath, docPath string) *Router {
	return &Router{
		handlerGroup: httpserver.NewHandlerGroup(),
		spec:         spec.NewRouter(),
		specPath:     specPath,
		docPath:      docPath,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handlerGroup.ServeHTTP(w, req)
}

func (r *Router) OpenAPI() {
	if err := r.spec.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Error(err.Error())
	}
}
