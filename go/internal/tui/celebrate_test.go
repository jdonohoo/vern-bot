package tui

import (
	"testing"
	"time"
)

func TestCelebrationConverges(t *testing.T) {
	var c CelebrationModel
	cmd := c.Start("pipeline", 120)
	if cmd == nil {
		t.Fatal("Start returned nil cmd")
	}
	if !c.active {
		t.Fatal("expected active after Start")
	}

	// Run 90 ticks — text spring should settle
	for i := 0; i < 90; i++ {
		cmd = c.Update(celebrateTickMsg{})
	}
	if !c.settled {
		t.Errorf("text spring not settled after 90 ticks, pos=%.2f vel=%.2f", c.pos, c.vel)
	}
}

func TestCelebrationEphemeralExpires(t *testing.T) {
	var c CelebrationModel
	c.StartEphemeral("historian", 120)

	// Run enough ticks to settle
	for i := 0; i < 90; i++ {
		c.Update(celebrateTickMsg{})
	}
	if !c.settled {
		t.Fatal("text spring not settled after 90 ticks")
	}

	// Force settledAt to the past to simulate timeout
	c.settledAt = time.Now().Add(-4 * time.Second)

	// Mark all particles as settled too
	for i := range c.particles {
		c.particles[i].settled = true
	}

	cmd := c.Update(celebrateTickMsg{})
	if cmd != nil {
		t.Error("expected nil cmd after ephemeral expiry")
	}
	if c.active {
		t.Error("expected inactive after ephemeral expiry")
	}
}

func TestCelebrationViewEmpty(t *testing.T) {
	var c CelebrationModel
	if v := c.View(); v != "" {
		t.Errorf("expected empty view when inactive, got %q", v)
	}
	if h := c.Height(); h != 0 {
		t.Errorf("expected height 0 when inactive, got %d", h)
	}
}

func TestCelebrationNarrowTerminal(t *testing.T) {
	var c CelebrationModel
	// Width <= text length + 5 should snap to settled
	cmd := c.Start("pipeline", 10)
	if cmd != nil {
		t.Error("expected nil cmd for narrow terminal (no animation)")
	}
	if !c.settled {
		t.Error("expected settled for narrow terminal")
	}
	if !c.active {
		t.Error("expected active even for narrow terminal")
	}
	if c.View() == "" {
		t.Error("expected non-empty view for narrow terminal")
	}
}

func TestCelebrationStartZeroWidth(t *testing.T) {
	var c CelebrationModel
	cmd := c.Start("pipeline", 0)
	if cmd != nil {
		t.Error("expected nil cmd for zero width")
	}
	if c.active {
		t.Error("expected inactive for zero width")
	}
}

func TestCelebrationConfettiHeight(t *testing.T) {
	var c CelebrationModel
	c.Start("pipeline", 120)
	h := c.Height()
	if h != 5 { // 4 canvas + 1 text
		t.Errorf("expected height 5 for confetti, got %d", h)
	}
}

func TestCelebrationSparkleHeight(t *testing.T) {
	var c CelebrationModel
	c.StartEphemeral("historian", 120)
	h := c.Height()
	if h != 1 { // sparkle = text line only
		t.Errorf("expected height 1 for sparkle, got %d", h)
	}
}

func TestCelebrationParticlesSettle(t *testing.T) {
	var c CelebrationModel
	c.Start("vernhole", 120)

	// Run 120 ticks — all particles should converge
	for i := 0; i < 120; i++ {
		c.Update(celebrateTickMsg{})
	}

	allSettled := true
	for _, p := range c.particles {
		if p.active && !p.settled {
			allSettled = false
			break
		}
	}
	if !allSettled {
		t.Error("not all particles settled after 120 ticks")
	}
}
