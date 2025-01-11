package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	dumpconfig "github.com/MenD32/Tempest/pkg/dump/config"
	requestconfig "github.com/MenD32/Tempest/pkg/request/config"
	responseconfig "github.com/MenD32/Tempest/pkg/response/config"
	"github.com/MenD32/Tempest/pkg/runner"
)

var (
	inputFile    string
	outputFile   string
	host         string
	requestType  string
	responseType string
	outputFormat string
)

var rootCmd = &cobra.Command{
	Short: "Tempest is a benchmarking tool for HTTP Servers",
	Long:  `Tempest is a benchmarking tool for HTTP Servers, with a specialization in AI/ML model serving.`,
	Run: func(cmd *cobra.Command, args []string) {

		Config := runner.Config{
			Host:         host,
			InputFile:    inputFile,
			OutputFile:   outputFile,
			InputType:    requestconfig.RequestFactoryType(requestType),
			ResponseType: responseconfig.ResponseBuilderType(responseType),
			OutputType:   dumpconfig.OutputType(outputFormat),
		}

		CompletedConfig := Config.Complete()
		runner.NewRunner(*CompletedConfig).Run()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file for benchmarking")
	rootCmd.MarkFlagRequired("input")

	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Input file for benchmarking")
	rootCmd.MarkFlagRequired("output")

	rootCmd.Flags().StringVar(&host, "host", "h", "Input file for benchmarking")
	rootCmd.MarkFlagRequired("host")

	rootCmd.Flags().StringVar(&requestType, "request-type", "shakespeare", "Request type (shakespeare)")
	rootCmd.Flags().StringVar(&responseType, "response-type", "openai", "Response format")
	rootCmd.Flags().StringVar(&outputFormat, "output-format", "json", "Output format (json or csv)")
}
