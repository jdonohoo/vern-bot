package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/embedded"
	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/jdonohoo/vern-bot/go/internal/packs"
	"github.com/spf13/cobra"
)

var modelsJSON bool

var modelsCmd = &cobra.Command{
	Use:   "models [llm]",
	Short: "Show persona tiers and the models each CLI offers",
	Long: `Show how persona tiers (deep, balanced, fast) map to a model and effort
per LLM, then list the models each installed CLI offers (codex and agy can
list theirs). Change the mapping with "tiers" in your config, for example:

  "tiers": {
    "deep": {"models": {"claude": "opus", "codex": "gpt-6-astra"}, "effort": "high"}
  }`,
	Args: cobra.MaximumNArgs(1),
	RunE: runModels,
}

func init() {
	modelsCmd.Flags().BoolVar(&modelsJSON, "json", false, "Print as JSON")
	rootCmd.AddCommand(modelsCmd)
}

type modelsProvider struct {
	Name      string   `json:"name"`
	Installed bool     `json:"installed"`
	Models    []string `json:"models,omitempty"`
	Hint      string   `json:"hint,omitempty"`
}

type modelsReport struct {
	Tiers     map[string]llm.Tier `json:"tiers"`
	Personas  map[string]string   `json:"personas"`
	Providers []modelsProvider    `json:"providers"`
}

func runModels(cmd *cobra.Command, args []string) error {
	agentsDir := resolveAgentsDir()
	cfg := config.Load(filepath.Dir(agentsDir))

	providers := llm.Providers()
	if len(args) == 1 {
		p, ok := llm.Lookup(args[0])
		if !ok {
			return fmt.Errorf("unknown LLM: %s (valid: %s)", args[0], strings.Join(llm.Names(), ", "))
		}
		providers = []*llm.Provider{p}
	}

	rep := modelsReport{Tiers: cfg.Tiers, Personas: personaTiers(agentsDir)}
	for _, p := range providers {
		mp := modelsProvider{Name: p.Name, Installed: p.Available(), Hint: p.ModelsHint}
		if mp.Installed {
			if ids, raw, ok := p.ListModels(); ok {
				mp.Models = ids
				if len(ids) == 0 && mp.Hint == "" {
					mp.Hint = "no model ids found in: " + firstLines(raw, 5)
				}
			} else if mp.Hint == "" && len(p.ModelsArgs) > 0 {
				mp.Hint = "could not list models (" + p.Binary + " " + strings.Join(p.ModelsArgs, " ") + " failed)"
			}
		}
		rep.Providers = append(rep.Providers, mp)
	}

	if modelsJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(rep)
	}

	fmt.Println("Tiers (model per LLM; blank = the CLI's default):")
	for _, name := range llm.TierNames {
		t := rep.Tiers[name]
		var cells []string
		for _, p := range llm.Names() {
			m := t.Models[p]
			if m == "" {
				m = "default"
			}
			cells = append(cells, p+"="+m)
		}
		effort := t.Effort
		if effort == "" {
			effort = "default"
		}
		fmt.Printf("  %-9s %s · effort=%s\n", name, strings.Join(cells, " "), effort)
	}

	byTier := map[string][]string{}
	for persona, tier := range rep.Personas {
		if tier == "" {
			tier = "(none)"
		}
		byTier[tier] = append(byTier[tier], persona)
	}
	fmt.Println("\nPersonas by tier:")
	for _, tier := range append(append([]string{}, llm.TierNames...), "(none)") {
		if ps := byTier[tier]; len(ps) > 0 {
			sort.Strings(ps)
			fmt.Printf("  %-9s %s\n", tier, strings.Join(ps, ", "))
		}
	}

	for _, mp := range rep.Providers {
		fmt.Printf("\n%s:\n", mp.Name)
		switch {
		case !mp.Installed:
			fmt.Println("  not installed")
		case len(mp.Models) > 0:
			for _, id := range mp.Models {
				fmt.Printf("  %s\n", id)
			}
		case mp.Hint != "":
			fmt.Printf("  %s\n", mp.Hint)
		}
	}
	return nil
}

// personaTiers maps every known persona to its tier.
func personaTiers(agentsDir string) map[string]string {
	out := map[string]string{}
	names := map[string]bool{}
	if entries, err := os.ReadDir(agentsDir); err == nil {
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				names[strings.TrimSuffix(e.Name(), ".md")] = true
			}
		}
	}
	if len(names) == 0 {
		for _, n := range embeddedAgentNames() {
			names[n] = true
		}
	}
	for n := range packs.PackAgentFiles() {
		names[n] = true
	}
	for n := range names {
		out[n] = llm.PersonaTier(agentsDir, n)
	}
	return out
}

func embeddedAgentNames() []string { return embedded.ListAgents() }

// firstLines returns up to n non-empty lines of s, joined with " / ".
func firstLines(s string, n int) string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
			if len(out) == n {
				break
			}
		}
	}
	if len(out) == 0 {
		return "(empty output)"
	}
	return strings.Join(out, " / ")
}
