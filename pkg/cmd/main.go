package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/MenD32/Tempest/pkg/client"
	"github.com/MenD32/Tempest/pkg/responses"
)

var (
	inputFilePath  string = "../Shakespeare/temp/test.json"
	outputFilePath string = "./temp/output.json"
	host           string = "http://localhost:8000"
)

var rootCmd = &cobra.Command{
	Short: "Tempest is a benchmarking tool for HTTP Servers",
	Long:  `Tempest is a benchmarking tool for HTTP Servers, with a specialization in AI/ML model serving.`,
	Run: func(cmd *cobra.Command, args []string) {

		requests, err := client.ShakespeareRequestFactory(
			inputFilePath,
			host,
		)
		if err != nil {
			fmt.Printf("Error creating requests: %v\n", err)
			os.Exit(1)
		}

		baseclient := client.NewClient(
			responses.OpenAIResponseFactory,
		)

		responses := client.Run(baseclient, requests)

		dumper := client.FileDumper{
			FilePath: outputFilePath,
		}

		dumper.Dump(responses)
	},
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
