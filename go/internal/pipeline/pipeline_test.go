package pipeline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsFailedOutput(t *testing.T) {
	dir := t.TempDir()

	// Missing file = failed
	if !IsFailedOutput(filepath.Join(dir, "nonexistent.md")) {
		t.Error("missing file should be failed")
	}

	// Empty file = failed
	emptyFile := filepath.Join(dir, "empty.md")
	os.WriteFile(emptyFile, nil, 0644)
	if !IsFailedOutput(emptyFile) {
		t.Error("empty file should be failed")
	}

	// Failure marker = failed
	failFile := filepath.Join(dir, "failed.md")
	os.WriteFile(failFile, []byte("# STEP FAILED\n\nSome error"), 0644)
	if !IsFailedOutput(failFile) {
		t.Error("failure marker should be failed")
	}

	// Real content = not failed
	goodFile := filepath.Join(dir, "good.md")
	os.WriteFile(goodFile, []byte("# Analysis\n\nSome real content here"), 0644)
	if IsFailedOutput(goodFile) {
		t.Error("good file should not be failed")
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Initial Analysis", "initial-analysis"},
		{"Chaos Check", "chaos-check"},
		{"Hello  World!!!", "hello-world"},
		{"architect", "architect"},
	}
	for _, tt := range tests {
		got := Slugify(tt.input)
		if got != tt.want {
			t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
