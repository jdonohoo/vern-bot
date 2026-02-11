package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Vern brand colors — matches docs/index.html site palette
	colorPrimary   = lipgloss.Color("#bc8cff") // purple (brand)
	colorSecondary = lipgloss.Color("#39d2c0") // teal (selection/accent)
	colorAccent    = lipgloss.Color("#58a6ff") // blue (links/buttons)
	colorSuccess   = lipgloss.Color("#3fb950") // green (OK states)
	colorDanger    = lipgloss.Color("#f85149") // red (errors)
	colorWarning   = lipgloss.Color("#d29922") // yellow (warnings)
	colorMuted     = lipgloss.Color("#8b949e") // gray (muted text)
	colorCyan      = lipgloss.Color("#39d2c0") // teal (same as secondary)
	colorOrange    = lipgloss.Color("#f0883e") // orange (flair)
	colorMagenta   = lipgloss.Color("#f778ba") // magenta (flair)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	menuSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(colorSecondary).
				Bold(true)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			MarginTop(1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorMuted)

	stepOKStyle = lipgloss.NewStyle().
			Foreground(colorSuccess)

	stepFailStyle = lipgloss.NewStyle().
			Foreground(colorDanger)

	stepWarningStyle = lipgloss.NewStyle().
				Foreground(colorWarning)

	llmStyle = lipgloss.NewStyle().
			Foreground(colorCyan)

	headerBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 2).
			MarginBottom(1)

	updateStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	// Log line styles for running views
	logHeaderStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	logStepStyle = lipgloss.NewStyle().
			Foreground(colorAccent)

	logOKStyle = lipgloss.NewStyle().
			Foreground(colorSuccess)

	logFailStyle = lipgloss.NewStyle().
			Foreground(colorDanger).
			Bold(true)

	logWarnStyle = lipgloss.NewStyle().
			Foreground(colorWarning)

	logSynthStyle = lipgloss.NewStyle().
			Foreground(colorMagenta)

	logDimStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	// Panel styles for split-panel running views
	statusPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Padding(0, 1)

	logPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMuted).
			Padding(0, 1)

	panelTitleStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			MarginBottom(1)

	// Help bar styles
	helpKeyStyle = lipgloss.NewStyle().Foreground(colorSecondary)

	helpDescStyle = lipgloss.NewStyle().Foreground(colorMuted)

	helpSepStyle = lipgloss.NewStyle().Foreground(colorMuted)

	helpBarStyle = lipgloss.NewStyle().MarginTop(1)

	// Progress bar styles
	progressStyle = lipgloss.NewStyle().MarginBottom(1)

	// Max/min content width for responsive layout.
	// contentWidth() computes the actual width for the current terminal.
	maxContentWidth = 160
	minContentWidth = 40

	// Minimum width to trigger split-panel layout
	splitPanelMinWidth = 80

	// Lines consumed by screen title + help bar chrome outside huh forms.
	// Each screen renders ~4 lines of title/spacing above the form, and
	// the App renders ~2 lines of help bar below.
	formChromeLines = 6
)

// contentWidth returns the usable content width for the given terminal width.
// Caps at maxContentWidth for readability, floors at minContentWidth.
func contentWidth(termWidth int) int {
	if termWidth <= 0 {
		return maxContentWidth
	}
	// Leave 4 columns of breathing room on each side for centering
	w := termWidth - 8
	if w < minContentWidth {
		w = minContentWidth
	}
	if w > maxContentWidth {
		w = maxContentWidth
	}
	return w
}

// formHeight returns the height available for huh forms, subtracting
// the screen title and help bar chrome from the raw terminal height.
func formHeight(termHeight int) int {
	h := termHeight - formChromeLines
	if h < 10 {
		h = 10
	}
	return h
}

// textareaLines returns the number of lines for a textarea based on terminal height.
// Scales between 3 (tiny terminal) and 10 (spacious terminal).
func textareaLines(termHeight int) int {
	if termHeight <= 0 {
		return 6
	}
	// Form chrome takes roughly 10-12 lines (title, borders, help, navigation)
	available := termHeight - 14
	lines := available / 3 // give about a third of remaining space to textarea
	if lines < 3 {
		return 3
	}
	if lines > 10 {
		return 10
	}
	return lines
}

// Pipeline step alternating color styles
var stepColors = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(colorAccent),   // blue
	lipgloss.NewStyle().Foreground(colorSecondary), // teal
}

// renderLogLine applies color to a pipeline/vernhole log line based on content.
func renderLogLine(line string) string {
	upper := strings.ToUpper(line)

	switch {
	case strings.HasPrefix(line, "==="):
		return logHeaderStyle.Render(line)
	case strings.HasPrefix(line, "[step] "):
		// Pipeline step list: "[step] N. Name → llm" — alternate colors
		display := strings.TrimPrefix(line, "[step] ")
		// Parse step number for alternating
		idx := 0
		if len(display) > 0 && display[0] >= '1' && display[0] <= '9' {
			idx = int(display[0]-'0') - 1
		}
		style := stepColors[idx%len(stepColors)]
		return style.Render("  " + display)
	case strings.HasPrefix(line, ">>>"):
		if strings.Contains(upper, "FAILED") || strings.Contains(upper, "SKIPPING") {
			return logFailStyle.Render(line)
		}
		if strings.Contains(upper, "SYNTHESIZ") {
			return logSynthStyle.Render(line)
		}
		if strings.Contains(upper, "SPLITTING") || strings.Contains(upper, "ORACLE") {
			return logSynthStyle.Render(line)
		}
		return logStepStyle.Render(line)
	case strings.Contains(upper, "HISTORIAN COMPLETE") || strings.Contains(upper, "HISTORIAN INDEX ALREADY EXISTS"):
		return logOKStyle.Render(line)
	case strings.Contains(upper, "OK (") || strings.Contains(upper, "FALLBACK SUCCEEDED"):
		return logOKStyle.Render(line)
	case strings.Contains(upper, "FAILED") || strings.Contains(upper, "ERROR"):
		return logFailStyle.Render(line)
	case strings.Contains(upper, "WARNING") || strings.Contains(upper, "TIMEOUT"):
		return logWarnStyle.Render(line)
	case strings.Contains(upper, "RETRY"):
		return logWarnStyle.Render(line)
	case strings.HasPrefix(line, "Summoning"):
		return logDimStyle.Render(line)
	default:
		return logDimStyle.Render(line)
	}
}

// isFilteredLogLine returns true for log lines that should not appear in the TUI activity log.
// These are either redundant with the static header or noise from LLM subprocess stderr.
func isFilteredLogLine(line string) bool {
	upper := strings.ToUpper(line)

	// Banner lines (redundant with static header)
	if strings.HasPrefix(line, "=== VERN") {
		return true
	}
	if strings.HasPrefix(line, "Prompt:") || strings.HasPrefix(line, "Pipeline:") || strings.HasPrefix(line, "Retries:") {
		return true
	}

	// Gemini CLI stderr noise
	if strings.Contains(upper, "YOLO MODE IS ENABLED") {
		return true
	}
	if strings.Contains(upper, "LOADED CACHED CREDENTIALS") {
		return true
	}
	if strings.Contains(upper, "HOOK REGISTRY INITIALIZED") {
		return true
	}
	if strings.Contains(upper, "ALL TOOL CALLS WILL BE AUTOMATICALLY APPROVED") {
		return true
	}

	// Empty/whitespace lines
	if strings.TrimSpace(line) == "" {
		return true
	}

	return false
}

// renderLogPanel renders the scrolling log lines into a bordered panel.
func renderLogPanel(stepLog []string, width, maxLines int) string {
	start := 0
	if len(stepLog) > maxLines {
		start = len(stepLog) - maxLines
	}
	var b strings.Builder
	for _, line := range stepLog[start:] {
		b.WriteString(renderLogLine(line) + "\n")
	}
	return logPanelStyle.Width(width).Render(
		panelTitleStyle.Render("Activity Log") + "\n" + b.String(),
	)
}

// stripMarkdown removes common markdown formatting for clean terminal display.
func stripMarkdown(content string) string {
	content = strings.ReplaceAll(content, "**", "")
	content = strings.ReplaceAll(content, "## ", "")
	content = strings.ReplaceAll(content, "# ", "")
	content = strings.ReplaceAll(content, "~~", "")
	var filtered []string
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" && len(filtered) > 0 && filtered[len(filtered)-1] == "" {
			continue
		}
		if strings.Count(trimmed, "-") > len(trimmed)/2 && strings.Contains(trimmed, "|") {
			continue
		}
		filtered = append(filtered, line)
	}
	return strings.Join(filtered, "\n")
}

// renderStatusPanel renders the pipeline status content into a bordered panel.
func renderStatusPanel(content string, width int) string {
	return statusPanelStyle.Width(width).Render(
		panelTitleStyle.Render("Pipeline Status") + "\n" + stripMarkdown(content),
	)
}

// renderVernStatusPanel renders VernHole progress into a bordered panel.
func renderVernStatusPanel(council string, llmMode string, vernsCompleted int, totalVerns int, width int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Council:  %s\n", llmStyle.Render(council)))
	b.WriteString(fmt.Sprintf("LLM Mode: %s\n\n", llmStyle.Render(llmMode)))

	b.WriteString(fmt.Sprintf("Progress: %s\n", llmStyle.Render(fmt.Sprintf("%d/%d Verns", vernsCompleted, totalVerns))))

	return statusPanelStyle.Width(width).Render(
		panelTitleStyle.Render("VernHole Status") + "\n" + b.String(),
	)
}
