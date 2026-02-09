package persona

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	// Create a temp agent file
	dir := t.TempDir()
	content := `---
name: mighty
description: MightyVern / Codex Vern - Raw computational power.
model: opus
color: blue
---

You are MightyVern. You wield the power of Codex.

PERSONALITY:
- Powerful and thorough
`
	path := filepath.Join(dir, "mighty.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != "mighty" {
		t.Errorf("name: got %q, want %q", p.Name, "mighty")
	}
	if p.Model != "opus" {
		t.Errorf("model: got %q, want %q", p.Model, "opus")
	}
	if p.Color != "blue" {
		t.Errorf("color: got %q, want %q", p.Color, "blue")
	}
	if p.Body == "" {
		t.Error("body should not be empty")
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	content := `---
name: yolo
description: YOLO Vern - No guardrails.
model: sonnet
---

Full send.
`
	if err := os.WriteFile(filepath.Join(dir, "yolo.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p, err := Load(dir, "yolo")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "yolo" {
		t.Errorf("name: got %q, want %q", p.Name, "yolo")
	}
}

func TestModelToLLM(t *testing.T) {
	tests := []struct {
		model string
		want  string
	}{
		{"opus", "claude"},
		{"sonnet", "claude"},
		{"haiku", "claude"},
		{"unknown", "claude"},
	}
	for _, tt := range tests {
		got := ModelToLLM(tt.model)
		if got != tt.want {
			t.Errorf("ModelToLLM(%q) = %q, want %q", tt.model, got, tt.want)
		}
	}
}

func TestShortDescription(t *testing.T) {
	tests := []struct {
		desc string
		want string
	}{
		{"MightyVern / Codex Vern - Raw computational power. Comprehensive solutions.", "Raw computational power"},
		{"YOLO Vern - No guardrails.", "No guardrails"},
		{"Simple description", "Simple description"},
	}
	for _, tt := range tests {
		got := ShortDescription(tt.desc)
		if got != tt.want {
			t.Errorf("ShortDescription(%q) = %q, want %q", tt.desc, got, tt.want)
		}
	}
}
