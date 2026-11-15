package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jdonohoo/vern-bot/go/internal/llm"
)

// A config written by v2 (gemini and copilot) keeps working.
func TestNormalizeLegacyConfig(t *testing.T) {
	legacy := `{
  "llm_mode": "mixed_gemini_fallback",
  "llms": {"claude": true, "codex": false, "gemini": true, "copilot": true},
  "llm_modes": {
    "mixed_claude_fallback": {"fallback": {"codex": "claude", "gemini": "claude", "copilot": "claude"}, "synthesis_llm": "claude"},
    "mixed_gemini_fallback": {"fallback": {"claude": "gemini"}, "synthesis_llm": "gemini"},
    "mixed_copilot_fallback": {"fallback": {"claude": "copilot"}, "synthesis_llm": "copilot"},
    "single_llm": {"override_llm": "gemini", "fallback": {}}
  },
  "discovery_pipelines": {"default": [{"step": 1, "name": "Chaos", "persona": "yolo", "llm": "gemini"}, {"step": 2, "name": "X", "persona": "great", "llm": "copilot"}]}
}`
	path := filepath.Join(t.TempDir(), "config.json")
	os.WriteFile(path, []byte(legacy), 0644)
	cfg, err := loadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.LLMMode != "mixed_antigravity_fallback" {
		t.Errorf("llm_mode = %q", cfg.LLMMode)
	}
	if _, ok := cfg.LLMModes["mixed_copilot_fallback"]; ok {
		t.Error("copilot mode should be dropped")
	}
	if cfg.GetSynthesisLLM() != "antigravity" {
		t.Errorf("synthesis = %q", cfg.GetSynthesisLLM())
	}
	if fb := cfg.GetFallbackLLM("gemini"); fb != "" {
		t.Errorf("antigravity is the safety net in this mode, fallback = %q", fb)
	}
	if fb := cfg.GetFallbackLLM("claude"); fb != "antigravity" {
		t.Errorf("claude fallback = %q", fb)
	}

	cfg.LLMMode = "mixed_claude_fallback"
	if fb := cfg.GetFallbackLLM("gemini"); fb != "claude" {
		t.Errorf("gemini (antigravity) fallback = %q", fb)
	}
	// Every other provider falls back to claude, including providers added
	// after the config was written (grok).
	if fb := cfg.LLMModes["mixed_claude_fallback"].Fallback; len(fb) != len(llm.Names())-1 || fb["antigravity"] != "claude" || fb["grok"] != "claude" {
		t.Errorf("claude mode fallback = %v", fb)
	}

	cfg.LLMMode = "single_llm"
	if cfg.GetOverrideLLM() != "antigravity" {
		t.Errorf("override = %q", cfg.GetOverrideLLM())
	}

	steps := cfg.Pipelines["default"]
	if steps[0].LLM != "antigravity" || steps[1].LLM != "claude" {
		t.Errorf("steps = %q, %q", steps[0].LLM, steps[1].LLM)
	}
	if _, ok := cfg.LLMs["copilot"]; ok || !cfg.LLMs["antigravity"] || len(cfg.LLMs) != 3 { // claude, codex, antigravity
		t.Errorf("llms = %v", cfg.LLMs)
	}
}

func TestRetiredModeFlag(t *testing.T) {
	cfg := hardcodedDefaults()
	cfg.LLMMode = "mixed_gemini_fallback" // e.g. from --llm-mode
	if cfg.GetSynthesisLLM() != "antigravity" {
		t.Errorf("retired mode name should resolve, synthesis = %q", cfg.GetSynthesisLLM())
	}
}

func TestDefaultModesCoverProviders(t *testing.T) {
	modes := defaultLLMModes()
	for _, name := range ModeNames() {
		if _, ok := modes[name]; !ok {
			t.Errorf("missing default mode %s", name)
		}
	}
	if len(modes) != len(ModeNames()) {
		t.Errorf("modes = %d, names = %d", len(modes), len(ModeNames()))
	}
}

func TestNormalizeTiers(t *testing.T) {
	cfg := &Config{
		Tiers: map[string]llm.Tier{
			"Deep": {Models: map[string]string{"gemini": "gemini-4-argon-high", "claude": ""}, Effort: "HIGH"},
			"fast": {Effort: "warp"},
		},
		Pipelines: map[string][]PipelineStep{
			"default": {{Step: 1, LLM: "g", Effort: "Max", Tier: "FAST"}, {Step: 2, LLM: "claude", Effort: "nope"}},
		},
	}
	normalize(cfg)

	deep := cfg.Tiers["deep"]
	if deep.Models["antigravity"] != "gemini-4-argon-high" || deep.Effort != "high" {
		t.Errorf("deep tier not normalized: %+v", deep)
	}
	if m, ok := deep.Models["claude"]; !ok || m != "" {
		t.Errorf("an explicit empty model should keep claude on its default: %+v", deep)
	}
	if cfg.Tiers["fast"].Effort != "" {
		t.Errorf("unknown effort should be dropped: %+v", cfg.Tiers["fast"])
	}
	if cfg.Tiers["balanced"].Models["claude"] != "sonnet" {
		t.Errorf("missing tiers should get defaults: %+v", cfg.Tiers)
	}
	steps := cfg.Pipelines["default"]
	if steps[0].Effort != "max" || steps[0].Tier != "fast" || steps[0].LLM != "antigravity" {
		t.Errorf("step not normalized: %+v", steps[0])
	}
	if steps[1].Effort != "" {
		t.Errorf("unknown step effort should be dropped: %+v", steps[1])
	}
}

func TestLoadSetsRunPolicy(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)
	os.MkdirAll(filepath.Join(dir, ".config", "vern"), 0755)
	os.WriteFile(filepath.Join(dir, ".config", "vern", "config.json"),
		[]byte(`{"allow_writes": true, "tiers": {"deep": {"models": {"codex": "gpt-6-astra"}}}}`), 0644)
	defer llm.SetPolicy(llm.Policy{})

	Load(dir)
	p := llm.CurrentPolicy()
	if !p.AllowWrites {
		t.Error("allow_writes should reach the run policy")
	}
	if p.Tiers["deep"].Models["codex"] != "gpt-6-astra" || p.Tiers["deep"].Models["claude"] != "opus" {
		t.Errorf("tiers should reach the run policy: %+v", p.Tiers)
	}
}
