package openapiserver

import (
	"encoding/json"
	"net/http"

	"github.com/oaswrap/spec/option"
	"github.com/platforma-dev/platforma/log"
)

type Handler[Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any] func(w *ResponseWriter[ResponseHeaders, ResponseBody], r *Request[Query, RequestHeaders, RequestBody])

func Get[Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody]) {
	// Prepare open api spec
	opts := []option.OperationOption{
		option.Request(new(Query)),
		option.Request(new(RequestHeaders)),
		option.Request(new(RequestBody)),
	}
	for statusCode, respModel := range resps {
		opts = append(opts, option.Response(statusCode, respModel))
	}

	group.spec.Get(pattern, opts...)

	// Add handler logic to mux
	group.handlerGroup.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// Convert http request to user request
		request := &Request[Query, RequestHeaders, RequestBody]{
			httpRequest: r,
		}
		// Query
		var query Query
		mapToStruct(r.URL.Query(), "query", &query)
		request.Query = query

		// Headers
		var requestHeaders RequestHeaders
		mapToStruct(r.Header, "header", &requestHeaders)
		request.Headers = requestHeaders

		// Call user handle
		writer := ResponseWriter[ResponseHeaders, ResponseBody]{}
		handler(&writer, request)

		// Headers
		headers := mapFromStruct[map[string][]string](writer.Headers, "header")
		for name, values := range headers {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// Status code
		if writer.StatusCode == 0 {
			writer.StatusCode = http.StatusOK
		}
		w.WriteHeader(writer.StatusCode)

		// Body
		if writer.bodySet {
			if err := json.NewEncoder(w).Encode(writer.body); err != nil {
				log.Error("failed to encode body", "error", err)
			}
		}
	})
}
