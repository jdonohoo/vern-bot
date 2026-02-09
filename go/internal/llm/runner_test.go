package llm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveLLM(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"claude", "claude"},
		{"c", "claude"},
		{"Claude", "claude"},
	}
	for _, tt := range tests {
		got := resolveLLM(tt.input)
		if got != tt.want {
			t.Errorf("resolveLLM(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestLoadPersonaContext(t *testing.T) {
	dir := t.TempDir()
	content := `---
name: mighty
description: MightyVern
model: opus
---

You are MightyVern. You wield the power of Codex.
`
	os.WriteFile(filepath.Join(dir, "mighty.md"), []byte(content), 0644)

	ctx := loadPersonaContext(dir, "mighty")
	if ctx == "" {
		t.Error("persona context should not be empty")
	}
	if !containsStr(ctx, "=== PERSONA ===") {
		t.Error("should contain PERSONA markers")
	}
	if !containsStr(ctx, "You are MightyVern") {
		t.Error("should contain persona body")
	}
	// Should NOT contain frontmatter
	if containsStr(ctx, "model: opus") {
		t.Error("should not contain frontmatter")
	}
}

func TestLoadPersonaContextMissing(t *testing.T) {
	ctx := loadPersonaContext("/nonexistent", "missing")
	if ctx != "" {
		t.Error("missing persona should return empty string")
	}
}

func TestExitCodeFromErr(t *testing.T) {
	if exitCodeFromErr(nil) != 0 {
		t.Error("nil error should return 0")
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
