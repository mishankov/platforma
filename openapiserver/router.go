package openapiserver

import (
	"net/http"

	"github.com/oaswrap/spec"
	specui "github.com/oaswrap/spec-ui"
	"github.com/platforma-dev/platforma/httpserver"
)

type Router struct {
	handlerGroup *httpserver.HandlerGroup
	spec         spec.Generator
	specPath     string // OpenAPI specifications path
	docPath      string // OpenAPI interactive documentation path
}

func NewRouter(specPath, docPath string) *Router {
	hg := httpserver.NewHandlerGroup()
	sp := spec.NewRouter()

	if specPath != "" {
		openapiHandler := specui.NewHandler(
			specui.WithDocsPath(docPath),
			specui.WithSpecPath(specPath),
			specui.WithSpecGenerator(sp),
			specui.WithScalar(),
		)

		hg.Handle(openapiHandler.SpecPath(), openapiHandler.Spec())
		hg.Handle(openapiHandler.DocsPath(), openapiHandler.Docs())
	}

	return &Router{
		handlerGroup: hg,
		spec:         sp,
		specPath:     specPath,
		docPath:      docPath,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handlerGroup.ServeHTTP(w, req)
}
