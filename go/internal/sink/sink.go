// Package sink sends a discovery run's tasks somewhere after the run.
//
// Sinks:
//
//	vts             the VTS files in output/vts/ (always written; the default)
//	json[:<path>]   output/tasks.json, copied to <path> when given
//	beads           creates Beads issues with the br CLI (like vern tobeads --apply)
//	exec:<command>  runs <command> through the shell with tasks.json on stdin
//
// An exec sink also gets VERN_TASKS_FILE, VERN_RUN_DIR and VERN_MANIFEST_FILE
// in its environment, runs in the current directory, and its output is
// saved to output/sink.log. This is how tools vern knows nothing about
// (an issue tracker, a studio's work board) receive the tasks.
package sink

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
	"github.com/jdonohoo/vern-bot/go/internal/tasks"
)

// LogFile is where an exec sink's output is saved, inside output/.
const LogFile = "sink.log"

// Options configures one sink run.
type Options struct {
	Spec         string // vts, json, json:<path>, beads, exec:<command>
	DiscoveryDir string
	Timeout      time.Duration // exec only; default 10 minutes
	// Beads imports output/vts into Beads. The CLI provides it so this
	// package does not depend on the br wiring.
	Beads func(vtsDir string) error
	// Stderr receives an exec sink's live output (default os.Stderr).
	Stderr io.Writer
}

// Names lists the sink kinds for help text.
var Names = []string{"vts", "json", "beads", "exec:<command>"}

// Validate checks a sink spec without running it.
func Validate(spec string) error {
	kind, arg := split(spec)
	switch kind {
	case "", "vts", "beads":
		return nil
	case "json":
		return nil
	case "exec":
		if strings.TrimSpace(arg) == "" {
			return fmt.Errorf("exec sink needs a command: exec:<command>")
		}
		return nil
	}
	return fmt.Errorf("unknown sink %q (valid: %s)", spec, strings.Join(Names, ", "))
}

func split(spec string) (kind, arg string) {
	spec = strings.TrimSpace(spec)
	if i := strings.Index(spec, ":"); i >= 0 {
		return strings.ToLower(spec[:i]), spec[i+1:]
	}
	return strings.ToLower(spec), ""
}

// Run sends the run's tasks to the sink and returns what happened.
func Run(o Options) pipeline.SinkInfo {
	kind, arg := split(o.Spec)
	if kind == "" {
		kind = "vts"
	}
	info := pipeline.SinkInfo{Sink: o.Spec, Status: "ok"}
	if info.Sink == "" {
		info.Sink = "vts"
	}
	if err := Validate(o.Spec); err != nil {
		return fail(info, err)
	}

	outputDir := filepath.Join(o.DiscoveryDir, "output")
	tasksFile := filepath.Join(outputDir, tasks.FileName)
	doc, err := tasks.Load(outputDir)
	if err != nil || len(doc.Tasks) == 0 {
		info.Status = "skipped"
		info.Detail = "the run produced no tasks"
		return info
	}
	// Runs from before tasks.json existed: write it from the VTS files.
	if _, statErr := os.Stat(tasksFile); statErr != nil {
		if werr := doc.Write(tasksFile); werr != nil {
			return fail(info, werr)
		}
	}

	switch kind {
	case "vts":
		info.Detail = fmt.Sprintf("%d VTS files in output/vts", len(doc.Tasks))
	case "json":
		info.Detail = "output/" + tasks.FileName
		if arg != "" {
			data, rerr := os.ReadFile(tasksFile)
			if rerr != nil {
				return fail(info, rerr)
			}
			if mkErr := os.MkdirAll(filepath.Dir(arg), 0755); mkErr != nil {
				return fail(info, mkErr)
			}
			if werr := os.WriteFile(arg, data, 0644); werr != nil {
				return fail(info, werr)
			}
			info.Detail = arg
		}
	case "beads":
		if o.Beads == nil {
			return fail(info, fmt.Errorf("beads sink is not available here"))
		}
		if berr := o.Beads(filepath.Join(outputDir, "vts")); berr != nil {
			return fail(info, berr)
		}
		info.Detail = fmt.Sprintf("%d tasks sent to Beads", len(doc.Tasks))
	case "exec":
		return runExec(o, info, arg, tasksFile)
	}
	return info
}

func runExec(o Options, info pipeline.SinkInfo, command, tasksFile string) pipeline.SinkInfo {
	timeout := o.Timeout
	if timeout == 0 {
		timeout = 10 * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		return fail(info, err)
	}
	runDir, _ := filepath.Abs(o.DiscoveryDir)
	absTasks, _ := filepath.Abs(tasksFile)
	cmd.Stdin = bytes.NewReader(data)
	cmd.Env = append(os.Environ(),
		"VERN_TASKS_FILE="+absTasks,
		"VERN_RUN_DIR="+runDir,
		"VERN_MANIFEST_FILE="+filepath.Join(runDir, "output", pipeline.ManifestFile),
	)
	live := o.Stderr
	if live == nil {
		live = os.Stderr
	}
	var log bytes.Buffer
	cmd.Stdout = io.MultiWriter(&log, live)
	cmd.Stderr = io.MultiWriter(&log, live)

	runErr := cmd.Run()
	logPath := filepath.Join(filepath.Dir(tasksFile), LogFile)
	os.WriteFile(logPath, log.Bytes(), 0644)
	info.Log = "output/" + LogFile

	if ctx.Err() == context.DeadlineExceeded {
		info.ExitCode = 124
		return fail(info, fmt.Errorf("exec sink timed out after %s", timeout))
	}
	if runErr != nil {
		if ee, ok := runErr.(*exec.ExitError); ok {
			info.ExitCode = ee.ExitCode()
		} else {
			info.ExitCode = 1
		}
		return fail(info, fmt.Errorf("exec sink failed: %v", runErr))
	}
	info.Detail = fmt.Sprintf("sent tasks to: %s", command)
	return info
}

func fail(info pipeline.SinkInfo, err error) pipeline.SinkInfo {
	info.Status = "failed"
	info.Detail = err.Error()
	if info.ExitCode == 0 {
		info.ExitCode = 1
	}
	return info
}
