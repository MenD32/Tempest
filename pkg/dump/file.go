package dump

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MenD32/Tempest/pkg/response"
)

type FileDumper struct {
	FilePath             string
	DumpFormatterFactory DumpFormatterFactory
}

func (fd FileDumper) Dump(responses []response.Response) error {
	f, err := os.Create(fd.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	metrics := []response.Metrics{}
	for _, res := range responses {
		metric, err := res.Metrics()
		if err != nil {
			return fmt.Errorf("failed to get metrics: %w", err)
		}
		metrics = append(metrics, *metric)
	}

	data, err := fd.DumpFormatterFactory(metrics)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

type DumpFormatterFactory func([]response.Metrics) ([]byte, error)

func DumpJSON(metrics []response.Metrics) ([]byte, error) {
	return json.Marshal(metrics)
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

func DumpCSV(metrics []response.Metrics) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	rows := MetricsToCSV(metrics)
	if err := writer.WriteAll(rows); err != nil {
		return nil, err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
