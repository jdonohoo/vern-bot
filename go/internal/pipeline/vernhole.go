package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jdonohoo/vern-bot/go/internal/council"
	"github.com/jdonohoo/vern-bot/go/internal/llm"
)

// VernHoleOptions configures a VernHole run.
type VernHoleOptions struct {
	Ctx          context.Context // optional: cancelled on TUI quit
	Idea         string
	OutputDir    string
	Council      string
	Count        int
	Context      string       // path to context file
	AgentsDir    string
	Timeout      int          // seconds
	SynthesisLLM string       // LLM for synthesis step (default: claude)
	OverrideLLM  string       // override all Vern LLMs (single_llm mode)
	OnLog        func(string) // optional callback for progress lines
}

// VernHoleResult holds per-Vern results.
type VernHoleResult struct {
	Index      int
	Vern       council.Vern
	Output     string
	ExitCode   int
	Succeeded  bool
	OutputFile string
}

// vernOutput prints to stdout in CLI mode, or routes to OnLog in TUI mode.
func vernOutput(opts *VernHoleOptions, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if opts.OnLog != nil {
		trimmed := strings.TrimSpace(msg)
		if trimmed != "" {
			opts.OnLog(trimmed)
		}
	} else {
		fmt.Print(msg)
	}
}

// RunVernHole executes a VernHole session with parallel Vern execution.
func RunVernHole(opts VernHoleOptions) error {
	roster := council.ScanRoster(opts.AgentsDir)

	// Determine council
	tierName := opts.Council
	if tierName == "" && opts.Count > 0 {
		tierName = fmt.Sprintf("%d", opts.Count)
	}
	if tierName == "" {
		tierName = "random"
	}

	selected, councilName := council.ResolveCouncil(tierName, roster, 3)
	numVerns := len(selected)

	// Load context
	contextBlock := ""
	if opts.Context != "" {
		data, err := os.ReadFile(opts.Context)
		if err == nil && len(data) > 0 {
			contextBlock = "\n\n=== PRIOR DISCOVERY PLAN ===\nThe following plan was synthesised from a full Vern Discovery Pipeline run on this idea. Use it as context, but bring your own unique perspective. Challenge it, build on it, tear it apart — whatever your persona demands.\n\n" + string(data) + "\n\n=== END PRIOR DISCOVERY PLAN ==="
			vernOutput(&opts, "Context loaded from: %s\n\n", opts.Context)
		}
	}

	// Banner
	vernOutput(&opts, "=== WELCOME TO THE VERNHOLE ===\n")
	vernOutput(&opts, "Idea: %s\n", opts.Idea)
	vernOutput(&opts, "Roster: %d Verns available\n", len(roster))
	if councilName != "" {
		vernOutput(&opts, "Council: %s\n", councilName)
	}
	vernOutput(&opts, "Summoning %d Verns...\n\n", numVerns)

	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 1200
	}

	// Run all Verns in parallel
	results := make([]VernHoleResult, numVerns)
	var wg sync.WaitGroup

	for i, v := range selected {
		wg.Add(1)
		go func(idx int, vern council.Vern) {
			defer wg.Done()

			outputFile := filepath.Join(opts.OutputDir, fmt.Sprintf("%02d-%s.md", idx+1, vern.ID))
			prompt := fmt.Sprintf("Analyze this idea from your unique perspective. Be true to your persona.\n\nOriginal idea: %s%s", opts.Idea, contextBlock)

			// Apply single_llm override if set
			vernLLM := vern.LLM
			if opts.OverrideLLM != "" {
				vernLLM = opts.OverrideLLM
			}

			vernOutput(&opts, ">>> Vern %d/%d: %s (%s)\n", idx+1, numVerns, vern.Desc, vernLLM)

			result, err := llm.Run(llm.RunOptions{
				Ctx:        opts.Ctx,
				LLM:        vernLLM,
				Prompt:     prompt,
				OutputFile: outputFile,
				Persona:    vern.ID,
				Timeout:    time.Duration(timeout) * time.Second,
				AgentsDir:  opts.AgentsDir,
			})

			r := VernHoleResult{
				Index:      idx,
				Vern:       vern,
				OutputFile: outputFile,
			}

			if err == nil && result.ExitCode == 0 && result.Output != "" {
				r.Output = result.Output
				r.ExitCode = 0
				r.Succeeded = true
			} else {
				exitCode := 1
				if result != nil {
					exitCode = result.ExitCode
				}
				r.ExitCode = exitCode
				vernOutput(&opts, "    WARNING: %s failed (exit %d) — excluding from synthesis\n", vern.ID, exitCode)
			}

			results[idx] = r
		}(i, v)
	}

	wg.Wait()

	// Collect results
	var allOutputs strings.Builder
	succeededCount := 0
	var failedVerns []string

	for _, r := range results {
		if r.Succeeded {
			allOutputs.WriteString(fmt.Sprintf("\n\n=== %s ===\n%s", r.Vern.Desc, r.Output))
			succeededCount++
		} else {
			failedVerns = append(failedVerns, r.Vern.ID)
		}
	}

	if len(failedVerns) > 0 {
		vernOutput(&opts, "\n>>> Failed Verns: %s\n\n", strings.Join(failedVerns, " "))
	}

	// Synthesis
	if succeededCount > 0 {
		vernOutput(&opts, ">>> Synthesizing the chaos (%d/%d Verns succeeded)...\n", succeededCount, numVerns)

		missingNote := ""
		if len(failedVerns) > 0 {
			missingNote = fmt.Sprintf("\n\nNOTE: The following Verns failed and their perspectives are missing from this synthesis: %s\nConsider what perspectives might be absent and note any gaps.", strings.Join(failedVerns, " "))
		}

		synthesisPrompt := fmt.Sprintf("Synthesize these diverse perspectives into actionable insights. Identify common themes, interesting contradictions, and recommended paths forward.\n\nORIGINAL IDEA: %s\n%s\nTHE VERNS HAVE SPOKEN:\n%s%s",
			opts.Idea, contextBlock, allOutputs.String(), missingNote)

		synthesisLLM := opts.SynthesisLLM
		if synthesisLLM == "" {
			synthesisLLM = "claude"
		}

		synthesisFile := filepath.Join(opts.OutputDir, "synthesis.md")
		_, err := llm.Run(llm.RunOptions{
			Ctx:        opts.Ctx,
			LLM:        synthesisLLM,
			Prompt:     synthesisPrompt,
			OutputFile: synthesisFile,
			Persona:    "vernhole-orchestrator",
			Timeout:    time.Duration(timeout) * time.Second,
			AgentsDir:  opts.AgentsDir,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nWARNING: Synthesis step failed\n")
		}
	} else {
		vernOutput(&opts, ">>> All Verns failed — skipping synthesis\n")
		vernOutput(&opts, "No perspectives to synthesize.\n")
	}

	// Summary
	vernOutput(&opts, "\n")
	vernOutput(&opts, "=== THE VERNHOLE HAS SPOKEN ===\n")
	vernOutput(&opts, "Files created in: %s\n", opts.OutputDir)
	vernOutput(&opts, "Succeeded: %d/%d\n", succeededCount, numVerns)
	if len(failedVerns) > 0 {
		vernOutput(&opts, "Failed: %s\n", strings.Join(failedVerns, " "))
	}

	if succeededCount == 0 {
		return fmt.Errorf("all Verns failed")
	}

	return nil
}
