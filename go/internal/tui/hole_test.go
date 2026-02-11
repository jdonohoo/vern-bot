package tui

import "testing"

func TestHoleUpdateProgress(t *testing.T) {
	m := NewHoleModel("/tmp", "/tmp/agents")

	tests := []struct {
		line          string
		wantTotal     int
		wantCompleted int
	}{
		// Verns are summoned in parallel — all appear quickly
		{">>> Vern 1/4: Architect (claude)", 4, 0},
		{">>> Vern 2/4: Engineer (codex)", 4, 0},
		{">>> Vern 3/4: Validator (gemini)", 4, 0},
		{">>> Vern 4/4: Oracle (claude)", 4, 0},
		// Completions come back in any order
		{"    OK (claude, 4800B, Vern 1/4)", 4, 1},
		{"    OK (codex, 3200B, Vern 2/4)", 4, 2},
		{"    FAILED (validator, exit 1, Vern 3/4) — excluding from synthesis", 4, 3},
		{"    OK (claude, 5100B, Vern 4/4)", 4, 4},
		// Synthesis lines don't change count
		{">>> Synthesizing the chaos (3/4 Verns succeeded)...", 4, 4},
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
