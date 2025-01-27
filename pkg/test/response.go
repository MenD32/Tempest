package test

import (
	"time"

	"github.com/MenD32/Tempest/pkg/response"
)

type TestResponse struct {
	Sent time.Time
}

func (tr TestResponse) Metrics() (*response.Metrics) {
	return &response.Metrics{
		Sent: time.Now(),
		Body: []byte(""),
		Metrics: map[string]interface{}{
			"test": "test",
		},
	}
}

func (tr TestResponse) Verify() error {
	return nil
}
