package llm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	ExitTimeout = 124 // GNU timeout convention
)

// RunOptions configures an LLM subprocess invocation.
type RunOptions struct {
	LLM        string // claude, codex, gemini
	Prompt     string
	OutputFile string // optional: write output to this file
	Persona    string // optional: persona context to inject
	Timeout    time.Duration
	WorkingDir string // for codex --cd
	AgentsDir  string // path to agents/ for persona loading
}

// Result holds the output of an LLM run.
type Result struct {
	Output   string
	ExitCode int
	TimedOut bool
	LLMUsed  string
	Duration time.Duration
}

// Run spawns an LLM subprocess with timeout and process group management.
func Run(opts RunOptions) (*Result, error) {
	if opts.Timeout == 0 {
		opts.Timeout = 20 * time.Minute
	}

	// Resolve LLM and check availability, fall back to claude
	llm := resolveLLM(opts.LLM)

	// Build persona context
	personaContext := ""
	if opts.Persona != "" && opts.AgentsDir != "" {
		personaContext = loadPersonaContext(opts.AgentsDir, opts.Persona)
	}

	// Text-only directive for claude/gemini
	textOnly := "IMPORTANT: Output your complete analysis as plain text to stdout. Do NOT create, write, or modify any files. Do NOT use any file-writing tools. Just output your analysis directly as text.\n\n"

	// Sign-off / dad joke
	var dadJoke string
	if personaContext != "" {
		dadJoke = "\n\n---\nSIGN-OFF REMINDER: Follow your persona's sign-off instructions above (dad joke in your style). After the joke, add your persona attribution on a new line starting with '-- ' followed by your persona name and a witty tag that fits your character. Examples: '-- MightyVern *mic drop*', '-- NyQuil Vern zzz...', '-- YOLO Vern 🚀', '-- Architect Vern (measure twice, deploy once)'. This is mandatory."
	} else {
		dadJoke = "\n\n---\nSIGN-OFF: You MUST end your response with a dad joke followed by a persona attribution. Format as a horizontal rule, your joke, then '-- [Your Name]' with a witty sign-off. This is mandatory — it's the law."
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	start := time.Now()

	var cmd *exec.Cmd
	var output []byte
	var err error

	switch llm {
	case "claude":
		fullPrompt := textOnly + personaContext + opts.Prompt + dadJoke
		cmd = exec.CommandContext(ctx, "claude", "--dangerously-skip-permissions", "-p", fullPrompt)
		cmd.Env = append(os.Environ(), "NODE_OPTIONS=--max-old-space-size=32768")

	case "codex":
		codexPrefix := "IMPORTANT: You are acting as a PLANNING and ANALYSIS agent for a discovery pipeline. Write your complete analysis, implementation plan, and recommendations as a detailed markdown document. Do NOT create any code files, project scaffolding, or application code. Do NOT build anything. Your entire output should be a thorough written analysis — problem space, architecture, risks, recommendations — not a built project.\n\n"
		fullPrompt := codexPrefix + personaContext + opts.Prompt + dadJoke

		tmpFile, tmpErr := os.CreateTemp("", "vern-codex.*.md")
		if tmpErr != nil {
			return nil, fmt.Errorf("create temp file for codex: %w", tmpErr)
		}
		tmpPath := tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(tmpPath)

		codexDir := opts.WorkingDir
		if codexDir == "" {
			codexDir = "."
		}

		cmd = exec.CommandContext(ctx, "codex", "exec",
			"--dangerously-bypass-approvals-and-sandbox",
			"--skip-git-repo-check",
			"--cd", codexDir,
			"-o", tmpPath,
			fullPrompt,
		)
		// Codex writes to tmpPath; we'll read it after

		// Set process group so we can kill children
		setProcGroup(cmd)
		cmd.Stdout = nil // codex output goes to -o file
		cmd.Stderr = nil

		runErr := cmd.Run()
		duration := time.Since(start)

		exitCode := exitCodeFromErr(runErr)
		timedOut := ctx.Err() == context.DeadlineExceeded
		if timedOut {
			killProcessGroup(cmd)
			exitCode = ExitTimeout
			fmt.Fprintf(os.Stderr, "[vern-run] Timeout: codex exceeded %s limit\n", opts.Timeout)
		}

		// Read codex output
		data, _ := os.ReadFile(tmpPath)
		outputStr := string(data)

		result := &Result{
			Output:   outputStr,
			ExitCode: exitCode,
			TimedOut: timedOut,
			LLMUsed:  llm,
			Duration: duration,
		}

		return writeOutput(result, opts.OutputFile)

	case "gemini":
		fullPrompt := textOnly + personaContext + opts.Prompt + dadJoke
		cmd = exec.CommandContext(ctx, "gemini", "--yolo", fullPrompt)

	default:
		return nil, fmt.Errorf("unknown LLM: %s (valid: claude, codex, gemini)", llm)
	}

	// Set process group for non-codex LLMs
	setProcGroup(cmd)

	output, err = cmd.CombinedOutput()
	duration := time.Since(start)

	exitCode := exitCodeFromErr(err)
	timedOut := ctx.Err() == context.DeadlineExceeded
	if timedOut {
		killProcessGroup(cmd)
		exitCode = ExitTimeout
		fmt.Fprintf(os.Stderr, "[vern-run] Timeout: %s exceeded %s limit\n", llm, opts.Timeout)
	}

	result := &Result{
		Output:   string(output),
		ExitCode: exitCode,
		TimedOut: timedOut,
		LLMUsed:  llm,
		Duration: duration,
	}

	return writeOutput(result, opts.OutputFile)
}

// resolveLLM normalizes LLM names and falls back to claude if unavailable.
func resolveLLM(llm string) string {
	switch strings.ToLower(llm) {
	case "claude", "c":
		return "claude"
	case "codex", "x":
		if _, err := exec.LookPath("codex"); err != nil {
			fmt.Fprintf(os.Stderr, "[vern-run] Warning: codex CLI not found, falling back to claude\n")
			return "claude"
		}
		return "codex"
	case "gemini", "g":
		if _, err := exec.LookPath("gemini"); err != nil {
			fmt.Fprintf(os.Stderr, "[vern-run] Warning: gemini CLI not found, falling back to claude\n")
			return "claude"
		}
		return "gemini"
	default:
		return llm
	}
}

func loadPersonaContext(agentsDir, persona string) string {
	path := agentsDir + "/" + persona + ".md"
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	// Strip YAML frontmatter (between first two --- lines)
	content := string(data)
	lines := strings.SplitN(content, "\n", -1)
	dashCount := 0
	bodyStartLine := 0
	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			dashCount++
			if dashCount >= 2 {
				bodyStartLine = i + 1
				break
			}
		}
	}

	if dashCount >= 2 && bodyStartLine < len(lines) {
		body := strings.Join(lines[bodyStartLine:], "\n")
		body = strings.TrimSpace(body)
		if body != "" {
			return "=== PERSONA ===\n" + body + "\n=== END PERSONA ===\n\n"
		}
	}

	return ""
}

func exitCodeFromErr(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return 1
}

func writeOutput(result *Result, outputFile string) (*Result, error) {
	if result.Output == "" {
		fmt.Fprintf(os.Stderr, "[vern-run] Warning: LLM produced empty output (exit code: %d)\n", result.ExitCode)
	}

	if outputFile != "" {
		if result.Output != "" {
			if err := os.WriteFile(outputFile, []byte(result.Output), 0644); err != nil {
				return result, fmt.Errorf("write output file: %w", err)
			}
		} else {
			// Write empty file so downstream knows the step ran
			if err := os.WriteFile(outputFile, nil, 0644); err != nil {
				return result, fmt.Errorf("write empty output file: %w", err)
			}
		}
	}

	return result, nil
}
