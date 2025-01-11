package dump_test

import (
	"os"
	"testing"
	"time"

	"github.com/MenD32/Tempest/pkg/dump"
	"github.com/MenD32/Tempest/pkg/response"
	"github.com/stretchr/testify/assert"
)

type TestResponse struct{}

func (tr TestResponse) Metrics() (*response.Metrics, error) {
	return &response.Metrics{
		Sent: time.Now(),
		Body: []byte(""),
		Metrics: map[string]interface{}{
			"test": "test",
		},
	}, nil
}

func (tr TestResponse) Verify() error {
	return nil
}

func TestFileDumper_DumpJSON(t *testing.T) {
	filePath := "test_output.json"
	defer os.Remove(filePath)

	dumper := dump.FileDumper{
		FilePath:             filePath,
		DumpFormatterFactory: dump.DumpJSON,
	}

	responses := make([]response.Response, 100)

	err := dumper.Dump(responses)
	assert.NoError(t, err)

	_, err = os.Stat(filePath)
	assert.NoError(t, err)
}

func TestFileDumper_DumpCSV(t *testing.T) {
	filePath := "test_output.csv"
	defer os.Remove(filePath)

	dumper := dump.FileDumper{
		FilePath:             filePath,
		DumpFormatterFactory: dump.DumpCSV,
	}

	responses := make([]response.Response, 100)

	err := dumper.Dump(responses)
	assert.NoError(t, err)

	_, err = os.Stat(filePath)
	assert.NoError(t, err)
}
