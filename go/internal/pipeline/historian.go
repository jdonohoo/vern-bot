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
	LLMName    string // override LLM (default: gemini)
	OnLog      func(string)
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

	// Recursively scan target directory for readable files
	var fileContents strings.Builder
	fileCount := 0

	err := filepath.WalkDir(opts.TargetDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		// Skip the output file itself and prompt.md
		if name == "input-history.md" || name == "prompt.md" {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".md" && ext != ".txt" && ext != ".json" && ext != ".yaml" && ext != ".yml" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		// Use relative path so the LLM sees the folder structure
		relPath, relErr := filepath.Rel(opts.TargetDir, path)
		if relErr != nil {
			relPath = name
		}
		fileContents.WriteString(fmt.Sprintf("\n\n=== %s ===\n%s\n", relPath, string(data)))
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

	logFn(fmt.Sprintf("Indexing %d files from %s using %s...", fileCount, opts.TargetDir, llmName))

	// Build the Historian prompt
	prompt := fmt.Sprintf(`You are Historian Vern. You read everything. Your job is to ingest the following collection of documents — including files from subdirectories — and produce a structured index called input-history.md.

The documents below were crawled recursively from the target directory and its subfolders. File paths are relative to the root directory (e.g. "subfolder/file.md"). Preserve the folder structure in your index.

Your output MUST be a navigable concept map with:
1. A table of contents listing every source document, organized by directory structure
2. Major themes and concepts identified across all documents and subfolders, with source references (relative-path/filename + section/heading)
3. Key decisions, requirements, constraints, and open questions tagged as [DECISION], [REQUIREMENT], [OPEN QUESTION], [CONTRADICTION]
4. Cross-references between related concepts across different documents and subdirectories
5. Any contradictions or gaps between documents

Format as clean, structured markdown with headers, sub-headers, and bullet hierarchies.
Every claim must include a source reference in (path/filename:section) format.

Here are %d documents to index:
%s`, fileCount, fileContents.String())

	start := time.Now()

	result, err := llm.Run(llm.RunOptions{
		Ctx:       opts.Ctx,
		LLM:       llmName,
		Prompt:    prompt,
		Persona:   "historian",
		Timeout:   time.Duration(opts.Timeout) * time.Second,
		AgentsDir: opts.AgentsDir,
	})
	if err != nil {
		return nil, fmt.Errorf("historian LLM call failed: %w", err)
	}

	duration := time.Since(start)

	if result.ExitCode != 0 {
		return nil, fmt.Errorf("historian LLM exited with code %d", result.ExitCode)
	}

	output := result.Output
	if strings.TrimSpace(output) == "" {
		return nil, fmt.Errorf("historian produced empty output")
	}

	// Write output file
	outputFile := opts.OutputFile
	if outputFile == "" {
		outputFile = filepath.Join(opts.TargetDir, "input-history.md")
	}
	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		return nil, fmt.Errorf("write input-history.md: %w", err)
	}

	logFn(fmt.Sprintf("Wrote %s (%d chars)", outputFile, len(output)))

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
