package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MenD32/Shakespeare/pkg/trace"
)

func ShakespeareRequestFactory(shakespeareFilePath string, host string) ([]Request, error) {

	traceLog := trace.TraceLog{}

	shakespeareFile, err := os.Open(shakespeareFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening trace log: %v", err)
	}
	defer shakespeareFile.Close()

	err = json.NewDecoder(shakespeareFile).Decode(&traceLog)
	if err != nil {
		return nil, fmt.Errorf("error decoding trace log: %v", err)
	}

	requests := make([]Request, 0)
	for _, trace := range traceLog {
		httpreq, err := http.NewRequest(
			trace.Method,
			getUrlString(trace, host),
			bytes.NewReader(trace.Body),
		)
		if err != nil {
			return nil, err
		}

		for k, v := range trace.Headers {
			httpreq.Header.Add(k, v)
		}

		req := NewRequest(
			trace.Delay,
			*httpreq,
		)

		requests = append(requests, req)
	}

	return requests, nil
}

func getUrlString(t trace.TraceLogRequest, host string) string {
	return fmt.Sprintf("%s%s", host, t.Path)
}
