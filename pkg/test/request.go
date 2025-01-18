package test

import (
	"net/http"
	"net/url"
	"time"
)

type TestRequest struct {
	delay time.Duration
	req   *http.Request
}

func NewTestRequest(serverurl string, delay time.Duration) TestRequest {
	u, _ := url.Parse(serverurl)
	return TestRequest{
		delay: delay,
		req: &http.Request{
			Method: http.MethodGet,
			URL:    u,
		},
	}
}

func (r TestRequest) Delay() time.Duration {
	return r.delay
}

func (r TestRequest) HTTPRequest() *http.Request {
	return r.req
}
