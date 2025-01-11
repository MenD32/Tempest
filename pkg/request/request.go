package request

import (
	"net/http"
	"time"
)

type Request interface {
	Delay() time.Duration
	HTTPRequest() *http.Request
}

type request struct {
	http.Request

	delay time.Duration
}

func (r *request) Delay() time.Duration {
	return r.delay
}

func (r *request) HTTPRequest() *http.Request {
	return &r.Request
}

func NewRequest(delay time.Duration, req http.Request) Request {
	return &request{
		delay:   delay,
		Request: req,
	}
}

type RequestFactory func(string, string) ([]Request, error)
