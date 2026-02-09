package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	configJSON := `{
		"version": "1.6.0",
		"timeout_seconds": 600,
		"max_retries": 2,
		"pipeline_mode": "default",
		"discovery_pipelines": {
			"default": [
				{
					"step": 1,
					"name": "Test Step",
					"persona": "mighty",
					"llm": "codex",
					"context_mode": "prompt_only",
					"prompt_prefix": "Test prompt"
				}
			]
		},
		"vernhole": {
			"default_council": "hammers",
			"min": 3
		}
	}`
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte(configJSON), 0644)

	cfg, err := loadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TimeoutSeconds != 600 {
		t.Errorf("timeout: got %d, want 600", cfg.TimeoutSeconds)
	}
	if cfg.MaxRetries != 2 {
		t.Errorf("max_retries: got %d, want 2", cfg.MaxRetries)
	}
	if cfg.VernHole.DefaultCouncil != "hammers" {
		t.Errorf("council: got %q, want %q", cfg.VernHole.DefaultCouncil, "hammers")
	}

	steps := cfg.GetPipeline("default")
	if len(steps) != 1 {
		t.Fatalf("steps: got %d, want 1", len(steps))
	}
	if steps[0].Name != "Test Step" {
		t.Errorf("step name: got %q, want %q", steps[0].Name, "Test Step")
	}
}

func TestLegacyConfig(t *testing.T) {
	dir := t.TempDir()
	configJSON := `{
		"version": "1.0.0",
		"discovery_pipeline": [
			{
				"step": 1,
				"name": "Legacy Step",
				"persona": "mighty",
				"llm": "claude",
				"context_mode": "prompt_only",
				"prompt_prefix": "Legacy prompt"
			}
		]
	}`
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte(configJSON), 0644)

	cfg, err := loadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	steps := cfg.GetPipeline("default")
	if len(steps) != 1 {
		t.Fatalf("steps: got %d, want 1", len(steps))
	}
	if steps[0].Name != "Legacy Step" {
		t.Errorf("step name: got %q, want %q", steps[0].Name, "Legacy Step")
	}
}

func TestHardcodedDefaults(t *testing.T) {
	cfg := hardcodedDefaults()
	if cfg.TimeoutSeconds != 1200 {
		t.Errorf("timeout: got %d, want 1200", cfg.TimeoutSeconds)
	}

	defaultSteps := cfg.GetPipeline("default")
	if len(defaultSteps) != 5 {
		t.Errorf("default steps: got %d, want 5", len(defaultSteps))
	}

	expandedSteps := cfg.GetPipeline("expanded")
	if len(expandedSteps) != 7 {
		t.Errorf("expanded steps: got %d, want 7", len(expandedSteps))
	}
}

func TestGetPipelineFallback(t *testing.T) {
	cfg := &Config{
		Pipelines: map[string][]PipelineStep{
			"default": {{Step: 1, Name: "Only Default"}},
		},
	}

	// Request non-existent pipeline, should fall back to default
	steps := cfg.GetPipeline("nonexistent")
	if len(steps) != 1 || steps[0].Name != "Only Default" {
		t.Errorf("fallback failed: got %v", steps)
	}
}

func TestLoadWithProjectRoot(t *testing.T) {
	// Load should work with the actual project config.default.json
	cfg := Load("/Users/justin/projects/jdonohoo/vern-bot")
	if cfg == nil {
		t.Fatal("Load returned nil")
	}
	defaultSteps := cfg.GetPipeline("default")
	if len(defaultSteps) < 3 {
		t.Errorf("expected at least 3 default steps, got %d", len(defaultSteps))
	}
}
