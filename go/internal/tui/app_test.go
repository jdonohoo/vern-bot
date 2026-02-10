package tui

import "testing"

func TestContentWidth(t *testing.T) {
	tests := []struct {
		termWidth int
		want      int
	}{
		{0, maxContentWidth},   // zero → default
		{-1, maxContentWidth},  // negative → default
		{200, maxContentWidth}, // huge terminal → capped at max
		{168, maxContentWidth}, // 168-8=160 → exactly max
		{108, 100},             // 108-8=100
		{80, 72},               // 80-8=72
		{60, 52},               // 60-8=52
		{48, minContentWidth},  // 48-8=40 → exactly min
		{30, minContentWidth},  // 30-8=22 → floored to min
	}
	for _, tt := range tests {
		got := contentWidth(tt.termWidth)
		if got != tt.want {
			t.Errorf("contentWidth(%d) = %d, want %d", tt.termWidth, got, tt.want)
		}
	}
}

func TestTextareaLines(t *testing.T) {
	tests := []struct {
		termHeight int
		want       int
	}{
		{0, 6},   // zero → default
		{-1, 6},  // negative → default
		{20, 3},  // (20-14)/3=2 → floored to 3
		{30, 5},  // (30-14)/3=5
		{40, 8},  // (40-14)/3=8
		{60, 10}, // (60-14)/3=15 → capped at 10
	}
	for _, tt := range tests {
		got := textareaLines(tt.termHeight)
		if got != tt.want {
			t.Errorf("textareaLines(%d) = %d, want %d", tt.termHeight, got, tt.want)
		}
	}
}

func TestContentWidthBounds(t *testing.T) {
	// Verify that for any reasonable terminal width, we stay within bounds
	for w := 1; w <= 300; w++ {
		got := contentWidth(w)
		if got < minContentWidth {
			t.Errorf("contentWidth(%d) = %d, below min %d", w, got, minContentWidth)
		}
		if got > maxContentWidth {
			t.Errorf("contentWidth(%d) = %d, above max %d", w, got, maxContentWidth)
		}
	}
}

func TestTextareaLinesBounds(t *testing.T) {
	// Verify that for any terminal height, lines stays in [3, 10]
	for h := 1; h <= 200; h++ {
		got := textareaLines(h)
		if got < 3 {
			t.Errorf("textareaLines(%d) = %d, below min 3", h, got)
		}
		if got > 10 {
			t.Errorf("textareaLines(%d) = %d, above max 10", h, got)
		}
	}
}

func TestIsNewer(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"2.2.0", "2.1.0", true},
		{"2.1.0", "2.2.0", false},
		{"2.1.0", "2.1.0", false},
		{"3.0.0", "2.9.9", true},
		{"2.10.0", "2.9.0", true},
		{"2.1.1", "2.1.0", true},
		{"2.1.0", "2.1.1", false},
		{"1.0.0", "2.0.0", false},
		{"10.0.0", "9.0.0", true},
		{"2.1.0.1", "2.1.0", true},
		{"2.1.0", "2.1.0.1", false},
	}

	for _, tt := range tests {
		got := isNewer(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("isNewer(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
