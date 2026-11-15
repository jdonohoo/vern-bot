package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/packs"
	"github.com/spf13/cobra"
)

var packsJSON, pipelinesJSON bool

var packsCmd = &cobra.Command{
	Use:   "packs",
	Short: "List the packs in use and the personas and pipelines they add",
	Long: `List the packs vern found, in search order (first wins): --pack flags, the
project's .vern/ folder and the packs its config.json lists, VERN_PATH, and
~/.config/vern/packs/*. A pack is a folder with agents/*.md (personas),
pipelines/*.json (discovery pipelines) and an optional pack.json. See
docs/packs.md.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		type packOut struct {
			packs.Pack
			Personas  []string `json:"personas"`
			Pipelines []string `json:"pipelines"`
		}
		var out struct {
			Project *packs.Project `json:"project,omitempty"`
			Packs   []packOut      `json:"packs"`
		}
		out.Project = packs.ActiveProject()
		pls, _ := packs.Pipelines()
		for _, p := range packs.Active() {
			po := packOut{Pack: p, Personas: []string{}, Pipelines: []string{}}
			files, _ := filepath.Glob(filepath.Join(p.Dir, "agents", "*.md"))
			for _, f := range files {
				po.Personas = append(po.Personas, strings.TrimSuffix(filepath.Base(f), ".md"))
			}
			for _, pl := range pls {
				if pl.Pack == p.Name {
					po.Pipelines = append(po.Pipelines, pl.Name)
				}
			}
			out.Packs = append(out.Packs, po)
		}
		if packsJSON {
			return printJSON(out)
		}
		if out.Project != nil {
			fmt.Printf("Project config: %s\n", out.Project.Dir)
			if out.Project.Sink != "" {
				fmt.Printf("  sink: %s\n", out.Project.Sink)
			}
			if out.Project.Pipeline != "" {
				fmt.Printf("  pipeline: %s\n", out.Project.Pipeline)
			}
			fmt.Println()
		}
		if len(out.Packs) == 0 {
			fmt.Println("No packs found. Add one with --pack <dir>, VERN_PATH, a project .vern/ folder, or ~/.config/vern/packs/<name>/.")
			return nil
		}
		for _, p := range out.Packs {
			fmt.Printf("%s (%s) %s\n", p.Name, p.Source, p.Dir)
			if p.Description != "" {
				fmt.Printf("  %s\n", p.Description)
			}
			if len(p.Personas) > 0 {
				fmt.Printf("  personas:  %s\n", strings.Join(p.Personas, ", "))
			}
			if len(p.Pipelines) > 0 {
				fmt.Printf("  pipelines: %s\n", strings.Join(p.Pipelines, ", "))
			}
		}
		return nil
	},
}

var pipelinesCmd = &cobra.Command{
	Use:   "pipelines",
	Short: "List the discovery pipelines you can pass to --pipeline",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		agentsDir := resolveAgentsDir()
		cfg := config.Load(filepath.Dir(agentsDir))
		infos := make([]config.PipelineInfo, 0, len(cfg.PipelineInfo))
		for _, name := range cfg.PipelineNames() {
			info := cfg.PipelineInfo[name]
			if info.Name == "" {
				info = config.PipelineInfo{Name: name, Kind: "software", Steps: len(cfg.Pipelines[name])}
			}
			infos = append(infos, info)
		}
		sort.Slice(infos, func(i, j int) bool { return infos[i].Name < infos[j].Name })
		if pipelinesJSON {
			return printJSON(infos)
		}
		for _, info := range infos {
			def := ""
			if info.Name == cfg.PipelineMode {
				def = " (default)"
			}
			src := "config"
			if info.Pack != "" {
				src = "pack " + info.Pack
			}
			fmt.Printf("%-16s %d steps · %s · %s%s\n", info.Name, info.Steps, info.Kind, src, def)
			if info.Description != "" {
				fmt.Printf("                 %s\n", info.Description)
			}
		}
		return nil
	},
}

func init() {
	packsCmd.Flags().BoolVar(&packsJSON, "json", false, "Print as JSON")
	pipelinesCmd.Flags().BoolVar(&pipelinesJSON, "json", false, "Print as JSON")
	rootCmd.AddCommand(packsCmd, pipelinesCmd)
}
