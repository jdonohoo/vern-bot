package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type logEntry struct {
	Time          string `json:"time"`
	LLMRequested  string `json:"llm_requested"`
	LLMUsed       string `json:"llm_used"`
	ExitCode      int    `json:"exit_code"`
	TimedOut      bool   `json:"timed_out"`
	DurationMs    int64  `json:"duration_ms"`
	Error         string `json:"error,omitempty"`
	OutputFile    string `json:"output_file,omitempty"`
	OutputBytes   int    `json:"output_bytes"`
	PromptPreview string `json:"prompt_preview"`
}

// logRun appends a JSONL entry to ~/.config/vern/logs/vern.log.
// Disabled if VERN_LOG=0.
func logRun(opts RunOptions, llmRequested string, result *Result, err error) {
	if os.Getenv("VERN_LOG") == "0" {
		return
	}

	logDir := filepath.Join(configDir(), "logs")
	if mkErr := os.MkdirAll(logDir, 0755); mkErr != nil {
		return
	}

	entry := logEntry{
		Time:          time.Now().UTC().Format(time.RFC3339),
		LLMRequested:  llmRequested,
		PromptPreview: truncatePrompt(opts.Prompt, 200),
	}

	if opts.OutputFile != "" {
		entry.OutputFile = opts.OutputFile
	}

	if result != nil {
		entry.LLMUsed = result.LLMUsed
		entry.ExitCode = result.ExitCode
		entry.TimedOut = result.TimedOut
		entry.DurationMs = result.Duration.Milliseconds()
		entry.OutputBytes = len(result.Output)
	}

	if err != nil {
		entry.Error = err.Error()
	}

	data, jsonErr := json.Marshal(entry)
	if jsonErr != nil {
		return
	}
	data = append(data, '\n')

	logPath := filepath.Join(logDir, "vern.log")
	f, fErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fErr != nil {
		return
	}
	defer f.Close()

	f.Write(data)
}

func truncatePrompt(prompt string, max int) string {
	if len(prompt) <= max {
		return prompt
	}
	return prompt[:max] + "..."
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[vern-log] Warning: cannot determine home directory: %v\n", err)
		return "/tmp/vern"
	}
	return filepath.Join(home, ".config", "vern")
}
