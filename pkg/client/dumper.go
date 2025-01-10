package client

import (
	"encoding/json"
	"os"
)

type Dumper interface {
	Dump([]Response)
}

type FileDumper struct {
	FilePath string
}

func (fd *FileDumper) Dump(responses []Response) error {
	f, err := os.Create(fd.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	metrics := []Metrics{}
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
