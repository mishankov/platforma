package main

import (
	"net/http"

	"github.com/platforma-dev/platforma/openapiserver"
)

type myReq struct {
	Name      []string `query:"name"`
	UserAgent []string `header:"User-Agent"`
}

type myRespHeaders struct {
	XMen        []string `header:"X-Men"`
	ContentType string   `header:"Content-Type"`
}

type errorRespBody struct {
	ErrorMessage string `json:"errorMessage,omitempty"`
}

type successRespBody struct {
	Data []string `json:"data,omitempty"`
}

type myRespBody struct {
	errorRespBody
	successRespBody
}

type myRequest = openapiserver.Request[myReq]
type myRespWriter = *openapiserver.ResponseWriter[myRespHeaders, myRespBody]

func main() {
	router := openapiserver.NewRouter("/docs/openapi.yml", "/docs")

	resps := map[int]any{
		http.StatusOK: struct {
			myRespHeaders
			successRespBody
		}{},
		http.StatusCreated: struct {
			myRespHeaders
			successRespBody
		}{},
		http.StatusBadRequest: struct {
			myRespHeaders
			errorRespBody
		}{},
	}

	helloGroup := openapiserver.NewGroup(router, "")

	openapiserver.Get(helloGroup, resps, "/hello", func(w myRespWriter, r myRequest) {
		w.Headers.XMen = r.Data.Name
		w.Headers.ContentType = "application/json"

		if r.Data.Name[0] == "xavier" {
			w.StatusCode = http.StatusBadRequest
			w.SetBody(myRespBody{errorRespBody: errorRespBody{ErrorMessage: "superhero banned"}})

			return
		}

		w.SetBody(myRespBody{successRespBody: successRespBody{Data: r.Data.Name}})
	})

	// openapiserver.Put(
	// 	helloGroup, resps, "/hello/{id}",
	// 	func(w myRespWriter, r *openapiserver.Request[myQuery, myReqHeaders, any]) {
	// 		w.Headers.XMen = r.Query.Name
	// 		w.Headers.ContentType = "application/json"

	// 		if r.Query.Name[0] == "xavier" {
	// 			w.StatusCode = http.StatusBadRequest
	// 			w.SetBody(myRespBody{errorRespBody: errorRespBody{ErrorMessage: "superhero banned"}})

	// 			return
	// 		}

	// 		w.SetBody(myRespBody{successRespBody: successRespBody{Data: r.Query.Name}})
	// 	})

	router.OpenAPI()

	http.ListenAndServe(":8080", router)
}
