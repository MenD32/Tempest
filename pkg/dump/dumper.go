package dump

import (
	"encoding/json"
	"os"

	"github.com/MenD32/Tempest/pkg/response"
)

type Dumper interface {
	Dump([]response.Response)
}

type FileDumper struct {
	FilePath string
}

func (fd *FileDumper) Dump(responses []response.Response) error {
	f, err := os.Create(fd.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	metrics := []response.Metrics{}
	for _, res := range responses {
		metrics = append(metrics, res.Metrics())
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
