package main

import (
	"fmt"
	"os"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
	"github.com/spf13/cobra"
)

var oracleCmd = &cobra.Command{
	Use:   "oracle",
	Short: "Post-hoc Oracle operations (consult and apply)",
	Long: `Oracle: consult the Oracle for a vision on VernHole output, or apply an existing vision to rewrite VTS tasks.

Subcommands:
  consult  Generate an Oracle vision from VernHole synthesis + VTS tasks
  apply    Have Architect Vern rewrite VTS tasks based on an Oracle vision`,
}

var oracleConsultCmd = &cobra.Command{
	Use:   `consult "<idea>"`,
	Short: "Generate an Oracle vision from VernHole synthesis + VTS tasks",
	Long: `Consult the Oracle to review VernHole synthesis and VTS tasks, and produce
a structured vision document with recommended modifications.

The synthesis directory should contain a synthesis.md file (from a VernHole run).
The VTS directory should contain vts-*.md task files (optional).`,
	Args: cobra.ExactArgs(1),
	RunE: runOracleConsult,
}

var oracleApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply an Oracle vision to rewrite VTS tasks",
	Long: `Have Architect Vern rewrite VTS task files based on an Oracle vision document.

The vision file is the oracle-vision.md produced by 'vern oracle consult'.
The VTS directory should contain the vts-*.md task files to be rewritten.`,
	RunE: runOracleApply,
}

var (
	oracleSynthesisDir string
	oracleVTSDir       string
	oracleOutputFile   string
	oracleVisionFile   string
	oracleLLMMode      string
	oracleSingleLLM    string
)

func init() {
	// consult flags
	oracleConsultCmd.Flags().StringVar(&oracleSynthesisDir, "synthesis-dir", ".", "Directory containing synthesis.md")
	oracleConsultCmd.Flags().StringVar(&oracleVTSDir, "vts-dir", "", "Directory containing vts-*.md files (optional)")
	oracleConsultCmd.Flags().StringVarP(&oracleOutputFile, "output", "o", "", "Output file for oracle-vision.md (default: <synthesis-dir>/oracle-vision.md)")
	oracleConsultCmd.Flags().StringVar(&oracleLLMMode, "llm-mode", "", "LLM fallback mode")
	oracleConsultCmd.Flags().StringVar(&oracleSingleLLM, "single-llm", "", "Use a single LLM")

	// apply flags
	oracleApplyCmd.Flags().StringVar(&oracleVisionFile, "vision-file", "", "Path to oracle-vision.md (required)")
	oracleApplyCmd.Flags().StringVar(&oracleVTSDir, "vts-dir", "", "Directory containing vts-*.md files (required)")
	oracleApplyCmd.Flags().StringVarP(&oracleOutputFile, "output", "o", "", "Output file for architect breakdown (default: <vts-dir>/../oracle-architect-breakdown.md)")
	oracleApplyCmd.Flags().StringVar(&oracleLLMMode, "llm-mode", "", "LLM fallback mode")
	oracleApplyCmd.Flags().StringVar(&oracleSingleLLM, "single-llm", "", "Use a single LLM")
	oracleApplyCmd.MarkFlagRequired("vision-file")
	oracleApplyCmd.MarkFlagRequired("vts-dir")

	oracleCmd.AddCommand(oracleConsultCmd)
	oracleCmd.AddCommand(oracleApplyCmd)
	rootCmd.AddCommand(oracleCmd)
}

func resolveOracleLLM() string {
	agentsDir := resolveAgentsDir()
	projectRoot := ""
	if agentsDir != "agents" {
		projectRoot = agentsDir[:len(agentsDir)-len("/agents")]
	}

	cfg := config.Load(projectRoot)

	if oracleSingleLLM != "" {
		return oracleSingleLLM
	}
	if oracleLLMMode != "" {
		cfg.LLMMode = oracleLLMMode
	}
	return cfg.GetSynthesisLLM()
}

func resolveOracleConfig() *config.Config {
	agentsDir := resolveAgentsDir()
	projectRoot := ""
	if agentsDir != "agents" {
		projectRoot = agentsDir[:len(agentsDir)-len("/agents")]
	}
	return config.Load(projectRoot)
}

func runOracleConsult(cmd *cobra.Command, args []string) error {
	idea := args[0]
	agentsDir := resolveAgentsDir()
	synthesisLLM := resolveOracleLLM()

	cfg := resolveOracleConfig()
	timeout := cfg.GetOracleTimeout()
	if envTimeout := os.Getenv("VERN_TIMEOUT"); envTimeout != "" {
		var t int
		if _, err := fmt.Sscanf(envTimeout, "%d", &t); err == nil {
			timeout = t
		}
	}

	err := pipeline.RunOracleConsult(pipeline.OracleConsultOptions{
		Idea:         idea,
		SynthesisDir: oracleSynthesisDir,
		VTSDir:       oracleVTSDir,
		OutputFile:   oracleOutputFile,
		AgentsDir:    agentsDir,
		SynthesisLLM: synthesisLLM,
		Timeout:      timeout,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func runOracleApply(cmd *cobra.Command, args []string) error {
	agentsDir := resolveAgentsDir()
	synthesisLLM := resolveOracleLLM()

	cfg := resolveOracleConfig()
	timeout := cfg.GetOracleApplyTimeout()
	if envTimeout := os.Getenv("VERN_TIMEOUT"); envTimeout != "" {
		var t int
		if _, err := fmt.Sscanf(envTimeout, "%d", &t); err == nil {
			timeout = t
		}
	}

	err := pipeline.RunOracleApply(pipeline.OracleApplyOptions{
		VisionFile:   oracleVisionFile,
		VTSDir:       oracleVTSDir,
		OutputFile:   oracleOutputFile,
		AgentsDir:    agentsDir,
		SynthesisLLM: synthesisLLM,
		Timeout:      timeout,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return nil
}
