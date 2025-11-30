package openapiserver

import (
	"net/http"
)

// Request represents an HTTP request with typed data.
type Request[T any] struct {
	HTTPRequest *http.Request
	Data        *T
}
