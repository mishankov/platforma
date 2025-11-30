package openapiserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oaswrap/spec/option"
	"github.com/platforma-dev/platforma/log"
)

// Handler defines a function type for handling HTTP requests with typed request and response parameters.
type Handler[RequestType, ResponseHeaders, ResponseBody any] func(ctx context.Context, w *ResponseWriter[ResponseHeaders, ResponseBody], r Request[RequestType])

// Get registers a GET route handler.
func Get[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	Handle(group, resps, http.MethodGet, pattern, handler)
}

// Head registers a HEAD route handler.
func Head[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	Handle(group, resps, http.MethodHead, pattern, handler)
}

// Post registers a POST route handler.
func Post[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	Handle(group, resps, http.MethodPost, pattern, handler)
}

// Put registers a PUT route handler.
func Put[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	Handle(group, resps, http.MethodPut, pattern, handler)
}

// Patch registers a PATCH route handler.
func Patch[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	Handle(group, resps, http.MethodPatch, pattern, handler)
}

// Delete registers a DELETE route handler.
func Delete[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	Handle(group, resps, http.MethodDelete, pattern, handler)
}

// Connect registers a CONNECT route handler.
func Connect[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	Handle(group, resps, http.MethodConnect, pattern, handler)
}

// Options registers an OPTIONS route handler.
func Options[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	Handle(group, resps, http.MethodOptions, pattern, handler)
}

// Trace registers a TRACE route handler.
func Trace[RequestType, ResponseHeaders, ResponseBody any](group *Group, resps map[int]any, pattern string, handler Handler[RequestType, ResponseHeaders, ResponseBody]) {
	Handle(group, resps, http.MethodTrace, pattern, handler)
}

// Handle registers a route handler for a specific HTTP method.
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
	group.handlerGroup.HandleFunc(method+" "+pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Convert http request to user request
		request := &Request[RequestType]{
			HTTPRequest: r,
			Data:        new(RequestType),
		}

		// Path
		if err := pathToStruct(r, request.Data); err != nil {
			log.ErrorContext(ctx, "failed to parse path parameters", "error", err)
		}

		// Query
		if err := mapToStruct(r.URL.Query(), "query", request.Data); err != nil {
			log.ErrorContext(ctx, "failed to parse query parameters", "error", err)
		}

		// Headers
		if err := mapToStruct(r.Header, "header", request.Data); err != nil {
			log.ErrorContext(ctx, "failed to parse headers", "error", err)
		}

		// Body
		if err := json.NewDecoder(r.Body).Decode(request.Data); err != nil {
			log.ErrorContext(ctx, "failed to decode body", "error", err)
		}

		// Call user handle
		writer := ResponseWriter[ResponseHeaders, ResponseBody]{}
		handler(ctx, &writer, *request)

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
				log.ErrorContext(ctx, "failed to encode body", "error", err)
			}
		}
	})
}
