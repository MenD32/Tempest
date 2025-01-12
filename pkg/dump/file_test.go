package dump_test

import (
	"os"
	"testing"

	"github.com/MenD32/Tempest/pkg/dump"
	"github.com/MenD32/Tempest/pkg/response"
	"github.com/MenD32/Tempest/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestFileDumper_DumpJSON(t *testing.T) {
	filePath := "test_output.json"
	defer os.Remove(filePath)

	dumper := dump.FileDumper{
		FilePath:             filePath,
		DumpFormatterFactory: dump.DumpJSON,
	}

	responses := make([]response.Response, 100)
	for i := range responses {
		responses[i] = test.TestResponse{}
	}

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
	for i := range responses {
		responses[i] = test.TestResponse{}
	}

	err := dumper.Dump(responses)
	assert.NoError(t, err)

	_, err = os.Stat(filePath)
	assert.NoError(t, err)
}
