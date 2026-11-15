package main

import (
	"fmt"
	"os"

	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "vern",
	Short:   "Vern-Bot CLI — multi-LLM discovery pipeline",
	Long:    "Vern CLI orchestrates multi-LLM discovery pipelines, VernHole councils, and task management.",
	Version: version,
}

var allowWrites bool

func init() {
	rootCmd.PersistentFlags().BoolVar(&allowWrites, "allow-writes", false,
		"Let the LLM CLIs edit files and run commands (default: read-only; also VERN_ALLOW_WRITES=1 or allow_writes in config)")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		pipeline.Version = version
		if allowWrites {
			// An env var reaches every llm.Run, including the TUI's.
			os.Setenv("VERN_ALLOW_WRITES", "1")
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
