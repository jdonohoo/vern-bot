package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <llm> <prompt>",
	Short: "Run a single LLM subprocess",
	Long: `Spawn an LLM subprocess with timeout and optional persona context.

Supported LLMs: claude (c), codex (x), antigravity (a; also agy, gemini), grok (k)

The CLI runs read-only unless --allow-writes is given. --model is passed to
that LLM only (a fallback to claude uses claude's default); --effort is one
of low, medium, high, xhigh, max. Without them the persona's tier decides
(see "vern models").

--json prints the output with the LLM, model, effort, exit code and the
usage the CLI reported, instead of the bare output.

Exit codes:
  0    Success
  1    Usage error / unknown LLM / the CLI failed
  124  Timeout (GNU timeout convention)
  other  the CLI's own exit code`,
	Args: cobra.ExactArgs(2),
	RunE: runRun,
}

var (
	runOutputFile string
	runPersona    string
	runTimeout    int
	runModel      string
	runEffort     string
	runTier       string
	runJSON       bool
)

func init() {
	runCmd.Flags().StringVarP(&runOutputFile, "output", "o", "", "File to save output to")
	runCmd.Flags().StringVarP(&runPersona, "persona", "p", "", "Persona name (loads agents/{persona}.md)")
	runCmd.Flags().IntVarP(&runTimeout, "timeout", "t", 0, "Timeout in seconds (default: VERN_TIMEOUT or 1200)")
	runCmd.Flags().StringVarP(&runModel, "model", "m", "", "Model id for this LLM (see vern models)")
	runCmd.Flags().StringVar(&runEffort, "effort", "", "Reasoning effort: low, medium, high, xhigh, max")
	runCmd.Flags().StringVar(&runTier, "tier", "", "Persona tier to use: deep, balanced, fast")
	runCmd.Flags().BoolVar(&runJSON, "json", false, "Print the result as JSON (output, llm, model, exit code, usage)")
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

	// Load config for its run policy (read-only, tiers).
	config.Load(filepath.Dir(agentsDir))

	opts := llm.RunOptions{
		LLM:        llmName,
		Prompt:     prompt,
		OutputFile: runOutputFile,
		Persona:    runPersona,
		Timeout:    time.Duration(timeout) * time.Second,
		WorkingDir: os.Getenv("VERN_WORKING_DIR"),
		AgentsDir:  agentsDir,
		Models:     llm.ModelsFor(llmName, runModel),
		Effort:     strings.ToLower(runEffort),
		Tier:       strings.ToLower(runTier),
		Plain:      true,
	}

	result, err := llm.Run(opts)
	if err != nil {
		return err
	}

	if runJSON {
		out := struct {
			LLM        string     `json:"llm"`
			Model      string     `json:"model,omitempty"`
			Effort     string     `json:"effort,omitempty"`
			ReadOnly   bool       `json:"read_only"`
			ExitCode   int        `json:"exit_code"`
			TimedOut   bool       `json:"timed_out"`
			DurationMS int64      `json:"duration_ms"`
			Usage      *llm.Usage `json:"usage,omitempty"`
			Output     string     `json:"output"`
			Error      string     `json:"error,omitempty"`
		}{
			LLM: result.LLMUsed, Model: result.Model, Effort: result.Effort, ReadOnly: result.ReadOnly,
			ExitCode: result.ExitCode, TimedOut: result.TimedOut, DurationMS: result.Duration.Milliseconds(),
			Output: result.Output,
		}
		if !result.Usage.IsZero() {
			out.Usage = &result.Usage
		}
		if result.ExitCode != 0 {
			out.Error = llm.FirstLine(result.Stderr)
		}
		if err := printJSON(out); err != nil {
			return err
		}
		if result.ExitCode != 0 {
			os.Exit(result.ExitCode)
		}
		return nil
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
