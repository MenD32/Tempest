package config

import (
	"github.com/MenD32/Tempest/pkg/dump"
)

type DumpFactoryConfig interface {
	GetName() string
	GetDumper() dump.Dumper
}

type OutputType string

const (
	JSONOutputType OutputType = "JSON"
	CSVOutputType  OutputType = "CSV"
)

func DumperFactory(outputType OutputType, filepath string) dump.Dumper {
	switch outputType {
	case JSONOutputType:
		return dump.FileDumper{
			FilePath:             filepath,
			DumpFormatterFactory: dump.DumpJSON,
		}
	case CSVOutputType:
		return dump.FileDumper{
			FilePath:             filepath,
			DumpFormatterFactory: dump.DumpCSV,
		}
	}
	return nil
}
