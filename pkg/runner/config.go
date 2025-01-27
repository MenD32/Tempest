package runner

import (
	"fmt"

	"github.com/MenD32/Tempest/pkg/dump"
	"github.com/MenD32/Tempest/pkg/request"
	"github.com/MenD32/Tempest/pkg/response"

	dumpconfig "github.com/MenD32/Tempest/pkg/dump/config"
	requestconfig "github.com/MenD32/Tempest/pkg/request/config"
	responseconfig "github.com/MenD32/Tempest/pkg/response/config"
)

type Config struct {
	Host       string `json:"host,omitempty"`
	InputFile  string `json:"input_file,omitempty"`
	OutputFile string `json:"output_file,omitempty"`

	InputType    requestconfig.RequestFactoryType   `json:"input_format,omitempty"`
	ResponseType responseconfig.ResponseBuilderType `json:"request_type,omitempty"`
	OutputType   dumpconfig.OutputType              `json:"output_format,omitempty"`
}

type CompletedConfig struct {
	Host       string
	InputFile  string
	OutputFile string

	RequestFactory  request.RequestFactory
	ResponseBuilder response.ResponseBuilder
	Dumper          dump.Dumper

	LogLevel int
}

func (c *Config) Complete() (*CompletedConfig, error) {

	if c.Host == "" {
		return nil, fmt.Errorf("host is required")
	}

	if c.InputFile == "" {
		return nil, fmt.Errorf("input file is required")
	}

	if c.OutputFile == "" {
		return nil, fmt.Errorf("output file is required")
	}

	requestFactory := requestconfig.RequestFactoryFactory(c.InputType)
	if requestFactory == nil {
		return nil, fmt.Errorf("invalid input type: %s", c.InputType)
	}

	responseBuilder := responseconfig.ResponseBuilderFactory(c.ResponseType)
	if responseBuilder == nil {
		return nil, fmt.Errorf("invalid response type: %s", c.ResponseType)
	}

	dumper := dumpconfig.DumperFactory(c.OutputType, c.OutputFile)
	if dumper == nil {
		return nil, fmt.Errorf("invalid output type: %s", c.OutputType)
	}

	return &CompletedConfig{
		Host:       c.Host,
		InputFile:  c.InputFile,
		OutputFile: c.OutputFile,

		RequestFactory:  requestFactory,
		ResponseBuilder: responseBuilder,
		Dumper:          dumper,
	}, nil
}
