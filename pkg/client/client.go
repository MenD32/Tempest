package client

import (
	"net/http"
	"sync"
	"time"

	"github.com/MenD32/Tempest/pkg/request"
	"github.com/MenD32/Tempest/pkg/response"
	"k8s.io/klog/v2"
)

const (
	COMPUTE_OFFSET = 6 * time.Millisecond
)

type Client interface {
	Send(request.Request, chan<- response.Response, *sync.WaitGroup)
	LogLevel() klog.Level
}

func Run(c Client, requests []request.Request) []response.Response {
	var traceWaitGroup sync.WaitGroup
	var requestWaitGroup sync.WaitGroup
	var requestChan = make(chan request.Request, len(requests))
	var responseChan = make(chan response.Response, len(requests))

	requestTimings := make(map[request.Request]time.Time, len(requests))

	klog.Infof("Indexing %d requests", len(requests))
	expectedStartTime := time.Now().Add(1 * time.Second)
	for _, req := range requests {
		requestTimings[req] = expectedStartTime.Add(req.Delay()).Add(COMPUTE_OFFSET)
	}

	time.Sleep(time.Until(expectedStartTime))
	klog.Infof("Starting benchmark")

	for req, calltime := range requestTimings {
		traceWaitGroup.Add(1)
		go func() {
			defer traceWaitGroup.Done()
			time.Sleep(time.Until(calltime))
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

	klog.Info("Finished sending requests, Waiting for responses...")

	requestWaitGroup.Wait()
	close(responseChan)

	klog.V(c.LogLevel()).Info("All responses received.")

	var responses []response.Response
	for res := range responseChan {
		responses = append(responses, res)
	}
	return responses
}

type client struct {
	respFactory func(*http.Response, time.Time) (response.Response, error)
	loglevel    int
}

func (client *client) Send(req request.Request, resChan chan<- response.Response, wg *sync.WaitGroup) {
	defer wg.Done()

	sent := time.Now()
	httpresp, _ := http.DefaultClient.Do(req.HTTPRequest())
	// if err != nil {
	// 	klog.Errorf("Error sending request: %v\n", err)
	// 	return
	// }

	resp, err := client.respFactory(httpresp, sent)
	if err != nil {
		klog.Errorf("Error creating response: %v\n", err)
		return
	}

	resChan <- resp
}

func (client *client) LogLevel() klog.Level {
	return klog.Level(client.loglevel)
}

func NewDefaultClient(respFactory func(*http.Response, time.Time) (response.Response, error), loglevel int) Client {
	return &client{respFactory: respFactory, loglevel: loglevel}
}
