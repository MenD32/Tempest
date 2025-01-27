package response

import (
	"net/http"
	"time"
)

type Response interface {
	Metrics() (*Metrics)
	Verify() error
}

type ResponseBuilder func(*http.Response, time.Time) (Response, error)

type ErrorResponse struct {
	Sent time.Time
	Err  error
}

func (er ErrorResponse) Metrics() (*Metrics) {
	return &Metrics{
		Sent: er.Sent,
		Error: er.Err,
	}
}

func (er ErrorResponse) Verify() error {
	return nil
}