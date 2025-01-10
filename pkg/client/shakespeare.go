package client

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/MenD32/Shakespeare/pkg/trace"
)

func ShakespeareRequestFactory(shakespeareFilePath string, host string, respFactory func(*http.Response, time.Time) (Response, error)) ([]Request, error) {

	traceLog, err := trace.Load(shakespeareFilePath)
	if err != nil {
		return nil, err
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
			respFactory,
		)

		requests = append(requests, req)
	}

	return requests, nil
}

func getUrlString(t trace.TraceLogRequest, host string) string {
	return fmt.Sprintf("%s%s", host, t.Path)
}
