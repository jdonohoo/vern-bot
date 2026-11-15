package main

import (
	"fmt"

	"github.com/jdonohoo/vern-bot/go/internal/report"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report <discovery-dir>",
	Short: "Render a discovery run as one shareable HTML page",
	Long: `Render a discovery run (pipeline steps, tasks, VernHole and Oracle) as a
single self-contained HTML file, saved as <discovery-dir>/output/report.html.
Runs from before run.json existed still render from their files.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := report.Write(args[0])
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	},
}

func init() { rootCmd.AddCommand(reportCmd) }
