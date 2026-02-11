package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/jdonohoo/vern-bot/go/internal/vts"
)

// OracleConsultOptions configures a standalone Oracle vision run.
type OracleConsultOptions struct {
	Ctx          context.Context
	Idea         string // Original idea/prompt
	SynthesisDir string // Dir containing synthesis.md
	VTSDir       string // Dir containing vts-*.md files (optional)
	OutputFile   string // Where to write oracle-vision.md
	AgentsDir    string
	SynthesisLLM string
	Timeout      int
	OnLog        func(string)
}

// OracleApplyOptions configures a standalone Oracle apply run.
type OracleApplyOptions struct {
	Ctx          context.Context
	VisionFile   string // Path to oracle-vision.md
	VTSDir       string // Dir containing vts-*.md files to rewrite
	OutputFile   string // Where to write oracle-architect-breakdown.md (optional, defaults to VTSDir/../oracle-architect-breakdown.md)
	AgentsDir    string
	SynthesisLLM string
	Timeout      int
	OnLog        func(string)
}

// oracleLog prints to stdout in CLI mode, or routes to OnLog callback.
func oracleLog(onLog func(string), format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if onLog != nil {
		trimmed := strings.TrimSpace(msg)
		if trimmed != "" {
			onLog(trimmed)
		}
	} else {
		fmt.Print(msg)
	}
}

// RunOracleConsult generates an Oracle vision from VernHole synthesis + VTS tasks.
func RunOracleConsult(opts OracleConsultOptions) error {
	synthesisFile := filepath.Join(opts.SynthesisDir, "synthesis.md")
	data, err := os.ReadFile(synthesisFile)
	if err != nil {
		return fmt.Errorf("no synthesis file found at %s: %w", synthesisFile, err)
	}

	oracleLog(opts.OnLog, "=== CONSULTING THE ORACLE ===\n")
	oracleLog(opts.OnLog, "Reading the VernHole synthesis and VTS tasks...\n")

	// Build VTS task contents
	var vtsIndex, vtsContents strings.Builder
	if opts.VTSDir != "" {
		entries, _ := os.ReadDir(opts.VTSDir)
		for _, e := range entries {
			if !strings.HasPrefix(e.Name(), "vts-") || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			vtsIndex.WriteString(fmt.Sprintf("\n- %s", e.Name()))
			vtsData, _ := os.ReadFile(filepath.Join(opts.VTSDir, e.Name()))
			vtsContents.WriteString(fmt.Sprintf("\n\n=== %s ===\n%s", e.Name(), string(vtsData)))
		}
	}

	oraclePrompt := fmt.Sprintf(`You are Oracle Vern. The ancient seer who reads the patterns in the Vern council's chaos.

Review these VTS tasks in light of the Vern council's synthesis. Recommend modifications: new tasks to add, tasks to modify, tasks to remove, dependency changes, complexity reassessments, and missing acceptance criteria.

Output as a structured vision document with these sections:
# Oracle Vision

## Summary
Brief overview of recommended changes.

## New Tasks
(Use ### TASK N+1: Title format — same format as the architect breakdown so it can be parsed by the VTS post-processor)

## Modified Tasks
### VTS-NNN: New Title (was: Old Title)
**Changes:** What changed and why
**Description:** ...
**Acceptance Criteria:** ...
**Complexity:** ...
**Dependencies:** ...

## Removed Tasks
- VTS-NNN: Reason for removal

## Dependency Changes
- Describe any dependency modifications

## Risk Assessment
Remaining risks after your recommended modifications.

ORIGINAL IDEA:
%s

VTS TASK INDEX:
%s

VTS TASK FILES:
%s

VERNHOLE SYNTHESIS:
%s`, opts.Idea, vtsIndex.String(), vtsContents.String(), string(data))

	synthesisLLM := opts.SynthesisLLM
	if synthesisLLM == "" {
		synthesisLLM = "claude"
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 1200
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	outputFile := opts.OutputFile
	if outputFile == "" {
		outputFile = filepath.Join(opts.SynthesisDir, "oracle-vision.md")
	}

	// Ensure output directory exists
	os.MkdirAll(filepath.Dir(outputFile), 0755)

	result, err := llm.Run(llm.RunOptions{
		Ctx:        ctx,
		LLM:        synthesisLLM,
		Prompt:     oraclePrompt,
		OutputFile: outputFile,
		Persona:    "oracle",
		Timeout:    time.Duration(timeout) * time.Second,
		AgentsDir:  opts.AgentsDir,
	})

	if err != nil || result.ExitCode != 0 {
		return fmt.Errorf("oracle consult failed (exit %d)", result.ExitCode)
	}

	oracleLog(opts.OnLog, "\nOracle vision written to: %s\n", outputFile)
	return nil
}

// RunOracleApply has Architect Vern rewrite VTS tasks based on Oracle's vision.
func RunOracleApply(opts OracleApplyOptions) error {
	oracleData, err := os.ReadFile(opts.VisionFile)
	if err != nil {
		return fmt.Errorf("no oracle vision file found at %s: %w", opts.VisionFile, err)
	}

	oracleLog(opts.OnLog, "=== ORACLE APPLYING VISION ===\n")
	oracleLog(opts.OnLog, "Architect Vern is rewriting VTS tasks...\n")

	// Build existing VTS contents
	var vtsContents strings.Builder
	entries, _ := os.ReadDir(opts.VTSDir)
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "vts-") || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, _ := os.ReadFile(filepath.Join(opts.VTSDir, e.Name()))
		vtsContents.WriteString(fmt.Sprintf("\n\n=== %s ===\n%s", e.Name(), string(data)))
	}

	architectPrompt := fmt.Sprintf(`You are Architect Vern. The Oracle has spoken. Apply the Oracle's vision to produce an updated task breakdown.

Read the Oracle's vision and the existing VTS tasks. Produce a complete, updated task breakdown incorporating the Oracle's recommendations (new tasks, modified tasks, removed tasks, dependency changes).

Format each task with an h3 header exactly like this: ### TASK 1: Title Here. Include for each task: **Description:** what needs to be done, **Acceptance Criteria:** bullet list, **Complexity:** S|M|L|XL, **Dependencies:** Task N references or None, **Files:** list of files likely touched.

ORACLE'S VISION:
%s

EXISTING VTS TASKS:
%s

Produce the complete updated task breakdown. Include ALL tasks (not just changed ones).`, string(oracleData), vtsContents.String())

	synthesisLLM := opts.SynthesisLLM
	if synthesisLLM == "" {
		synthesisLLM = "claude"
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 1200
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	outputFile := opts.OutputFile
	if outputFile == "" {
		outputFile = filepath.Join(filepath.Dir(opts.VTSDir), "oracle-architect-breakdown.md")
	}

	// Ensure output directory exists
	os.MkdirAll(filepath.Dir(outputFile), 0755)

	result, err := llm.Run(llm.RunOptions{
		Ctx:        ctx,
		LLM:        synthesisLLM,
		Prompt:     architectPrompt,
		OutputFile: outputFile,
		Persona:    "architect",
		Timeout:    time.Duration(timeout) * time.Second,
		AgentsDir:  opts.AgentsDir,
	})

	if err != nil || result.ExitCode != 0 || IsFailedOutput(outputFile) {
		return fmt.Errorf("oracle apply failed")
	}

	// Clear old VTS files and re-process
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "vts-") && strings.HasSuffix(e.Name(), ".md") {
			os.Remove(filepath.Join(opts.VTSDir, e.Name()))
		}
	}

	oracleLog(opts.OnLog, "\n>>> Re-splitting updated architect breakdown into VTS task files...\n")

	// Process VTS from updated breakdown
	archData, err := os.ReadFile(outputFile)
	if err != nil {
		return fmt.Errorf("read updated breakdown: %w", err)
	}

	tasks, header, footer := vts.ParseArchitectOutput(string(archData))
	if len(tasks) > 0 {
		if err := vts.WriteVTSFiles(tasks, opts.VTSDir, "oracle", filepath.Base(outputFile), opts.OnLog); err != nil {
			oracleLog(opts.OnLog, "  Error writing VTS files: %v\n", err)
		}
		if err := vts.WriteSummary(tasks, outputFile, header, footer, "", opts.OnLog); err != nil {
			oracleLog(opts.OnLog, "  Error writing summary: %v\n", err)
		}
	}

	oracleLog(opts.OnLog, "\nOracle's vision applied. Updated VTS files in: %s\n", opts.VTSDir)
	return nil
}
