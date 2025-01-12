package client_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/MenD32/Tempest/pkg/client"
	"github.com/MenD32/Tempest/pkg/response"
	"github.com/MenD32/Tempest/pkg/test"
)

func TestNewDefaultClient(t *testing.T) {
	respFactory := func(resp *http.Response, sent time.Time) (response.Response, error) {
		return test.TestResponse{}, nil
	}

	c := client.NewDefaultClient(respFactory)
	if c == nil {
		t.Fatalf("Expected non-nil client")
	}
}
