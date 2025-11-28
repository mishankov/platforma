package httpserverv2

import (
	"encoding/json"
	"net/http"

	"github.com/platforma-dev/platforma/log"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

type Router struct {
	mux      http.ServeMux
	specPath string // OpenAPI specifications path
	docPath  string // OpenAPI interactive documentation path
}

func NewRouter(specPath, docPath string) *Router {
	return &Router{
		mux:      *http.NewServeMux(),
		specPath: specPath,
		docPath:  docPath,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

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

type Handler[Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any] func(w *ResponseWriter[ResponseHeaders, ResponseBody], r *Request[Query, RequestHeaders, RequestBody])

func Get[Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody any](router *Router, pattern string, handler Handler[Query, RequestHeaders, RequestBody, ResponseHeaders, ResponseBody]) {
	// Prepare open api spec
	r := spec.NewRouter()

	// Add routes
	v1 := r.Group("")

	v1.Get(pattern,
		option.Summary("User login"),
		option.Request(new(Query)),
		option.Request(new(RequestHeaders)),
		option.Request(new(RequestBody)),
		option.Response(200, new(ResponseHeaders)),
		option.Response(201, new(ResponseBody)),
	)

	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Error(err.Error())
	}

	router.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
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

func main() {
	type myQuery struct {
		Name string `query:"name"`
	}

	type myRespHeaders struct {
		XMen string `header:"X-Men"`
	}

	type myRequest = Request[myQuery, any, any]
	type myRespWriter = *ResponseWriter[myRespHeaders, any]

	router := &Router{}
	Get(router, "/hey", func(w myRespWriter, r *myRequest) {
		w.Headers.XMen = r.Query.Name
	})
}
