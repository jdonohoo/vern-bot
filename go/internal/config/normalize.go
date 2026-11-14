package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/llm"
)

// modeAliases maps retired LLM mode names to their replacements.
var modeAliases = map[string]string{
	"mixed_gemini_fallback":  "mixed_antigravity_fallback",
	"mixed_copilot_fallback": "mixed_claude_fallback",
}

// CanonicalMode returns the current name for an LLM mode, mapping retired
// names (mixed_gemini_fallback, mixed_copilot_fallback) to their replacements.
func CanonicalMode(mode string) string {
	if to, ok := modeAliases[mode]; ok {
		return to
	}
	return mode
}

// ModeNames returns the built-in LLM mode names in display order.
func ModeNames() []string {
	names := []string{}
	for _, n := range llm.Names() {
		names = append(names, "mixed_"+n+"_fallback")
	}
	return append(names, "single_llm")
}

// canonicalLLM maps an LLM name or alias to its provider name. Retired
// providers (copilot) become claude; unknown names are kept so a run can
// report them.
func canonicalLLM(name string) string {
	if name == "" {
		return ""
	}
	if c := llm.Canonical(name); c != "" {
		return c
	}
	if _, gone := llm.IsRetired(name); gone {
		return "claude"
	}
	return name
}

// normalize rewrites provider and mode names in a loaded config to the
// current set, so configs written before v3 (gemini, copilot) keep working.
func normalize(cfg *Config) {
	for mode, steps := range cfg.Pipelines {
		for i := range steps {
			steps[i].LLM = canonicalLLM(steps[i].LLM)
		}
		cfg.Pipelines[mode] = steps
	}

	if cfg.LLMs != nil {
		llms := map[string]bool{}
		for name, enabled := range cfg.LLMs {
			c := llm.Canonical(name)
			if c == "" {
				continue // retired or unknown
			}
			llms[c] = llms[c] || enabled
		}
		cfg.LLMs = llms
	}

	modes := map[string]LLMModeConfig{}
	for name, m := range cfg.LLMModes {
		if strings.HasPrefix(name, "mixed_") {
			if _, retired := modeAliases[name]; retired {
				continue // replaced by the built-in mode below
			}
		}
		fb := map[string]string{}
		for from, to := range m.Fallback {
			if c := llm.Canonical(from); c != "" {
				if t := canonicalLLM(to); t != c {
					fb[c] = t
				}
			}
		}
		m.Fallback = fb
		m.SynthesisLLM = canonicalLLM(m.SynthesisLLM)
		m.OverrideLLM = canonicalLLM(m.OverrideLLM)
		modes[name] = m
	}
	for name, m := range defaultLLMModes() {
		if _, ok := modes[name]; !ok {
			modes[name] = m
		}
	}
	cfg.LLMModes = modes

	if to, ok := modeAliases[cfg.LLMMode]; ok {
		fmt.Fprintf(os.Stderr, "[vern] llm_mode %s is retired; using %s\n", cfg.LLMMode, to)
		cfg.LLMMode = to
	}
}
