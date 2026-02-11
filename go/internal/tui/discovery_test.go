package tui

import "testing"

func TestDiscoveryUpdateProgress(t *testing.T) {
	m := NewDiscoveryModel("/tmp", "/tmp/agents")
	m.totalSteps = 5

	tests := []struct {
		line            string
		wantCompleted   int
		wantProgressPct float64
	}{
		{">>> Running step 1: Architect", 0, 0.0},
		{"    OK (2.3s)", 1, 0.2},
		{">>> Running step 2: Engineer", 1, 0.2},
		{"    FAILED: timeout", 1, 0.2},
		{"    RETRY 1/3...", 1, 0.2},
		{"    FALLBACK SUCCEEDED (3.1s)", 2, 0.4},
		{">>> Running step 3: Validator", 2, 0.4},
		{"    OK (1.8s)", 3, 0.6},
		{">>> Running step 4: Consolidator", 3, 0.6},
		{"    OK (2.1s)", 4, 0.8},
		{">>> Running step 5: Reviewer", 4, 0.8},
		{"    OK (1.5s)", 5, 1.0},
		{">>> VernHole starting", 5, 1.0},
		{"    OK (5.0s)", 6, 1.0}, // overflow capped at 1.0
	}

	for _, tt := range tests {
		m.updateProgress(tt.line)
		if m.stepsCompleted != tt.wantCompleted {
			t.Errorf("after line %q: stepsCompleted = %d, want %d",
				tt.line, m.stepsCompleted, tt.wantCompleted)
		}
		pct := float64(m.stepsCompleted) / float64(m.totalSteps)
		if pct > 1 {
			pct = 1
		}
		if pct != tt.wantProgressPct {
			t.Errorf("after line %q: progress percent = %.1f, want %.1f",
				tt.line, pct, tt.wantProgressPct)
		}
	}
}
