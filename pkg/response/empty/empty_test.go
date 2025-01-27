package empty

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmptyResponseBuilder(t *testing.T) {
	sentTime := time.Now()
	resp := &http.Response{}
	emptyResp, err := EmptyResponseBuilder(resp, sentTime)

	assert.NoError(t, err)
	assert.Equal(t, sentTime, emptyResp.(EmptyResponse).Sent)
}

func TestEmptyResponse_Metrics(t *testing.T) {
	sentTime := time.Now()
	emptyResp := EmptyResponse{Sent: sentTime}
	metrics := emptyResp.Metrics()

	assert.Equal(t, sentTime, metrics.Sent)
}

func TestEmptyResponse_Verify(t *testing.T) {
	emptyResp := EmptyResponse{}
	err := emptyResp.Verify()

	assert.NoError(t, err)
}
