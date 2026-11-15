package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jdonohoo/vern-bot/go/internal/packs"
	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
	"github.com/jdonohoo/vern-bot/go/internal/sink"
	"github.com/jdonohoo/vern-bot/go/internal/tobeads"
	"github.com/spf13/cobra"
)

var (
	sinkSpec string
	sinkJSON bool
)

var sinkCmd = &cobra.Command{
	Use:   "sink <discovery-dir>",
	Short: "Send a finished run's tasks to a task sink",
	Long: `Send the tasks of a finished discovery run somewhere else. Discovery does
this itself with --sink; use this to resend, or to try another sink.

Sinks:
  vts             the VTS files in output/vts/ (always written)
  json[:<path>]   output/tasks.json, copied to <path> when given
  beads           create Beads issues with the br CLI (like vern tobeads --apply)
  exec:<command>  run <command> through the shell with tasks.json on stdin,
                  and VERN_TASKS_FILE, VERN_RUN_DIR, VERN_MANIFEST_FILE set;
                  its output is saved to output/sink.log

The default sink comes from "sink" in the project's .vern/config.json.

Exit codes: 0 ok or nothing to send, 1 the sink failed.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		info := runSink(args[0], resolveSinkSpec(sinkSpec))
		if sinkJSON {
			if err := printJSON(info); err != nil {
				return err
			}
		} else {
			fmt.Printf("Sink %s: %s", info.Sink, info.Status)
			if info.Detail != "" {
				fmt.Printf(" (%s)", info.Detail)
			}
			fmt.Println()
		}
		if info.Status == "failed" {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	sinkCmd.Flags().StringVar(&sinkSpec, "sink", "", "vts, json[:<path>], beads or exec:<command> (default: the project's sink, else vts)")
	sinkCmd.Flags().BoolVar(&sinkJSON, "json", false, "Print the result as JSON")
	rootCmd.AddCommand(sinkCmd)
}

// resolveSinkSpec returns the flag value, else the project's sink, else vts.
func resolveSinkSpec(flag string) string {
	if flag != "" {
		return flag
	}
	if p := packs.ActiveProject(); p != nil && p.Sink != "" {
		return p.Sink
	}
	return "vts"
}

// runSink sends a run's tasks and records the result in its run.json.
func runSink(discoveryDir, spec string) pipeline.SinkInfo {
	info := sink.Run(sink.Options{
		Spec:         spec,
		DiscoveryDir: discoveryDir,
		Beads:        beadsSink,
	})
	pipeline.UpdateManifest(discoveryDir, func(m *pipeline.Manifest) { m.Sink = &info })
	return info
}

// beadsSink imports a VTS folder into the nearest Beads repo.
func beadsSink(vtsDir string) error {
	abs, err := filepath.Abs(vtsDir)
	if err != nil {
		return err
	}
	beadsDir := findBeadsDir(abs)
	if beadsDir == "" {
		cwd, _ := os.Getwd()
		beadsDir = findBeadsDir(cwd)
	}
	if beadsDir == "" {
		return fmt.Errorf("no .beads/ folder found above %s or the current directory (run br init)", abs)
	}
	res, err := tobeads.Run(tobeads.ImportOptions{VTSDir: abs, BeadsDir: beadsDir, Apply: true},
		&tobeads.RealBrRunner{WorkDir: beadsDir})
	if err != nil {
		return err
	}
	if res != nil && res.Failed > 0 {
		return fmt.Errorf("%d of the tasks failed to import: %s", res.Failed, strings.Join(res.Errors, "; "))
	}
	return nil
}
