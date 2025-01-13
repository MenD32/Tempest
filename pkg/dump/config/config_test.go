package config_test

import (
	"reflect"
	"testing"

	"github.com/MenD32/Tempest/pkg/dump"
	"github.com/MenD32/Tempest/pkg/dump/config"
)

func TestDumpFactoryConfig(t *testing.T) {
	tests := []struct {
		name       string
		outputType config.OutputType
		filepath   string
		expected   dump.Dumper
	}{
		{
			name:       "JSON",
			outputType: config.JSONOutputType,
			filepath:   "",
			expected: dump.FileDumper{
				DumpFormatterFactory: dump.DumpJSON,
				FilePath:             "",
			},
		},
		{
			name:       "CSV",
			outputType: config.CSVOutputType,
			filepath:   "",
			expected: dump.FileDumper{
				DumpFormatterFactory: dump.DumpCSV,
				FilePath:             "",
			},
		},
		{
			name:       "Invalid",
			outputType: "Invalid",
			filepath:   "",
			expected:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := config.DumperFactory(tt.outputType, tt.filepath)
			if !(reflect.TypeOf(actual) == reflect.TypeOf(tt.expected)) {
				t.Errorf("Expected %v, got %v", tt.expected, actual)
				return
			}
			if fd, ok := actual.(dump.FileDumper); ok {
				expected_fd, _ := tt.expected.(dump.FileDumper)
				if reflect.ValueOf(fd.DumpFormatterFactory) != reflect.ValueOf(expected_fd.DumpFormatterFactory) {
					t.Errorf("Expected %v, got %v", expected_fd.DumpFormatterFactory, fd.DumpFormatterFactory)
				}
			}
		})
	}

}
