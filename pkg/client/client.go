package client

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var start time.Time

type Client interface {
	Send(Request, chan<- Response, *sync.WaitGroup)
}

func Run(c Client, requests []Request) []Response {
	var traceWaitGroup sync.WaitGroup
	var requestWaitGroup sync.WaitGroup
	var requestChan = make(chan Request, len(requests))
	var responseChan = make(chan Response, len(requests))

	start = time.Now()

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

type client struct {
	respFactory func(*http.Response, time.Time) (Response, error)
}

func (client *client) Send(req Request, resChan chan<- Response, wg *sync.WaitGroup) {
	defer wg.Done()

	sent := time.Now()
	fmt.Printf("Sent request at %v\n", time.Since(start))
	httpresp, err := http.DefaultClient.Do(req.HTTPRequest())
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}

	resp, err := client.respFactory(httpresp, sent)
	if err != nil {
		fmt.Printf("Error creating response: %v\n", err)
		return
	}

	resChan <- resp
}

func NewClient(respFactory func(*http.Response, time.Time) (Response, error)) Client {
	return &client{respFactory: respFactory}
}
