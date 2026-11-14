package tui

import (
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/llm"
)

func TestLLMModeOptionsCount(t *testing.T) {
	if len(LLMModeOptions) != len(config.ModeNames()) {
		t.Errorf("expected %d LLM mode options, got %d", len(config.ModeNames()), len(LLMModeOptions))
	}
}

func TestSingleLLMOptionsCount(t *testing.T) {
	if len(SingleLLMOptions) != len(llm.Names()) {
		t.Errorf("expected %d single LLM options, got %d", len(llm.Names()), len(SingleLLMOptions))
	}
}

func TestCouncilOptionsCount(t *testing.T) {
	if len(CouncilOptions) != 7 {
		t.Errorf("expected 7 council options, got %d", len(CouncilOptions))
	}
}

func TestVernHoleOptionsCount(t *testing.T) {
	if len(VernHoleOptions) != 8 {
		t.Errorf("expected 8 VernHole options (7 councils + none), got %d", len(VernHoleOptions))
	}
}

func TestPipelineOptionsCount(t *testing.T) {
	if len(PipelineOptions) != 2 {
		t.Errorf("expected 2 pipeline options, got %d", len(PipelineOptions))
	}
}

func TestOutputPathOptionsCount(t *testing.T) {
	if len(OutputPathOptions) != 2 {
		t.Errorf("expected 2 output path options, got %d", len(OutputPathOptions))
	}
}

func noDuplicateValues(t *testing.T, name string, opts []huh.Option[string]) {
	t.Helper()
	seen := make(map[string]bool)
	for _, opt := range opts {
		if seen[opt.Value] {
			t.Errorf("%s: duplicate value %q", name, opt.Value)
		}
		seen[opt.Value] = true
	}
}

func TestNoDuplicateValues(t *testing.T) {
	noDuplicateValues(t, "LLMModeOptions", LLMModeOptions)
	noDuplicateValues(t, "SingleLLMOptions", SingleLLMOptions)
	noDuplicateValues(t, "CouncilOptions", CouncilOptions)
	noDuplicateValues(t, "PipelineOptions", PipelineOptions)
	noDuplicateValues(t, "OutputPathOptions", OutputPathOptions)
}

func TestVernHoleOptionsNoDuplicateKeys(t *testing.T) {
	seen := make(map[string]bool)
	for _, opt := range VernHoleOptions {
		if seen[opt.Key] {
			t.Errorf("VernHoleOptions: duplicate key %q", opt.Key)
		}
		seen[opt.Key] = true
	}
}

func TestCouncilLabel(t *testing.T) {
	label := councilLabel("full")
	if label == "" || label == "full" {
		t.Errorf("expected a display label for 'full', got %q", label)
	}

	label = councilLabel("unknown")
	if label != "unknown" {
		t.Errorf("expected 'unknown' for unknown council, got %q", label)
	}
}

func TestExpectedLLMValues(t *testing.T) {
	// Every built-in mode must be selectable, and nothing else.
	expected := map[string]bool{}
	for _, name := range config.ModeNames() {
		expected[name] = false
	}
	for _, opt := range LLMModeOptions {
		if _, ok := expected[opt.Value]; !ok {
			t.Errorf("unexpected LLM mode value %q", opt.Value)
		}
		expected[opt.Value] = true
	}
	for val, found := range expected {
		if !found {
			t.Errorf("missing expected LLM mode value %q", val)
		}
	}
}

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}

	tests := []struct {
		input string
		want  string
	}{
		{"~/projects", home + "/projects"},
		{"~/", home + "/"},
		{"./foo", "./foo"},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
		{"~notslash", "~notslash"},
		{"", ""},
	}

	for _, tt := range tests {
		got := expandHome(tt.input)
		if got != tt.want {
			t.Errorf("expandHome(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestExpandHomeContainsNoTilde(t *testing.T) {
	result := expandHome("~/some/path")
	if strings.Contains(result, "~") {
		t.Errorf("expandHome should resolve ~, got %q", result)
	}
}

func TestExpectedSingleLLMValues(t *testing.T) {
	// Every registered provider must be selectable, and nothing else.
	expected := map[string]bool{}
	for _, name := range llm.Names() {
		expected[name] = false
	}
	for _, opt := range SingleLLMOptions {
		if _, ok := expected[opt.Value]; !ok {
			t.Errorf("unexpected single LLM value %q", opt.Value)
		}
		expected[opt.Value] = true
	}
	for val, found := range expected {
		if !found {
			t.Errorf("missing expected single LLM value %q", val)
		}
	}
}
