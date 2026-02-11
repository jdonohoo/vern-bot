package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/jdonohoo/vern-bot/go/internal/generate"
)

type genState int

const (
	genStateForm genState = iota
	genStateRunning
	genStateDone
)

// genVals holds form-bound values on the heap so pointers survive
// bubbletea's value-copy semantics.
type genVals struct {
	name        string
	description string
	model       string
	color       string
	logCh       chan string
}

// GenerateModel handles the Generate Persona screen.
type GenerateModel struct {
	state       genState
	form        *huh.Form
	spinner     spinner.Model
	viewport    viewport.Model
	width       int
	height      int
	projectRoot string
	agentsDir   string
	vals        *genVals

	running bool
	stepLog []string
	err     error
	summary string
}

// ModelOptions are the model options for persona generation (determines which LLM runs the persona).
var ModelOptions = []huh.Option[string]{
	huh.NewOption("Auto (let AI decide) (Recommended)", ""),
	huh.NewOption("Claude Opus — deep thinker, thorough analysis", "opus"),
	huh.NewOption("Claude Sonnet — fast and scrappy", "sonnet"),
	huh.NewOption("Claude Haiku — brief and minimal", "haiku"),
	huh.NewOption("Gemini 3 — 2M context window, large-scale analysis", "gemini-3"),
	huh.NewOption("Gemini Pro — deep reasoning", "gemini-pro"),
	huh.NewOption("Gemini Flash — speed-optimized", "gemini-flash"),
	huh.NewOption("Codex — raw computational power", "codex"),
	huh.NewOption("Codex Mini — lighter and faster", "codex-mini"),
	huh.NewOption("Copilot — code-focused assistance", "copilot"),
	huh.NewOption("Copilot GPT-4 — GPT-4 backbone", "copilot-gpt4"),
}

// ColorOptions are the color choices for persona TUI display.
var ColorOptions = []huh.Option[string]{
	huh.NewOption("Auto (let AI decide) (Recommended)", ""),
	huh.NewOption("Red", "red"),
	huh.NewOption("Blue", "blue"),
	huh.NewOption("Green", "green"),
	huh.NewOption("Yellow", "yellow"),
	huh.NewOption("Pink", "pink"),
	huh.NewOption("Cyan", "cyan"),
	huh.NewOption("Orange", "orange"),
	huh.NewOption("Purple", "purple"),
	huh.NewOption("Gray", "gray"),
}

func NewGenerateModel(projectRoot, agentsDir string) GenerateModel {
	vals := &genVals{
		model: "",
		color: "",
	}

	m := GenerateModel{
		state:       genStateForm,
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

func (m *GenerateModel) SetSize(w, h int) {
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

func (m *GenerateModel) buildForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Persona name").
				Placeholder("nihilist").
				Value(&v.name).
				Validate(func(s string) error {
					return generate.ValidateName(s, m.projectRoot)
				}),
		),
		huh.NewGroup(
			huh.NewText().
				Title("Description / personality prompt").
				Placeholder("Describe the persona's personality and purpose...").
				Lines(3).
				CharLimit(500).
				Value(&v.description).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("description is required")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Model / LLM").
				Options(ModelOptions...).
				Height(10).
				Value(&v.model),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Color").
				Options(ColorOptions...).
				Height(8).
				Value(&v.color),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m GenerateModel) Init() tea.Cmd {
	return m.form.Init()
}

type genDoneMsg struct {
	summary string
	err     error
}

type genLogMsg struct {
	line string
}

func (m GenerateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.running {
			if m.state == genStateDone {
				return m, backToMenu
			}
			return m, backToMenu
		}

	case genDoneMsg:
		m.state = genStateDone
		m.running = false
		m.summary = msg.summary
		m.err = msg.err
		m.initDoneViewport()
		return m, tea.DisableMouse

	case genLogMsg:
		m.stepLog = append(m.stepLog, msg.line)
		return m, m.waitForLog()
	}

	switch m.state {
	case genStateForm:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted {
			m.state = genStateRunning
			m.running = true
			m.stepLog = nil
			m.vals.logCh = make(chan string, 50)
			return m, tea.Batch(m.spinner.Tick, m.startGenerate(), m.waitForLog())
		}
		if m.form.State == huh.StateAborted {
			return m, backToMenu
		}
		return m, cmd

	case genStateRunning:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case genStateDone:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "q":
				return m, backToMenu
			case "r":
				if m.err != nil {
					m.state = genStateForm
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

func (m *GenerateModel) initDoneViewport() {
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
		content.WriteString("\n\n")
		content.WriteString(logDimStyle.Render("Press [r] to retry, or [q/esc] to return to menu."))
	} else {
		content.WriteString(stepOKStyle.Render("Persona generated successfully!"))
		content.WriteString("\n\n")
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

func (m GenerateModel) startGenerate() tea.Cmd {
	return func() tea.Msg {
		v := m.vals
		defer close(v.logCh)

		repoRoot := m.projectRoot

		var summaryLines []string

		logFn := func(msg string) {
			select {
			case v.logCh <- msg:
			default:
			}
		}

		err := generate.Run(generate.Options{
			Name:        v.name,
			Description: v.description,
			Model:       v.model,
			Color:       v.color,
			LLM:         "claude",
			RepoRoot:    repoRoot,
			LogFunc:     logFn,
		})

		if err != nil {
			return genDoneMsg{err: err}
		}

		summaryLines = append(summaryLines,
			fmt.Sprintf("Created agents/%s.md", v.name),
			fmt.Sprintf("Created commands/%s.md", v.name),
			fmt.Sprintf("Created skills/%s/SKILL.md", v.name),
			"Updated commands/v.md",
			"Updated commands/help.md",
			"Updated go/internal/embedded/embedded_test.go",
			"Updated go/internal/council/selection.go",
			"Regenerated embedded assets",
		)

		return genDoneMsg{summary: strings.Join(summaryLines, "\n")}
	}
}

func (m GenerateModel) waitForLog() tea.Cmd {
	return func() tea.Msg {
		if m.vals.logCh == nil {
			return nil
		}
		line, ok := <-m.vals.logCh
		if !ok {
			return nil
		}
		return genLogMsg{line: line}
	}
}

// Cancel aborts any running generation.
func (m GenerateModel) Cancel() {
	// Generation uses context.Background internally;
	// cancellation is handled by the bubbletea framework via ctrl+c
}

func (m GenerateModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Generate Persona"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Create a new Vern persona using AI."))
	b.WriteString("\n\n")

	switch m.state {
	case genStateForm:
		b.WriteString(m.form.View())

	case genStateRunning:
		b.WriteString(fmt.Sprintf("%s %s %s...\n", m.spinner.View(), logHeaderStyle.Render("Generating"), llmStyle.Render(m.vals.name)))
		b.WriteString(fmt.Sprintf("  Description: %s\n", logDimStyle.Render(m.vals.description)))
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

	case genStateDone:
		b.WriteString(m.viewport.View())
	}

	return b.String()
}
