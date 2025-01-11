package client

import (
	"net/http"
	"sync"
	"time"

	"github.com/MenD32/Tempest/pkg/request"
	"github.com/MenD32/Tempest/pkg/response"
	"k8s.io/klog/v2"
)

type Client interface {
	Send(request.Request, chan<- response.Response, *sync.WaitGroup)
}

func Run(c Client, requests []request.Request) []response.Response {
	var traceWaitGroup sync.WaitGroup
	var requestWaitGroup sync.WaitGroup
	var requestChan = make(chan request.Request, len(requests))
	var responseChan = make(chan response.Response, len(requests))

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

	var responses []response.Response
	for res := range responseChan {
		responses = append(responses, res)
	}
	return responses
}

type client struct {
	respFactory func(*http.Response, time.Time) (response.Response, error)
}

func (client *client) Send(req request.Request, resChan chan<- response.Response, wg *sync.WaitGroup) {
	defer wg.Done()

	sent := time.Now()
	httpresp, err := http.DefaultClient.Do(req.HTTPRequest())
	if err != nil {
		klog.Errorf("Error sending request: %v\n", err)
		return
	}

	resp, err := client.respFactory(httpresp, sent)
	if err != nil {
		klog.Errorf("Error creating response: %v\n", err)
		return
	}

	resChan <- resp
}

func NewDefaultClient(respFactory func(*http.Response, time.Time) (response.Response, error)) Client {
	return &client{respFactory: respFactory}
}
