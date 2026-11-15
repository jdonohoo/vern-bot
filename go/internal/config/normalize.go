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
			steps[i].Tier = strings.ToLower(steps[i].Tier)
			steps[i].Effort = strings.ToLower(steps[i].Effort)
			if !llm.ValidEffort(steps[i].Effort) {
				fmt.Fprintf(os.Stderr, "[vern] pipeline %s step %d: unknown effort %q ignored\n", mode, steps[i].Step, steps[i].Effort)
				steps[i].Effort = ""
			}
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
		// A built-in mixed mode covers every provider, including ones added
		// after the config was written.
		if safety := strings.TrimSuffix(strings.TrimPrefix(name, "mixed_"), "_fallback"); llm.Canonical(safety) == safety && name == "mixed_"+safety+"_fallback" {
			for _, other := range llm.Names() {
				if _, ok := fb[other]; !ok && other != safety {
					fb[other] = safety
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

	cfg.Tiers = normalizeTiers(cfg.Tiers)

	if to, ok := modeAliases[cfg.LLMMode]; ok {
		fmt.Fprintf(os.Stderr, "[vern] llm_mode %s is retired; using %s\n", cfg.LLMMode, to)
		cfg.LLMMode = to
	}
}

// normalizeTiers lower-cases tier names, canonicalizes provider keys and
// fills in the default tiers and models the config leaves out. An empty
// model string ("") keeps that provider on its CLI default, so a config
// can opt out of a default mapping.
func normalizeTiers(in map[string]llm.Tier) map[string]llm.Tier {
	out := map[string]llm.Tier{}
	for name, t := range in {
		models := map[string]string{}
		for p, m := range t.Models {
			if c := canonicalLLM(p); c != "" {
				models[c] = m
			}
		}
		effort := strings.ToLower(t.Effort)
		if !llm.ValidEffort(effort) {
			fmt.Fprintf(os.Stderr, "[vern] tier %s: unknown effort %q ignored\n", name, t.Effort)
			effort = ""
		}
		out[strings.ToLower(name)] = llm.Tier{Models: models, Effort: effort}
	}
	// Defaults fill in every tier and provider the config leaves out.
	for name, def := range llm.DefaultTiers() {
		t, ok := out[name]
		if !ok {
			out[name] = def
			continue
		}
		for p, m := range def.Models {
			if _, set := t.Models[p]; !set {
				t.Models[p] = m
			}
		}
		out[name] = t
	}
	return out
}
