package tui

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
)

type histState int

const (
	histStateForm histState = iota
	histStateRunning
	histStateDone
)

// histVals holds form-bound values on the heap so pointers survive
// bubbletea's value-copy semantics.
type histVals struct {
	directory string
	logCh     chan string
	cancel    context.CancelFunc
}

// HistorianModel handles the Historian screen.
type HistorianModel struct {
	state       histState
	form        *huh.Form
	spinner     spinner.Model
	viewport    viewport.Model
	celebration CelebrationModel
	width       int
	height      int
	projectRoot string
	agentsDir   string
	vals        *histVals

	running bool
	stepLog []string
	err     error
	summary string
}

func NewHistorianModel(projectRoot, agentsDir string) HistorianModel {
	vals := &histVals{}

	m := HistorianModel{
		state:       histStateForm,
		projectRoot: projectRoot,
		agentsDir:   agentsDir,
		vals:        vals,
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	m.spinner = s

	m.form = m.buildForm()

	return m
}

func (m *HistorianModel) SetSize(w, h int) {
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

func (m *HistorianModel) buildForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Directory to index").
				Placeholder("./my-discovery/input").
				Value(&v.directory).
				Validate(func(s string) error {
					s = strings.TrimSpace(s)
					if s == "" {
						return fmt.Errorf("directory is required")
					}
					s = expandHome(s)
					info, err := os.Stat(s)
					if err != nil {
						return fmt.Errorf("directory not found: %s", s)
					}
					if !info.IsDir() {
						return fmt.Errorf("not a directory: %s", s)
					}
					return nil
				}),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m HistorianModel) Init() tea.Cmd {
	return m.form.Init()
}

type histDoneMsg struct {
	summary string
	err     error
}

type histLogMsg struct {
	line string
}

func (m HistorianModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.running {
			if m.state == histStateDone {
				return m, backToMenu
			}
			return m, backToMenu
		}

	case celebrateTickMsg:
		cmd := m.celebration.Update(msg)
		return m, cmd

	case histDoneMsg:
		m.state = histStateDone
		m.running = false
		m.summary = msg.summary
		m.err = msg.err
		var celebCmd tea.Cmd
		if msg.err == nil {
			celebCmd = m.celebration.Start("historian", m.width)
		}
		m.initDoneViewport()
		return m, tea.Batch(tea.DisableMouse, celebCmd)

	case histLogMsg:
		m.stepLog = append(m.stepLog, msg.line)
		return m, m.waitForLog()
	}

	switch m.state {
	case histStateForm:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted {
			m.state = histStateRunning
			m.running = true
			m.stepLog = nil
			m.vals.logCh = make(chan string, 50)
			return m, tea.Batch(m.spinner.Tick, m.startHistorian(), m.waitForLog())
		}
		if m.form.State == huh.StateAborted {
			return m, backToMenu
		}
		return m, cmd

	case histStateRunning:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case histStateDone:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "q":
				return m, backToMenu
			case "r":
				if m.err != nil {
					m.state = histStateForm
					m.err = nil
					m.summary = ""
					m.stepLog = nil
					m.form = m.buildForm()
					return m, tea.Batch(m.form.Init(), tea.EnableMouseCellMotion)
				}
			}
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *HistorianModel) initDoneViewport() {
	cw := contentWidth(m.width)
	vpHeight := m.height - 6
	if m.err == nil {
		vpHeight -= m.celebration.Height()
	}
	if vpHeight < 5 {
		vpHeight = 5
	}
	m.viewport = viewport.New(cw, vpHeight)
	m.viewport.MouseWheelEnabled = true

	var content strings.Builder
	if m.err != nil {
		content.WriteString(stepFailStyle.Render("Error: " + m.err.Error()))
		content.WriteString("\n\n")
		content.WriteString(logDimStyle.Render("Press [r] to retry, or [q/esc] to return to menu."))
	} else {
		content.WriteString(m.summary)
	}

	if len(m.stepLog) > 0 {
		content.WriteString("\n\n")
		content.WriteString(logHeaderStyle.Render("Activity Log"))
		content.WriteString("\n")
		for _, line := range m.stepLog {
			content.WriteString(renderLogLine(line))
			content.WriteString("\n")
		}
	}

	m.viewport.SetContent(content.String())
}

func (m HistorianModel) startHistorian() tea.Cmd {
	return func() tea.Msg {
		v := m.vals
		defer close(v.logCh)

		ctx, cancel := context.WithCancel(context.Background())
		v.cancel = cancel

		logFn := func(msg string) {
			select {
			case v.logCh <- msg:
			default:
			}
		}

		cfg := config.Load(m.projectRoot)

		result, err := pipeline.RunHistorian(pipeline.HistorianOptions{
			Ctx:       ctx,
			TargetDir: expandHome(v.directory),
			AgentsDir: m.agentsDir,
			Timeout:   cfg.GetHistorianTimeout(),
			OnLog:     logFn,
		})

		if err != nil {
			return histDoneMsg{err: err}
		}

		var summaryLines []string
		summaryLines = append(summaryLines,
			fmt.Sprintf("Output: %s", result.OutputFile),
			fmt.Sprintf("Size: %d chars", result.CharCount),
			fmt.Sprintf("Duration: %s", result.Duration.Round(100*time.Millisecond)),
			fmt.Sprintf("LLM: %s", result.LLMUsed),
		)
		if result.FellBack {
			summaryLines = append(summaryLines,
				"",
				fmt.Sprintf("WARNING: Gemini not installed -- used %s as fallback.", result.LLMUsed),
				"The Historian needs Gemini's 2M context window for large inputs.",
				"Index may be incomplete. Install Gemini CLI for best results.",
			)
		}

		return histDoneMsg{summary: strings.Join(summaryLines, "\n")}
	}
}

func (m HistorianModel) waitForLog() tea.Cmd {
	return func() tea.Msg {
		if m.vals.logCh == nil {
			return nil
		}
		line, ok := <-m.vals.logCh
		if !ok {
			return nil
		}
		return histLogMsg{line: line}
	}
}

// Cancel aborts any running historian operation.
func (m HistorianModel) Cancel() {
	if m.vals != nil && m.vals.cancel != nil {
		m.vals.cancel()
	}
}

func (m HistorianModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Historian"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Index a directory into a structured concept map."))
	b.WriteString("\n\n")

	switch m.state {
	case histStateForm:
		b.WriteString(m.form.View())

	case histStateRunning:
		b.WriteString(fmt.Sprintf("%s %s...\n", m.spinner.View(), logHeaderStyle.Render("Historian is reading everything")))
		b.WriteString(fmt.Sprintf("  Directory: %s\n", logDimStyle.Render(m.vals.directory)))
		b.WriteString("\n")

		if len(m.stepLog) > 0 {
			maxLines := m.height - 14
			if maxLines < 4 {
				maxLines = 4
			}
			start := 0
			if len(m.stepLog) > maxLines {
				start = len(m.stepLog) - maxLines
			}
			for _, line := range m.stepLog[start:] {
				b.WriteString("  " + renderLogLine(line) + "\n")
			}
		}

	case histStateDone:
		if cv := m.celebration.View(); cv != "" {
			b.WriteString(cv)
			b.WriteString("\n")
		}
		b.WriteString(m.viewport.View())
	}

	return b.String()
}
