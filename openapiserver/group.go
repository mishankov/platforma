package openapiserver

import (
	"github.com/oaswrap/spec"
	"github.com/platforma-dev/platforma/httpserver"
)

type Group struct {
	spec         spec.Router
	handlerGroup *httpserver.HandlerGroup
}

func NewGroup(router *Router, pattern string) *Group {
	hg := httpserver.NewHandlerGroup()
	router.handlerGroup.HandleGroup(pattern, hg)
	group := router.spec.Group(pattern)

	return &Group{
		spec:         group,
		handlerGroup: hg,
	}
}
