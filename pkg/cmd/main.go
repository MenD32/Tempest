package main

import (
	"fmt"
	"os"

	"github.com/MenD32/Tempest/pkg/runner"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	dumpconfig "github.com/MenD32/Tempest/pkg/dump/config"
	requestconfig "github.com/MenD32/Tempest/pkg/request/config"
	responseconfig "github.com/MenD32/Tempest/pkg/response/config"
)

var (
	inputFile    string
	outputFile   string
	host         string
	requestType  string
	responseType string
	outputFormat string
)

var (
	Version        = "dev"
	CommitHash     = "none"
	BuildTimestamp = "unknown"
)

func BuildVersion() string {
	return fmt.Sprintf("%s-%s (%s)", Version, CommitHash, BuildTimestamp)
}

var rootCmd = &cobra.Command{
	Short: "Tempest is a benchmarking tool for HTTP Servers",
	Long:  `Tempest is a benchmarking tool for HTTP Servers, with a specialization in AI/ML model serving.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func RunCmd(cmd *cobra.Command, args []string) {
	Config := runner.Config{
		Host:         host,
		InputFile:    inputFile,
		OutputFile:   outputFile,
		InputType:    requestconfig.RequestFactoryType(requestType),
		ResponseType: responseconfig.ResponseBuilderType(responseType),
		OutputType:   dumpconfig.OutputType(outputFormat),
	}

	CompletedConfig, err := Config.Complete()
	if err != nil {
		klog.Errorf("Error completing config: %s\n", err)
		os.Exit(1)
	}
	err = runner.NewRunner(*CompletedConfig).Run()
	if err != nil {
		klog.Errorf("Runtime error: %s\n", err)
		os.Exit(1)
	}
}

func VersionCmd(cmd *cobra.Command, args []string) {
	klog.Infof("Tempest version: %s\n", BuildVersion())
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the benchmark",
	Long:  `Run the benchmark using the specified input and output files, host, request type, response type, and output format.`,
	Run:   RunCmd,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: `Get the version of Tempest.`,
	Long:  `Get the version of Tempest.`,
	Run:   VersionCmd,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		klog.Errorf("%s\n", err)
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

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}
