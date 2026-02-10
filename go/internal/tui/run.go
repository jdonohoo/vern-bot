package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/jdonohoo/vern-bot/go/internal/llm"
)

type runState int

const (
	runStateForm runState = iota
	runStateRunning
	runStateDone
)

// runVals holds form-bound values on the heap so pointers survive
// bubbletea's value-copy semantics.
type runVals struct {
	llmName string
	prompt  string
	cancel  context.CancelFunc
}

// RunModel handles single LLM runs.
type RunModel struct {
	state   runState
	form    *huh.Form
	spinner spinner.Model
	width   int
	height  int

	agentsDir string
	vals      *runVals

	// Result
	running bool
	output  string
	err     error
}

func NewRunModel(agentsDir string) RunModel {
	vals := &runVals{
		llmName: "claude",
	}

	m := RunModel{
		state:     runStateForm,
		agentsDir: agentsDir,
		vals:      vals,
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
	if m.form != nil {
		m.form.WithWidth(cw).WithHeight(h)
	}
}

func (m *RunModel) buildForm() *huh.Form {
	lines := textareaLines(m.height)
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which LLM?").
				Options(SingleLLMOptions...).
				Height(len(SingleLLMOptions)+1).
				Value(&v.llmName),
		),
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
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(m.height)
}

func (m RunModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m RunModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.running {
			return m, backToMenu
		}

	case runDoneMsg:
		m.state = runStateDone
		m.running = false
		m.output = msg.output
		m.err = msg.err
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
			return m, tea.Batch(m.spinner.Tick, m.startRun())
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
			if keyMsg.String() == "enter" || keyMsg.String() == "q" || keyMsg.String() == "esc" {
				return m, backToMenu
			}
		}
	}

	return m, nil
}

type runDoneMsg struct {
	output string
	err    error
}

func (m RunModel) startRun() tea.Cmd {
	return func() tea.Msg {
		v := m.vals

		ctx, cancel := context.WithCancel(context.Background())
		v.cancel = cancel
		defer cancel()

		result, err := llm.Run(llm.RunOptions{
			Ctx:       ctx,
			LLM:       v.llmName,
			Prompt:    v.prompt,
			Timeout:   20 * time.Minute,
			AgentsDir: m.agentsDir,
		})
		if err != nil {
			return runDoneMsg{err: err}
		}
		return runDoneMsg{output: result.Output}
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
	b.WriteString("\n\n")

	switch m.state {
	case runStateForm:
		b.WriteString(m.form.View())

	case runStateRunning:
		b.WriteString(fmt.Sprintf("%s %s %s...\n", m.spinner.View(), logHeaderStyle.Render("Running"), llmStyle.Render(m.vals.llmName)))

	case runStateDone:
		if m.err != nil {
			b.WriteString(stepFailStyle.Render("Error: " + m.err.Error()))
		} else {
			b.WriteString(stepOKStyle.Render("Complete!"))
			b.WriteString("\n\n")
			output := m.output
			maxLines := 30
			if m.height > 40 {
				maxLines = m.height - 10
			}
			lines := strings.Split(output, "\n")
			if len(lines) > maxLines {
				output = strings.Join(lines[:maxLines], "\n") + "\n... (truncated)"
			}
			b.WriteString(output)
		}
		b.WriteString("\n\n")
		b.WriteString(subtitleStyle.Render("Press Enter or q to return to menu"))
	}

	return b.String()
}
