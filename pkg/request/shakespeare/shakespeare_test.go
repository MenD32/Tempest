package shakespeare_test

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/MenD32/Tempest/pkg/request"
	"github.com/MenD32/Tempest/pkg/request/shakespeare"
)

func requestInPEACE(method string, host string, body []byte) *http.Request {

	httpreq, err := http.NewRequest(
		method,
		host,
		bytes.NewReader([]byte(body)),
	)
	if err != nil {
		panic(err)
	}

	return httpreq
}

func TestShakespeareRequestFactory(t *testing.T) {

	tests := []struct {
		name        string
		fileContent string
		host        string
		expected    []request.Request
		err         bool
	}{
		{
			name:        "Valid",
			fileContent: `[{"delay": 0,"method": "GET","path": "test","headers": {},"body": ""}]`,
			host:        "http://localhost:8080",
			expected: []request.Request{
				request.NewRequest(
					0,
					*requestInPEACE("GET", "http://localhost:8080/test", []byte("")),
				),
			},
			err: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "testfile")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.fileContent)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			result, err := shakespeare.ShakespeareRequestFactory(tmpfile.Name(), tt.host)
			if (err != nil) != tt.err {
				t.Errorf("Expected error %v, got %v", tt.err, err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
