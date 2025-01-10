package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Short: "Tempest is a benchmarking tool for HTTP Servers",
		Long:  `Tempest is a benchmarking tool for HTTP Servers, with a specialization in AI/ML model serving.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			fmt.Println("Hello, Cobra!")
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
