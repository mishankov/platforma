package openapiserver

type ResponseWriter[Headers, Body any] struct {
	StatusCode int
	Headers    Headers
	bodySet    bool
	body       Body
}

func (w *ResponseWriter[Headers, Body]) SetBody(b Body) {
	w.body = b
	w.bodySet = true
}
