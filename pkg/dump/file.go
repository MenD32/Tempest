package dump

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/MenD32/Tempest/pkg/response"
)

type FileDumper struct {
	FilePath             string
	DumpFormatterFactory DumpFormatterFactory

	StartedAt time.Time
}

func (fd FileDumper) Dump(responses []response.Response) error {
	f, err := os.Create(fd.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	metrics := []response.Metrics{}
	for _, res := range responses {
		if res == nil {
			return fmt.Errorf("response is nil")
		}
		metric := res.Metrics()
		metrics = append(metrics, *metric)
	}

	dd := NewDumpData(metrics, fd.StartedAt)

	data, err := fd.DumpFormatterFactory(dd)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

type DumpFormatterFactory func(DumpData) ([]byte, error)

func DumpJSON(dd DumpData) ([]byte, error) {
	return json.Marshal(
		struct {
			Metrics  []response.Metrics `json:"metrics"`
			Metadata struct {
				Count     int       `json:"count"`
				StartTime time.Time `json:"start_time"`
			} `json:"metadata"`
		}{
			Metrics: dd.Metrics,
			Metadata: struct {
				Count     int       `json:"count"`
				StartTime time.Time `json:"start_time"`
			}{
				Count:     len(dd.Metrics),
				StartTime: dd.Metadata.StartTime,
			},
		},
	)
}

func MetricsToCSV(metrics []response.Metrics) [][]string {
	columns := []string{
		"sent",
		"body",
	}
	for key := range metrics[0].Metrics {
		columns = append(columns, key)
	}

	rows := [][]string{}
	for _, metric := range metrics {
		row := []string{
			metric.Sent.String(),
			string(metric.Body),
		}
		for _, value := range metric.Metrics {
			row = append(row, fmt.Sprintf("%v", value))
		}
		rows = append(rows, row)
	}

	return rows
}

func DumpCSV(dd DumpData) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	rows := MetricsToCSV(dd.Metrics)
	if err := writer.WriteAll(rows); err != nil {
		return nil, err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
