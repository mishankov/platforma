package main

import (
	"net/http"

	httpserverv2 "github.com/platforma-dev/platforma/httpserver-v2"
)

func main() {
	router := httpserverv2.NewRouter("/openapi.yml", "/docs")

	type myQuery struct {
		Name []string `query:"name"`
	}

	type myReqHeaders struct {
		UserAgent []string `header:"User-Agent"`
	}

	type myRespHeaders struct {
		XMen        string `header:"X-Men"`
		ContentType string `header:"Content-Type"`
	}

	type myRespBody struct {
		Data string `json:"data"`
	}

	type myRequest = httpserverv2.Request[myQuery, myReqHeaders, any]
	type myRespWriter = *httpserverv2.ResponseWriter[myRespHeaders, myRespBody]

	httpserverv2.Get(router, "/hello", func(w myRespWriter, r *myRequest) {
		w.Headers.XMen = r.Query.Name[0]
		w.Headers.ContentType = "application/json"

		w.SetBody(myRespBody{Data: "Hi, " + r.Query.Name[1]})
	})

	http.ListenAndServe(":8080", router)
}
