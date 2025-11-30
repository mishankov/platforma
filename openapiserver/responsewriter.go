package openapiserver

// ResponseWriter provides a typed interface for writing HTTP responses.
type ResponseWriter[Headers, Body any] struct {
	StatusCode int
	Headers    Headers
	bodySet    bool
	body       Body
}

// SetBody sets the response body.
func (w *ResponseWriter[Headers, Body]) SetBody(b Body) {
	w.body = b
	w.bodySet = true
}
