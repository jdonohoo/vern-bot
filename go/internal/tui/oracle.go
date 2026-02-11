package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/jdonohoo/vern-bot/go/internal/config"
	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
)

type oracleState int

const (
	oracleStateForm oracleState = iota
	oracleStateRunning
	oracleStateDone
)

// OracleOperationOptions are the selectable operations.
var OracleOperationOptions = []huh.Option[string]{
	huh.NewOption("Consult the Oracle (generate vision)", "consult"),
	huh.NewOption("Apply Oracle's Vision (rewrite VTS tasks)", "apply"),
	huh.NewOption("Run VernHole on existing output", "vernhole"),
}

// oracleVals holds form-bound values on the heap so pointers survive
// bubbletea's value-copy semantics.
type oracleVals struct {
	operation    string
	synthDir     string
	vtsDir       string
	visionFile   string
	contextFile  string
	idea         string
	council      string
	llmMode      string
	singleLLM    string
	confirm      bool
	cancel       context.CancelFunc
	logCh        chan string
}

// OracleModel handles the Oracle wizard.
type OracleModel struct {
	state       oracleState
	form        *huh.Form
	spinner     spinner.Model
	progress    progress.Model
	viewport    viewport.Model
	projectRoot string
	agentsDir   string
	width       int
	height      int
	vals        *oracleVals

	// Execution state
	running        bool
	stepLog        []string
	vernsCompleted int
	totalVerns     int
	err            error
	statusMsg      string
}

func NewOracleModel(projectRoot, agentsDir string) OracleModel {
	vals := &oracleVals{
		operation: "consult",
		council:   "full",
		llmMode:   "mixed_claude_fallback",
		singleLLM: "claude",
		confirm:   true,
	}

	m := OracleModel{
		state:       oracleStateForm,
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
func (m *OracleModel) SetSize(w, h int) {
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

func (m *OracleModel) buildForm() *huh.Form {
	lines := textareaLines(m.height)
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
		// Step 1: Operation
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Options(OracleOperationOptions...).
				Height(len(OracleOperationOptions)+1).
				Value(&v.operation),
		),
		// Step 2a: Synthesis directory (consult only)
		huh.NewGroup(
			huh.NewInput().
				Title("Path to synthesis directory").
				Description("Directory containing synthesis.md from a VernHole run").
				Placeholder("./vernhole/").
				Value(&v.synthDir).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("synthesis directory is required")
					}
					return nil
				}),
		).WithHideFunc(func() bool { return v.operation != "consult" }),
		// Step 2b: VTS directory (consult + apply)
		huh.NewGroup(
			huh.NewInput().
				Title("Path to VTS directory").
				Description("Directory containing vts-*.md task files").
				Placeholder("./discovery/output/vts/").
				Value(&v.vtsDir).
				Validate(func(s string) error {
					if v.operation == "apply" && strings.TrimSpace(s) == "" {
						return fmt.Errorf("VTS directory is required for apply")
					}
					return nil
				}),
		).WithHideFunc(func() bool { return v.operation == "vernhole" }),
		// Step 2c: Oracle vision file (apply only)
		huh.NewGroup(
			huh.NewInput().
				Title("Path to oracle-vision.md").
				Description("Oracle vision file from a previous consult").
				Placeholder("./oracle-vision.md").
				Value(&v.visionFile).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("oracle vision file is required")
					}
					return nil
				}),
		).WithHideFunc(func() bool { return v.operation != "apply" }),
		// Step 2d: Context file (vernhole on output)
		huh.NewGroup(
			huh.NewInput().
				Title("Path to context file").
				Description("File to feed as context (e.g. consolidation.md or any discovery output)").
				Placeholder("./discovery/output/consolidation.md").
				Value(&v.contextFile).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("context file is required")
					}
					return nil
				}),
		).WithHideFunc(func() bool { return v.operation != "vernhole" }),
		// Step 3: Council (vernhole on output only)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which council do you want to summon?").
				Options(CouncilOptions...).
				Height(len(CouncilOptions)+1).
				Value(&v.council),
		).WithHideFunc(func() bool { return v.operation != "vernhole" }),
		// Step 4: LLM mode
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("LLM Mode").
				Options(LLMModeOptions...).
				Height(len(LLMModeOptions)+1).
				Value(&v.llmMode),
		),
		// Step 4b: Single LLM (conditional)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which LLM should run all steps?").
				Options(SingleLLMOptions...).
				Height(len(SingleLLMOptions)+1).
				Value(&v.singleLLM),
		).WithHideFunc(func() bool { return v.llmMode != "single_llm" }),
		// Step 5: Idea (consult + vernhole)
		huh.NewGroup(
			huh.NewText().
				Title("Enter the original idea / prompt").
				Placeholder("Describe the idea that was explored...").
				Lines(lines).
				CharLimit(2000).
				Value(&v.idea).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("idea is required")
					}
					return nil
				}),
		).WithHideFunc(func() bool { return v.operation == "apply" }),
		// Step 6: Confirm
		huh.NewGroup(
			huh.NewNote().
				Title("Review Configuration").
				DescriptionFunc(func() string {
					return m.confirmSummary()
				}, &v.operation),
			huh.NewConfirm().
				Title("Proceed?").
				Affirmative("Yes, run it").
				Negative("Cancel").
				Value(&v.confirm),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m OracleModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m OracleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.running {
			return m, backToMenu
		}

	case oracleDoneMsg:
		m.state = oracleStateDone
		m.running = false
		m.err = msg.err
		m.initDoneViewport()
		return m, tea.DisableMouse

	case oracleLogMsg:
		m.stepLog = append(m.stepLog, msg.line)
		m.updateProgress(msg.line)
		return m, m.waitForLog()

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case oracleStatusClearMsg:
		m.statusMsg = ""
		return m, nil
	}

	switch m.state {
	case oracleStateForm:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted {
			if !m.vals.confirm {
				return m, backToMenu
			}
			m.state = oracleStateRunning
			m.running = true
			m.vals.logCh = make(chan string, 100)
			return m, tea.Batch(m.spinner.Tick, m.startOracle(), m.waitForLog())
		}
		if m.form.State == huh.StateAborted {
			return m, backToMenu
		}
		return m, cmd

	case oracleStateRunning:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case oracleStateDone:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "q":
				return m, backToMenu
			case "c":
				text := m.resultContent()
				if text == "" {
					break
				}
				if err := copyToClipboard(text); err != nil {
					m.statusMsg = "Copy failed: " + err.Error()
				} else {
					m.statusMsg = "Copied to clipboard!"
				}
				return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
					return oracleStatusClearMsg{}
				})
			}
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *OracleModel) updateProgress(line string) {
	upper := strings.ToUpper(line)
	if strings.HasPrefix(line, ">>> Vern ") {
		m.totalVerns++
	}
	if strings.Contains(upper, "OK (") || strings.Contains(upper, "FALLBACK SUCCEEDED") {
		m.vernsCompleted++
	}
}

func (m OracleModel) resultContent() string {
	v := m.vals
	switch v.operation {
	case "consult":
		outFile := m.consultOutputFile()
		if data, err := os.ReadFile(outFile); err == nil {
			return string(data)
		}
	case "apply":
		dir := expandHome(v.vtsDir)
		entries, _ := os.ReadDir(dir)
		var b strings.Builder
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "vts-") && strings.HasSuffix(e.Name(), ".md") {
				b.WriteString(fmt.Sprintf("=== %s ===\n", e.Name()))
				if data, err := os.ReadFile(filepath.Join(dir, e.Name())); err == nil {
					b.Write(data)
					b.WriteString("\n\n")
				}
			}
		}
		return b.String()
	case "vernhole":
		synthPath := filepath.Join("./vernhole/", "synthesis.md")
		if data, err := os.ReadFile(synthPath); err == nil {
			return string(data)
		}
	}
	return ""
}

func (m OracleModel) consultOutputFile() string {
	dir := expandHome(m.vals.synthDir)
	return filepath.Join(dir, "oracle-vision.md")
}

func (m *OracleModel) initDoneViewport() {
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
	} else {
		switch m.vals.operation {
		case "consult":
			content.WriteString(stepOKStyle.Render("Oracle has spoken!"))
			content.WriteString(fmt.Sprintf("\nVision written to: %s\n\n", logDimStyle.Render(m.consultOutputFile())))
			if data, err := os.ReadFile(m.consultOutputFile()); err == nil {
				content.WriteString(logHeaderStyle.Render("Oracle Vision"))
				content.WriteString("\n")
				content.WriteString(string(data))
				content.WriteString("\n")
			}
		case "apply":
			content.WriteString(stepOKStyle.Render("Oracle's vision applied!"))
			content.WriteString(fmt.Sprintf("\nUpdated VTS files in: %s\n\n", logDimStyle.Render(expandHome(m.vals.vtsDir))))
		case "vernhole":
			content.WriteString(stepOKStyle.Render("The VernHole has spoken!"))
			content.WriteString("\nFiles created in: ./vernhole/\n\n")
			synthPath := filepath.Join("./vernhole/", "synthesis.md")
			if data, err := os.ReadFile(synthPath); err == nil {
				content.WriteString(logHeaderStyle.Render("Synthesis"))
				content.WriteString("\n")
				content.WriteString(string(data))
				content.WriteString("\n")
			}
		}
	}

	// Show full log
	if len(m.stepLog) > 0 {
		content.WriteString("\n")
		content.WriteString(logHeaderStyle.Render("Activity Log"))
		content.WriteString("\n")
		for _, line := range m.stepLog {
			content.WriteString(renderLogLine(line))
			content.WriteString("\n")
		}
	}

	m.viewport.SetContent(content.String())
}

func (m OracleModel) llmModeLabel() string {
	if m.vals.llmMode == "single_llm" {
		return "single_llm (" + m.vals.singleLLM + ")"
	}
	return m.vals.llmMode
}

func (m OracleModel) operationLabel() string {
	for _, opt := range OracleOperationOptions {
		if opt.Value == m.vals.operation {
			return opt.Key
		}
	}
	return m.vals.operation
}

func (m OracleModel) confirmSummary() string {
	v := m.vals
	label := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render
	val := lipgloss.NewStyle().Foreground(colorSecondary).Render
	dim := lipgloss.NewStyle().Foreground(colorMuted).Italic(true).Render

	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("Operation:"), val(m.operationLabel())))

	switch v.operation {
	case "consult":
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("Synthesis:"), dim(v.synthDir)))
		if v.vtsDir != "" {
			b.WriteString(fmt.Sprintf("  %s  %s\n", label("VTS Dir:"), dim(v.vtsDir)))
		}
		idea := v.idea
		if len(idea) > 80 {
			idea = idea[:77] + "..."
		}
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("Idea:"), dim(idea)))
	case "apply":
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("Vision:"), dim(v.visionFile)))
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("VTS Dir:"), dim(v.vtsDir)))
	case "vernhole":
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("Context:"), dim(v.contextFile)))
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("Council:"), val(councilLabel(v.council))))
		idea := v.idea
		if len(idea) > 80 {
			idea = idea[:77] + "..."
		}
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("Idea:"), dim(idea)))
	}

	b.WriteString("\n  " + lipgloss.NewStyle().Foreground(colorMuted).Render(strings.Repeat("─", 40)) + "\n\n")
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("LLM Mode:"), val(m.llmModeLabel())))

	return b.String()
}

type oracleDoneMsg struct {
	err error
}

type oracleLogMsg struct {
	line string
}

type oracleStatusClearMsg struct{}

func (m OracleModel) startOracle() tea.Cmd {
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

		onLog := func(line string) {
			select {
			case v.logCh <- line:
			default:
			}
		}

		var err error

		switch v.operation {
		case "consult":
			err = pipeline.RunOracleConsult(pipeline.OracleConsultOptions{
				Ctx:          ctx,
				Idea:         v.idea,
				SynthesisDir: expandHome(v.synthDir),
				VTSDir:       expandHome(v.vtsDir),
				AgentsDir:    m.agentsDir,
				SynthesisLLM: synthesisLLM,
				Timeout:      1200,
				OnLog:        onLog,
			})

		case "apply":
			err = pipeline.RunOracleApply(pipeline.OracleApplyOptions{
				Ctx:          ctx,
				VisionFile:   expandHome(v.visionFile),
				VTSDir:       expandHome(v.vtsDir),
				AgentsDir:    m.agentsDir,
				SynthesisLLM: synthesisLLM,
				Timeout:      1200,
				OnLog:        onLog,
			})

		case "vernhole":
			os.MkdirAll("./vernhole/", 0755)
			err = pipeline.RunVernHole(pipeline.VernHoleOptions{
				Ctx:          ctx,
				Idea:         v.idea,
				OutputDir:    "./vernhole/",
				Council:      v.council,
				Context:      expandHome(v.contextFile),
				AgentsDir:    m.agentsDir,
				Timeout:      1200,
				SynthesisLLM: synthesisLLM,
				OverrideLLM:  overrideLLM,
				OnLog:        onLog,
			})
		}

		return oracleDoneMsg{err: err}
	}
}

func (m OracleModel) waitForLog() tea.Cmd {
	return func() tea.Msg {
		if m.vals.logCh == nil {
			return nil
		}
		line, ok := <-m.vals.logCh
		if !ok {
			return nil
		}
		return oracleLogMsg{line: line}
	}
}

// Cancel aborts any running oracle goroutine.
func (m OracleModel) Cancel() {
	if m.vals != nil && m.vals.cancel != nil {
		m.vals.cancel()
	}
}

func (m OracleModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Post-Processing"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Run Oracle or VernHole on existing discovery output."))
	b.WriteString("\n\n")

	switch m.state {
	case oracleStateForm:
		b.WriteString(m.form.View())

	case oracleStateRunning:
		opLabel := m.operationLabel()
		b.WriteString(fmt.Sprintf("%s %s\n", m.spinner.View(), logHeaderStyle.Render(opLabel+"...")))
		b.WriteString(fmt.Sprintf("  LLM Mode: %s\n\n", logDimStyle.Render(m.llmModeLabel())))

		// Progress bar for vernhole operation
		if m.vals.operation == "vernhole" && m.totalVerns > 0 {
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

		if m.vals.operation == "vernhole" && cw >= splitPanelMinWidth && len(m.stepLog) > 3 {
			leftW := cw * 2 / 5
			rightW := cw - leftW - 1

			left := renderVernStatusPanel(m.vals.council, m.llmModeLabel(), m.vernsCompleted, m.totalVerns, leftW, availHeight)
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

	case oracleStateDone:
		if m.statusMsg != "" {
			b.WriteString(stepOKStyle.Render(m.statusMsg) + "\n")
		}
		b.WriteString(m.viewport.View())
	}

	return b.String()
}
