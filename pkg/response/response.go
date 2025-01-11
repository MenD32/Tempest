package response

import (
	"net/http"
	"time"
)

type Response interface {
	Metrics() Metrics
}

type ResponseBuilder func(*http.Response, time.Time) (Response, error)
