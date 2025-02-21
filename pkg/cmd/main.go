package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/MenD32/Tempest/pkg/runner"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	dumpconfig "github.com/MenD32/Tempest/pkg/dump/config"
	requestconfig "github.com/MenD32/Tempest/pkg/request/config"
	responseconfig "github.com/MenD32/Tempest/pkg/response/config"
)

var (
	// runner config
	inputFile    string
	outputFile   string
	host         string
	requestType  string
	responseType string
	outputFormat string

	// debug config
	enablePprof bool
	pprofFile   string
	enableTrace bool
	traceFile   string
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

	klog.Infof("%s\n", responseType)

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

	f, _ := os.Create("temp/cpu.prof")
	defer f.Close()

	t, _ := os.Create("temp/trace.out")
	defer t.Close()

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	trace.Start(t)
	defer trace.Stop()

	if err := rootCmd.Execute(); err != nil {
		klog.Errorf("%s\n", err)
		os.Exit(1)
	}
}

func init() {
	runCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file for benchmarking")
	runCmd.MarkFlagRequired("input")

	runCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for benchmark results")
	runCmd.MarkFlagRequired("output")

	runCmd.Flags().StringVar(&host, "host", "", "Host to send requests to")
	runCmd.MarkFlagRequired("host")

	runCmd.Flags().StringVar(&requestType, "request-type", "Shakespeare", "Request type (shakespeare)")
	runCmd.Flags().StringVar(&responseType, "response-type", "openai", "Response format")
	runCmd.Flags().StringVar(&outputFormat, "output-format", "JSON", "Output format (json or csv)")

	runCmd.Flags().BoolVar(&enablePprof, "enable-pprof", false, "Enable pprof profiling")
	runCmd.Flags().StringVar(&pprofFile, "pprof-file", "cpu.prof", "File to write pprof output to")
	runCmd.Flags().BoolVar(&enableTrace, "enable-trace", false, "Enable pprof profiling")
	runCmd.Flags().StringVar(&traceFile, "trace-file", "trace.out", "File to write trace output to")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}
