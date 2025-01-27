package dump_test

import (
	"testing"
	"time"

	"github.com/MenD32/Tempest/pkg/dump"
	"github.com/MenD32/Tempest/pkg/response"
)

func TestNewDumpData(t *testing.T) {
	startTime := time.Now()
	metrics := []response.Metrics{
		{Sent: time.Now().Add(100 * time.Millisecond)},
		{Sent: time.Now().Add(200 * time.Millisecond)},
	}

	dumpData := dump.NewDumpData(metrics, startTime)

	if len(dumpData.Metrics) != len(metrics) {
		t.Fatalf("Expected %d metrics, got %d", len(metrics), len(dumpData.Metrics))
	}

	if dumpData.Metadata.Count != len(metrics) {
		t.Fatalf("Expected count to be %d, got %d", len(metrics), dumpData.Metadata.Count)
	}

	if !dumpData.Metadata.StartTime.Equal(startTime) {
		t.Fatalf("Expected start time to be %v, got %v", startTime, dumpData.Metadata.StartTime)
	}
}
