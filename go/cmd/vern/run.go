package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <llm> <prompt>",
	Short: "Run a single LLM subprocess",
	Long: `Spawn an LLM subprocess with timeout and optional persona context.

Supported LLMs: claude (c), codex (x), gemini (g), copilot (p)

Exit codes:
  0    Success
  1    Usage error / unknown LLM
  124  Timeout (GNU timeout convention)`,
	Args: cobra.ExactArgs(2),
	RunE: runRun,
}

var (
	runOutputFile string
	runPersona    string
	runTimeout    int
)

func init() {
	runCmd.Flags().StringVarP(&runOutputFile, "output", "o", "", "File to save output to")
	runCmd.Flags().StringVarP(&runPersona, "persona", "p", "", "Persona name (loads agents/{persona}.md)")
	runCmd.Flags().IntVarP(&runTimeout, "timeout", "t", 0, "Timeout in seconds (default: VERN_TIMEOUT or 1200)")
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) error {
	llmName := args[0]
	prompt := args[1]

	// Resolve timeout: flag > env > default
	timeout := 1200
	if envTimeout := os.Getenv("VERN_TIMEOUT"); envTimeout != "" {
		fmt.Sscanf(envTimeout, "%d", &timeout)
	}
	if runTimeout > 0 {
		timeout = runTimeout
	}

	// Find agents dir relative to binary
	agentsDir := resolveAgentsDir()

	opts := llm.RunOptions{
		LLM:        llmName,
		Prompt:     prompt,
		OutputFile: runOutputFile,
		Persona:    runPersona,
		Timeout:    time.Duration(timeout) * time.Second,
		WorkingDir: os.Getenv("VERN_WORKING_DIR"),
		AgentsDir:  agentsDir,
	}

	result, err := llm.Run(opts)
	if err != nil {
		return err
	}

	// Print output to stdout (tee behavior when output file is set)
	if result.Output != "" {
		fmt.Print(result.Output)
	}

	if result.ExitCode != 0 {
		os.Exit(result.ExitCode)
	}

	return nil
}

// resolveAgentsDir finds the agents/ directory relative to the binary or project root.
func resolveAgentsDir() string {
	// Try relative to binary: binary is at go/bin/vern or go/cmd/vern/vern
	exe, err := os.Executable()
	if err == nil {
		exe, _ = filepath.EvalSymlinks(exe)
		// Binary at {root}/go/bin/vern → agents at {root}/agents/
		root := filepath.Dir(filepath.Dir(filepath.Dir(exe)))
		candidate := filepath.Join(root, "agents")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		// Binary at {root}/go/cmd/vern/vern → agents at {root}/agents/
		root = filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(exe))))
		candidate = filepath.Join(root, "agents")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}

	// Fallback: look for VERN_ROOT env
	if vernRoot := os.Getenv("VERN_ROOT"); vernRoot != "" {
		return filepath.Join(vernRoot, "agents")
	}

	// Last resort: relative to cwd
	return "agents"
}
