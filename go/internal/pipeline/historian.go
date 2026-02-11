package pipeline

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/llm"
)

// HistorianOptions configures a Historian invocation.
type HistorianOptions struct {
	Ctx        context.Context
	TargetDir  string // Directory to index (input/ for pipeline, user-specified for standalone)
	OutputFile string // Where to write input-history.md
	PromptFile string // Optional: prompt.md to update with reference to the index
	AgentsDir  string
	Timeout    int // seconds
	LLMName     string // override LLM (default: gemini)
	QuietStderr bool   // suppress stderr (TUI mode)
	OnLog       func(string)
}

// HistorianResult holds the outcome of a Historian run.
type HistorianResult struct {
	OutputFile string
	CharCount  int
	Duration   time.Duration
	LLMUsed    string
	FellBack   bool // true if gemini wasn't available
}

// RunHistorian scans a directory, builds a prompt from its contents, calls an LLM
// (preferring gemini for its large context window), and writes input-history.md.
func RunHistorian(opts HistorianOptions) (*HistorianResult, error) {
	if opts.Ctx == nil {
		opts.Ctx = context.Background()
	}
	if opts.Timeout == 0 {
		opts.Timeout = 600 // 10 minutes
	}

	// Resolve absolute path for the target directory
	absDir, absErr := filepath.Abs(opts.TargetDir)
	if absErr != nil {
		absDir = opts.TargetDir
	}

	// Lightweight walk: count files without reading contents
	fileCount := 0
	err := filepath.WalkDir(opts.TargetDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		if name == "input-history.md" || name == "prompt.md" {
			return nil
		}
		fileCount++
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk target directory: %w", err)
	}

	if fileCount == 0 {
		return nil, fmt.Errorf("no indexable files found in %s", opts.TargetDir)
	}

	// Resolve LLM: prefer gemini, fall back if unavailable
	llmName, fellBack := resolveHistorianLLM(opts.LLMName)

	logFn := opts.OnLog
	if logFn == nil {
		logFn = func(string) {}
	}

	if fellBack {
		logFn(fmt.Sprintf("WARNING: Gemini is not installed. Falling back to %s.", llmName))
		logFn("The Historian is designed for Gemini's 2M token context window.")
		logFn(fmt.Sprintf("Using %s may FAIL or produce incomplete results on large input directories.", llmName))
		logFn("For best results: install the Gemini CLI (https://ai.google.dev/gemini-api/docs/downloads/cli)")
	}

	logFn(fmt.Sprintf("Indexing %d files from %s using %s (LLM will crawl directory)...", fileCount, absDir, llmName))

	// Build the Historian prompt — instructs the LLM to crawl the directory itself
	prompt := fmt.Sprintf(`You are Historian Vern. You read everything. Your job is to produce a DEEP, EXHAUSTIVE index — not a summary.

## TASK

Recursively read and index ALL files in this directory:
  %s

There are approximately %d indexable files. Read ALL of them. Skip only input-history.md and prompt.md.

## STEP 1: DISCOVER THE FILE TREE (do this FIRST)

Before reading any files, run a recursive directory listing to discover the full structure:
  ls -R
or list subdirectories individually. This gives you the complete roadmap of every file and folder.

Use ONLY paths from your directory listing when reading files. Do NOT guess at directory names or file paths. If a path fails, re-list that directory to find the correct name.

## STEP 2: READ AND INDEX

Work through the file listing systematically. Read files individually or in small batches. Compile your index incrementally as you go.

## PURPOSE

This index (input-history.md) will be consumed by OTHER LLMs with SMALLER context windows. They will NOT read the original files. Your index IS their only source of truth. If you leave something out, it is lost. If you summarize too aggressively, downstream LLMs will make decisions based on incomplete information.

You have a 2M token context window. USE IT. Produce a large, detailed, thorough index. Err on the side of too much detail, not too little.

## OUTPUT REQUIREMENTS

### 1. File Inventory
- List EVERY file found, organized by directory structure
- For each file: relative path, file type, approximate size/length, and a 1-2 sentence description of what it contains

### 2. Deep Content Index (the core of your output)
For each file, produce a DETAILED breakdown — not a summary. Include:
- **Section-by-section content** with source references (relative-path/filename:section-heading or line range)
- **Code snippets**: When you encounter code, function signatures, API endpoints, config blocks, SQL schemas, or implementation details — INCLUDE THEM VERBATIM in fenced code blocks with the source reference. Do not describe code in prose when you can show it.
- **Tables and lists**: If the source contains tables, data lists, comparison matrices, or structured data — REPRODUCE THEM in your index. These are high-value reference material.
- **Specific numbers, metrics, thresholds, limits, config values** — capture them exactly, not approximately
- **Names, identifiers, URLs, file paths, version numbers** referenced in the documents — these are facts that downstream LLMs will need to look up

### 3. Semantic Tags
Tag key items inline as they appear:
- [DECISION] — choices that were made and their rationale
- [REQUIREMENT] — hard constraints, must-haves, acceptance criteria
- [OPEN QUESTION] — unresolved items, things marked as TBD or TODO
- [CONTRADICTION] — places where documents disagree with each other
- [RISK] — identified risks, concerns, failure modes
- [DEPENDENCY] — external dependencies, blockers, prerequisites

### 4. Cross-References
- Link related concepts across different files (e.g., "see also: path/other-file.md:section")
- Identify when multiple files discuss the same topic and note where they agree/disagree

### 5. Contradictions and Gaps
- Dedicated section listing contradictions between documents with exact source references for both sides
- Note gaps — topics that seem important but have no coverage, or questions raised but never answered

## FORMAT RULES

- Use clean structured markdown with headers, sub-headers, and bullet hierarchies
- EVERY claim, fact, or data point MUST include a source reference: (relative-path/filename:section) or (relative-path/filename:L42) for specific lines
- Prefer showing over telling — include the actual content (code, tables, lists, quotes) rather than describing it
- Use fenced code blocks with language tags when including code
- Do NOT editorialize or add your own opinions about the content — index what IS there
- Do NOT truncate, abbreviate, or skip files because the output is getting long — completeness is the entire point`, absDir, fileCount)

	start := time.Now()

	result, err := llm.Run(llm.RunOptions{
		Ctx:           opts.Ctx,
		LLM:           llmName,
		Prompt:        prompt,
		Persona:       "historian",
		Timeout:       time.Duration(opts.Timeout) * time.Second,
		AgentsDir:     opts.AgentsDir,
		AllowFileRead: true,
		WorkingDir:    absDir,
		QuietStderr:   opts.QuietStderr,
	})
	if err != nil {
		return nil, fmt.Errorf("historian LLM call failed: %w", err)
	}

	duration := time.Since(start)

	if result.ExitCode != 0 {
		return nil, fmt.Errorf("historian LLM exited with code %d", result.ExitCode)
	}

	// Resolve output file path before checking content sources
	outputFile := opts.OutputFile
	if outputFile == "" {
		outputFile = filepath.Join(opts.TargetDir, "input-history.md")
	}

	output := result.Output

	// Agentic LLMs (especially Gemini --yolo) may write the output file directly
	// instead of returning content on stdout. Check both sources.
	if strings.TrimSpace(output) == "" {
		if data, readErr := os.ReadFile(outputFile); readErr == nil && len(strings.TrimSpace(string(data))) > 0 {
			output = string(data)
			logFn("LLM wrote output file directly (stdout was empty) — preserving it")
		}
	}

	if strings.TrimSpace(output) == "" {
		return nil, fmt.Errorf("historian produced empty output (checked stdout and %s)", outputFile)
	}

	// Only write if the LLM didn't already write the file with this content
	existingData, _ := os.ReadFile(outputFile)
	if string(existingData) != output {
		if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
			return nil, fmt.Errorf("write input-history.md: %w", err)
		}
		logFn(fmt.Sprintf("Wrote %s (%d chars)", outputFile, len(output)))
	} else {
		logFn(fmt.Sprintf("Output file already written by LLM: %s (%d chars)", outputFile, len(output)))
	}

	// Update prompt.md if it exists
	if opts.PromptFile != "" {
		updatePromptFile(opts.PromptFile, logFn)
	}

	return &HistorianResult{
		OutputFile: outputFile,
		CharCount:  len(output),
		Duration:   duration,
		LLMUsed:    llmName,
		FellBack:   fellBack,
	}, nil
}

// updatePromptFile appends a reference to input-history.md if not already present.
func updatePromptFile(promptFile string, logFn func(string)) {
	data, err := os.ReadFile(promptFile)
	if err != nil {
		return
	}

	content := string(data)
	if strings.Contains(content, "input-history.md") {
		return // already has reference
	}

	addition := "\n\n## Additional Context\n\nSee `input-history.md` in this folder — it contains a structured index of all input materials with source references. Read this file first for an overview of the full input corpus.\n"
	if err := os.WriteFile(promptFile, []byte(content+addition), 0644); err != nil {
		logFn(fmt.Sprintf("Warning: could not update prompt.md: %v", err))
		return
	}
	logFn("Updated prompt.md with input-history.md reference")
}

// resolveHistorianLLM picks the LLM for the Historian.
// Prefers gemini for its 2M context window. Falls back to config fallback or claude.
func resolveHistorianLLM(override string) (llmName string, fellBack bool) {
	if override != "" && override != "gemini" {
		return override, false
	}

	if _, err := exec.LookPath("gemini"); err == nil {
		return "gemini", false
	}

	// Gemini not available — fall back
	cfg := config.Load("")
	fallback := cfg.GetFallbackLLM("gemini")
	if fallback == "" || fallback == "gemini" {
		fallback = "claude"
	}
	return fallback, true
}
