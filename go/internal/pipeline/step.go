package pipeline

import (
	"regexp"
	"strings"
)

// StepResult tracks the outcome of a pipeline step.
type StepResult struct {
	StepNum     int    `json:"step"`
	Name        string `json:"name"`
	OutputFile  string `json:"output_file"`
	Status      string `json:"status"` // "ok", "failed", "skipped"
	ExitCode    int    `json:"exit_code"`
	Attempts    int    `json:"attempts"`
	LLMUsed     string `json:"llm_used"`
	DurationMS  int64  `json:"duration_ms"`
	OutputBytes int64  `json:"output_bytes"`
}

// IsFailedOutput checks if a file is a failure marker or empty/missing.
func IsFailedOutput(path string) bool {
	data, err := readFileIfExists(path)
	if err != nil || len(data) == 0 {
		return true
	}
	// Check for failure marker written by this script
	firstLine := strings.SplitN(string(data), "\n", 2)[0]
	return strings.HasPrefix(firstLine, "# STEP FAILED")
}

// Slugify converts a string to a URL-friendly slug.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
