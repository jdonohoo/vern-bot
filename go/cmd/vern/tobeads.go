package main

import (
	"fmt"
	"os"

	"github.com/jdonohoo/vern-bot/go/internal/tobeads"
	"github.com/spf13/cobra"
)

var tobeadsApply bool
var tobeadsSync bool
var tobeadsBeadsDir string

var tobeadsCmd = &cobra.Command{
	Use:   "tobeads <vts-directory>",
	Short: "Import VTS tasks into Beads",
	Long:  "One-way importer: reads VTS task files and creates Beads issues via br CLI. Dry-run by default.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vtsDir := args[0]

		// Verify directory exists
		info, err := os.Stat(vtsDir)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("VTS directory not found: %s", vtsDir)
		}

		runner := &tobeads.RealBrRunner{
			WorkDir: tobeadsBeadsDir,
		}

		opts := tobeads.ImportOptions{
			VTSDir:   vtsDir,
			BeadsDir: tobeadsBeadsDir,
			Apply:    tobeadsApply,
			Sync:     tobeadsSync,
		}

		result, err := tobeads.Run(opts, runner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			// Exit code 1 for preflight failures
			os.Exit(1)
		}

		if result != nil && result.Failed > 0 {
			os.Exit(2)
		}

		return nil
	},
}

func init() {
	tobeadsCmd.Flags().BoolVar(&tobeadsApply, "apply", false, "Actually create issues (default: dry-run)")
	tobeadsCmd.Flags().BoolVar(&tobeadsSync, "sync", false, "Run br sync --flush-only after apply")
	tobeadsCmd.Flags().StringVar(&tobeadsBeadsDir, "beads-dir", "", "Target Beads repo directory")
	rootCmd.AddCommand(tobeadsCmd)
}
