package openapiserver

import (
	"net/http"

	"github.com/oaswrap/spec"
)

type Group struct {
	spec spec.Router
	mux  *http.ServeMux
}

func NewGroup(router *Router, pattern string) *Group {
	group := router.spec.Group(pattern)
	return &Group{
		spec: group,
		mux:  http.NewServeMux(),
	}
}
