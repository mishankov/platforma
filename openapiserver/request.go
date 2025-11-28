package openapiserver

import (
	"encoding/json"
	"net/http"
)

type Request[Query, Headers, Body any] struct {
	httpRequest *http.Request
	Query       Query
	Headers     Headers
	bodyDecoded bool
	body        *Body
}

func (r *Request[Query, Headers, Body]) Body() (*Body, error) {
	if r.bodyDecoded {
		return r.body, nil
	}

	if err := json.NewDecoder(r.httpRequest.Body).Decode(r.body); err != nil {
		return nil, err
	}
	r.bodyDecoded = true

	return r.body, nil
}
