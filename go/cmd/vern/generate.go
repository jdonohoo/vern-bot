package main

import (
	"fmt"
	"os"

	"github.com/jdonohoo/vern-bot/go/internal/generate"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate <name> <description>",
	Short: "Generate a new Vern persona using AI",
	Long: `Generate a new Vern persona by having an LLM design the personality,
command, and skill files based on a name and description.

Creates 3 files and updates 4 registration files:
  agents/{name}.md          - Persona definition
  commands/{name}.md        - Plugin command
  skills/{name}/SKILL.md    - Plugin skill
  commands/v.md             - Router aliases
  commands/help.md          - Help text
  embedded_test.go          - Test expectations
  council/selection.go      - Hardcoded roster`,
	Args: cobra.ExactArgs(2),
	RunE: runGenerate,
}

var (
	genModel    string
	genColor    string
	genLLM      string
	genDryRun   bool
	genNoUpdate bool
)

func init() {
	generateCmd.Flags().StringVarP(&genModel, "model", "m", "", "Override model (opus/sonnet/haiku/gemini-3/gemini-pro/gemini-flash/codex/codex-mini/copilot/copilot-gpt4)")
	generateCmd.Flags().StringVar(&genColor, "color", "", "Override TUI color")
	generateCmd.Flags().StringVar(&genLLM, "llm", "claude", "LLM for generation")
	generateCmd.Flags().BoolVar(&genDryRun, "dry-run", false, "Print generated content without writing")
	generateCmd.Flags().BoolVar(&genNoUpdate, "no-update", false, "Write core files only, skip registration updates")
	rootCmd.AddCommand(generateCmd)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	name := args[0]
	description := args[1]

	repoRoot, err := generate.DetectRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		fmt.Fprintf(os.Stderr, "Set VERN_ROOT env var or run from within the vern-bot repo.\n")
		os.Exit(1)
	}

	opts := generate.Options{
		Name:        name,
		Description: description,
		Model:       genModel,
		Color:       genColor,
		LLM:         genLLM,
		DryRun:      genDryRun,
		NoUpdate:    genNoUpdate,
		RepoRoot:    repoRoot,
	}

	if err := generate.Run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return nil
}
