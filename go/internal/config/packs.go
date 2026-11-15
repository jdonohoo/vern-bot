package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/jdonohoo/vern-bot/go/internal/packs"
)

// PipelineInfo describes where a pipeline came from.
type PipelineInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// Kind is "software" (the default: tasks with files and acceptance
	// criteria) or "creative" (deliverables such as outlines and beats).
	Kind  string `json:"kind"`
	Pack  string `json:"pack,omitempty"` // "" for config pipelines
	Steps int    `json:"steps"`
}

// applyPacks adds the active packs' pipelines. A pack never replaces a
// pipeline the config already defines.
func applyPacks(cfg *Config) {
	if cfg.Pipelines == nil {
		cfg.Pipelines = map[string][]PipelineStep{}
	}
	cfg.PipelineInfo = map[string]PipelineInfo{}
	for name, steps := range cfg.Pipelines {
		cfg.PipelineInfo[name] = PipelineInfo{Name: name, Kind: "software", Steps: len(steps)}
	}

	pls, errs := packs.Pipelines()
	for _, err := range errs {
		fmt.Fprintf(os.Stderr, "[vern] pack pipeline skipped: %v\n", err)
	}
	for _, pl := range pls {
		if _, exists := cfg.Pipelines[pl.Name]; exists {
			continue
		}
		var steps []PipelineStep
		if err := json.Unmarshal(pl.Steps, &steps); err != nil || len(steps) == 0 {
			fmt.Fprintf(os.Stderr, "[vern] pack %s: pipeline %s has no valid steps (%s)\n", pl.Pack, pl.Name, pl.File)
			continue
		}
		cfg.Pipelines[pl.Name] = normalizeSteps(pl.Name, steps)
		kind := pl.Kind
		if kind == "" {
			kind = "software"
		}
		cfg.PipelineInfo[pl.Name] = PipelineInfo{Name: pl.Name, Description: pl.Description, Kind: kind, Pack: pl.Pack, Steps: len(steps)}
	}
}

// applyProject overlays a project's .vern/config.json.
func applyProject(cfg *Config, proj *packs.Project) {
	if proj == nil {
		return
	}
	if proj.LLMMode != "" {
		cfg.LLMMode = CanonicalMode(proj.LLMMode)
	}
	if proj.Pipeline != "" {
		cfg.PipelineMode = proj.Pipeline
	}
	if proj.AllowWrites != nil {
		cfg.AllowWrites = *proj.AllowWrites
	}
	if len(proj.Tiers) > 0 {
		var tiers map[string]llm.Tier
		if err := json.Unmarshal(proj.Tiers, &tiers); err != nil {
			fmt.Fprintf(os.Stderr, "[vern] %s: tiers ignored: %v\n", proj.Dir, err)
		} else {
			for name, t := range normalizeTiers(tiers) {
				base := cfg.Tiers[name]
				merged := llm.Tier{Models: map[string]string{}, Effort: base.Effort}
				for p, m := range base.Models {
					merged.Models[p] = m
				}
				if _, set := tiers[name]; set {
					for p, m := range tiers[name].Models {
						if c := canonicalLLM(p); c != "" {
							merged.Models[c] = m
						}
					}
					if t.Effort != "" {
						merged.Effort = t.Effort
					}
				}
				cfg.Tiers[name] = merged
			}
		}
	}
	cfg.ProjectDir = proj.Dir
}

// PipelineNames lists the pipelines in this config, sorted.
func (c *Config) PipelineNames() []string {
	names := make([]string, 0, len(c.Pipelines))
	for n := range c.Pipelines {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// HasPipeline reports whether a pipeline with this name exists.
func (c *Config) HasPipeline(name string) bool {
	_, ok := c.Pipelines[name]
	return ok
}
