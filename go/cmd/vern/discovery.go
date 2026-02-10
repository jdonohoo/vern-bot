package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
	"github.com/spf13/cobra"
)

var discoveryCmd = &cobra.Command{
	Use:   "discovery <idea> [discovery_dir]",
	Short: "Run the full Vern Discovery Pipeline",
	Long: `Run the multi-LLM discovery pipeline on an idea.

Flags:
  --batch              Non-interactive mode (for skill/plugin use)
  --skip-input         Don't read input/ files as context
  --vernhole N         Run VernHole with N verns after pipeline
  --vernhole-council   Use a named council tier
  --oracle             Run Oracle Vern after VernHole
  --oracle-apply       Auto-apply Oracle's vision via Architect Vern
  --expanded           Use expanded pipeline
  --extra-context FILE Add extra context file (repeatable)
  --resume-from N      Resume pipeline from step N
  --max-retries N      Max retry attempts per step
  --llm-mode MODE      LLM fallback mode (mixed_claude_fallback, mixed_codex_fallback, etc.)
  --single-llm LLM     Use a single LLM for all steps`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runDiscovery,
}

var (
	discBatch         bool
	discSkipInput     bool
	discVernhole      int
	discCouncil       string
	discOracle        bool
	discOracleApply   bool
	discExpanded      bool
	discExtraContext   []string
	discResumeFrom    int
	discMaxRetries    int
	discLLMMode       string
	discSingleLLM     string
)

func init() {
	discoveryCmd.Flags().BoolVar(&discBatch, "batch", false, "Non-interactive mode")
	discoveryCmd.Flags().BoolVar(&discSkipInput, "skip-input", false, "Don't read input/ files as context")
	discoveryCmd.Flags().IntVar(&discVernhole, "vernhole", 0, "Run VernHole with N verns after pipeline")
	discoveryCmd.Flags().StringVar(&discCouncil, "vernhole-council", "", "Use a named VernHole council tier")
	discoveryCmd.Flags().BoolVar(&discOracle, "oracle", false, "Run Oracle Vern after VernHole")
	discoveryCmd.Flags().BoolVar(&discOracleApply, "oracle-apply", false, "Auto-apply Oracle's vision")
	discoveryCmd.Flags().BoolVar(&discExpanded, "expanded", false, "Use expanded pipeline")
	discoveryCmd.Flags().StringArrayVar(&discExtraContext, "extra-context", nil, "Extra context files (repeatable)")
	discoveryCmd.Flags().IntVar(&discResumeFrom, "resume-from", 0, "Resume pipeline from step N")
	discoveryCmd.Flags().IntVar(&discMaxRetries, "max-retries", 0, "Max retry attempts per step")
	discoveryCmd.Flags().StringVar(&discLLMMode, "llm-mode", "", "LLM fallback mode (mixed_claude_fallback, mixed_codex_fallback, mixed_gemini_fallback, mixed_copilot_fallback, single_llm)")
	discoveryCmd.Flags().StringVar(&discSingleLLM, "single-llm", "", "Use a single LLM for all steps (shorthand for --llm-mode single_llm)")
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
		Expanded:          discExpanded,
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

	return pipeline.Run(opts)
}
