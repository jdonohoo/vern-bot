package tui

import "testing"

func TestHoleUpdateProgress(t *testing.T) {
	m := NewHoleModel("/tmp", "/tmp/agents")

	tests := []struct {
		line            string
		wantTotal       int
		wantCompleted   int
	}{
		{">>> Vern 1: Architect (claude)", 1, 0},
		{"    OK (2.3s)", 1, 1},
		{">>> Vern 2: Engineer (codex)", 2, 1},
		{"    FALLBACK SUCCEEDED (3.1s)", 2, 2},
		{">>> Vern 3: Validator (gemini)", 3, 2},
		{"    OK (1.8s)", 3, 3},
		{">>> Vern 4: Oracle (claude)", 4, 3},
		{"    OK (2.1s)", 4, 4},
		{">>> SYNTHESIZING council output", 4, 4},
		{"    OK (1.5s)", 4, 5}, // synthesis counts as completion
	}

	for _, tt := range tests {
		m.updateProgress(tt.line)
		if m.totalVerns != tt.wantTotal {
			t.Errorf("after line %q: totalVerns = %d, want %d",
				tt.line, m.totalVerns, tt.wantTotal)
		}
		if m.vernsCompleted != tt.wantCompleted {
			t.Errorf("after line %q: vernsCompleted = %d, want %d",
				tt.line, m.vernsCompleted, tt.wantCompleted)
		}
	}
}
