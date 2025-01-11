package runner

import (
	"fmt"

	"github.com/MenD32/Tempest/pkg/client"
	"github.com/MenD32/Tempest/pkg/dump"
	"github.com/MenD32/Tempest/pkg/request"
	"github.com/MenD32/Tempest/pkg/request/shakespeare"
	"github.com/MenD32/Tempest/pkg/response/openai"
)

type Runner struct {
	config CompletedConfig
}

func NewRunner(config CompletedConfig) *Runner {
	return &Runner{
		config: config,
	}
}

func (r *Runner) Run() error {

	var requestFactory request.RequestFactory = shakespeare.ShakespeareRequestFactory

	requests, err := requestFactory(
		r.config.InputFile,
		r.config.Host,
	)
	if err != nil {
		return fmt.Errorf("error creating requests: %w", err)
	}

	baseclient := client.NewDefaultClient(
		openai.OpenAIResponseBuilder,
	)

	responses := client.Run(baseclient, requests)

	dumper := dump.FileDumper{
		FilePath:             r.config.OutputFile,
		DumpFormatterFactory: dump.DumpJSON,
	}

	dumper.Dump(responses)
	return nil
}
