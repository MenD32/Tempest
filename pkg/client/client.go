package client

import (
	"net/http"
	"sync"
	"time"

	"github.com/MenD32/Tempest/pkg/dump"
	"github.com/MenD32/Tempest/pkg/request"
	"github.com/MenD32/Tempest/pkg/response"
	"k8s.io/klog/v2"
)

const (
	PRERUN_DELAY = 2 * time.Second // 1 second is too short but 3 seconds is too much
)

func (c client) Run(requests []request.Request) ([]response.Response, dump.Metadata) {
	var traceWaitGroup sync.WaitGroup
	var responseChan = make(chan response.Response, len(requests))

	requestTimings := make(map[request.Request]time.Time, len(requests))

	klog.Infof("Indexing %d requests", len(requests))
	expectedStartTime := time.Now().Add(PRERUN_DELAY)
	for _, req := range requests {
		// The more threads that run concurrently, the lower the accuracy of the sleep until goroutine.
		// To mend this, we substract a small amount of time from the expected start time. relative to the number for threads (approximately index)
		klog.Infof("Request will be sent at %v", expectedStartTime.Add(req.Delay()))
		requestTimings[req] = expectedStartTime.Add(req.Delay())
	}

	for req, calltime := range requestTimings {
		traceWaitGroup.Add(1)
		go func() {
			defer traceWaitGroup.Done()
			time.Sleep(time.Until(calltime))
			c.Send(req, responseChan)
		}()
	}

	time.Sleep(time.Until(expectedStartTime))
	klog.Infof("Starting requests at %v", expectedStartTime)

	traceWaitGroup.Wait()
	close(responseChan)

	klog.Info("All responses received.")

	var responses []response.Response
	for res := range responseChan {
		responses = append(responses, res)
	}
	return responses, dump.Metadata{Count: len(responses), StartTime: expectedStartTime}
}

type client struct {
	respFactory func(*http.Response, time.Time) (response.Response, error)
	loglevel    int
	failOnError bool
}

func (client *client) Send(req request.Request, resChan chan<- response.Response) {

	sent := time.Now()

	klog.Infof("Sending request... %s", req.HTTPRequest().URL)

	httpresp, err := http.DefaultClient.Do(req.HTTPRequest())
	if err != nil {
		klog.Errorf("Error sending request: %v\n", err)
		resChan <- response.ErrorResponse{
			Sent: sent,
			Err:  err,
		}
		return
	}

	resp, err := client.respFactory(httpresp, sent)
	if err != nil {
		klog.Errorf("Error creating response: %v\n", err)
		resChan <- response.ErrorResponse{
			Sent: sent,
			Err:  err,
		}
		return
	}

	resChan <- resp
}

func (client *client) LogLevel() klog.Level {
	return klog.Level(client.loglevel)
}

func NewDefaultClient(respFactory func(*http.Response, time.Time) (response.Response, error), loglevel int, failOnError bool) *client {
	return &client{respFactory: respFactory, loglevel: loglevel, failOnError: failOnError}
}
