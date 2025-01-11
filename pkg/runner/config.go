package runner

import (
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
}

func (c *Config) Complete() *CompletedConfig {
	return &CompletedConfig{
		Host:       c.Host,
		InputFile:  c.InputFile,
		OutputFile: c.OutputFile,

		RequestFactory:  requestconfig.RequestFactoryFactory(c.InputType),
		ResponseBuilder: responseconfig.ResponseBuilderFactory(c.ResponseType),
		Dumper:          dumpconfig.DumperFactory(c.OutputType, c.OutputFile),
	}
}
