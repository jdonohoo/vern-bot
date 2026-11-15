package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
	"github.com/jdonohoo/vern-bot/go/internal/report"
	"github.com/jdonohoo/vern-bot/go/internal/sink"
	"github.com/spf13/cobra"
)

var discoveryCmd = &cobra.Command{
	Use:   "discovery <idea> [discovery_dir]",
	Short: "Run the full Vern Discovery Pipeline",
	Long: `Run the multi-LLM discovery pipeline on an idea.

Flags:
  --batch              Non-interactive mode (for skill/plugin use)
  --skip-input         Don't read input/ files as context
  --skip-historian     Skip the Historian pre-step
  --vernhole N         Run VernHole with N verns after pipeline
  --vernhole-council   Use a named council tier
  --oracle             Run Oracle Vern after VernHole
  --oracle-apply       Auto-apply Oracle's vision via Architect Vern
  --expanded           Use expanded pipeline
  --pipeline NAME      Use a named pipeline from your config or a pack (see vern pipelines)
  --extra-context FILE Add extra context file (repeatable)
  --resume-from N      Resume pipeline from step N
  --max-retries N      Max retry attempts per step
  --llm-mode MODE      LLM fallback mode (mixed_claude_fallback, mixed_codex_fallback, etc.)
  --single-llm LLM     Use a single LLM for all steps
  --sink SINK          Send the tasks on after the run: vts (default), json[:<path>], beads, exec:<command> (see vern sink)
  --json               Print the run manifest (output/run.json) on stdout; progress goes to stderr

Every run writes output/run.json (steps, LLMs, models, timings, usage),
output/tasks.json (the tasks as JSON, schema v1: schemas/tasks.v1.json) and
output/report.html (the whole run as one page; rebuild with vern report).

Exit codes with --json:
  0  every step succeeded and tasks were produced
  1  the run could not start, or every step failed
  3  partial: some steps failed, VernHole failed, no tasks were produced, or the sink failed`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runDiscovery,
}

var (
	discBatch         bool
	discSkipInput     bool
	discSkipHistorian bool
	discVernhole      int
	discCouncil       string
	discOracle        bool
	discOracleApply   bool
	discExpanded      bool
	discExtraContext  []string
	discResumeFrom    int
	discMaxRetries    int
	discLLMMode       string
	discSingleLLM     string
	discJSON          bool
	discPipeline      string
	discSink          string
)

func init() {
	discoveryCmd.Flags().BoolVar(&discBatch, "batch", false, "Non-interactive mode")
	discoveryCmd.Flags().BoolVar(&discSkipInput, "skip-input", false, "Don't read input/ files as context")
	discoveryCmd.Flags().BoolVar(&discSkipHistorian, "skip-historian", false, "Skip the Historian pre-step (indexing input files)")
	discoveryCmd.Flags().IntVar(&discVernhole, "vernhole", 0, "Run VernHole with N verns after pipeline")
	discoveryCmd.Flags().StringVar(&discCouncil, "vernhole-council", "", "Use a named VernHole council tier")
	discoveryCmd.Flags().BoolVar(&discOracle, "oracle", false, "Run Oracle Vern after VernHole")
	discoveryCmd.Flags().BoolVar(&discOracleApply, "oracle-apply", false, "Auto-apply Oracle's vision")
	discoveryCmd.Flags().BoolVar(&discExpanded, "expanded", false, "Use expanded pipeline")
	discoveryCmd.Flags().StringArrayVar(&discExtraContext, "extra-context", nil, "Extra context files (repeatable)")
	discoveryCmd.Flags().IntVar(&discResumeFrom, "resume-from", 0, "Resume pipeline from step N")
	discoveryCmd.Flags().IntVar(&discMaxRetries, "max-retries", 0, "Max retry attempts per step")
	discoveryCmd.Flags().StringVar(&discLLMMode, "llm-mode", "", "LLM fallback mode (mixed_claude_fallback, mixed_codex_fallback, mixed_antigravity_fallback, mixed_grok_fallback, single_llm)")
	discoveryCmd.Flags().StringVar(&discSingleLLM, "single-llm", "", "Use a single LLM for all steps (shorthand for --llm-mode single_llm)")
	discoveryCmd.Flags().BoolVar(&discJSON, "json", false, "Print the run manifest as JSON on stdout (progress on stderr)")
	discoveryCmd.Flags().StringVar(&discPipeline, "pipeline", "", "Named pipeline from config or a pack (see vern pipelines)")
	discoveryCmd.Flags().StringVar(&discSink, "sink", "", "Where tasks go after the run: vts (default), json, beads, exec:<command> (default from .vern/config.json)")
	rootCmd.AddCommand(discoveryCmd)
}

func runDiscovery(cmd *cobra.Command, args []string) error {
	idea := args[0]

	// Discovery dir: arg or derived from idea
	discoveryDir := ""
	if len(args) > 1 {
		discoveryDir = args[1]
	}
	if discoveryDir == "" {
		slug := strings.ToLower(idea)
		slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
		slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
		slug = strings.Trim(slug, "-")
		if len(slug) > 50 {
			slug = slug[:50]
		}
		discoveryDir = "./discovery/" + slug
	}

	// --oracle-apply implies --oracle
	if discOracleApply {
		discOracle = true
	}

	// Resolve timeout
	timeout := 0
	if envTimeout := os.Getenv("VERN_TIMEOUT"); envTimeout != "" {
		fmt.Sscanf(envTimeout, "%d", &timeout)
	}

	agentsDir := resolveAgentsDir()

	// Find project root (parent of agents dir)
	projectRoot := ""
	if agentsDir != "agents" {
		// agentsDir is absolute
		projectRoot = agentsDir[:len(agentsDir)-len("/agents")]
	}

	opts := pipeline.Options{
		Idea:              idea,
		DiscoveryDir:      discoveryDir,
		BatchMode:         discBatch,
		ReadInput:         !discSkipInput,
		SkipHistorian:     discSkipHistorian,
		Expanded:          discExpanded,
		Pipeline:          discPipeline,
		ResumeFrom:        discResumeFrom,
		MaxRetries:        discMaxRetries,
		VernHoleCount:     discVernhole,
		VernHoleCouncil:   discCouncil,
		OracleFlag:        discOracle,
		OracleApplyFlag:   discOracleApply,
		ExtraContextFiles: discExtraContext,
		AgentsDir:         agentsDir,
		ProjectRoot:       projectRoot,
		Timeout:           timeout,
		LLMMode:           discLLMMode,
		SingleLLM:         discSingleLLM,
	}

	spec := resolveSinkSpec(discSink)
	if err := sink.Validate(spec); err != nil {
		return err
	}

	if !discJSON {
		err := pipeline.Run(opts)
		if err != nil {
			return err
		}
		if spec != "vts" {
			info := runSink(discoveryDir, spec)
			fmt.Printf("Sink %s: %s (%s)\n", info.Sink, info.Status, info.Detail)
		}
		if path, rerr := report.Write(discoveryDir); rerr == nil {
			fmt.Printf("Report: %s\n", path)
		}
		return nil
	}

	// --json: progress to stderr, the manifest alone on stdout.
	opts.BatchMode = true
	opts.OnLog = func(line string) { fmt.Fprintln(os.Stderr, line) }
	m, err := pipeline.RunWithManifest(opts)
	if err != nil {
		return err
	}
	if spec != "vts" && m.Tasks != nil {
		info := runSink(discoveryDir, spec)
		m.Sink = &info
		if info.Status == "failed" && m.Status == "ok" {
			m.Status = "partial"
		}
	}
	report.Write(discoveryDir)
	if err := printJSON(m); err != nil {
		return err
	}
	switch m.Status {
	case "ok":
		return nil
	case "failed":
		os.Exit(1)
	default:
		os.Exit(3)
	}
	return nil
}
