package main

import (
	"github.com/jdonohoo/vern-bot/go/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive Bubble Tea TUI",
	Long:  "Start the Vern-Bot terminal UI for interactive discovery pipelines, VernHole councils, and settings.",
	RunE:  runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func runTUI(cmd *cobra.Command, args []string) error {
	agentsDir := resolveAgentsDir()

	// Find project root
	projectRoot := ""
	if agentsDir != "agents" && len(agentsDir) > len("/agents") {
		projectRoot = agentsDir[:len(agentsDir)-len("/agents")]
	}

	return tui.Run(projectRoot, agentsDir, version)
}
