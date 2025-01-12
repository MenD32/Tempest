package test

import (
	"time"

	"github.com/MenD32/Tempest/pkg/response"
)

type TestResponse struct{}

func (tr TestResponse) Metrics() (*response.Metrics, error) {
	return &response.Metrics{
		Sent: time.Now(),
		Body: []byte(""),
		Metrics: map[string]interface{}{
			"test": "test",
		},
	}, nil
}

func (tr TestResponse) Verify() error {
	return nil
}
