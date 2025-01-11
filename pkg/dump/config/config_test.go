package config_test

import (
	"testing"

	"github.com/MenD32/Tempest/pkg/dump"
	"github.com/MenD32/Tempest/pkg/dump/config"
)

func TestDumperFactory(t *testing.T) {
	tests := []struct {
		outputType config.OutputType
		filepath   string
		expected   dump.Dumper
	}{
		{
			outputType: config.JSONOutputType,
			filepath:   "test.json",
			expected: dump.FileDumper{
				FilePath:             "test.json",
				DumpFormatterFactory: dump.DumpJSON,
			},
		},
		{
			outputType: config.CSVOutputType,
			filepath:   "test.csv",
			expected: dump.FileDumper{
				FilePath:             "test.csv",
				DumpFormatterFactory: dump.DumpJSON,
			},
		},
		{
			outputType: "INVALID",
			filepath:   "test.invalid",
			expected:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.outputType), func(t *testing.T) {
			result := config.DumperFactory(tt.outputType, tt.filepath)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}