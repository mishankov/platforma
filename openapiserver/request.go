package openapiserver

import (
	"net/http"
)

type Request[T any] struct {
	HttpRequest *http.Request
	Data        *T
}
