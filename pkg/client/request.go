package client

import (
	"net/http"
	"sync"
	"time"
)

type Request interface {
	Delay() time.Duration
	Send(chan<- Response, *sync.WaitGroup) (Response, error)
}

type request struct {
	http.Request

	delay       time.Duration
	respFactory func(*http.Response, time.Time) (Response, error)
}

func (r *request) Delay() time.Duration {
	return r.delay
}

func (r *request) Send(responseChan chan<- Response, requestWaitGroup *sync.WaitGroup) (Response, error) {
	defer requestWaitGroup.Done()

	sent := time.Now()
	resp, err := http.DefaultClient.Do(&r.Request)
	if err != nil {
		return nil, err
	}

	return r.respFactory(resp, sent)
}

func NewRequest(delay time.Duration, req http.Request, respFactory func(*http.Response, time.Time) (Response, error)) Request {
	return &request{
		delay:       delay,
		Request:     req,
		respFactory: respFactory,
	}
}
