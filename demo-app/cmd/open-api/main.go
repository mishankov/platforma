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

type myRequest = httpserverv2.Request[myQuery, myReqHeaders, any]
type myRespWriter = *httpserverv2.ResponseWriter[myRespHeaders, myRespBody]

func main() {
	router := httpserverv2.NewRouter("/docs/openapi.yml", "/docs")

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

	httpserverv2.Get(router, resps, "/hello", func(w myRespWriter, r *myRequest) {
		w.Headers.XMen = r.Query.Name
		w.Headers.ContentType = "application/json"

		if r.Query.Name[0] == "xavier" {
			w.StatusCode = http.StatusBadRequest
			w.SetBody(myRespBody{errorRespBody: errorRespBody{ErrorMessage: "superhero banned"}})

			return
		}

		w.SetBody(myRespBody{successRespBody: successRespBody{Data: r.Query.Name}})
	})

	http.ListenAndServe(":8080", router)
}
