package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/jdonohoo/vern-bot/go/internal/council"
	"github.com/jdonohoo/vern-bot/go/internal/llm"
)

type runState int

const (
	runStateForm runState = iota
	runStateRunning
	runStateDone
)

const maxRunAttempts = 3

// runVals holds form-bound values on the heap so pointers survive
// bubbletea's value-copy semantics.
type runVals struct {
	llmName    string
	persona    string
	outputPath string
	customPath string
	prompt     string
	cancel     context.CancelFunc
	logCh      chan string
}

// RunModel handles single LLM runs.
type RunModel struct {
	state       runState
	form        *huh.Form
	spinner     spinner.Model
	viewport    viewport.Model
	width       int
	height      int
	projectRoot string
	agentsDir   string
	vals        *runVals

	// Result
	running   bool
	stepLog   []string
	output    string
	stderr    string
	err       error
	statusMsg string
}

func NewRunModel(projectRoot, agentsDir string) RunModel {
	vals := &runVals{
		llmName:    "claude",
		outputPath: "none",
	}

	m := RunModel{
		state:       runStateForm,
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

// SetSize updates the terminal dimensions and resizes the form.
func (m *RunModel) SetSize(w, h int) {
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

func (m *RunModel) buildForm() *huh.Form {
	lines := textareaLines(m.height)
	w := contentWidth(m.width)
	v := m.vals

	// Build persona options dynamically from roster
	personaOpts := []huh.Option[string]{
		huh.NewOption("None (raw prompt)", ""),
	}
	for _, vern := range council.ScanRoster(m.agentsDir) {
		personaOpts = append(personaOpts, huh.NewOption(
			fmt.Sprintf("%s — %s", vern.ID, vern.Desc), vern.ID,
		))
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which LLM?").
				Options(SingleLLMOptions...).
				Height(len(SingleLLMOptions)+1).
				Value(&v.llmName),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Vern Persona").
				Options(personaOpts...).
				Height(8).
				Value(&v.persona),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Output location").
				Options(RunOutputOptions...).
				Height(len(RunOutputOptions)+1).
				Value(&v.outputPath),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Custom output path").
				Placeholder("./output/").
				Value(&v.customPath),
		).WithHideFunc(func() bool { return v.outputPath != "custom" }),
		huh.NewGroup(
			huh.NewText().
				Title("Enter your prompt").
				Placeholder("Describe what you want...").
				Lines(lines).
				CharLimit(2000).
				Value(&v.prompt).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("prompt is required")
					}
					return nil
				}),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m RunModel) Init() tea.Cmd {
	return m.form.Init()
}

// runDoneMsg signals the LLM run completed (with or without error).
type runDoneMsg struct {
	output string
	stderr string
	err    error
}

// runLogMsg delivers a log line to the TUI during execution.
type runLogMsg struct {
	line string
}

// runStatusClearMsg clears the temporary status flash.
type runStatusClearMsg struct{}

func (m RunModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.running {
			if m.state == runStateDone {
				return m, backToMenu
			}
			return m, backToMenu
		}

	case runDoneMsg:
		m.state = runStateDone
		m.running = false
		m.output = msg.output
		m.stderr = msg.stderr
		m.err = msg.err
		m.initDoneViewport()
		return m, tea.DisableMouse

	case runLogMsg:
		m.stepLog = append(m.stepLog, msg.line)
		return m, m.waitForLog()

	case runStatusClearMsg:
		m.statusMsg = ""
		return m, nil
	}

	switch m.state {
	case runStateForm:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted {
			m.state = runStateRunning
			m.running = true
			m.stepLog = nil
			m.vals.logCh = make(chan string, 50)
			return m, tea.Batch(m.spinner.Tick, m.startRun(), m.waitForLog())
		}
		if m.form.State == huh.StateAborted {
			return m, backToMenu
		}
		return m, cmd

	case runStateRunning:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case runStateDone:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "q":
				return m, backToMenu
			case "r":
				if m.err != nil {
					// Reset to form state, keeping prompt filled in
					m.state = runStateForm
					m.err = nil
					m.output = ""
					m.stepLog = nil
					m.form = m.buildForm()
					return m, tea.Batch(m.form.Init(), tea.EnableMouseCellMotion)
				}
			case "c":
				text := m.output
				if text == "" {
					break
				}
				if err := copyToClipboard(text); err != nil {
					m.statusMsg = "Copy failed: " + err.Error()
				} else {
					m.statusMsg = "Copied to clipboard!"
				}
				return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
					return runStatusClearMsg{}
				})
			}
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *RunModel) initDoneViewport() {
	cw := contentWidth(m.width)
	vpHeight := m.height - 6
	if vpHeight < 5 {
		vpHeight = 5
	}
	m.viewport = viewport.New(cw, vpHeight)
	m.viewport.MouseWheelEnabled = true

	var content strings.Builder
	if m.err != nil {
		content.WriteString(stepFailStyle.Render("Error: " + m.err.Error()))
		content.WriteString("\n")
		if m.stderr != "" {
			content.WriteString("\n")
			content.WriteString(logDimStyle.Render(m.stderr))
			content.WriteString("\n")
		}
		content.WriteString("\n")
		// Show retry hint
		content.WriteString(logDimStyle.Render("Press [r] to retry with a different LLM, or [q/esc] to return to menu."))
	} else {
		content.WriteString(stepOKStyle.Render("Complete!"))
		// Show output file path if one was written
		outFile := m.resolveOutputFile()
		if outFile != "" {
			content.WriteString(fmt.Sprintf("\nOutput saved to: %s", logDimStyle.Render(outFile)))
		}
		content.WriteString("\n\n")
		content.WriteString(m.output)
	}

	// Show activity log if there were retries
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

// resolveOutputFile computes the output file path from form values (without side effects).
func (m RunModel) resolveOutputFile() string {
	v := m.vals
	if v.outputPath == "none" {
		return ""
	}
	dir := "./"
	if v.outputPath == "custom" && v.customPath != "" {
		dir = expandHome(v.customPath)
	}
	name := "tui-run.md"
	if v.persona != "" {
		name = v.persona + "-vern-tui-run.md"
	}
	return filepath.Join(dir, name)
}

func (m RunModel) startRun() tea.Cmd {
	return func() tea.Msg {
		v := m.vals
		defer close(v.logCh)

		ctx, cancel := context.WithCancel(context.Background())
		v.cancel = cancel
		defer cancel()

		outputFile := m.resolveOutputFile()
		if outputFile != "" {
			dir := filepath.Dir(outputFile)
			os.MkdirAll(dir, 0755)
		}

		var lastErr error
		var lastOutput string
		var lastStderr string

		for attempt := 1; attempt <= maxRunAttempts; attempt++ {
			logLine := fmt.Sprintf(">>> Attempt %d/%d — Running %s", attempt, maxRunAttempts, v.llmName)
			if v.persona != "" {
				logLine += fmt.Sprintf(" (persona: %s)", v.persona)
			}
			select {
			case v.logCh <- logLine:
			default:
			}

			result, err := llm.Run(llm.RunOptions{
				Ctx:        ctx,
				LLM:        v.llmName,
				Prompt:     v.prompt,
				Persona:    v.persona,
				OutputFile: outputFile,
				Timeout:    20 * time.Minute,
				AgentsDir:  m.agentsDir,
			})

			if err == nil && result != nil && result.ExitCode == 0 && result.Output != "" {
				okLine := fmt.Sprintf("    OK (%s)", result.Duration.Round(100*time.Millisecond))
				select {
				case v.logCh <- okLine:
				default:
				}
				return runDoneMsg{output: result.Output}
			}

			lastErr = err
			if result != nil {
				lastOutput = result.Output
				if result.Stderr != "" {
					lastStderr = result.Stderr
				}
			}

			failLine := "    FAILED"
			if err != nil {
				failLine += ": " + err.Error()
			} else if result != nil && result.Output == "" {
				failLine += ": empty output"
			}
			select {
			case v.logCh <- failLine:
			default:
			}

			if attempt < maxRunAttempts {
				retryLine := fmt.Sprintf("    RETRY %d/%d...", attempt, maxRunAttempts-1)
				select {
				case v.logCh <- retryLine:
				default:
				}
			}
		}

		if lastErr == nil {
			lastErr = fmt.Errorf("all %d attempts failed (empty or error output)", maxRunAttempts)
		}
		return runDoneMsg{err: lastErr, output: lastOutput, stderr: lastStderr}
	}
}

func (m RunModel) waitForLog() tea.Cmd {
	return func() tea.Msg {
		if m.vals.logCh == nil {
			return nil
		}
		line, ok := <-m.vals.logCh
		if !ok {
			return nil
		}
		return runLogMsg{line: line}
	}
}

// Cancel aborts any running LLM goroutine.
func (m RunModel) Cancel() {
	if m.vals != nil && m.vals.cancel != nil {
		m.vals.cancel()
	}
}

func (m RunModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Single LLM Run"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Run a one-shot prompt against any LLM, optionally with a Vern persona."))
	b.WriteString("\n\n")

	switch m.state {
	case runStateForm:
		b.WriteString(m.form.View())

	case runStateRunning:
		b.WriteString(fmt.Sprintf("%s %s %s...\n", m.spinner.View(), logHeaderStyle.Render("Running"), llmStyle.Render(m.vals.llmName)))
		if m.vals.persona != "" {
			b.WriteString(fmt.Sprintf("  Persona: %s\n", llmStyle.Render(m.vals.persona)))
		}
		outFile := m.resolveOutputFile()
		if outFile != "" {
			b.WriteString(fmt.Sprintf("  Output: %s\n", logDimStyle.Render(outFile)))
		}
		b.WriteString("\n")

		// Show activity log
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

	case runStateDone:
		if m.statusMsg != "" {
			b.WriteString(stepOKStyle.Render(m.statusMsg) + "\n")
		}
		b.WriteString(m.viewport.View())
	}

	return b.String()
}
