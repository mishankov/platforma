package main

import (
	"net/http"

	httpserverv2 "github.com/platforma-dev/platforma/httpserver-v2"
)

type myQuery struct {
	Name []string `query:"name"`
}

type myReqHeaders struct {
	UserAgent []string `header:"User-Agent"`
}

type myRespHeaders struct {
	XMen        []string `header:"X-Men"`
	ContentType string   `header:"Content-Type"`
}

type myRespBody struct {
	Data []string `json:"data"`
}

type errorResp struct {
	ErrorMessage string `json:"errorMessage"`
}

type myRequest = httpserverv2.Request[myQuery, myReqHeaders, any]
type myRespWriter = *httpserverv2.ResponseWriter[myRespHeaders, myRespBody]

func main() {
	router := httpserverv2.NewRouter("/docs/openapi.yml", "/docs")

	resps := map[int]any{
		http.StatusOK: struct {
			myRespHeaders
			myRespBody
		}{},
		http.StatusCreated: struct {
			myRespHeaders
			myRespBody
		}{},
		http.StatusBadRequest: struct {
			myRespHeaders
			errorResp
		}{},
	}

	httpserverv2.Get(router, resps, "/hello", func(w myRespWriter, r *myRequest) {
		w.Headers.XMen = r.Query.Name
		w.Headers.ContentType = "application/json"

		if r.Query.Name[0] == "xavier" {
			w.StatusCode = http.StatusBadRequest
			w.SetBody(errorResp{ErrorMessage: "banned superhero"})

			return
		}

		w.SetBody(myRespBody{Data: r.Query.Name})
	})

	http.ListenAndServe(":8080", router)
}
