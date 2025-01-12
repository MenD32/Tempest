package client_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MenD32/Tempest/pkg/client"
	"github.com/MenD32/Tempest/pkg/request"
	"github.com/MenD32/Tempest/pkg/response"
	"github.com/MenD32/Tempest/pkg/test"
)

const (
	MAX_DRIFT               = 50 * time.Millisecond
	MAX_CONCURRENT_REQUESTS = 1000
	REQUEST_COUNT           = 100
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

func TestRunWithServer(t *testing.T) {
	// Start a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world")
	}))
	defer ts.Close()

	respFactory := func(resp *http.Response, sent time.Time) (response.Response, error) {
		return test.TestResponse{}, nil
	}

	c := client.NewDefaultClient(respFactory)

	requests := []request.Request{}
	for i := 0; i < REQUEST_COUNT; i++ {
		req := test.NewTestRequest(ts.URL, 0)
		requests = append(requests, req)
	}

	responses := client.Run(c, requests)

	if len(responses) != len(requests) {
		t.Fatalf("Expected %d responses, got %d", len(requests), len(responses))
	}
}

func TestRunWithServerAndHighConcurrentRequests(t *testing.T) {
	// Start a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		fmt.Fprintf(w, "Hello, world")
	}))
	defer ts.Close()

	respFactory := func(resp *http.Response, sent time.Time) (response.Response, error) {
		return test.TestResponse{Sent: sent}, nil
	}

	c := client.NewDefaultClient(respFactory)

	requests := []request.Request{}
	start := time.Now()
	for i := 0; i < MAX_CONCURRENT_REQUESTS; i++ {
		req := test.NewTestRequest(ts.URL, time.Millisecond*time.Duration(i))
		requests = append(requests, req)
	}

	responses := client.Run(c, requests)

	if len(responses) != len(requests) {
		t.Fatalf("Expected %d responses, got %d", len(requests), len(responses))
	}

	for i, res := range responses {
		drift := res.(test.TestResponse).Sent.Sub(start) - time.Millisecond*time.Duration(i)
		if drift > MAX_DRIFT {
			t.Fatalf("Expected drift to be less than %v, got %v", MAX_DRIFT, drift)
		}
	}
}
