package openapiserver

import (
	"github.com/oaswrap/spec"
	"github.com/platforma-dev/platforma/httpserver"
)

// Group represents a group of routes with a common path prefix.
type Group struct {
	spec         spec.Router
	handlerGroup *httpserver.HandlerGroup
}

// NewGroup creates a new route group with the specified pattern.
func NewGroup(router *Router, pattern string) *Group {
	hg := httpserver.NewHandlerGroup()
	router.handlerGroup.Handle(pattern, hg)
	group := router.spec.Group(pattern)

	return &Group{
		spec:         group,
		handlerGroup: hg,
	}
}
