package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

// expandHome replaces a leading ~/ with the user's home directory.
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return home + path[1:]
		}
	}
	return path
}

// LLMModeOptions are the selectable LLM mode options used across screens.
var LLMModeOptions = []huh.Option[string]{
	huh.NewOption("Mixed LLMs + Claude fallback (Recommended)", "mixed_claude_fallback"),
	huh.NewOption("Mixed LLMs + Codex fallback", "mixed_codex_fallback"),
	huh.NewOption("Mixed LLMs + Gemini fallback", "mixed_gemini_fallback"),
	huh.NewOption("Mixed LLMs + Copilot fallback", "mixed_copilot_fallback"),
	huh.NewOption("Single LLM", "single_llm"),
}

// SingleLLMOptions are the selectable LLMs for single-LLM mode.
var SingleLLMOptions = []huh.Option[string]{
	huh.NewOption("Claude", "claude"),
	huh.NewOption("Codex", "codex"),
	huh.NewOption("Gemini", "gemini"),
	huh.NewOption("Copilot", "copilot"),
}

// CouncilOptions are the VernHole council tier options.
var CouncilOptions = []huh.Option[string]{
	huh.NewOption("The Full Vern Experience (15) (Recommended)", "full"),
	huh.NewOption("Fate's Hand (random count, random selection)", "random"),
	huh.NewOption("The War Room (10-13)", "war"),
	huh.NewOption("The Round Table (6-9)", "round"),
	huh.NewOption("Max Conflict (6)", "conflict"),
	huh.NewOption("The Inner Circle (3-5)", "inner"),
	huh.NewOption("Council of the Three Hammers (3)", "hammers"),
}

// VernHoleOptions are council options plus a "No VernHole" option for discovery.
var VernHoleOptions = []huh.Option[string]{
	huh.NewOption("The Full Vern Experience (15) (Recommended)", "full"),
	huh.NewOption("Fate's Hand (random count, random selection)", "random"),
	huh.NewOption("The War Room (10-13)", "war"),
	huh.NewOption("The Round Table (6-9)", "round"),
	huh.NewOption("Max Conflict (6)", "conflict"),
	huh.NewOption("The Inner Circle (3-5)", "inner"),
	huh.NewOption("Council of the Three Hammers (3)", "hammers"),
	huh.NewOption("No VernHole, just the pipeline", ""),
}

// PipelineOptions are the discovery pipeline mode options.
var PipelineOptions = []huh.Option[string]{
	huh.NewOption("Default (5-step) (Recommended)", "default"),
	huh.NewOption("Expanded (7-step)", "expanded"),
}

// OracleApplyOptions are the Oracle vision handling options.
var OracleApplyOptions = []huh.Option[string]{
	huh.NewOption("Vision only — just generate the Oracle's analysis", "vision"),
	huh.NewOption("Auto-apply — Architect Vern rewrites tasks from the vision", "apply"),
}

// OutputPathOptions are the output path choices.
var OutputPathOptions = []huh.Option[string]{
	huh.NewOption("Current directory (Recommended)", "default"),
	huh.NewOption("Custom path", "custom"),
}

// validateName rejects path-traversal characters in a discovery folder name.
func validateName(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.Contains(s, "..") || strings.Contains(s, "/") || strings.Contains(s, "\\") {
		return fmt.Errorf("name cannot contain '..', '/' or '\\'")
	}
	return nil
}

// councilLabel returns the display label for a council value.
func councilLabel(council string) string {
	for _, opt := range CouncilOptions {
		if opt.Value == council {
			return opt.Key
		}
	}
	return council
}
