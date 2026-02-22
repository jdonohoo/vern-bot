package tui

import (
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/harmonica"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// celebrateTickMsg drives the spring animation loop.
type celebrateTickMsg struct{}

// Phase message pools — randomly selected on Start/StartEphemeral.
var phaseMessages = map[string][]string{
	"historian": {
		"The Historian has read everything. Everything.",
		"Ancient knowledge, freshly indexed.",
		"All files consumed. The Historian is satisfied.",
	},
	"pipeline": {
		"The gauntlet is complete.",
		"All Verns survived the pipeline. Barely.",
		"Discovery complete. Opinions delivered.",
	},
	"vernhole": {
		"The VernHole has spoken!",
		"Chaos processed. Results may vary.",
		"The council rests. For now.",
	},
	"oracle": {
		"The Oracle sees all.",
		"Vision delivered. Dad jokes included.",
		"The Oracle has spoken. Mostly clearly.",
	},
}

// Particle characters and colors for celebration effects.
var particleChars = []rune{'*', '+', '.', '~'}

var particleColors = []lipgloss.Color{
	colorPrimary,
	colorSecondary,
	colorAccent,
	colorSuccess,
	colorWarning,
	colorOrange,
	colorMagenta,
}

const (
	celebrateFPS         = 30
	celebrateInterval    = time.Second / celebrateFPS
	celebrateSettleThresh = 0.5
	ephemeralTimeout     = 3 * time.Second
	confettiCanvasHeight = 4
)

type particle struct {
	char    rune
	col     int
	targetY float64
	yPos    float64
	yVel    float64
	xOffset float64
	xVel    float64
	xTarget float64
	spring  harmonica.Spring
	color   lipgloss.Color
	delay   int
	tick    int
	active  bool
	settled bool
}

// CelebrationModel provides spring-animated celebration banners with particle effects.
type CelebrationModel struct {
	active    bool
	ephemeral bool
	text      string
	style     lipgloss.Style
	spring    harmonica.Spring
	pos       float64
	vel       float64
	target    float64
	width     int
	settled   bool
	settledAt time.Time

	particles []particle
	canvas    int // 0 for sparkle (ephemeral), confettiCanvasHeight for confetti (persistent)
}

// Start begins a persistent celebration (confetti rain + text slide-in).
func (c *CelebrationModel) Start(phase string, termWidth int) tea.Cmd {
	if termWidth <= 0 {
		return nil
	}
	c.Reset()
	c.active = true
	c.ephemeral = false
	c.text = pickPhaseMessage(phase)
	c.style = lipgloss.NewStyle().Bold(true).Foreground(colorMagenta)
	c.width = termWidth
	c.canvas = confettiCanvasHeight
	c.target = 2.0

	// Narrow terminal: snap to settled
	if termWidth <= len(c.text)+5 {
		c.pos = c.target
		c.vel = 0
		c.settled = true
		c.settledAt = time.Now()
		return nil
	}

	c.pos = float64(termWidth)
	c.vel = 0
	c.spring = harmonica.NewSpring(harmonica.FPS(celebrateFPS), 6.0, 0.5)

	// Confetti particles: 12-16, spread across canvas
	nParticles := 12 + rand.Intn(5)
	c.particles = make([]particle, nParticles)
	for i := range c.particles {
		c.particles[i] = particle{
			char:    particleChars[rand.Intn(len(particleChars))],
			col:     rand.Intn(termWidth),
			targetY: float64(rand.Intn(confettiCanvasHeight)),
			yPos:    -1.0,
			yVel:    0,
			spring:  harmonica.NewSpring(harmonica.FPS(celebrateFPS), 4.0, 0.7),
			color:   particleColors[rand.Intn(len(particleColors))],
			delay:   rand.Intn(3), // stagger by 0-2 ticks
			active:  false,
		}
	}

	return c.tick()
}

// StartEphemeral begins an ephemeral celebration (sparkle burst + text slide-in).
func (c *CelebrationModel) StartEphemeral(phase string, termWidth int) tea.Cmd {
	if termWidth <= 0 {
		return nil
	}
	c.Reset()
	c.active = true
	c.ephemeral = true
	c.text = pickPhaseMessage(phase)
	c.style = lipgloss.NewStyle().Bold(true).Foreground(colorMagenta)
	c.width = termWidth
	c.canvas = 0
	c.target = 2.0

	// Narrow terminal: snap to settled
	if termWidth <= len(c.text)+5 {
		c.pos = c.target
		c.vel = 0
		c.settled = true
		c.settledAt = time.Now()
		return nil
	}

	c.pos = float64(termWidth)
	c.vel = 0
	c.spring = harmonica.NewSpring(harmonica.FPS(celebrateFPS), 6.0, 0.5)

	// Sparkle particles: 6-8, scatter outward from text center
	nParticles := 6 + rand.Intn(3)
	c.particles = make([]particle, nParticles)
	textCenter := float64(len(c.text)) / 2.0
	for i := range c.particles {
		// Random target offset: ±10-20 chars from center
		sign := 1.0
		if rand.Intn(2) == 0 {
			sign = -1.0
		}
		xTarget := sign * (10.0 + float64(rand.Intn(11)))
		c.particles[i] = particle{
			char:    particleChars[rand.Intn(len(particleChars))],
			col:     int(textCenter),
			xOffset: 0,
			xVel:    0,
			xTarget: xTarget,
			spring:  harmonica.NewSpring(harmonica.FPS(celebrateFPS), 5.0, 0.4),
			color:   particleColors[rand.Intn(len(particleColors))],
			delay:   0,
			active:  true,
		}
	}

	return c.tick()
}

// Update advances all springs by one tick.
func (c *CelebrationModel) Update(msg tea.Msg) tea.Cmd {
	if !c.active {
		return nil
	}
	if _, ok := msg.(celebrateTickMsg); !ok {
		return nil
	}

	// Advance text spring
	if !c.settled {
		c.pos, c.vel = c.spring.Update(c.pos, c.vel, c.target)
		if abs(c.pos-c.target) < celebrateSettleThresh && abs(c.vel) < celebrateSettleThresh {
			c.pos = c.target
			c.vel = 0
			c.settled = true
			c.settledAt = time.Now()
		}
	}

	// Advance particle springs
	allSettled := true
	for i := range c.particles {
		p := &c.particles[i]
		if !p.active {
			p.tick++
			if p.tick >= p.delay {
				p.active = true
			}
			allSettled = false
			continue
		}
		if p.settled {
			continue
		}

		if c.ephemeral {
			// Sparkle: X offset spring
			p.xOffset, p.xVel = p.spring.Update(p.xOffset, p.xVel, p.xTarget)
			if abs(p.xOffset-p.xTarget) < celebrateSettleThresh && abs(p.xVel) < celebrateSettleThresh {
				p.settled = true
			} else {
				allSettled = false
			}
		} else {
			// Confetti: Y position spring
			p.yPos, p.yVel = p.spring.Update(p.yPos, p.yVel, p.targetY)
			if abs(p.yPos-p.targetY) < celebrateSettleThresh && abs(p.yVel) < celebrateSettleThresh {
				p.yPos = p.targetY
				p.settled = true
			} else {
				allSettled = false
			}
		}
	}

	// Check for ephemeral expiry
	if c.ephemeral && c.settled && allSettled && time.Since(c.settledAt) >= ephemeralTimeout {
		c.active = false
		return nil
	}

	// Persistent: fade particles after settling + timeout
	if !c.ephemeral && c.settled && allSettled && time.Since(c.settledAt) >= ephemeralTimeout {
		// Deactivate all particles (text stays)
		for i := range c.particles {
			c.particles[i].active = false
		}
		c.canvas = 0
		return nil
	}

	return c.tick()
}

// View renders the celebration. Returns "" when inactive.
func (c *CelebrationModel) View() string {
	if !c.active {
		return ""
	}

	var b strings.Builder

	if !c.ephemeral && c.canvas > 0 {
		// Confetti canvas: render particles on grid lines
		for row := 0; row < c.canvas; row++ {
			line := make([]byte, c.width)
			for i := range line {
				line[i] = ' '
			}
			for _, p := range c.particles {
				if !p.active || p.col < 0 || p.col >= c.width {
					continue
				}
				rounded := int(p.yPos + 0.5)
				if rounded == row {
					style := lipgloss.NewStyle().Foreground(p.color)
					// Replace char at position — for simplicity, render line then overlay
					line[p.col] = byte(p.char)
					_ = style // we'll render per-char below
				}
			}
			// Render with per-char coloring
			rendered := c.renderCanvasRow(row)
			b.WriteString(rendered)
			b.WriteString("\n")
		}
	}

	// Text line with optional sparkles
	pad := int(c.pos + 0.5)
	if pad < 0 {
		pad = 0
	}
	if pad > c.width {
		pad = c.width
	}

	if c.ephemeral {
		// Sparkle mode: particles on same line
		leftSparkles, rightSparkles := c.renderSparkles()
		b.WriteString(strings.Repeat(" ", pad))
		b.WriteString(leftSparkles)
		b.WriteString(c.style.Render(c.text))
		b.WriteString(rightSparkles)
	} else {
		b.WriteString(strings.Repeat(" ", pad))
		b.WriteString(c.style.Render(c.text))
	}

	return b.String()
}

// Height returns the number of lines the celebration occupies.
func (c *CelebrationModel) Height() int {
	if !c.active {
		return 0
	}
	return 1 + c.canvas
}

// Reset clears all celebration state.
func (c *CelebrationModel) Reset() {
	c.active = false
	c.ephemeral = false
	c.text = ""
	c.pos = 0
	c.vel = 0
	c.settled = false
	c.particles = nil
	c.canvas = 0
}

func (c *CelebrationModel) tick() tea.Cmd {
	return tea.Tick(celebrateInterval, func(time.Time) tea.Msg {
		return celebrateTickMsg{}
	})
}

// renderCanvasRow renders one row of the confetti canvas with per-particle coloring.
func (c *CelebrationModel) renderCanvasRow(row int) string {
	type coloredChar struct {
		ch    rune
		color lipgloss.Color
	}
	chars := make(map[int]coloredChar)
	for _, p := range c.particles {
		if !p.active || p.col < 0 || p.col >= c.width {
			continue
		}
		rounded := int(p.yPos + 0.5)
		if rounded == row {
			chars[p.col] = coloredChar{ch: p.char, color: p.color}
		}
	}

	if len(chars) == 0 {
		return strings.Repeat(" ", c.width)
	}

	var b strings.Builder
	for col := 0; col < c.width; col++ {
		if cc, ok := chars[col]; ok {
			style := lipgloss.NewStyle().Foreground(cc.color)
			b.WriteString(style.Render(string(cc.ch)))
		} else {
			b.WriteByte(' ')
		}
	}
	return b.String()
}

// renderSparkles returns left and right sparkle strings for ephemeral mode.
func (c *CelebrationModel) renderSparkles() (string, string) {
	var left, right strings.Builder
	textCenter := float64(len(c.text)) / 2.0

	type sparkle struct {
		offset float64
		char   rune
		color  lipgloss.Color
	}
	var leftS, rightS []sparkle

	for _, p := range c.particles {
		if !p.active {
			continue
		}
		if p.xOffset < 0 {
			leftS = append(leftS, sparkle{offset: p.xOffset, char: p.char, color: p.color})
		} else {
			rightS = append(rightS, sparkle{offset: p.xOffset, char: p.char, color: p.color})
		}
	}

	// Render left sparkles (positioned relative to text start)
	if len(leftS) > 0 {
		maxDist := 0.0
		for _, s := range leftS {
			if -s.offset > maxDist {
				maxDist = -s.offset
			}
		}
		width := int(maxDist + 1)
		buf := make([]rune, width)
		for i := range buf {
			buf[i] = ' '
		}
		for _, s := range leftS {
			pos := width - 1 + int(s.offset+0.5)
			if pos >= 0 && pos < width {
				buf[pos] = s.char
			}
		}
		// Render with color
		for i, ch := range buf {
			if ch == ' ' {
				left.WriteByte(' ')
				continue
			}
			// Find matching sparkle color
			for _, s := range leftS {
				pos := width - 1 + int(s.offset+0.5)
				if pos == i {
					style := lipgloss.NewStyle().Foreground(s.color)
					left.WriteString(style.Render(string(ch)))
					break
				}
			}
		}
	}

	// Render right sparkles
	if len(rightS) > 0 {
		maxDist := 0.0
		for _, s := range rightS {
			if s.offset > maxDist {
				maxDist = s.offset
			}
		}
		width := int(maxDist+1) + 1
		_ = textCenter
		buf := make([]rune, width)
		for i := range buf {
			buf[i] = ' '
		}
		for _, s := range rightS {
			pos := int(s.offset + 0.5)
			if pos >= 0 && pos < width {
				buf[pos] = s.char
			}
		}
		for i, ch := range buf {
			if ch == ' ' {
				right.WriteByte(' ')
				continue
			}
			for _, s := range rightS {
				pos := int(s.offset + 0.5)
				if pos == i {
					style := lipgloss.NewStyle().Foreground(s.color)
					right.WriteString(style.Render(string(ch)))
					break
				}
			}
		}
	}

	return left.String(), right.String()
}

func pickPhaseMessage(phase string) string {
	msgs, ok := phaseMessages[phase]
	if !ok || len(msgs) == 0 {
		return "Done!"
	}
	return msgs[rand.Intn(len(msgs))]
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
