package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/llm"
	"github.com/jdonohoo/vern-bot/go/internal/vts"
)

// Options configures a pipeline run.
type Options struct {
	Idea              string
	DiscoveryDir      string
	BatchMode         bool
	ReadInput         bool
	Expanded          bool
	ResumeFrom        int
	MaxRetries        int
	VernHoleCount     int
	VernHoleCouncil   string
	OracleFlag        bool
	OracleApplyFlag   bool
	ExtraContextFiles []string
	AgentsDir         string
	ProjectRoot       string
	Timeout           int // seconds
}

// Pipeline orchestrates the discovery pipeline execution.
type Pipeline struct {
	opts          Options
	cfg           *config.Config
	steps         []config.PipelineStep
	results       []StepResult
	fullPrompt    string
	consolidation string // path to consolidation file
	logFile       *os.File
	statusPath    string    // path to pipeline-status.md
	startTime     time.Time // pipeline start time
	mode          string    // pipeline mode (default/expanded)
}

// Run executes the full discovery pipeline.
func Run(opts Options) error {
	cfg := config.Load(opts.ProjectRoot)

	// Override timeout from config if not set
	if opts.Timeout == 0 {
		opts.Timeout = cfg.TimeoutSeconds
	}
	if opts.MaxRetries <= 0 {
		if cfg.MaxRetries > 0 {
			opts.MaxRetries = cfg.MaxRetries
		} else {
			opts.MaxRetries = 1
		}
	}

	// Determine pipeline mode
	mode := cfg.PipelineMode
	if opts.Expanded {
		mode = "expanded"
	}

	steps := cfg.GetPipeline(mode)

	p := &Pipeline{
		opts:    opts,
		cfg:     cfg,
		steps:   steps,
		results: make([]StepResult, len(steps)),
	}

	return p.execute(mode)
}

func (p *Pipeline) execute(mode string) error {
	opts := p.opts

	// Create discovery directory structure
	inputDir := filepath.Join(opts.DiscoveryDir, "input")
	outputDir := filepath.Join(opts.DiscoveryDir, "output")
	vtsDir := filepath.Join(outputDir, "vts")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		return fmt.Errorf("create input dir: %w", err)
	}
	if err := os.MkdirAll(vtsDir, 0755); err != nil {
		return fmt.Errorf("create vts dir: %w", err)
	}

	// Initialize pipeline log
	logPath := filepath.Join(outputDir, "pipeline.log")
	var err error
	p.logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open pipeline log: %w", err)
	}
	defer p.logFile.Close()

	// Initialize status file and timing
	p.statusPath = filepath.Join(outputDir, "pipeline-status.md")
	p.startTime = time.Now()
	p.mode = mode

	// Banner
	fmt.Println("=== VERN DISCOVERY PIPELINE ===")
	fmt.Printf("Idea: %s\n", opts.Idea)
	fmt.Printf("Discovery folder: %s\n", opts.DiscoveryDir)
	fmt.Println("Output: Vern Task Spec (VTS)")
	fmt.Printf("Pipeline mode: %s (%d steps)\n", mode, len(p.steps))
	if opts.BatchMode {
		fmt.Println("Mode: batch (non-interactive)")
	}
	if opts.ResumeFrom > 0 {
		fmt.Printf("Resuming from: step %d\n", opts.ResumeFrom)
	}
	fmt.Printf("Max retries: %d\n", opts.MaxRetries)
	fmt.Printf("Timeout: %ds per step\n", opts.Timeout)
	fmt.Println()

	// Print pipeline steps
	fmt.Println("Pipeline:")
	for _, step := range p.steps {
		fmt.Printf("  %d. %s → %s\n", step.Step, step.Name, step.LLM)
	}
	fmt.Println()

	p.log("=== Pipeline started ===")
	p.log("Idea: %s", opts.Idea)
	p.log("Mode: %s (%d steps)", mode, len(p.steps))
	if opts.ResumeFrom > 0 {
		p.log("Resuming from step %d", opts.ResumeFrom)
	}

	// Save prompt
	promptFile := filepath.Join(inputDir, "prompt.md")
	if _, err := os.Stat(promptFile); os.IsNotExist(err) {
		os.WriteFile(promptFile, []byte("# Discovery Prompt\n\n"+opts.Idea+"\n"), 0644)
		fmt.Printf("Saved prompt to %s\n", promptFile)
	}

	// Build context from input files
	inputContext := p.buildInputContext(inputDir)

	// Extra context files
	for _, f := range opts.ExtraContextFiles {
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Printf("Warning: Extra context file not found: %s\n", f)
			continue
		}
		inputContext += fmt.Sprintf("\n\n=== %s ===\n%s\n", filepath.Base(f), string(data))
		fmt.Printf("Added extra context: %s\n", f)
	}

	// Build full prompt
	p.fullPrompt = opts.Idea
	if inputContext != "" {
		p.fullPrompt += "\n\n=== INPUT MATERIALS ===\n" + inputContext + "\n=== END INPUT MATERIALS ==="
	}

	// Set working dir for codex
	os.Setenv("VERN_WORKING_DIR", opts.DiscoveryDir)

	// Execute pipeline steps
	failedSteps := []int{}

	for idx, step := range p.steps {
		stepNum := step.Step
		outputFile := filepath.Join(outputDir, fmt.Sprintf("%02d-%s-%s.md",
			stepNum, Slugify(step.Persona), Slugify(step.Name)))

		// Resume logic
		if opts.ResumeFrom > 0 && stepNum < opts.ResumeFrom {
			if !IsFailedOutput(outputFile) {
				fmt.Printf("\n>>> Pass %d/%d: %s — SKIPPED (resuming, output exists)\n", stepNum, len(p.steps), step.Name)
				p.log("Step %d (%s): SKIPPED (resume, output exists)", stepNum, step.Name)
				p.results[idx] = StepResult{
					StepNum:    stepNum,
					Name:       step.Name,
					OutputFile: outputFile,
					Status:     "skipped",
				}
				// Track consolidation file
				if step.ContextMode == "all_previous" {
					p.consolidation = outputFile
				}
				continue
			}
			fmt.Printf("\n>>> Pass %d/%d: %s — re-running (no valid output for resume)\n", stepNum, len(p.steps), step.Name)
			p.log("Step %d (%s): re-running (no valid output for resume)", stepNum, step.Name)
		}

		fmt.Printf("\n>>> Pass %d/%d: %s (%s)\n", stepNum, len(p.steps), step.Name, step.LLM)

		// Build prompt based on context mode
		runPrompt := p.buildStepPrompt(step, idx)

		// Track consolidation file
		if step.ContextMode == "all_previous" {
			p.consolidation = outputFile
		}

		// Retry loop with claude fallback
		succeeded := false
		originalLLM := step.LLM
		retryLLM := step.LLM
		retryPersona := step.Persona
		fellBack := false
		totalAttempts := opts.MaxRetries + 1
		var lastExitCode int
		var attemptCount int
		var duration time.Duration

		for attempt := 1; attempt <= totalAttempts; attempt++ {
			if attempt > 1 {
				fmt.Printf("    Retry %d/%d for step %d (%s) with %s...\n", attempt-1, opts.MaxRetries, stepNum, step.Name, retryLLM)
				p.log("Step %d (%s): retry %d/%d with %s", stepNum, step.Name, attempt-1, opts.MaxRetries, retryLLM)
			}

			result, _ := llm.Run(llm.RunOptions{
				LLM:        retryLLM,
				Prompt:     runPrompt,
				OutputFile:  outputFile,
				Persona:    retryPersona,
				Timeout:    time.Duration(opts.Timeout) * time.Second,
				AgentsDir:  opts.AgentsDir,
			})

			lastExitCode = result.ExitCode
			attemptCount = attempt
			duration = result.Duration

			if result.ExitCode == 0 && !IsFailedOutput(outputFile) {
				succeeded = true
				break
			}

			// On timeout with non-claude LLM, fall back to claude immediately (don't waste more timeout cycles)
			if result.ExitCode == llm.ExitTimeout && retryLLM != "claude" {
				fmt.Printf("    Timeout on %s — falling back to claude\n", retryLLM)
				p.log("Step %d (%s): timeout on %s after %s, falling back to claude", stepNum, step.Name, retryLLM, result.Duration.Truncate(time.Second))
				retryLLM = "claude"
				fellBack = true
			}
		}

		// Claude fallback: if all retries failed with a non-claude LLM, try claude as final safety net
		if !succeeded && retryLLM != "claude" {
			fmt.Printf("    %s failed after %d attempt(s) — falling back to claude\n", originalLLM, totalAttempts)
			p.log("Step %d (%s): %s FAILED after %d attempt(s) (last exit %d), falling back to claude",
				stepNum, step.Name, originalLLM, totalAttempts, lastExitCode)

			retryLLM = "claude"
			fellBack = true

			result, _ := llm.Run(llm.RunOptions{
				LLM:        "claude",
				Prompt:     runPrompt,
				OutputFile:  outputFile,
				Persona:    retryPersona,
				Timeout:    time.Duration(opts.Timeout) * time.Second,
				AgentsDir:  opts.AgentsDir,
			})

			attemptCount++
			lastExitCode = result.ExitCode
			duration = result.Duration

			if result.ExitCode == 0 && !IsFailedOutput(outputFile) {
				succeeded = true
				fmt.Printf("    Claude fallback succeeded\n")
				p.log("Step %d (%s): claude fallback OK after %s failed", stepNum, step.Name, originalLLM)
			} else {
				p.log("Step %d (%s): claude fallback also FAILED (exit %d)", stepNum, step.Name, result.ExitCode)
			}
		}

		if succeeded {
			outputBytes := fileSize(outputFile)
			if fellBack {
				fmt.Printf("    OK (%s→claude fallback, %d bytes, %s)\n", originalLLM, outputBytes, duration.Truncate(time.Second))
				p.log("Step %d (%s): OK via claude fallback (original=%s, attempt %d, %d bytes)", stepNum, step.Name, originalLLM, attemptCount, outputBytes)
			} else {
				fmt.Printf("    OK (%s, %d bytes, %s)\n", retryLLM, outputBytes, duration.Truncate(time.Second))
				p.log("Step %d (%s): OK (exit %d, attempt %d, llm=%s, %d bytes)", stepNum, step.Name, lastExitCode, attemptCount, retryLLM, outputBytes)
			}
			p.results[idx] = StepResult{
				StepNum:     stepNum,
				Name:        step.Name,
				OutputFile:  outputFile,
				Status:      "ok",
				ExitCode:    lastExitCode,
				Attempts:    attemptCount,
				LLMUsed:     retryLLM,
				OriginalLLM: originalLLM,
				FellBack:    fellBack,
				DurationMS:  duration.Milliseconds(),
				OutputBytes: outputBytes,
			}
		} else {
			fmt.Printf("    FAILED after %d attempts including claude fallback (last exit: %d)\n", attemptCount, lastExitCode)
			p.log("Step %d (%s): FAILED (exit %d, %d attempts, original=%s, final=%s)", stepNum, step.Name, lastExitCode, attemptCount, originalLLM, retryLLM)

			// Write failure marker
			failureContent := fmt.Sprintf("# STEP FAILED\n\nStep %d (%s) failed after %d attempt(s).\nOriginal LLM: %s\nFinal LLM: %s (fallback)\nLast exit code: %d\n\nRe-run with: --resume-from %d\n",
				stepNum, step.Name, attemptCount, originalLLM, retryLLM, lastExitCode, stepNum)
			os.WriteFile(outputFile, []byte(failureContent), 0644)

			failedSteps = append(failedSteps, stepNum)
			p.results[idx] = StepResult{
				StepNum:     stepNum,
				Name:        step.Name,
				OutputFile:  outputFile,
				Status:      "failed",
				ExitCode:    lastExitCode,
				Attempts:    attemptCount,
				LLMUsed:     retryLLM,
				OriginalLLM: originalLLM,
				FellBack:    fellBack,
				DurationMS:  duration.Milliseconds(),
				OutputBytes: 0,
			}
		}

		// Write structured JSON log entry + update status file
		p.logJSON(p.results[idx])
		p.writeStatus("running", failedSteps)
	}

	// Pipeline summary
	if len(failedSteps) > 0 {
		fmt.Printf("\nWARNING: %d step(s) failed: %v\n", len(failedSteps), failedSteps)
		fmt.Println("Re-run failed steps with: --resume-from <N>")
		p.log("Pipeline completed with %d failure(s): steps %v", len(failedSteps), failedSteps)
	} else {
		p.log("Pipeline completed successfully (%d/%d steps)", len(p.steps), len(p.steps))
	}
	p.writeStatus("pipeline_complete", failedSteps)

	// VTS post-processing
	lastResult := p.results[len(p.results)-1]
	if lastResult.Status == "failed" || IsFailedOutput(lastResult.OutputFile) {
		fmt.Println("\n>>> Skipping VTS post-processing (architect step failed)")
		fmt.Printf("    Re-run with: --resume-from %d\n", len(p.steps))
		p.log("VTS post-processing: SKIPPED (architect step failed)")
	} else {
		p.processVTS(lastResult.OutputFile, vtsDir, "discovery")
	}

	// Directory structure output
	fmt.Println()
	fmt.Println("=== DISCOVERY COMPLETE ===")
	fmt.Printf("Files created in: %s\n", opts.DiscoveryDir)
	p.printDirectoryStructure()

	// VernHole
	if len(failedSteps) > 0 {
		fmt.Printf("\n>>> Skipping VernHole (pipeline had %d failed step(s): %v)\n", len(failedSteps), failedSteps)
		p.log("VernHole: SKIPPED (%d pipeline step(s) failed)", len(failedSteps))
	} else if opts.VernHoleCount > 0 || opts.VernHoleCouncil != "" {
		p.runVernHole()
	}

	return nil
}

func (p *Pipeline) buildStepPrompt(step config.PipelineStep, idx int) string {
	switch step.ContextMode {
	case "prompt_only":
		return step.PromptPrefix + "\n\n" + p.fullPrompt

	case "previous":
		prevOutput := ""
		if idx > 0 {
			prevFile := p.results[idx-1].OutputFile
			if prevFile != "" && !IsFailedOutput(prevFile) {
				data, _ := os.ReadFile(prevFile)
				prevOutput = string(data)
			}
		}
		return step.PromptPrefix + "\n\nORIGINAL REQUEST:\n" + p.fullPrompt + "\n\nPREVIOUS ANALYSIS:\n" + prevOutput

	case "all_previous":
		var allContext strings.Builder
		for j := 0; j < idx; j++ {
			r := p.results[j]
			if r.OutputFile != "" && !IsFailedOutput(r.OutputFile) {
				data, _ := os.ReadFile(r.OutputFile)
				allContext.WriteString(fmt.Sprintf("\n\n%s: %s", p.steps[j].Name, string(data)))
			}
		}
		return step.PromptPrefix + "\n\nORIGINAL REQUEST:\n" + p.fullPrompt + allContext.String()

	case "consolidation":
		consolOutput := ""
		if p.consolidation != "" && !IsFailedOutput(p.consolidation) {
			data, _ := os.ReadFile(p.consolidation)
			consolOutput = string(data)
		}
		return step.PromptPrefix + "\n\nORIGINAL REQUEST:\n" + p.fullPrompt + "\n\nMaster plan: " + consolOutput

	default:
		// Fallback: treat like previous
		prevOutput := ""
		if idx > 0 {
			prevFile := p.results[idx-1].OutputFile
			if prevFile != "" && !IsFailedOutput(prevFile) {
				data, _ := os.ReadFile(prevFile)
				prevOutput = string(data)
			}
		}
		return step.PromptPrefix + "\n\nORIGINAL REQUEST:\n" + p.fullPrompt + "\n\nPREVIOUS ANALYSIS:\n" + prevOutput
	}
}

func (p *Pipeline) buildInputContext(inputDir string) string {
	if !p.opts.ReadInput {
		fmt.Println("Skipping input files (--skip-input).")
		return ""
	}

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return ""
	}

	var context strings.Builder
	count := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".md" && ext != ".txt" && ext != ".json" && ext != ".yaml" && ext != ".yml" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(inputDir, entry.Name()))
		if err != nil {
			continue
		}
		context.WriteString(fmt.Sprintf("\n\n=== %s ===\n%s\n", entry.Name(), string(data)))
		count++
	}

	if count > 0 {
		fmt.Printf("Loaded %d input files as context.\n", count)
	}

	return context.String()
}

func (p *Pipeline) processVTS(architectFile string, vtsDir string, source string) {
	data, err := os.ReadFile(architectFile)
	if err != nil {
		return
	}

	fmt.Println("\n>>> Splitting architect breakdown into VTS task files...")

	tasks, header, footer := vts.ParseArchitectOutput(string(data))
	if len(tasks) == 0 {
		fmt.Println("  No tasks found in architect breakdown, skipping split")
		return
	}

	if err := vts.WriteVTSFiles(tasks, vtsDir, source, filepath.Base(architectFile)); err != nil {
		fmt.Printf("  Error writing VTS files: %v\n", err)
		return
	}

	if err := vts.WriteSummary(tasks, architectFile, header, footer, ""); err != nil {
		fmt.Printf("  Error writing summary: %v\n", err)
	}
}

func (p *Pipeline) runVernHole() {
	opts := p.opts

	// Find consolidation file
	consolFile := p.consolidation
	if consolFile == "" {
		// Fallback: find by pattern
		outputDir := filepath.Join(opts.DiscoveryDir, "output")
		entries, _ := os.ReadDir(outputDir)
		for _, e := range entries {
			if strings.Contains(strings.ToLower(e.Name()), "consolidation") && strings.HasSuffix(e.Name(), ".md") {
				consolFile = filepath.Join(outputDir, e.Name())
				break
			}
		}
	}

	fmt.Println()
	fmt.Println("=== ENTERING THE VERNHOLE ===")
	fmt.Println("Feeding original idea + master plan into the chaos...")

	vernholeDir := filepath.Join(opts.DiscoveryDir, "vernhole")
	os.MkdirAll(vernholeDir, 0755)

	err := RunVernHole(VernHoleOptions{
		Idea:      opts.Idea,
		OutputDir: vernholeDir,
		Council:   opts.VernHoleCouncil,
		Count:     opts.VernHoleCount,
		Context:   consolFile,
		AgentsDir: opts.AgentsDir,
		Timeout:   opts.Timeout,
	})
	if err != nil {
		fmt.Printf("\nWARNING: VernHole failed: %v\n", err)
		p.log("VernHole: FAILED (%v)", err)
		p.writeStatus("complete_vernhole_failed", nil)
		return
	}

	p.log("VernHole: OK")

	// Oracle integration
	if opts.OracleFlag {
		p.runOracle(vernholeDir)
	}

	p.writeStatus("complete", nil)
}

func (p *Pipeline) runOracle(vernholeDir string) {
	opts := p.opts
	synthesisFile := filepath.Join(vernholeDir, "synthesis.md")

	data, err := os.ReadFile(synthesisFile)
	if err != nil {
		fmt.Println(">>> Skipping Oracle (no synthesis file)")
		return
	}

	fmt.Println()
	fmt.Println("=== CONSULTING THE ORACLE ===")
	fmt.Println("Reading the VernHole synthesis and VTS tasks...")

	// Build VTS task contents
	vtsDir := filepath.Join(opts.DiscoveryDir, "output", "vts")
	var vtsIndex, vtsContents strings.Builder

	entries, _ := os.ReadDir(vtsDir)
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "vts-") || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		vtsIndex.WriteString(fmt.Sprintf("\n- %s", e.Name()))
		vtsData, _ := os.ReadFile(filepath.Join(vtsDir, e.Name()))
		vtsContents.WriteString(fmt.Sprintf("\n\n=== %s ===\n%s", e.Name(), string(vtsData)))
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
%s`, p.fullPrompt, vtsIndex.String(), vtsContents.String(), string(data))

	oracleVisionFile := filepath.Join(opts.DiscoveryDir, "oracle-vision.md")
	result, err := llm.Run(llm.RunOptions{
		LLM:        "claude",
		Prompt:     oraclePrompt,
		OutputFile:  oracleVisionFile,
		Persona:    "oracle",
		Timeout:    time.Duration(opts.Timeout) * time.Second,
		AgentsDir:  opts.AgentsDir,
	})

	if err != nil || result.ExitCode != 0 {
		fmt.Println("\nWARNING: Oracle step failed")
		p.log("Oracle: FAILED")
		return
	}

	fmt.Printf("\nOracle vision written to: %s\n", oracleVisionFile)
	p.log("Oracle: OK")

	// Auto-apply: Architect Vern rewrites VTS based on Oracle's vision
	if opts.OracleApplyFlag {
		p.applyOracleVision(oracleVisionFile)
	}
}

func (p *Pipeline) applyOracleVision(oracleVisionFile string) {
	opts := p.opts

	fmt.Println()
	fmt.Println("=== APPLYING THE ORACLE'S VISION ===")
	fmt.Println("Architect Vern is rewriting VTS tasks...")

	oracleData, _ := os.ReadFile(oracleVisionFile)

	// Build existing VTS contents
	vtsDir := filepath.Join(opts.DiscoveryDir, "output", "vts")
	var vtsContents strings.Builder
	entries, _ := os.ReadDir(vtsDir)
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "vts-") || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		data, _ := os.ReadFile(filepath.Join(vtsDir, e.Name()))
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

	updatedFile := filepath.Join(opts.DiscoveryDir, "output", "oracle-architect-breakdown.md")
	result, err := llm.Run(llm.RunOptions{
		LLM:        "claude",
		Prompt:     architectPrompt,
		OutputFile:  updatedFile,
		Persona:    "architect",
		Timeout:    time.Duration(opts.Timeout) * time.Second,
		AgentsDir:  opts.AgentsDir,
	})

	if err != nil || result.ExitCode != 0 || IsFailedOutput(updatedFile) {
		fmt.Println("\nWARNING: Oracle apply step failed")
		p.log("Oracle apply: FAILED")
		return
	}

	// Clear old VTS files and re-process
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "vts-") && strings.HasSuffix(e.Name(), ".md") {
			os.Remove(filepath.Join(vtsDir, e.Name()))
		}
	}

	fmt.Println("\n>>> Re-splitting updated architect breakdown into VTS task files...")
	p.processVTS(updatedFile, vtsDir, "oracle")

	fmt.Printf("\nOracle's vision applied. Updated VTS files in: %s\n", vtsDir)
	p.log("Oracle apply: OK")
}

func (p *Pipeline) log(format string, args ...interface{}) {
	if p.logFile == nil {
		return
	}
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(p.logFile, "[%s] %s\n", timestamp, msg)
	p.logFile.Sync()
}

func (p *Pipeline) logJSON(result StepResult) {
	if p.logFile == nil {
		return
	}
	entry := struct {
		Time       string `json:"time"`
		StepResult `json:",inline"`
	}{
		Time:       time.Now().Format(time.RFC3339),
		StepResult: result,
	}
	data, _ := json.Marshal(entry)
	fmt.Fprintf(p.logFile, "%s\n", string(data))
	p.logFile.Sync()
}

func (p *Pipeline) printDirectoryStructure() {
	fmt.Printf("\nStructure:\n")
	fmt.Printf("  %s/\n", p.opts.DiscoveryDir)
	fmt.Printf("  ├── input/\n")
	inputDir := filepath.Join(p.opts.DiscoveryDir, "input")
	entries, _ := os.ReadDir(inputDir)
	for _, e := range entries {
		fmt.Printf("  │   ├── %s\n", e.Name())
	}
	fmt.Printf("  └── output/\n")
	outputDir := filepath.Join(p.opts.DiscoveryDir, "output")
	entries, _ = os.ReadDir(outputDir)
	for _, e := range entries {
		fmt.Printf("      ├── %s\n", e.Name())
	}
}

func readFileIfExists(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	return os.ReadFile(path)
}

// fileSize returns the size of a file in bytes, or 0 if it doesn't exist.
func fileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// writeStatus writes a human-readable pipeline-status.md file.
// This file is designed to be read by Claude Code to report progress to the user.
func (p *Pipeline) writeStatus(phase string, failedSteps []int) {
	if p.statusPath == "" {
		return
	}

	elapsed := time.Since(p.startTime).Truncate(time.Second)

	var b strings.Builder
	b.WriteString("# Pipeline Status\n\n")
	b.WriteString(fmt.Sprintf("**Phase:** %s\n", phase))
	b.WriteString(fmt.Sprintf("**Mode:** %s (%d steps)\n", p.mode, len(p.steps)))
	b.WriteString(fmt.Sprintf("**Elapsed:** %s\n", elapsed))
	b.WriteString(fmt.Sprintf("**Started:** %s\n", p.startTime.Format("2006-01-02 15:04:05")))

	if p.opts.ResumeFrom > 0 {
		b.WriteString(fmt.Sprintf("**Resumed from:** step %d\n", p.opts.ResumeFrom))
	}
	b.WriteString("\n")

	// Step results table
	b.WriteString("## Pipeline Steps\n\n")
	b.WriteString("| Step | Name | LLM | Status | Duration | Size |\n")
	b.WriteString("|------|------|-----|--------|----------|------|\n")

	completedSteps := 0
	for _, r := range p.results {
		if r.Name == "" {
			continue // not yet run
		}

		status := r.Status
		switch status {
		case "ok":
			completedSteps++
			if r.FellBack {
				status = fmt.Sprintf("ok (fallback: %s→claude)", r.OriginalLLM)
			}
		case "failed":
			status = fmt.Sprintf("FAILED (exit %d, %d attempts)", r.ExitCode, r.Attempts)
		}

		dur := ""
		if r.DurationMS > 0 {
			dur = (time.Duration(r.DurationMS) * time.Millisecond).Truncate(time.Second).String()
		}

		size := ""
		if r.OutputBytes > 0 {
			if r.OutputBytes > 1024 {
				size = fmt.Sprintf("%.1fKB", float64(r.OutputBytes)/1024)
			} else {
				size = fmt.Sprintf("%dB", r.OutputBytes)
			}
		}

		llmCol := r.LLMUsed
		if llmCol == "" {
			llmCol = "-"
		}
		if r.FellBack {
			llmCol = fmt.Sprintf("~~%s~~ → claude", r.OriginalLLM)
		}

		b.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s |\n",
			r.StepNum, r.Name, llmCol, status, dur, size))
	}

	b.WriteString(fmt.Sprintf("\n**Progress:** %d/%d steps complete\n", completedSteps, len(p.steps)))

	if len(failedSteps) > 0 {
		b.WriteString(fmt.Sprintf("\n**Failed steps:** %v\n", failedSteps))
		b.WriteString(fmt.Sprintf("**Resume command:** `--resume-from %d`\n", failedSteps[0]))
	}

	// VernHole info
	if phase == "complete" || phase == "complete_vernhole_failed" {
		b.WriteString("\n## VernHole\n\n")
		if phase == "complete_vernhole_failed" {
			b.WriteString("**Status:** FAILED\n")
		} else {
			vernholeDir := filepath.Join(p.opts.DiscoveryDir, "vernhole")
			entries, err := os.ReadDir(vernholeDir)
			if err == nil && len(entries) > 0 {
				b.WriteString("**Status:** complete\n")
				b.WriteString(fmt.Sprintf("**Council:** %s\n", p.opts.VernHoleCouncil))
				b.WriteString("**Files:**\n")
				for _, e := range entries {
					if strings.HasSuffix(e.Name(), ".md") {
						size := fileSize(filepath.Join(vernholeDir, e.Name()))
						b.WriteString(fmt.Sprintf("- %s (%.1fKB)\n", e.Name(), float64(size)/1024))
					}
				}
			}
		}
	}

	// Oracle info
	oracleFile := filepath.Join(p.opts.DiscoveryDir, "oracle-vision.md")
	if size := fileSize(oracleFile); size > 0 {
		b.WriteString(fmt.Sprintf("\n## Oracle\n\n**Status:** complete\n**File:** oracle-vision.md (%.1fKB)\n", float64(size)/1024))
	}

	// VTS info
	vtsDir := filepath.Join(p.opts.DiscoveryDir, "output", "vts")
	entries, err := os.ReadDir(vtsDir)
	if err == nil {
		vtsCount := 0
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "vts-") && strings.HasSuffix(e.Name(), ".md") {
				vtsCount++
			}
		}
		if vtsCount > 0 {
			b.WriteString(fmt.Sprintf("\n## VTS Tasks\n\n**Count:** %d task files in `output/vts/`\n", vtsCount))
		}
	}

	os.WriteFile(p.statusPath, []byte(b.String()), 0644)
}
