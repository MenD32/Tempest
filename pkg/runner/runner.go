package runner

import (
	"fmt"

	"github.com/MenD32/Tempest/pkg/client"
	"github.com/MenD32/Tempest/pkg/dump"
	"github.com/MenD32/Tempest/pkg/request"
	"github.com/MenD32/Tempest/pkg/request/shakespeare"
	"k8s.io/klog/v2"
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
		r.config.ResponseBuilder,
		client.NewRecommendedClientConfig(),
	)

	responses, metadata := baseclient.Run(requests)

	if len(responses) != len(requests) {
		klog.Warningf("expected %d responses, got %d", len(requests), len(responses))
	}

	dumper := dump.FileDumper{
		FilePath:             r.config.OutputFile,
		DumpFormatterFactory: dump.DumpJSON,
		StartedAt:            metadata.StartTime,
	}

	dumper.Dump(responses)
	return nil
}
