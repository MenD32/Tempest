package empty

import (
	"net/http"
	"time"

	"github.com/MenD32/Tempest/pkg/response"
)

type EmptyResponse struct {
	Sent time.Time
}

func EmptyResponseBuilder(resp *http.Response, sent time.Time) (response.Response, error) {
	return EmptyResponse{Sent: sent}, nil
}

func (r EmptyResponse) Metrics() (*response.Metrics) {
	return &response.Metrics{
		Sent: r.Sent,
	}
}

func (r EmptyResponse) Verify() error {
	return nil
}
