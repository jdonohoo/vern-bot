package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the vern-bot configuration.
type Config struct {
	Version        string                      `json:"version"`
	TimeoutSeconds int                         `json:"timeout_seconds"`
	MaxRetries     int                         `json:"max_retries"`
	PipelineMode   string                      `json:"pipeline_mode"`
	Pipelines      map[string][]PipelineStep   `json:"discovery_pipelines"`
	LLMs           map[string]bool             `json:"llms"`
	VernHole       VernHoleConfig              `json:"vernhole"`

	// Backward compat: old config format
	LegacyPipeline []PipelineStep              `json:"discovery_pipeline"`
}

// PipelineStep defines a single step in a discovery pipeline.
type PipelineStep struct {
	Step         int    `json:"step"`
	Name         string `json:"name"`
	Persona      string `json:"persona"`
	LLM          string `json:"llm"`
	ContextMode  string `json:"context_mode"`
	PromptPrefix string `json:"prompt_prefix"`
}

// VernHoleConfig holds VernHole-specific settings.
type VernHoleConfig struct {
	DefaultCouncil string `json:"default_council"`
	Min            int    `json:"min"`
}

// Load reads configuration using the 3-tier chain:
//  1. ~/.claude/vern-bot-config.json (user config)
//  2. {projectRoot}/config.default.json (project defaults)
//  3. Hardcoded defaults
func Load(projectRoot string) *Config {
	// Tier 1: user config
	userConfig := filepath.Join(os.Getenv("HOME"), ".claude", "vern-bot-config.json")
	if cfg, err := loadFile(userConfig); err == nil {
		return cfg
	}

	// Tier 2: project default config
	defaultConfig := filepath.Join(projectRoot, "config.default.json")
	if cfg, err := loadFile(defaultConfig); err == nil {
		return cfg
	}

	// Tier 3: hardcoded defaults
	return hardcodedDefaults()
}

func loadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	// Backward compat: migrate legacy discovery_pipeline to discovery_pipelines
	if cfg.Pipelines == nil && cfg.LegacyPipeline != nil {
		cfg.Pipelines = map[string][]PipelineStep{
			"default": cfg.LegacyPipeline,
		}
		cfg.PipelineMode = "default"
	}

	// Apply defaults for missing fields
	if cfg.TimeoutSeconds == 0 {
		cfg.TimeoutSeconds = 1200
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 2
	}
	if cfg.PipelineMode == "" {
		cfg.PipelineMode = "default"
	}
	if cfg.VernHole.Min == 0 {
		cfg.VernHole.Min = 3
	}
	if cfg.VernHole.DefaultCouncil == "" {
		cfg.VernHole.DefaultCouncil = "random"
	}

	return cfg, nil
}

// GetPipeline returns the pipeline steps for the given mode.
// Falls back to "default" if the mode doesn't exist.
func (c *Config) GetPipeline(mode string) []PipelineStep {
	if mode == "" {
		mode = c.PipelineMode
	}
	if steps, ok := c.Pipelines[mode]; ok {
		return steps
	}
	if steps, ok := c.Pipelines["default"]; ok {
		return steps
	}
	return hardcodedDefaults().Pipelines["default"]
}

func hardcodedDefaults() *Config {
	return &Config{
		Version:        "1.6.0",
		TimeoutSeconds: 1200,
		MaxRetries:     2,
		PipelineMode:   "default",
		VernHole: VernHoleConfig{
			DefaultCouncil: "random",
			Min:            3,
		},
		LLMs: map[string]bool{
			"claude": true,
			"codex":  true,
			"gemini": true,
		},
		Pipelines: map[string][]PipelineStep{
			"default": {
				{Step: 1, Name: "Initial Analysis", Persona: "mighty", LLM: "codex", ContextMode: "prompt_only",
					PromptPrefix: "You are MightyVern. Analyze this idea and provide comprehensive initial analysis including: problem space, technical requirements, proposed architecture, unknowns and risks."},
				{Step: 2, Name: "Refinement", Persona: "great", LLM: "claude", ContextMode: "previous",
					PromptPrefix: "You are Vernile the Great. Review and refine this analysis. Identify gaps, add architectural considerations, consider maintainability and elegance."},
				{Step: 3, Name: "Chaos Check", Persona: "yolo", LLM: "gemini", ContextMode: "previous",
					PromptPrefix: "You are YOLO Vern. Challenge and stress-test this plan. What could go wrong? What unconventional approaches exist? No sacred cows."},
				{Step: 4, Name: "Consolidation", Persona: "mighty", LLM: "codex", ContextMode: "all_previous",
					PromptPrefix: "You are MightyVern. Synthesize all inputs into a master plan. Merge insights, resolve contradictions, create unified vision, prioritize features."},
				{Step: 5, Name: "Architect Breakdown", Persona: "architect", LLM: "claude", ContextMode: "consolidation",
					PromptPrefix: "You are Architect Vern. Break down this master plan into actionable Vern Task Spec (VTS) tasks. Format each task with an h3 header exactly like this: ### TASK 1: Title Here. Include for each task: **Description:** what needs to be done, **Acceptance Criteria:** bullet list, **Complexity:** S|M|L|XL, **Dependencies:** Task N references or None, **Files:** list of files likely touched. Think in systems. Consider failure modes. Make it maintainable."},
			},
			"expanded": {
				{Step: 1, Name: "Initial Analysis", Persona: "mighty", LLM: "codex", ContextMode: "prompt_only",
					PromptPrefix: "You are MightyVern. Analyze this idea and provide comprehensive initial analysis including: problem space, technical requirements, proposed architecture, unknowns and risks."},
				{Step: 2, Name: "Refinement", Persona: "great", LLM: "claude", ContextMode: "previous",
					PromptPrefix: "You are Vernile the Great. Review and refine this analysis. Identify gaps, add architectural considerations, consider maintainability and elegance."},
				{Step: 3, Name: "Reality Check", Persona: "mediocre", LLM: "claude", ContextMode: "previous",
					PromptPrefix: "You are Vern the Mediocre. Reality-check this plan. What's over-engineered? What can be simplified? Where is cleverness hiding complexity? Cut the fluff, keep what ships."},
				{Step: 4, Name: "Chaos Check", Persona: "yolo", LLM: "gemini", ContextMode: "previous",
					PromptPrefix: "You are YOLO Vern. Challenge and stress-test this plan. What could go wrong? What unconventional approaches exist? No sacred cows."},
				{Step: 5, Name: "MVP Lens", Persona: "startup", LLM: "claude", ContextMode: "previous",
					PromptPrefix: "You are Startup Vern. What's the MVP here? Cut scope ruthlessly. What can ship in week one? What's a nice-to-have disguised as a must-have? If you're not embarrassed by v1, you shipped too late."},
				{Step: 6, Name: "Consolidation", Persona: "mighty", LLM: "codex", ContextMode: "all_previous",
					PromptPrefix: "You are MightyVern. Synthesize all inputs into a master plan. Merge insights, resolve contradictions, create unified vision, prioritize features."},
				{Step: 7, Name: "Architect Breakdown", Persona: "architect", LLM: "claude", ContextMode: "consolidation",
					PromptPrefix: "You are Architect Vern. Break down this master plan into actionable Vern Task Spec (VTS) tasks. Format each task with an h3 header exactly like this: ### TASK 1: Title Here. Include for each task: **Description:** what needs to be done, **Acceptance Criteria:** bullet list, **Complexity:** S|M|L|XL, **Dependencies:** Task N references or None, **Files:** list of files likely touched. Think in systems. Consider failure modes. Make it maintainable."},
			},
		},
	}
}
