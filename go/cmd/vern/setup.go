package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive first-run configuration wizard",
	Long: `Configure vern-bot: detect available LLMs, set default LLM mode,
and save configuration for standalone use.

Config is saved to ~/.config/vern/config.json.`,
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) error {
	fmt.Println("=== VERN-BOT SETUP ===")
	fmt.Println()

	// Detect available LLMs
	llms := map[string]bool{
		"claude":  detectCLI("claude"),
		"codex":   detectCLI("codex"),
		"gemini":  detectCLI("gemini"),
		"copilot": detectCLI("copilot"),
	}

	fmt.Println("Detected LLMs:")
	for _, name := range []string{"claude", "codex", "gemini", "copilot"} {
		status := "not found"
		if llms[name] {
			status = "found"
		}
		fmt.Printf("  %s: %s\n", name, status)
	}
	fmt.Println()

	// Build config
	cfg := config.Load("")

	// Update LLMs based on detection
	cfg.LLMs = llms

	// Default LLM mode based on what's available
	cfg.LLMMode = "mixed_claude_fallback"
	if !llms["claude"] {
		// If no claude, pick first available as fallback
		for _, name := range []string{"codex", "gemini", "copilot"} {
			if llms[name] {
				cfg.LLMMode = "mixed_" + name + "_fallback"
				break
			}
		}
	}

	fmt.Printf("Default LLM mode: %s\n", cfg.LLMMode)
	fmt.Printf("Default pipeline: %s\n", cfg.PipelineMode)
	fmt.Println()

	// Save config
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "vern")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Printf("Config saved to: %s\n", configPath)
	fmt.Println()

	// Summary
	fmt.Println("=== VERN-BOT CONFIGURED ===")
	fmt.Println()
	fmt.Printf("LLM Mode: %s\n", cfg.LLMMode)
	fmt.Print("LLMs: ")
	for _, name := range []string{"claude", "codex", "gemini", "copilot"} {
		if llms[name] {
			fmt.Printf("%s ✓  ", name)
		} else {
			fmt.Printf("%s ✗  ", name)
		}
	}
	fmt.Println()
	fmt.Printf("Pipeline: %s\n", cfg.PipelineMode)
	fmt.Println()
	fmt.Println("Run 'vern tui' for interactive mode, or 'vern discovery' to start a pipeline.")

	return nil
}

func detectCLI(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
