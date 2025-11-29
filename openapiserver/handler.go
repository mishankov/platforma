package openapiserver

import (
	"encoding/json"
	"net/http"

	"github.com/oaswrap/spec/option"
	"github.com/platforma-dev/platforma/log"
)

type Handler[RequestType, ResponseHeaders, ResponseBody any] func(w *ResponseWriter[ResponseHeaders, ResponseBody], r Request[RequestType])

func Get[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	Handle(group, resps, http.MethodGet, pattern, handler)
}

// func Head[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody]) {
// 	Handle(group, resps, http.MethodHead, pattern, handler)
// }

// func Post[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody]) {
// 	Handle(group, resps, http.MethodPost, pattern, handler)
// }

// func Put[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody]) {
// 	Handle(group, resps, http.MethodPut, pattern, handler)
// }

// func Patch[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody]) {
// 	Handle(group, resps, http.MethodPatch, pattern, handler)
// }

// func Delete[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody]) {
// 	Handle(group, resps, http.MethodDelete, pattern, handler)
// }

// func Connect[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody]) {
// 	Handle(group, resps, http.MethodConnect, pattern, handler)
// }

// func Options[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[Path, Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody]) {
// 	Handle(group, resps, http.MethodOptions, pattern, handler)
// }

// func Trace[Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody]) {
// 	Handle(group, resps, http.MethodTrace, pattern, handler)
// }

func Handle[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, method string, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	// Prepare open api spec
	opts := []option.OperationOption{
		option.Request(new(RequestType)),
	}
	for statusCode, respModel := range resps {
		opts = append(opts, option.Response(statusCode, respModel))
	}

	group.spec.Add(method, pattern, opts...)

	// Add handler logic to mux
	group.handlerGroup.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// Convert http request to user request
		request := &Request[RequestType]{
			HttpRequest: r,
			Data:        new(RequestType),
		}
		// Query
		mapToStruct(r.URL.Query(), "query", request.Data)

		// Headers
		mapToStruct(r.Header, "header", request.Data)

		// Body
		if err := json.NewDecoder(r.Body).Decode(request.Data); err != nil {
			log.Error("failed to decode body", "error", err)
		}

		// Call user handle
		writer := ResponseWriter[ResponseHeaders, ResponseBody]{}
		handler(&writer, *request)

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
