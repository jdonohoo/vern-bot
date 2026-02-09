package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "vern",
	Short: "Vern-Bot CLI — multi-LLM discovery pipeline",
	Long:  "Vern CLI orchestrates multi-LLM discovery pipelines, VernHole councils, and task management.",
	Version: version,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
