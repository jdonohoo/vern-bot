package llm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jdonohoo/vern-bot/go/internal/embedded"
)

// Tier is a persona's weight class. Config maps each tier to a model per
// provider, so personas never name a vendor model.
type Tier struct {
	// Models maps a canonical provider name to the model id passed to it.
	// A provider with no entry runs on its CLI default.
	Models map[string]string `json:"models,omitempty"`
	// Effort is the reasoning effort for this tier ("" = CLI default).
	Effort string `json:"effort,omitempty"`
}

// Tier names.
const (
	TierDeep     = "deep"
	TierBalanced = "balanced"
	TierFast     = "fast"
)

// TierNames lists the tiers from heaviest to lightest.
var TierNames = []string{TierDeep, TierBalanced, TierFast}

// DefaultTiers maps tiers to Claude's model aliases and leaves every other
// provider on its CLI default.
func DefaultTiers() map[string]Tier {
	return map[string]Tier{
		TierDeep:     {Models: map[string]string{"claude": "opus"}},
		TierBalanced: {Models: map[string]string{"claude": "sonnet"}},
		TierFast:     {Models: map[string]string{"claude": "haiku"}},
	}
}

// Policy holds the process-wide run settings that come from config.
type Policy struct {
	// AllowWrites lets the CLIs' agents edit files and run commands.
	// The default (false) runs every CLI read-only where it supports it;
	// vern itself still writes the output files.
	AllowWrites bool
	// Tiers maps a tier name to its models and effort.
	Tiers map[string]Tier
}

var (
	policyMu sync.RWMutex
	policy   = Policy{Tiers: DefaultTiers()}
	warnOnce sync.Map
)

// SetPolicy replaces the process-wide policy. config.Load calls it.
func SetPolicy(p Policy) {
	if p.Tiers == nil {
		p.Tiers = DefaultTiers()
	}
	policyMu.Lock()
	policy = p
	policyMu.Unlock()
}

// CurrentPolicy returns the process-wide policy.
func CurrentPolicy() Policy {
	policyMu.RLock()
	defer policyMu.RUnlock()
	return policy
}

// WritesAllowed reports whether runs may modify files: the policy, or
// VERN_ALLOW_WRITES=1 (set by the --allow-writes flag).
func WritesAllowed() bool {
	if v := os.Getenv("VERN_ALLOW_WRITES"); v == "1" || strings.EqualFold(v, "true") {
		return true
	}
	return CurrentPolicy().AllowWrites
}

// PersonaTier returns a persona's tier: its `tier:` frontmatter, else one
// derived from its Claude Code `model:` (opus = deep, sonnet = balanced,
// haiku = fast). Returns "" when the persona is unknown or has neither.
func PersonaTier(agentsDir, persona string) string {
	if persona == "" {
		return ""
	}
	fm := personaFrontmatter(agentsDir, persona)
	if t := strings.ToLower(fm["tier"]); t != "" {
		return t
	}
	switch strings.ToLower(fm["model"]) {
	case "opus", "fable":
		return TierDeep
	case "sonnet":
		return TierBalanced
	case "haiku":
		return TierFast
	}
	return ""
}

// personaFrontmatter reads the key: value lines of a persona's frontmatter
// from agentsDir, falling back to the agents compiled into the binary.
func personaFrontmatter(agentsDir, persona string) map[string]string {
	out := map[string]string{}
	var content string
	if agentsDir != "" {
		if data, err := os.ReadFile(filepath.Join(agentsDir, persona+".md")); err == nil {
			content = string(data)
		}
	}
	if content == "" {
		if c, ok := embedded.GetAgent(persona); ok {
			content = c
		}
	}
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return out
	}
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "---" {
			break
		}
		if i := strings.Index(line, ":"); i > 0 {
			out[strings.TrimSpace(line[:i])] = strings.TrimSpace(line[i+1:])
		}
	}
	return out
}

// Selection is the model and effort a run resolved to.
type Selection struct {
	Tier   string
	Model  string
	Effort string
}

// selectModel resolves the model and effort for provider llm:
// explicit Models entry > tier mapping > CLI default, and explicit
// Effort > tier effort.
func selectModel(opts RunOptions, provider string) Selection {
	sel := Selection{Tier: strings.ToLower(opts.Tier)}
	if sel.Tier == "" {
		sel.Tier = PersonaTier(opts.AgentsDir, opts.Persona)
	}
	if t, ok := CurrentPolicy().Tiers[sel.Tier]; ok {
		sel.Model = t.Models[provider]
		sel.Effort = t.Effort
	}
	if m := opts.Models[provider]; m != "" {
		sel.Model = m
	}
	if opts.Effort != "" {
		sel.Effort = opts.Effort
	}
	return sel
}

// ModelsFor builds RunOptions.Models for a model meant for one provider.
// Returns nil when model is "".
func ModelsFor(provider, model string) map[string]string {
	if model == "" {
		return nil
	}
	name := Canonical(provider)
	if name == "" {
		name = strings.ToLower(provider)
	}
	return map[string]string{name: model}
}

func warnReadOnlyUnsupported(p *Provider) {
	if _, seen := warnOnce.LoadOrStore(p.Name, true); seen {
		return
	}
	fmt.Fprintf(os.Stderr, "[vern-run] Warning: %s has no read-only mode; it runs with write access\n", p.Name)
}
