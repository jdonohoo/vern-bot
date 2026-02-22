package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

		// Resolve to absolute path
		absVTSDir, err := filepath.Abs(vtsDir)
		if err != nil {
			return fmt.Errorf("resolve VTS path: %w", err)
		}

		// Verify directory exists
		info, err := os.Stat(absVTSDir)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("VTS directory not found: %s", vtsDir)
		}

		// Verify br CLI is on PATH
		if _, err := exec.LookPath("br"); err != nil {
			return fmt.Errorf("br CLI not found on PATH\n  Install: cargo install beads_rust\n  Or build from source and add to PATH")
		}

		// Resolve beads directory: explicit flag > walk up from VTS dir
		beadsDir := tobeadsBeadsDir
		if beadsDir == "" {
			beadsDir = findBeadsDir(absVTSDir)
		}
		if beadsDir != "" {
			beadsDir, _ = filepath.Abs(beadsDir)
		}

		if beadsDir == "" {
			return fmt.Errorf(".beads/ not found in ancestry of %s\n  Either run: br init --prefix <prefix> in the target project\n  Or specify: --beads-dir <path>", absVTSDir)
		}
		beadsDB := filepath.Join(beadsDir, ".beads")
		if info, err := os.Stat(beadsDB); err != nil || !info.IsDir() {
			return fmt.Errorf(".beads/ not found in %s\n  Run: br init --prefix <prefix>\n  in the target project directory first", beadsDir)
		}

		runner := &tobeads.RealBrRunner{
			WorkDir: beadsDir,
		}

		opts := tobeads.ImportOptions{
			VTSDir:   absVTSDir,
			BeadsDir: beadsDir,
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

// findBeadsDir walks up from startDir looking for a .beads/ directory.
func findBeadsDir(startDir string) string {
	dir := startDir
	for {
		candidate := filepath.Join(dir, ".beads")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "" // reached root
		}
		dir = parent
	}
}

func init() {
	tobeadsCmd.Flags().BoolVar(&tobeadsApply, "apply", false, "Actually create issues (default: dry-run)")
	tobeadsCmd.Flags().BoolVar(&tobeadsSync, "sync", false, "Run br sync --flush-only after apply")
	tobeadsCmd.Flags().StringVar(&tobeadsBeadsDir, "beads-dir", "", "Target Beads repo directory (auto-detected from VTS dir if omitted)")
	rootCmd.AddCommand(tobeadsCmd)
}
