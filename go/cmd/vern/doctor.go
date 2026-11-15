package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/spf13/cobra"
)

var doctorJSON bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the installed LLM CLIs, sign-ins and flags",
	Long: `Check every LLM CLI vern can drive, without spending tokens: whether it
is installed, its version, whether it is signed in (claude and codex can
say), and whether its help still lists every flag vern passes. A missing
flag usually means the CLI changed and vern needs an update.

Exit codes:
  0  claude (the fallback for every mode) is ready
  1  claude is missing, signed out or changed`,
	Args: cobra.NoArgs,
	RunE: runDoctor,
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Print the report as JSON")
	rootCmd.AddCommand(doctorCmd)
}

type doctorReport struct {
	Version     string       `json:"version"`
	Config      string       `json:"config"`
	AgentsDir   string       `json:"agents_dir"`
	AllowWrites bool         `json:"allow_writes"`
	LLMMode     string       `json:"llm_mode"`
	Providers   []llm.Status `json:"providers"`
	OK          bool         `json:"ok"`
}

func runDoctor(cmd *cobra.Command, args []string) error {
	agentsDir := resolveAgentsDir()
	cfg := config.Load(filepath.Dir(agentsDir))

	rep := doctorReport{
		Version:     version,
		Config:      cfg.SourcePath,
		AgentsDir:   agentsDir,
		AllowWrites: llm.WritesAllowed(),
		LLMMode:     cfg.LLMMode,
	}
	if rep.Config == "" {
		rep.Config = "built-in defaults"
	}
	if info, err := os.Stat(agentsDir); err != nil || !info.IsDir() {
		rep.AgentsDir = "built-in personas"
	}
	for _, p := range llm.Providers() {
		st := p.Check()
		rep.Providers = append(rep.Providers, st)
		if p.Name == "claude" {
			rep.OK = st.OK()
		}
	}

	if doctorJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(rep); err != nil {
			return err
		}
	} else {
		printDoctor(rep, cfg)
	}
	if !rep.OK {
		os.Exit(1)
	}
	return nil
}

func printDoctor(rep doctorReport, cfg *config.Config) {
	fmt.Printf("vern %s\n", rep.Version)
	fmt.Printf("  config   %s\n", rep.Config)
	fmt.Printf("  personas %s\n", rep.AgentsDir)
	fmt.Printf("  llm mode %s\n", rep.LLMMode)
	if rep.AllowWrites {
		fmt.Println("  runs     writes allowed (allow_writes, --allow-writes or VERN_ALLOW_WRITES)")
	} else {
		fmt.Println("  runs     read-only")
	}
	fmt.Println()

	for _, st := range rep.Providers {
		mark := "ok"
		switch {
		case !st.Installed:
			mark = "--"
		case !st.OK():
			mark = "!!"
		}
		line := fmt.Sprintf("[%s] %-12s", mark, st.Name)
		if !st.Installed {
			fmt.Printf("%s not installed (%s): %s\n", line, st.Binary, st.Install)
			continue
		}
		parts := []string{st.Path}
		if st.Version != "" {
			parts = append(parts, st.Version)
		}
		switch {
		case st.SignedIn == nil:
			parts = append(parts, "sign-in not checkable")
		case *st.SignedIn:
			parts = append(parts, "signed in")
		default:
			parts = append(parts, "NOT signed in")
		}
		fmt.Printf("%s %s\n", line, strings.Join(parts, " · "))
		if !st.ReadOnly && !rep.AllowWrites {
			fmt.Printf("     %s has no read-only mode; it runs with write access\n", st.Name)
		}
		if len(st.MissingFlags) > 0 {
			fmt.Printf("     help no longer lists %s: this CLI changed, vern may need an update\n", strings.Join(st.MissingFlags, ", "))
		}
		if st.Error != "" {
			fmt.Printf("     %s\n", st.Error)
		}
	}

	fmt.Println()
	if rep.OK {
		fmt.Println("Ready. Steps on a missing CLI fall back as the llm mode says.")
	} else {
		fmt.Println("claude is not ready, and it is the fallback for every mode.")
	}
}
