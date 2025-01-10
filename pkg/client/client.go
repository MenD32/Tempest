package client

import (
	"sync"
	"time"
)

type Client interface {
	Send(Request, chan<- Response, *sync.WaitGroup)
}

func Run(c Client, requests []Request) []Response {
	var traceWaitGroup sync.WaitGroup
	var requestWaitGroup sync.WaitGroup
	var requestChan = make(chan Request, len(requests))
	var responseChan = make(chan Response, len(requests))

	for _, req := range requests {
		traceWaitGroup.Add(1)
		go func() {
			defer traceWaitGroup.Done()

			time.Sleep(req.Delay())
			requestChan <- req
		}()
	}

	go func() {
		for req := range requestChan {
			requestWaitGroup.Add(1)
			go c.Send(req, responseChan, &requestWaitGroup)
		}
	}()

	traceWaitGroup.Wait()
	close(requestChan)

	requestWaitGroup.Wait()
	close(responseChan)

	var responses []Response
	for res := range responseChan {
		responses = append(responses, res)
	}
	return responses
}
