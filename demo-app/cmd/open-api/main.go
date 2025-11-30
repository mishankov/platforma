package main

import (
	"context"
	"net/http"

	"github.com/platforma-dev/platforma/openapiserver"
)

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

	helloGroup := openapiserver.NewGroup(router, "/")

	openapiserver.Get(
		helloGroup, resps, "/hello",
		func(_ context.Context, w myRespWriter, r openapiserver.Request[struct {
			Name      []string `query:"name"`
			UserAgent []string `header:"User-Agent"`
		}]) {
			w.Headers.XMen = r.Data.Name
			w.Headers.ContentType = "application/json"

			if r.Data.Name[0] == "xavier" {
				w.StatusCode = http.StatusBadRequest
				w.SetBody(myRespBody{errorRespBody: errorRespBody{ErrorMessage: "superhero banned"}})

				return
			}

			w.SetBody(myRespBody{successRespBody: successRespBody{Data: r.Data.Name}})
		})

	openapiserver.Put(
		helloGroup, resps, "/hello/{id}",
		func(_ context.Context, w myRespWriter, r openapiserver.Request[struct {
			Id        string   `path:"id"`
			Name      []string `query:"name"`
			UserAgent []string `header:"User-Agent"`
		}]) {
			w.Headers.XMen = r.Data.Name
			w.Headers.ContentType = "application/json"

			if r.Data.Name[0] == "xavier" {
				w.StatusCode = http.StatusBadRequest
				w.SetBody(myRespBody{errorRespBody: errorRespBody{ErrorMessage: "superhero banned"}})

				return
			}

			w.SetBody(myRespBody{successRespBody: successRespBody{Data: r.Data.Name}})
		})

	http.ListenAndServe(":8080", router)
}
