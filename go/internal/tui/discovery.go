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

	"github.com/jdonohoo/vern-bot/go/internal/pipeline"
)

type discoveryState int

const (
	discStateSetupForm discoveryState = iota
	discStateEditFiles
	discStateConfigForm
	discStateRunning
	discStateDone
)

// discoveryVals holds form-bound values on the heap so pointers survive
// bubbletea's value-copy semantics.
type discoveryVals struct {
	idea         string
	name         string
	outputPath   string
	customPath   string
	llmMode      string
	singleLLM    string
	pipeline     string
	runHistorian bool
	vernhole     string
	oracle       bool
	oracleApply  string
	confirm      bool
	cancel       context.CancelFunc
	logCh        chan string
}

// DiscoveryModel handles the discovery pipeline wizard.
type DiscoveryModel struct {
	state       discoveryState
	setupForm   *huh.Form
	configForm  *huh.Form
	spinner     spinner.Model
	progress    progress.Model
	viewport    viewport.Model
	projectRoot string
	agentsDir   string
	width       int
	height      int
	vals        *discoveryVals

	// Execution state
	running               bool
	stepLog               []string
	pipelineSteps         []string // step definitions captured from "[step] N. Name → llm"
	stepsCompleted        int
	totalSteps            int
	historianPhase        string // "pending", "running", "done", "failed", "skipped"
	currentStep           int    // which pipeline step number is currently running
	completedPipelineStep int    // highest completed pipeline step number
	statusContent         string // pipeline-status.md content for status panel
	err                   error
	setupErr              error
}

func NewDiscoveryModel(projectRoot, agentsDir string) DiscoveryModel {
	vals := &discoveryVals{
		outputPath:   "default",
		llmMode:      "mixed_claude_fallback",
		singleLLM:    "claude",
		pipeline:     "default",
		runHistorian: true,
		vernhole:     "full",
		oracle:       true,
		oracleApply:  "vision",
		confirm:      true,
	}

	m := DiscoveryModel{
		state:       discStateSetupForm,
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

	m.setupForm = m.buildSetupForm()
	m.configForm = m.buildConfigForm()

	return m
}

// SetSize updates the terminal dimensions and resizes forms accordingly.
func (m *DiscoveryModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	cw := contentWidth(w)
	fh := formHeight(h)
	if m.setupForm != nil {
		m.setupForm.WithWidth(cw).WithHeight(fh)
	}
	if m.configForm != nil {
		m.configForm.WithWidth(cw).WithHeight(fh)
	}
	m.viewport.Width = cw
	m.viewport.Height = h - 6
	if m.viewport.Height < 5 {
		m.viewport.Height = 5
	}
}

func (m *DiscoveryModel) buildSetupForm() *huh.Form {
	lines := textareaLines(m.height)
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Prompt").
				Description("This prompt will be sent to each step in the discovery pipeline").
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
			huh.NewInput().
				Title("Name for this discovery (folder name)").
				Placeholder("my-project").
				CharLimit(80).
				Value(&v.name).
				Validate(validateName),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Output location").
				Options(OutputPathOptions...).
				Value(&v.outputPath),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Custom output path").
				Placeholder("./discovery").
				Value(&v.customPath),
		).WithHideFunc(func() bool { return v.outputPath != "custom" }),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m *DiscoveryModel) buildConfigForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
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
				Title("Pipeline mode").
				Options(PipelineOptions...).
				Height(len(PipelineOptions)+1).
				Value(&v.pipeline),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Run Historian to index input files?").
				Description("Uses Gemini's 2M context to create a deep index for downstream steps").
				Affirmative("Yes (Recommended)").
				Negative("Skip").
				Value(&v.runHistorian),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Run VernHole council after pipeline?").
				Options(VernHoleOptions...).
				Height(len(VernHoleOptions)+1).
				Value(&v.vernhole),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Consult the Oracle after VernHole?").
				Description("Oracle Vern reads the council's chaos and finds the signal").
				Affirmative("Yes").
				Negative("No").
				Value(&v.oracle),
		).WithHideFunc(func() bool { return v.vernhole == "" }),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How should the Oracle's vision be handled?").
				Options(OracleApplyOptions...).
				Value(&v.oracleApply),
		).WithHideFunc(func() bool { return !v.oracle || v.vernhole == "" }),
		huh.NewGroup(
			huh.NewNote().
				Title("Review Configuration").
				DescriptionFunc(func() string {
					return m.confirmSummary()
				}, &v.idea),
			huh.NewConfirm().
				Title("Start discovery pipeline?").
				Affirmative("Yes, start pipeline").
				Negative("Cancel").
				Value(&v.confirm),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m DiscoveryModel) Init() tea.Cmd {
	return m.setupForm.Init()
}

func (m DiscoveryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.running {
			if m.state == discStateDone {
				return m, backToMenu
			}
			return m, backToMenu
		}

	case pipelineDoneMsg:
		m.state = discStateDone
		m.running = false
		m.err = msg.err
		m.readStatus()
		m.initDoneViewport()
		return m, tea.DisableMouse

	case pipelineLogMsg:
		line := msg.line

		// Capture step list lines for the outline (don't show in activity log)
		if strings.HasPrefix(line, "[step] ") {
			m.pipelineSteps = append(m.pipelineSteps, strings.TrimPrefix(line, "[step] "))
			return m, m.waitForLog()
		}

		// Filter banner/noise lines
		if isFilteredLogLine(line) {
			return m, m.waitForLog()
		}

		// Track historian state
		upper := strings.ToUpper(line)
		if strings.Contains(upper, "INVOKING THE ANCIENT SECRETS OF THE HISTORIAN") {
			m.historianPhase = "running"
		} else if strings.Contains(upper, "HISTORIAN COMPLETE") || strings.Contains(upper, "HISTORIAN INDEX ALREADY EXISTS") {
			m.historianPhase = "done"
		} else if strings.Contains(upper, "HISTORIAN FAILED") {
			m.historianPhase = "failed"
		}

		// Track current pipeline step from ">>> Pass N/M:" lines
		if strings.HasPrefix(line, ">>> Pass ") {
			rest := strings.TrimPrefix(line, ">>> Pass ")
			if idx := strings.Index(rest, "/"); idx > 0 {
				var n int
				fmt.Sscanf(rest[:idx], "%d", &n)
				if n > 0 {
					m.currentStep = n
				}
			}
		}

		m.stepLog = append(m.stepLog, line)
		m.updateProgress(line)
		return m, m.waitForLog()

	case discStatusTickMsg:
		if m.state == discStateRunning {
			m.readStatus()
			return m, m.statusTick()
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	switch m.state {
	case discStateSetupForm:
		form, cmd := m.setupForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.setupForm = f
		}
		if m.setupForm.State == huh.StateCompleted {
			m.setupErr = m.createInputFolder()
			m.state = discStateEditFiles
			return m, nil
		}
		if m.setupForm.State == huh.StateAborted {
			return m, backToMenu
		}
		return m, cmd

	case discStateEditFiles:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter":
				m.state = discStateConfigForm
				return m, m.configForm.Init()
			case "esc":
				return m, backToMenu
			}
		}
		return m, nil

	case discStateConfigForm:
		form, cmd := m.configForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.configForm = f
		}
		if m.configForm.State == huh.StateCompleted {
			if !m.vals.confirm {
				return m, backToMenu
			}
			m.state = discStateRunning
			m.running = true
			m.totalSteps = 5
			if m.vals.pipeline == "expanded" {
				m.totalSteps = 7
			}
			if m.vals.runHistorian {
				m.totalSteps++ // +1 for historian pre-step
				m.historianPhase = "pending"
			} else {
				m.historianPhase = "skipped"
			}
			m.vals.logCh = make(chan string, 100)
			return m, tea.Batch(m.spinner.Tick, m.startPipeline(), m.waitForLog(), m.statusTick())
		}
		if m.configForm.State == huh.StateAborted {
			return m, backToMenu
		}
		return m, cmd

	case discStateRunning:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case discStateDone:
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

func (m *DiscoveryModel) updateProgress(line string) {
	upper := strings.ToUpper(line)
	if strings.Contains(upper, "OK (") || strings.Contains(upper, "FALLBACK SUCCEEDED") {
		m.stepsCompleted++
		m.completedPipelineStep = m.currentStep
	} else if strings.Contains(upper, "HISTORIAN COMPLETE") ||
		strings.Contains(upper, "HISTORIAN INDEX ALREADY EXISTS") ||
		strings.Contains(upper, "HISTORIAN FAILED") {
		m.stepsCompleted++
	}
	pct := float64(m.stepsCompleted) / float64(m.totalSteps)
	if pct > 1 {
		pct = 1
	}
	m.progress.SetPercent(pct)
}

func (m *DiscoveryModel) initDoneViewport() {
	cw := contentWidth(m.width)
	vpHeight := m.height - 6
	if vpHeight < 5 {
		vpHeight = 5
	}
	m.viewport = viewport.New(cw, vpHeight)
	m.viewport.MouseWheelEnabled = true

	var content strings.Builder
	if m.err != nil {
		content.WriteString(stepFailStyle.Render("Pipeline failed: " + m.err.Error()))
		content.WriteString("\n\n")
	} else {
		content.WriteString(stepOKStyle.Render("Discovery pipeline complete!"))
		content.WriteString(fmt.Sprintf("\n\nOutput: %s\n", m.discoveryDir()))
		content.WriteString("\n")
	}

	if m.statusContent != "" {
		content.WriteString(logHeaderStyle.Render("Pipeline Status"))
		content.WriteString("\n")
		content.WriteString(stripMarkdown(m.statusContent))
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

func (m DiscoveryModel) outputBase() string {
	if m.vals.outputPath == "custom" && m.vals.customPath != "" {
		return expandHome(m.vals.customPath)
	}
	return "./discovery"
}

func (m DiscoveryModel) discoveryDir() string {
	return fmt.Sprintf("%s/%s", m.outputBase(), m.vals.name)
}

func (m DiscoveryModel) createInputFolder() error {
	inputDir := filepath.Join(m.discoveryDir(), "input")
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		return fmt.Errorf("create input folder: %w", err)
	}

	promptPath := filepath.Join(inputDir, "prompt.md")
	if _, err := os.Stat(promptPath); err == nil {
		return nil
	}

	content := fmt.Sprintf("# %s\n\n%s\n", m.vals.name, m.vals.idea)
	if err := os.WriteFile(promptPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("write prompt.md: %w", err)
	}
	return nil
}

func (m DiscoveryModel) confirmSummary() string {
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
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("Name:"), val(v.name)))
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("Output:"), dim(m.discoveryDir())))

	// Horizontal rule
	b.WriteString("\n  " + lipgloss.NewStyle().Foreground(colorMuted).Render(strings.Repeat("─", 40)) + "\n\n")

	// Config section
	mode := "Default (5-step)"
	if v.pipeline == "expanded" {
		mode = "Expanded (7-step)"
	}
	council := "None"
	if v.vernhole != "" {
		council = councilLabel(v.vernhole)
	}
	llmDesc := v.llmMode
	if v.llmMode == "single_llm" {
		llmDesc = "single_llm (" + v.singleLLM + ")"
	}

	hist := "Yes"
	if !v.runHistorian {
		hist = "Skipped"
	}

	b.WriteString(fmt.Sprintf("  %s    %s\n", label("Pipeline:"), val(mode)))
	b.WriteString(fmt.Sprintf("  %s    %s\n", label("LLM Mode:"), val(llmDesc)))
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("Historian:"), val(hist)))
	b.WriteString(fmt.Sprintf("  %s    %s\n", label("VernHole:"), val(council)))

	if v.vernhole != "" && v.oracle {
		oracleDesc := "Vision only"
		if v.oracleApply == "apply" {
			oracleDesc = "Auto-apply via Architect"
		}
		b.WriteString(fmt.Sprintf("  %s      %s\n", label("Oracle:"), val(oracleDesc)))
	}

	return b.String()
}

type pipelineDoneMsg struct {
	results []pipeline.StepResult
	err     error
}

type pipelineLogMsg struct {
	line string
}

type discStatusTickMsg time.Time

func (m DiscoveryModel) startPipeline() tea.Cmd {
	return func() tea.Msg {
		v := m.vals
		defer close(v.logCh)

		ctx, cancel := context.WithCancel(context.Background())
		v.cancel = cancel
		defer cancel()

		dir := m.discoveryDir()

		opts := pipeline.Options{
			Ctx:           ctx,
			Idea:          v.idea,
			DiscoveryDir:  dir,
			BatchMode:     true,
			ReadInput:     true,
			SkipHistorian: !v.runHistorian,
			Expanded:      v.pipeline == "expanded",
			AgentsDir:     m.agentsDir,
			ProjectRoot:   m.projectRoot,
			LLMMode:       v.llmMode,
			SingleLLM:     v.singleLLM,
			OnLog: func(line string) {
				select {
				case v.logCh <- line:
				default:
				}
			},
		}

		if v.vernhole != "" {
			opts.VernHoleCouncil = v.vernhole
			if v.oracle {
				opts.OracleFlag = true
				opts.OracleApplyFlag = v.oracleApply == "apply"
			}
		}

		err := pipeline.Run(opts)
		return pipelineDoneMsg{err: err}
	}
}

func (m DiscoveryModel) waitForLog() tea.Cmd {
	return func() tea.Msg {
		if m.vals.logCh == nil {
			return nil
		}
		line, ok := <-m.vals.logCh
		if !ok {
			return nil
		}
		return pipelineLogMsg{line: line}
	}
}

func (m DiscoveryModel) statusTick() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return discStatusTickMsg(t)
	})
}

func (m *DiscoveryModel) readStatus() {
	statusPath := filepath.Join(m.discoveryDir(), "output", "pipeline-status.md")
	data, err := os.ReadFile(statusPath)
	if err == nil {
		m.statusContent = string(data)
	}
}

// Cancel aborts any running pipeline goroutine.
func (m DiscoveryModel) Cancel() {
	if m.vals != nil && m.vals.cancel != nil {
		m.vals.cancel()
	}
}

func (m DiscoveryModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Discovery Pipeline"))
	b.WriteString("\n\n")

	switch m.state {
	case discStateSetupForm:
		b.WriteString(m.setupForm.View())

	case discStateEditFiles:
		dir := m.discoveryDir()
		inputDir := filepath.Join(dir, "input")
		promptPath := filepath.Join(inputDir, "prompt.md")

		if m.setupErr != nil {
			b.WriteString(stepFailStyle.Render("Failed to create folder: " + m.setupErr.Error()))
			b.WriteString("\n\n")
			b.WriteString(subtitleStyle.Render("Press Enter to continue anyway"))
			b.WriteString("\n")
			b.WriteString(subtitleStyle.Render("Press Esc to go back"))
		} else {
			b.WriteString(stepOKStyle.Render("Created discovery folder:"))
			b.WriteString("\n\n")
			b.WriteString(fmt.Sprintf("  %s\n", dir))
			b.WriteString(fmt.Sprintf("  %s\n", inputDir))
			b.WriteString(fmt.Sprintf("  %s\n", promptPath))
			b.WriteString("\n")
			b.WriteString("Your idea has been written to prompt.md.\n")
			b.WriteString("You can now:\n\n")
			b.WriteString("  1. Edit prompt.md to refine your idea\n")
			b.WriteString("  2. Add reference files to the input/ folder\n")
			b.WriteString("\n")
			b.WriteString(subtitleStyle.Render("Press Enter when ready to configure"))
			b.WriteString("\n")
			b.WriteString(subtitleStyle.Render("Press Esc to go back"))
		}

	case discStateConfigForm:
		b.WriteString(m.configForm.View())

	case discStateRunning:
		label := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render
		cw := contentWidth(m.width)

		b.WriteString(fmt.Sprintf("%s %s\n\n", m.spinner.View(), logHeaderStyle.Render("Running discovery pipeline...")))

		// Config line
		mode := "default"
		if m.vals.pipeline == "expanded" {
			mode = "expanded"
		}
		b.WriteString(fmt.Sprintf("  %s  %s  |  %s  %s  |  %s  %s\n",
			label("Pipeline:"), llmStyle.Render(mode),
			label("LLM:"), llmStyle.Render(m.vals.llmMode),
			label("Output:"), logDimStyle.Render(m.discoveryDir())))

		// Step outline — only show when there's enough vertical space
		outlineLines := len(m.pipelineSteps)
		if m.historianPhase != "skipped" {
			outlineLines++
		}
		// Need: 2 title + 1 config + 2 separators + outline + panels(8 min) + 2 progress = ~15 + outline
		showOutline := (m.height - 15 - outlineLines) >= 8

		if showOutline && (len(m.pipelineSteps) > 0 || m.historianPhase != "skipped") {
			b.WriteString("  " + logDimStyle.Render(strings.Repeat("─", 50)) + "\n")

			// Historian line
			if m.historianPhase != "skipped" {
				historianLabel := "Historian (pre-step)"
				switch m.historianPhase {
				case "done":
					b.WriteString("  " + stepOKStyle.Render("✓ "+historianLabel) + "\n")
				case "running":
					b.WriteString("  " + m.spinner.View() + " " + llmStyle.Render(historianLabel) + "\n")
				case "failed":
					b.WriteString("  " + stepFailStyle.Render("✗ "+historianLabel+" (failed)") + "\n")
				default: // pending
					b.WriteString("  " + logDimStyle.Render("· "+historianLabel) + "\n")
				}
			}

			// Pipeline steps
			for i, stepLine := range m.pipelineSteps {
				stepNum := i + 1
				if stepNum <= m.completedPipelineStep {
					b.WriteString("  " + stepOKStyle.Render("✓ "+stepLine) + "\n")
				} else if stepNum == m.currentStep {
					style := stepColors[(stepNum-1)%len(stepColors)]
					b.WriteString("  " + m.spinner.View() + " " + style.Render(stepLine) + "\n")
				} else {
					b.WriteString("  " + logDimStyle.Render("· "+stepLine) + "\n")
				}
			}
		} else {
			outlineLines = 0 // not shown, don't count
		}

		// Separator before panels/log
		b.WriteString("  " + logDimStyle.Render(strings.Repeat("─", 50)) + "\n")

		// Panels / activity log — reserve 2 lines for bottom progress bar
		availHeight := m.height - 10 - outlineLines
		if showOutline {
			availHeight -= 2 // separators around outline
		}
		if availHeight < 6 {
			availHeight = 6
		}

		if cw >= splitPanelMinWidth && m.statusContent != "" {
			leftW := cw * 2 / 5
			rightW := cw - leftW - 1

			left := renderStatusPanel(m.statusContent, leftW, availHeight)
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

		// Progress bar — pinned at bottom, always visible
		pct := float64(m.stepsCompleted) / float64(m.totalSteps)
		if pct > 1 {
			pct = 1
		}
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("  %s %d/%d steps",
			m.progress.ViewAs(pct),
			m.stepsCompleted, m.totalSteps))

	case discStateDone:
		b.WriteString(m.viewport.View())
	}

	return b.String()
}
