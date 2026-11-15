package main

import (
	"fmt"
	"os"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/packs"
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

var (
	allowWrites bool
	packFlags   []string
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&allowWrites, "allow-writes", false,
		"Let the LLM CLIs edit files and run commands (default: read-only; also VERN_ALLOW_WRITES=1 or allow_writes in config)")
	rootCmd.PersistentFlags().StringArrayVar(&packFlags, "pack", nil,
		"Pack directory with extra personas and pipelines (repeatable; also VERN_PATH and .vern/)")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		pipeline.Version = version
		cwd, _ := os.Getwd()
		if _, proj, err := packs.Setup(packFlags, cwd, config.HomeDir()); err != nil && proj != nil {
			fmt.Fprintf(os.Stderr, "[vern] %s/config.json ignored: %v\n", proj.Dir, err)
		}
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
