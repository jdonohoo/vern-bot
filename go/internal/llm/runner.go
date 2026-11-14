package llm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	ExitTimeout = 124 // GNU timeout convention
)

// RunOptions configures an LLM subprocess invocation.
type RunOptions struct {
	Ctx           context.Context // optional parent context for cancellation
	LLM           string          // provider name or alias: claude, codex, antigravity (gemini)
	Prompt        string
	OutputFile    string // optional: write output to this file
	Persona       string // optional: persona context to inject
	Timeout       time.Duration
	WorkingDir    string // working directory for the LLM subprocess
	AgentsDir     string // path to agents/ for persona loading
	AllowFileRead bool   // when true, permit the LLM to read files from the filesystem
	QuietStderr   bool   // when true, discard stderr (TUI mode — prevents display corruption)
}

// Result holds the output of an LLM run.
type Result struct {
	Output   string
	Stderr   string // last 2KB of subprocess stderr
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
	llmRequested := opts.LLM
	llm := resolveLLM(opts.LLM)

	// Build persona context
	personaContext := ""
	if opts.Persona != "" && opts.AgentsDir != "" {
		personaContext = loadPersonaContext(opts.AgentsDir, opts.Persona)
	}

	// Sign-off / dad joke
	var dadJoke string
	if personaContext != "" {
		dadJoke = "\n\n---\nSIGN-OFF REMINDER: Follow your persona's sign-off instructions above (dad joke in your style). After the joke, add your persona attribution on a new line starting with '-- ' followed by your persona name and a witty tag that fits your character. Examples: '-- MightyVern *mic drop*', '-- NyQuil Vern zzz...', '-- YOLO Vern 🚀', '-- Architect Vern (measure twice, deploy once)'. This is mandatory."
	} else {
		dadJoke = "\n\n---\nSIGN-OFF: You MUST end your response with a dad joke followed by a persona attribution. Format as a horizontal rule, your joke, then '-- [Your Name]' with a witty sign-off. This is mandatory — it's the law."
	}

	provider, ok := Lookup(llm)
	if !ok {
		return nil, fmt.Errorf("unknown LLM: %s (valid: %s)", llm, strings.Join(Names(), ", "))
	}

	fullPrompt := provider.Preamble(opts.AllowFileRead) + personaContext + opts.Prompt + dadJoke

	parent := opts.Ctx
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithTimeout(parent, opts.Timeout)
	defer cancel()

	inv := Invocation{WorkingDir: opts.WorkingDir, Timeout: opts.Timeout}
	if provider.WritesOutputFile {
		tmpFile, tmpErr := os.CreateTemp("", "vern-"+provider.Name+".*.md")
		if tmpErr != nil {
			return nil, fmt.Errorf("create temp file for %s: %w", provider.Name, tmpErr)
		}
		inv.OutputPath = tmpFile.Name()
		tmpFile.Close()
		defer os.Remove(inv.OutputPath)
	}

	start := time.Now()

	if provider.PromptAsArg {
		inv.Prompt = fullPrompt
		if len(fullPrompt) > maxArgPrompt() {
			promptFile, pErr := writePromptFile(fullPrompt)
			if pErr != nil {
				return nil, pErr
			}
			defer os.RemoveAll(filepath.Dir(promptFile))
			inv.PromptFile = promptFile
			inv.Prompt = promptFilePointer(promptFile)
		}
	}

	cmd := exec.CommandContext(ctx, provider.Binary, provider.Args(inv)...)
	if len(provider.Env) > 0 {
		cmd.Env = append(os.Environ(), provider.Env...)
	}
	// Providers that write their own output file get WorkingDir through their
	// args (codex --cd); the rest run in it so they can read files there.
	if opts.WorkingDir != "" && !provider.WritesOutputFile {
		cmd.Dir = opts.WorkingDir
	}

	// The prompt goes on stdin unless the CLI only takes it as an argument:
	// large prompts would overflow the OS command-line limits.
	if !provider.PromptAsArg {
		cmd.Stdin = strings.NewReader(fullPrompt)
	}

	// Capture stdout only — stderr passes through to the terminal (unless
	// QuietStderr) so the user sees progress, but it never pollutes the output.
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	if opts.QuietStderr {
		cmd.Stderr = &stderrBuf
	} else {
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	}

	// Own process group so a timeout kills the CLI's children too.
	setProcGroup(cmd)

	runErr := cmd.Run()
	duration := time.Since(start)

	exitCode := exitCodeFromErr(runErr)
	timedOut := ctx.Err() == context.DeadlineExceeded
	if timedOut {
		killProcessGroup(cmd)
		exitCode = ExitTimeout
		fmt.Fprintf(os.Stderr, "[vern-run] Timeout: %s exceeded %s limit\n", llm, opts.Timeout)
	}

	outputStr := stdoutBuf.String()
	if provider.WritesOutputFile {
		data, _ := os.ReadFile(inv.OutputPath)
		outputStr = string(data)
	}

	result := &Result{
		Output:   outputStr,
		Stderr:   truncStderr(stderrBuf.String(), 2048),
		ExitCode: exitCode,
		TimedOut: timedOut,
		LLMUsed:  llm,
		Duration: duration,
	}

	result, wErr := writeOutput(result, opts.OutputFile)
	logRun(opts, llmRequested, result, runErr, wErr)
	return result, wErr
}

// writePromptFile writes a prompt to a file in its own temp directory, so a
// CLI can be given read access to that directory alone.
func writePromptFile(prompt string) (string, error) {
	dir, err := os.MkdirTemp("", "vern-prompt-")
	if err != nil {
		return "", fmt.Errorf("create prompt dir: %w", err)
	}
	path := filepath.Join(dir, "prompt.md")
	if err := os.WriteFile(path, []byte(prompt), 0600); err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("write prompt file: %w", err)
	}
	return path, nil
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

// truncStderr trims stderr to the last max bytes, strips leading blank lines
// and non-printable characters.
func truncStderr(s string, max int) string {
	if s == "" {
		return ""
	}
	// Strip non-printable chars (keep tabs and newlines)
	cleaned := strings.Map(func(r rune) rune {
		if r == '\t' || r == '\n' || r == '\r' {
			return r
		}
		if r < 32 {
			return -1
		}
		return r
	}, s)
	// Keep last max bytes
	if len(cleaned) > max {
		cleaned = cleaned[len(cleaned)-max:]
	}
	// Strip leading blank lines
	cleaned = strings.TrimLeft(cleaned, "\n\r")
	return strings.TrimSpace(cleaned)
}

// FirstLine returns the first non-empty line of s, or empty string.
func FirstLine(s string) string {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}
	return ""
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
