package dump

import (
	"time"

	"github.com/MenD32/Tempest/pkg/response"
)

type Dumper interface {
	Dump([]response.Response) error
}

type DumpData struct {
	Metrics  []response.Metrics `json:"metrics"`
	Metadata Metadata           `json:"metadata"`
}

type Metadata struct {
	Count     int       `json:"count"`
	StartTime time.Time `json:"start_time"`
}

func NewDumpData(metrics []response.Metrics, starttime time.Time) DumpData {
	return DumpData{
		Metrics: metrics,
		Metadata: Metadata{
			Count:     len(metrics),
			StartTime: starttime,
		},
	}
}
