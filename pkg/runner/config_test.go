
package runner_test

import (
	"testing"

	"github.com/MenD32/Tempest/pkg/runner"
	requestconfig "github.com/MenD32/Tempest/pkg/request/config"
	responseconfig "github.com/MenD32/Tempest/pkg/response/config"
	dumpconfig "github.com/MenD32/Tempest/pkg/dump/config"
)

func TestRunnerConfigComplete(t *testing.T) {
	tests := []struct {
		name    string
		config  runner.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: runner.Config{
				Host:         "localhost",
				InputFile:    "input.txt",
				OutputFile:   "output.txt",
				InputType:    requestconfig.ShakespeareRequestFactoryType,
				ResponseType: responseconfig.EmptyResponseType,
				OutputType:   dumpconfig.JSONOutputType,
			},
			wantErr: false,
		},
		{
			name: "missing host",
			config: runner.Config{
				InputFile:    "input.txt",
				OutputFile:   "output.txt",
				InputType:    requestconfig.ShakespeareRequestFactoryType,
				ResponseType: responseconfig.EmptyResponseType,
				OutputType:   dumpconfig.JSONOutputType,
			},
			wantErr: true,
		},
		{
			name: "missing input file",
			config: runner.Config{
				Host:         "localhost",
				OutputFile:   "output.txt",
				InputType:    requestconfig.ShakespeareRequestFactoryType,
				ResponseType: responseconfig.EmptyResponseType,
				OutputType:   dumpconfig.JSONOutputType,
			},
			wantErr: true,
		},
		{
			name: "missing output file",
			config: runner.Config{
				Host:         "localhost",
				InputFile:    "input.txt",
				InputType:    requestconfig.ShakespeareRequestFactoryType,
				ResponseType: responseconfig.EmptyResponseType,
				OutputType:   dumpconfig.JSONOutputType,
			},
			wantErr: true,
		},
		{
			name: "missing input type",
			config: runner.Config{
				Host:         "localhost",
				InputFile:    "input.txt",
				OutputFile:   "output.txt",
				ResponseType: responseconfig.EmptyResponseType,
				OutputType:   dumpconfig.JSONOutputType,
			},
			wantErr: true,
		},
		{
			name: "invalid input type",
			config: runner.Config{
				Host:         "localhost",
				InputFile:    "input.txt",
				OutputFile:   "output.txt",
				InputType:    "invalid",
				ResponseType: responseconfig.EmptyResponseType,
				OutputType:   dumpconfig.JSONOutputType,
			},
			wantErr: true,
		},
		{
			name: "missing response type",
			config: runner.Config{
				Host:         "localhost",
				InputFile:    "input.txt",
				OutputFile:   "output.txt",
				InputType:    requestconfig.ShakespeareRequestFactoryType,
				OutputType:   dumpconfig.JSONOutputType,
			},
			wantErr: true,
		},
		{
			name: "invalid response type",
			config: runner.Config{
				Host:         "localhost",
				InputFile:    "input.txt",
				OutputFile:   "output.txt",
				InputType:    requestconfig.ShakespeareRequestFactoryType,
				ResponseType: "invalid",
				OutputType:   dumpconfig.JSONOutputType,
			},
			wantErr: true,
		},
		{
			name: "missing output type",
			config: runner.Config{
				Host:         "localhost",
				InputFile:    "input.txt",
				OutputFile:   "output.txt",
				InputType:    requestconfig.ShakespeareRequestFactoryType,
				ResponseType: responseconfig.EmptyResponseType,
			},
			wantErr: true,
		},
		{
			name: "invalid output type",
			config: runner.Config{
				Host:         "localhost",
				InputFile:    "input.txt",
				OutputFile:   "output.txt",
				InputType:    requestconfig.ShakespeareRequestFactoryType,
				ResponseType: responseconfig.EmptyResponseType,
				OutputType:   "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.config.Complete()
			if (err != nil) != tt.wantErr {
				t.Errorf("Complete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}