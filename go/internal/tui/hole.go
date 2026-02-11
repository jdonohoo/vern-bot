package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
)

type holeState int

const (
	holeStateForm holeState = iota
	holeStateRunning
	holeStateDone
)

// holeVals holds form-bound values on the heap so pointers survive
// bubbletea's value-copy semantics.
type holeVals struct {
	council    string
	llmMode    string
	singleLLM  string
	outputPath string
	customPath string
	idea       string
	confirm    bool
	cancel     context.CancelFunc
	logCh      chan string
}

// HoleModel handles the VernHole wizard.
type HoleModel struct {
	state       holeState
	form        *huh.Form
	spinner     spinner.Model
	progress    progress.Model
	viewport    viewport.Model
	projectRoot string
	agentsDir   string
	width       int
	height      int
	vals        *holeVals

	// Execution state
	running        bool
	stepLog        []string
	vernsCompleted int
	totalVerns     int
	err            error
}

func NewHoleModel(projectRoot, agentsDir string) HoleModel {
	vals := &holeVals{
		council:    "full",
		llmMode:    "mixed_claude_fallback",
		singleLLM:  "claude",
		outputPath: "default",
		confirm:    true,
	}

	m := HoleModel{
		state:       holeStateForm,
		projectRoot: projectRoot,
		agentsDir:   agentsDir,
		vals:        vals,
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	m.spinner = s

	p := progress.New(
		progress.WithSolidFill(string(colorSecondary)),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	m.progress = p

	m.form = m.buildForm()

	return m
}

// SetSize updates the terminal dimensions and resizes the form.
func (m *HoleModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	cw := contentWidth(w)
	fh := formHeight(h)
	if m.form != nil {
		m.form.WithWidth(cw).WithHeight(fh)
	}
	m.viewport.Width = cw
	m.viewport.Height = h - 6
	if m.viewport.Height < 5 {
		m.viewport.Height = 5
	}
}

func (m *HoleModel) buildForm() *huh.Form {
	lines := textareaLines(m.height)
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which council do you want to summon?").
				Options(CouncilOptions...).
				Height(len(CouncilOptions)+1).
				Value(&v.council),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("LLM Mode").
				Options(LLMModeOptions...).
				Height(len(LLMModeOptions)+1).
				Value(&v.llmMode),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which LLM should run all steps?").
				Options(SingleLLMOptions...).
				Height(len(SingleLLMOptions)+1).
				Value(&v.singleLLM),
		).WithHideFunc(func() bool { return v.llmMode != "single_llm" }),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Output location").
				Options(OutputPathOptions...).
				Height(len(OutputPathOptions)+1).
				Value(&v.outputPath),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Custom output path").
				Placeholder("./vernhole/").
				Value(&v.customPath),
		).WithHideFunc(func() bool { return v.outputPath != "custom" }),
		huh.NewGroup(
			huh.NewText().
				Title("Enter the prompt for the council").
				Placeholder("Describe your idea in detail...").
				Lines(lines).
				CharLimit(2000).
				Value(&v.idea).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("idea is required")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewNote().
				Title("Review Configuration").
				DescriptionFunc(func() string {
					return m.confirmSummary()
				}, &v.council),
			huh.NewConfirm().
				Title("Start VernHole?").
				Affirmative("Yes, summon the council").
				Negative("Cancel").
				Value(&v.confirm),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m HoleModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m HoleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.running {
			if m.state == holeStateDone {
				return m, backToMenu
			}
			return m, backToMenu
		}

	case holeDoneMsg:
		m.state = holeStateDone
		m.running = false
		m.err = msg.err
		m.initDoneViewport()
		return m, tea.DisableMouse

	case holeLogMsg:
		m.stepLog = append(m.stepLog, msg.line)
		m.updateProgress(msg.line)
		return m, m.waitForLog()

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	switch m.state {
	case holeStateForm:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted {
			if !m.vals.confirm {
				return m, backToMenu
			}
			m.state = holeStateRunning
			m.running = true
			m.vals.logCh = make(chan string, 100)
			return m, tea.Batch(m.spinner.Tick, m.startHole(), m.waitForLog())
		}
		if m.form.State == huh.StateAborted {
			return m, backToMenu
		}
		return m, cmd

	case holeStateRunning:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case holeStateDone:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "q" {
				return m, backToMenu
			}
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *HoleModel) updateProgress(line string) {
	upper := strings.ToUpper(line)
	// Count total Verns from ">>> Vern N:" lines
	if strings.HasPrefix(line, ">>> Vern ") {
		m.totalVerns++
	}
	if strings.Contains(upper, "OK (") || strings.Contains(upper, "FALLBACK SUCCEEDED") {
		m.vernsCompleted++
	}
}

func (m *HoleModel) initDoneViewport() {
	cw := contentWidth(m.width)
	vpHeight := m.height - 6
	if vpHeight < 5 {
		vpHeight = 5
	}
	m.viewport = viewport.New(cw, vpHeight)
	m.viewport.MouseWheelEnabled = true

	var content strings.Builder
	if m.err != nil {
		content.WriteString(stepFailStyle.Render("VernHole failed: " + m.err.Error()))
		content.WriteString("\n\n")
	} else {
		content.WriteString(stepOKStyle.Render("The VernHole has spoken!"))
		content.WriteString(fmt.Sprintf("\n\nFiles created in: %s\n", m.outputDir()))
		content.WriteString("\n")
	}

	// Try to read synthesis
	synthPath := filepath.Join(m.outputDir(), "synthesis.md")
	if data, err := os.ReadFile(synthPath); err == nil {
		content.WriteString(logHeaderStyle.Render("Synthesis"))
		content.WriteString("\n")
		content.WriteString(string(data))
		content.WriteString("\n")
	}

	// Show full log
	if len(m.stepLog) > 0 {
		content.WriteString(logHeaderStyle.Render("Activity Log"))
		content.WriteString("\n")
		for _, line := range m.stepLog {
			content.WriteString(renderLogLine(line))
			content.WriteString("\n")
		}
	}

	m.viewport.SetContent(content.String())
}

func (m HoleModel) outputDir() string {
	if m.vals.outputPath == "custom" && m.vals.customPath != "" {
		return expandHome(m.vals.customPath)
	}
	return "./vernhole/"
}

func (m HoleModel) llmModeLabel() string {
	if m.vals.llmMode == "single_llm" {
		return "single_llm (" + m.vals.singleLLM + ")"
	}
	return m.vals.llmMode
}

func (m HoleModel) confirmSummary() string {
	v := m.vals
	label := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render
	val := lipgloss.NewStyle().Foreground(colorSecondary).Render
	dim := lipgloss.NewStyle().Foreground(colorMuted).Italic(true).Render

	// Prompt section
	idea := v.idea
	if len(idea) > 120 {
		idea = idea[:117] + "..."
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("Prompt:"), dim(idea)))

	// Horizontal rule
	b.WriteString("\n  " + lipgloss.NewStyle().Foreground(colorMuted).Render(strings.Repeat("─", 40)) + "\n\n")

	// Config section
	b.WriteString(fmt.Sprintf("  %s   %s\n", label("Council:"), val(councilLabel(v.council))))
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("LLM Mode:"), val(m.llmModeLabel())))
	b.WriteString(fmt.Sprintf("  %s    %s\n", label("Output:"), dim(m.outputDir())))

	return b.String()
}

type holeDoneMsg struct {
	err error
}

type holeLogMsg struct {
	line string
}

func (m HoleModel) startHole() tea.Cmd {
	return func() tea.Msg {
		v := m.vals
		defer close(v.logCh)

		ctx, cancel := context.WithCancel(context.Background())
		v.cancel = cancel
		defer cancel()

		cfg := config.Load(m.projectRoot)

		synthesisLLM := cfg.GetSynthesisLLM()
		overrideLLM := cfg.GetOverrideLLM()

		if v.llmMode == "single_llm" && v.singleLLM != "" {
			overrideLLM = v.singleLLM
			synthesisLLM = v.singleLLM
		} else if v.llmMode != "" {
			cfg.LLMMode = v.llmMode
			synthesisLLM = cfg.GetSynthesisLLM()
			overrideLLM = cfg.GetOverrideLLM()
		}

		err := pipeline.RunVernHole(pipeline.VernHoleOptions{
			Ctx:          ctx,
			Idea:         v.idea,
			OutputDir:    m.outputDir(),
			Council:      v.council,
			AgentsDir:    m.agentsDir,
			Timeout:      1200,
			SynthesisLLM: synthesisLLM,
			OverrideLLM:  overrideLLM,
			OnLog: func(line string) {
				select {
				case v.logCh <- line:
				default:
				}
			},
		})
		return holeDoneMsg{err: err}
	}
}

func (m HoleModel) waitForLog() tea.Cmd {
	return func() tea.Msg {
		if m.vals.logCh == nil {
			return nil
		}
		line, ok := <-m.vals.logCh
		if !ok {
			return nil
		}
		return holeLogMsg{line: line}
	}
}

// Cancel aborts any running VernHole goroutine.
func (m HoleModel) Cancel() {
	if m.vals != nil && m.vals.cancel != nil {
		m.vals.cancel()
	}
}

func (m HoleModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("VernHole Council"))
	b.WriteString("\n\n")

	switch m.state {
	case holeStateForm:
		b.WriteString(m.form.View())

	case holeStateRunning:
		b.WriteString(fmt.Sprintf("%s %s\n", m.spinner.View(), logHeaderStyle.Render("Summoning the VernHole council...")))
		b.WriteString(fmt.Sprintf("  Output: %s\n\n", logDimStyle.Render(m.outputDir())))

		// Progress bar
		if m.totalVerns > 0 {
			pct := float64(m.vernsCompleted) / float64(m.totalVerns)
			if pct > 1 {
				pct = 1
			}
			b.WriteString(fmt.Sprintf("  %s %d/%d Verns\n\n",
				progressStyle.Render(m.progress.ViewAs(pct)),
				m.vernsCompleted, m.totalVerns))
		}

		cw := contentWidth(m.width)
		availHeight := m.height - 12
		if availHeight < 8 {
			availHeight = 8
		}

		if cw >= splitPanelMinWidth && len(m.stepLog) > 3 {
			leftW := cw * 2 / 5
			rightW := cw - leftW - 1

			left := renderVernStatusPanel(m.vals.council, m.llmModeLabel(), m.vernsCompleted, m.totalVerns, leftW)
			right := renderLogPanel(m.stepLog, rightW, availHeight)

			b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, left, right))
		} else {
			maxLines := availHeight
			start := 0
			if len(m.stepLog) > maxLines {
				start = len(m.stepLog) - maxLines
			}
			for _, line := range m.stepLog[start:] {
				b.WriteString("  " + renderLogLine(line) + "\n")
			}
		}

	case holeStateDone:
		b.WriteString(m.viewport.View())
	}

	return b.String()
}
