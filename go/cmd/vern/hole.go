package main

import (
	"fmt"
	"os"

	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
	"github.com/spf13/cobra"
)

var holeCmd = &cobra.Command{
	Use:   "hole <idea>",
	Short: "Summon random Vern personas for chaotic discovery",
	Long: `VernHole: summon a council of Vern personas to analyze your idea in parallel.

Council tiers:
  hammers  - Council of the Three Hammers (great, mediocre, ketamine)
  conflict - Max Conflict (startup, enterprise, yolo, paranoid, optimist, inverse)
  inner    - The Inner Circle (architect, inverse, paranoid + random fill)
  round    - The Round Table (mighty, yolo, startup, academic, enterprise + random fill)
  war      - The War Room (all round table + ux, retro, optimist, nyquil + random fill)
  full     - The Full Vern Experience (all summonable personas)
  random   - Fate's Hand (random count, random selection)`,
	Args: cobra.ExactArgs(1),
	RunE: runHole,
}

var (
	holeOutputDir string
	holeCouncil   string
	holeContext   string
	holeCount     int
)

func init() {
	holeCmd.Flags().StringVarP(&holeOutputDir, "output-dir", "d", ".", "Output directory for VernHole files")
	holeCmd.Flags().StringVar(&holeCouncil, "council", "", "Named council tier")
	holeCmd.Flags().StringVar(&holeContext, "context", "", "Path to context file (e.g. discovery master plan)")
	holeCmd.Flags().IntVarP(&holeCount, "count", "n", 0, "Number of Verns to summon (min 3)")
	rootCmd.AddCommand(holeCmd)
}

func runHole(cmd *cobra.Command, args []string) error {
	idea := args[0]
	agentsDir := resolveAgentsDir()

	// Resolve timeout
	timeout := 1200
	if envTimeout := os.Getenv("VERN_TIMEOUT"); envTimeout != "" {
		var t int
		if _, err := fmt.Sscanf(envTimeout, "%d", &t); err == nil {
			timeout = t
		}
	}

	err := pipeline.RunVernHole(pipeline.VernHoleOptions{
		Idea:      idea,
		OutputDir: holeOutputDir,
		Council:   holeCouncil,
		Count:     holeCount,
		Context:   holeContext,
		AgentsDir: agentsDir,
		Timeout:   timeout,
	})
	if err != nil {
		os.Exit(1)
	}
	return nil
}
