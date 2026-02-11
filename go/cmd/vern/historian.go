package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
	"github.com/spf13/cobra"
)

var historianCmd = &cobra.Command{
	Use:   "historian <directory>",
	Short: "Index a directory of input files into a structured concept map",
	Long: `Historian Vern reads everything. Point it at a directory and it will
ingest all readable files (.md, .txt, .json, .yaml, .yml), build a structured
concept map with source references, and write input-history.md.

Uses Gemini by default for its 2M context window. Falls back to the configured
fallback LLM if Gemini is not available.`,
	Args: cobra.ExactArgs(1),
	RunE: runHistorian,
}

var (
	historianLLM     string
	historianTimeout int
)

func init() {
	historianCmd.Flags().StringVar(&historianLLM, "llm", "", "Override LLM (default: gemini, falls back if unavailable)")
	historianCmd.Flags().IntVar(&historianTimeout, "timeout", 600, "Timeout in seconds (default: 600)")
	rootCmd.AddCommand(historianCmd)
}

func runHistorian(cmd *cobra.Command, args []string) error {
	targetDir := args[0]

	// Resolve to absolute path
	absDir, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("resolve directory: %w", err)
	}

	// Verify directory exists
	info, err := os.Stat(absDir)
	if err != nil {
		return fmt.Errorf("directory not found: %s", absDir)
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", absDir)
	}

	agentsDir := resolveAgentsDir()

	outputFile := filepath.Join(absDir, "input-history.md")

	// Check for prompt.md in the target directory
	promptFile := filepath.Join(absDir, "prompt.md")
	if _, err := os.Stat(promptFile); os.IsNotExist(err) {
		promptFile = "" // don't try to update if it doesn't exist
	}

	fmt.Printf("=== HISTORIAN VERN ===\n")
	fmt.Printf("Target: %s\n", absDir)
	fmt.Printf("Output: %s\n", outputFile)
	fmt.Println()

	result, err := pipeline.RunHistorian(pipeline.HistorianOptions{
		TargetDir:  absDir,
		OutputFile: outputFile,
		PromptFile: promptFile,
		AgentsDir:  agentsDir,
		Timeout:    historianTimeout,
		LLMName:    historianLLM,
		OnLog: func(msg string) {
			fmt.Printf("    %s\n", msg)
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("=== HISTORIAN COMPLETE ===\n")
	fmt.Printf("Output: %s\n", result.OutputFile)
	fmt.Printf("Size: %d chars\n", result.CharCount)
	fmt.Printf("Duration: %s\n", result.Duration.Round(100*time.Millisecond))
	fmt.Printf("LLM: %s\n", result.LLMUsed)
	if result.FellBack {
		fmt.Printf("\nWARNING: Gemini is not installed. Used %s as fallback.\n", result.LLMUsed)
		fmt.Printf("The Historian is designed for Gemini's 2M token context window.\n")
		fmt.Printf("Using %s may produce incomplete results on large input directories.\n", result.LLMUsed)
		fmt.Printf("Install the Gemini CLI for best results.\n")
	}

	return nil
}
