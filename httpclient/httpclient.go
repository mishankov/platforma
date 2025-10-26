package httpclient

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mishankov/platforma/log"
)

type Client struct {
	client *http.Client
}

func New(timeout time.Duration) *Client {
	c := &Client{&http.Client{
		Timeout: timeout,
	}}

	return c
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	log.DebugContext(req.Context(), "request started", "url", req.URL, "headers", maskedHeaders(req.Header))
	resp, err := c.client.Do(req)
	if err != nil {
		log.DebugContext(req.Context(), "request failed", "error", err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	log.DebugContext(req.Context(), "request made", "status", resp.Status, "headers", maskedHeaders(resp.Header))

	return resp, nil
}

func maskedHeaders(headers http.Header) http.Header {
	newHeaders := http.Header{}
	for hn, hvs := range headers {
		if strings.ToLower(hn) == "authorization" {
			for _, hv := range hvs {
				if strings.HasPrefix(strings.ToLower(hv), "basic ") || strings.HasPrefix(strings.ToLower(hv), "bearer ") {
					newHeaders.Add(hn, strings.Split(hv, " ")[0]+" ***")
					continue
				}
				newHeaders.Add(hn, hv)
			}
		} else {
			for _, hv := range hvs {
				newHeaders.Add(hn, hv)
			}
		}
	}

	return newHeaders
}
