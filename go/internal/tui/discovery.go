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

// vernStatus tracks an individual Vern's async state in the VernHole phase.
type vernStatus struct {
	num    int    // 1-based index
	desc   string // display name from log line (e.g. "Ketamine Vern")
	llm    string // e.g. "claude"
	status string // "summoned", "ok", "failed"
}

type discoveryState int

const (
	discStateSetupForm discoveryState = iota
	discStatePathForm
	discStateProjectSelect
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

// projectInfo describes an existing discovery project for the rerun selector.
type projectInfo struct {
	name    string
	path    string
	modTime time.Time
}

// DiscoveryModel handles the discovery pipeline wizard.
type DiscoveryModel struct {
	state       discoveryState
	rerun       bool
	setupForm   *huh.Form
	pathForm    *huh.Form
	projectForm *huh.Form
	configForm  *huh.Form
	spinner     spinner.Model
	progress    progress.Model
	viewport    viewport.Model
	celebration CelebrationModel
	projectRoot string
	agentsDir   string
	width       int
	height      int
	vals        *discoveryVals
	projects    []projectInfo
	selectedProject string

	// Execution state
	running               bool
	stepLog               []string
	pipelineSteps         []string // step definitions captured from "[step] N. Name → llm"
	stepsCompleted        int
	totalSteps            int
	historianPhase        string // "pending", "running", "done", "failed", "skipped"
	currentStep           int    // which pipeline step number is currently running
	completedPipelineStep int    // highest completed pipeline step number
	statusContent         string // pipeline-status.md content (read at completion for done view)
	statusMsg             string // transient feedback (e.g. "Copied to clipboard!")
	err                   error
	setupErr              error

	// Phase tracking — screen redraws per phase
	runningPhase  string // "pipeline", "vernhole", "oracle"
	phaseLogStart int    // index into stepLog where current phase started

	// VernHole phase tracking (async — all Verns run in parallel)
	vernRoster     map[int]*vernStatus // keyed by Vern number (1-based)
	vernsCompleted int
	totalVerns     int

	// Oracle phase tracking
	oracleStep string // "consult", "apply"
}

func NewDiscoveryModel(projectRoot, agentsDir string) DiscoveryModel {
	cfg := config.Load(projectRoot)

	outputPath := "default"
	customPath := ""
	if cfg.DefaultDiscoveryPath != "" {
		outputPath = "custom"
		customPath = cfg.DefaultDiscoveryPath
	}

	vals := &discoveryVals{
		outputPath:   outputPath,
		customPath:   customPath,
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

// NewRerunDiscoveryModel creates a DiscoveryModel in rerun mode that starts
// with a path selection form, then shows a project selector.
func NewRerunDiscoveryModel(projectRoot, agentsDir string) DiscoveryModel {
	cfg := config.Load(projectRoot)

	outputPath := "default"
	customPath := ""
	if cfg.DefaultDiscoveryPath != "" {
		outputPath = "custom"
		customPath = cfg.DefaultDiscoveryPath
	}

	vals := &discoveryVals{
		outputPath:   outputPath,
		customPath:   customPath,
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
		state:       discStatePathForm,
		rerun:       true,
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

	m.pathForm = m.buildPathForm()
	m.configForm = m.buildConfigForm()

	return m
}

// scanProjects finds existing discovery projects by looking for directories
// containing input/prompt.md under the given base path.
func scanProjects(base string) []projectInfo {
	base = expandHome(base)
	var projects []projectInfo

	entries, err := os.ReadDir(base)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		projPath := filepath.Join(base, e.Name())
		promptPath := filepath.Join(projPath, "input", "prompt.md")
		info, err := os.Stat(promptPath)
		if err != nil {
			continue
		}
		projects = append(projects, projectInfo{
			name:    e.Name(),
			path:    projPath,
			modTime: info.ModTime(),
		})
	}

	// Sort newest first
	for i := 0; i < len(projects); i++ {
		for j := i + 1; j < len(projects); j++ {
			if projects[j].modTime.After(projects[i].modTime) {
				projects[i], projects[j] = projects[j], projects[i]
			}
		}
	}

	return projects
}

// buildPathForm creates a form for selecting the discovery directory to scan.
func (m *DiscoveryModel) buildPathForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Discovery folder location").
				Description("Where are your discovery projects?").
				Options(OutputPathOptions...).
				Value(&v.outputPath),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Path to discovery folder").
				Placeholder("~/discovery").
				Value(&v.customPath),
		).WithHideFunc(func() bool { return v.outputPath != "custom" }),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

// buildProjectForm creates a huh.Select form for choosing an existing project.
func (m *DiscoveryModel) buildProjectForm() *huh.Form {
	if len(m.projects) == 0 {
		return nil
	}

	w := contentWidth(m.width)
	options := make([]huh.Option[string], len(m.projects))
	for i, p := range m.projects {
		label := fmt.Sprintf("%s  (%s)", p.name, p.modTime.Format("Jan 2 15:04"))
		options[i] = huh.NewOption(label, p.path)
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a project to rerun").
				Description("Choose an existing discovery project").
				Options(options...).
				Height(min(len(options)+1, 15)).
				Value(&m.selectedProject),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

// cleanProjectOutput removes previous pipeline output from a discovery project
// so a rerun starts fresh. Preserves the input/ directory (prompt.md, input-history.md, etc.).
func cleanProjectOutput(dir string) error {
	for _, sub := range []string{"output", "vernhole"} {
		p := filepath.Join(dir, sub)
		if _, err := os.Stat(p); err == nil {
			if err := os.RemoveAll(p); err != nil {
				return fmt.Errorf("remove %s: %w", sub, err)
			}
		}
	}
	oraclePath := filepath.Join(dir, "oracle-vision.md")
	if _, err := os.Stat(oraclePath); err == nil {
		if err := os.Remove(oraclePath); err != nil {
			return fmt.Errorf("remove oracle-vision.md: %w", err)
		}
	}
	return nil
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
	if m.pathForm != nil {
		m.pathForm.WithWidth(cw).WithHeight(fh)
	}
	if m.projectForm != nil {
		m.projectForm.WithWidth(cw).WithHeight(fh)
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
	if m.rerun && m.pathForm != nil {
		return m.pathForm.Init()
	}
	if m.setupForm != nil {
		return m.setupForm.Init()
	}
	return nil
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

	case celebrateTickMsg:
		cmd := m.celebration.Update(msg)
		return m, cmd

	case pipelineDoneMsg:
		m.state = discStateDone
		m.running = false
		m.err = msg.err
		m.readStatus()
		var celebCmd tea.Cmd
		if msg.err == nil {
			celebCmd = m.celebration.Start("pipeline", m.width)
		}
		m.initDoneViewport()
		return m, tea.Batch(tea.DisableMouse, celebCmd)

	case pipelineLogMsg:
		line := msg.line

		// Capture step list lines for the outline (don't show in activity log)
		if strings.HasPrefix(line, "[step] ") {
			m.pipelineSteps = append(m.pipelineSteps, strings.TrimPrefix(line, "[step] "))
			return m, m.waitForLog()
		}

		upper := strings.ToUpper(line)

		// Detect phase transitions (before filtering, so state updates even for filtered lines)
		if strings.Contains(upper, "ENTERING THE VERNHOLE") {
			m.runningPhase = "vernhole"
			m.phaseLogStart = len(m.stepLog)
		} else if strings.Contains(upper, "CONSULTING THE ORACLE") {
			// VernHole just finished — fire ephemeral celebration
			m.celebration.StartEphemeral("vernhole", m.width)
			m.runningPhase = "oracle"
			m.phaseLogStart = len(m.stepLog)
			m.oracleStep = "consult"
		} else if strings.Contains(upper, "ORACLE APPLYING VISION") {
			m.oracleStep = "apply"
		}

		// Filter banner/noise lines (after phase detection)
		if isFilteredLogLine(line) {
			return m, m.waitForLog()
		}

		// Pipeline phase: track historian and step progress
		if m.runningPhase == "pipeline" {
			if strings.Contains(upper, "INVOKING THE ANCIENT SECRETS OF THE HISTORIAN") {
				m.historianPhase = "running"
			} else if strings.Contains(upper, "HISTORIAN COMPLETE") || strings.Contains(upper, "HISTORIAN INDEX ALREADY EXISTS") {
				m.historianPhase = "done"
				m.celebration.StartEphemeral("historian", m.width)
			} else if strings.Contains(upper, "HISTORIAN: PROMPT ONLY") || strings.Contains(upper, "NOTHING TO INDEX") {
				m.historianPhase = "done"
			} else if strings.Contains(upper, "HISTORIAN FAILED") {
				m.historianPhase = "failed"
			}

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
		}

		// VernHole phase: track Vern personas (async — all run in parallel)
		if m.runningPhase == "vernhole" {
			if strings.HasPrefix(line, ">>> Vern ") {
				// Parse ">>> Vern N/Total: Description (llm)"
				rest := strings.TrimPrefix(line, ">>> Vern ")
				if num, total, desc, llmName, ok := parseVernLine(rest); ok {
					if m.vernRoster == nil {
						m.vernRoster = make(map[int]*vernStatus)
					}
					m.totalVerns = total
					m.vernRoster[num] = &vernStatus{
						num:    num,
						desc:   desc,
						llm:    llmName,
						status: "summoned",
					}
				}
			} else if strings.Contains(upper, "VERN ") && (strings.Contains(upper, "OK (") || strings.Contains(upper, "FAILED (")) {
				// Parse "    OK (llm, NB, Vern N/Total)" or "    FAILED (id, exit N, Vern N/Total)"
				if num := parseVernResultLine(line); num > 0 {
					if vs, ok := m.vernRoster[num]; ok {
						if strings.Contains(upper, "OK (") {
							vs.status = "ok"
						} else {
							vs.status = "failed"
						}
					}
				}
			}
		}

		m.stepLog = append(m.stepLog, line)
		m.updateProgress(line)
		return m, m.waitForLog()

	case discStatusClearMsg:
		m.statusMsg = ""
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	switch m.state {
	case discStatePathForm:
		form, cmd := m.pathForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.pathForm = f
		}
		if m.pathForm.State == huh.StateCompleted {
			m.projects = scanProjects(m.outputBase())
			m.projectForm = m.buildProjectForm()
			m.state = discStateProjectSelect
			if m.projectForm != nil {
				return m, m.projectForm.Init()
			}
			return m, nil
		}
		if m.pathForm.State == huh.StateAborted {
			return m, backToMenu
		}
		return m, cmd

	case discStateProjectSelect:
		if m.projectForm == nil {
			// No projects found — esc goes back
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				if keyMsg.String() == "esc" || keyMsg.String() == "enter" {
					return m, backToMenu
				}
			}
			return m, nil
		}
		form, cmd := m.projectForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.projectForm = f
		}
		if m.projectForm.State == huh.StateCompleted {
			// Load project info from selected path
			dir := m.selectedProject
			m.vals.name = filepath.Base(dir)
			promptPath := filepath.Join(dir, "input", "prompt.md")
			if data, err := os.ReadFile(promptPath); err == nil {
				m.vals.idea = string(data)
			}
			// Determine output path settings
			parent := filepath.Dir(dir)
			if absParent, err := filepath.Abs(parent); err == nil {
				absCwd, _ := filepath.Abs("./discovery")
				if absParent == absCwd {
					m.vals.outputPath = "default"
				} else {
					m.vals.outputPath = "custom"
					m.vals.customPath = parent
				}
			}
			m.state = discStateEditFiles
			return m, nil
		}
		if m.projectForm.State == huh.StateAborted {
			return m, backToMenu
		}
		return m, cmd

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
				// For rerun, re-read prompt.md in case user edited it externally
				if m.rerun {
					promptPath := filepath.Join(m.discoveryDir(), "input", "prompt.md")
					if data, err := os.ReadFile(promptPath); err == nil {
						m.vals.idea = string(data)
					}
				}
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
			if m.rerun {
				if err := cleanProjectOutput(m.discoveryDir()); err != nil {
					m.err = err
					m.state = discStateDone
					m.initDoneViewport()
					return m, tea.DisableMouse
				}
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
			m.runningPhase = "pipeline"
			m.vals.logCh = make(chan string, 100)
			return m, tea.Batch(m.spinner.Tick, m.startPipeline(), m.waitForLog())
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
			switch keyMsg.String() {
			case "q":
				return m, backToMenu
			case "c":
				text := m.resultContent()
				if text != "" {
					if err := copyToClipboard(text); err != nil {
						m.statusMsg = "Copy failed: " + err.Error()
					} else {
						m.statusMsg = "Copied to clipboard!"
					}
					return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
						return discStatusClearMsg{}
					})
				}
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

	switch m.runningPhase {
	case "vernhole":
		// Count completed Verns from roster (handles async completion)
		completed := 0
		for _, vs := range m.vernRoster {
			if vs.status == "ok" || vs.status == "failed" {
				completed++
			}
		}
		m.vernsCompleted = completed
	case "oracle":
		// Oracle has no granular progress
	default: // pipeline
		if strings.Contains(upper, "OK (") || strings.Contains(upper, "FALLBACK SUCCEEDED") {
			m.stepsCompleted++
			m.completedPipelineStep = m.currentStep
		} else if strings.Contains(upper, "HISTORIAN COMPLETE") ||
			strings.Contains(upper, "HISTORIAN INDEX ALREADY EXISTS") ||
			strings.Contains(upper, "HISTORIAN: PROMPT ONLY") ||
			strings.Contains(upper, "NOTHING TO INDEX") ||
			strings.Contains(upper, "HISTORIAN FAILED") {
			m.stepsCompleted++
		}
	}
}

func (m *DiscoveryModel) initDoneViewport() {
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
		content.WriteString(stepFailStyle.Render("Pipeline failed: " + m.err.Error()))
		content.WriteString("\n\n")
	} else {
		content.WriteString(fmt.Sprintf("Output: %s\n", m.discoveryDir()))
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

// resultContent builds the full pipeline output as plain text for clipboard copy.
func (m DiscoveryModel) resultContent() string {
	var b strings.Builder

	b.WriteString("=== DISCOVERY PIPELINE RESULTS ===\n")
	b.WriteString(fmt.Sprintf("Output: %s\n\n", m.discoveryDir()))

	if m.statusContent != "" {
		b.WriteString(stripMarkdown(m.statusContent))
		b.WriteString("\n")
	}

	if len(m.stepLog) > 0 {
		b.WriteString("--- Activity Log ---\n")
		for _, line := range m.stepLog {
			b.WriteString(line + "\n")
		}
	}

	return b.String()
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

	content := fmt.Sprintf("# %s\n\n## Prompt\n\n%s\n\n## Additional Context\n<!-- Add any extra context, constraints, or goals. -->\n", m.vals.name, m.vals.idea)
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

type discStatusClearMsg struct{}

func (m DiscoveryModel) startPipeline() tea.Cmd {
	return func() tea.Msg {
		v := m.vals
		defer close(v.logCh)

		ctx, cancel := context.WithCancel(context.Background())
		v.cancel = cancel
		defer cancel()

		dir := m.discoveryDir()

		singleLLM := ""
		if v.llmMode == "single_llm" {
			singleLLM = v.singleLLM
		}

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
			SingleLLM:     singleLLM,
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

func (m *DiscoveryModel) readStatus() {
	statusPath := filepath.Join(m.discoveryDir(), "output", "pipeline-status.md")
	data, err := os.ReadFile(statusPath)
	if err == nil {
		m.statusContent = string(data)
	}
}

// renderStepOutlinePanel renders the pipeline step outline in a bordered panel.
func (m DiscoveryModel) renderStepOutlinePanel(width, height int) string {
	var b strings.Builder

	// Historian line
	if m.historianPhase != "skipped" {
		historianLabel := "Historian (pre-step)"
		switch m.historianPhase {
		case "done":
			b.WriteString(stepOKStyle.Render("✓ "+historianLabel) + "\n")
		case "running":
			b.WriteString(m.spinner.View() + " " + llmStyle.Render(historianLabel) + "\n")
		case "failed":
			b.WriteString(stepFailStyle.Render("✗ "+historianLabel+" (failed)") + "\n")
		default:
			b.WriteString(logDimStyle.Render("· "+historianLabel) + "\n")
		}
	}

	// Pipeline steps
	for i, stepLine := range m.pipelineSteps {
		stepNum := i + 1
		if stepNum <= m.completedPipelineStep {
			b.WriteString(stepOKStyle.Render("✓ "+stepLine) + "\n")
		} else if stepNum == m.currentStep {
			style := stepColors[(stepNum-1)%len(stepColors)]
			b.WriteString(m.spinner.View() + " " + style.Render(stepLine) + "\n")
		} else {
			b.WriteString(logDimStyle.Render("· "+stepLine) + "\n")
		}
	}

	return statusPanelStyle.Width(width).Height(height).Render(
		panelTitleStyle.Render("Pipeline Steps") + "\n" + b.String(),
	)
}

// renderVernOutlinePanel renders VernHole persona status in a bordered panel.
// Verns run in parallel, so status is independent per-Vern (not sequential).
func (m DiscoveryModel) renderVernOutlinePanel(width, height int) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Council:  %s\n", llmStyle.Render(councilLabel(m.vals.vernhole))))
	b.WriteString(fmt.Sprintf("LLM Mode: %s\n", llmStyle.Render(m.vals.llmMode)))
	if m.totalVerns > 0 {
		b.WriteString(fmt.Sprintf("Progress: %s\n", llmStyle.Render(fmt.Sprintf("%d/%d complete", m.vernsCompleted, m.totalVerns))))
	}
	b.WriteString("\n")

	// Show Verns sorted by number, with per-Vern status
	for i := 1; i <= m.totalVerns; i++ {
		vs, ok := m.vernRoster[i]
		if !ok {
			continue
		}
		label := fmt.Sprintf("%s (%s)", vs.desc, vs.llm)
		switch vs.status {
		case "ok":
			b.WriteString(stepOKStyle.Render("✓ "+label) + "\n")
		case "failed":
			b.WriteString(stepFailStyle.Render("✗ "+label) + "\n")
		default: // "summoned" — running in parallel
			b.WriteString(logDimStyle.Render(m.spinner.View()+" "+label) + "\n")
		}
	}

	return statusPanelStyle.Width(width).Height(height).Render(
		panelTitleStyle.Render("VernHole Council") + "\n" + b.String(),
	)
}

// renderOracleOutlinePanel renders Oracle phase status in a bordered panel.
func (m DiscoveryModel) renderOracleOutlinePanel(width, height int) string {
	var b strings.Builder

	consultLabel := "Consult the Oracle"
	applyLabel := "Apply via Architect"

	switch m.oracleStep {
	case "apply":
		b.WriteString(stepOKStyle.Render("✓ "+consultLabel) + "\n")
		b.WriteString(m.spinner.View() + " " + llmStyle.Render(applyLabel) + "\n")
	case "consult":
		b.WriteString(m.spinner.View() + " " + llmStyle.Render(consultLabel) + "\n")
		if m.vals.oracleApply == "apply" {
			b.WriteString(logDimStyle.Render("· "+applyLabel) + "\n")
		}
	default:
		b.WriteString(logDimStyle.Render("· "+consultLabel) + "\n")
		if m.vals.oracleApply == "apply" {
			b.WriteString(logDimStyle.Render("· "+applyLabel) + "\n")
		}
	}

	return statusPanelStyle.Width(width).Height(height).Render(
		panelTitleStyle.Render("Oracle Vision") + "\n" + b.String(),
	)
}

// parseVernLine parses "N/Total: Description (llm)" from a ">>> Vern " line.
func parseVernLine(rest string) (num, total int, desc, llmName string, ok bool) {
	// Format: "3/15: Needs 6 meetings and a committee first (claude)"
	slashIdx := strings.Index(rest, "/")
	if slashIdx < 1 {
		return
	}
	colonIdx := strings.Index(rest, ":")
	if colonIdx < slashIdx {
		return
	}
	fmt.Sscanf(rest[:slashIdx], "%d", &num)
	fmt.Sscanf(rest[slashIdx+1:colonIdx], "%d", &total)
	if num == 0 || total == 0 {
		return
	}

	afterColon := strings.TrimSpace(rest[colonIdx+1:])
	// Extract LLM from trailing "(llm)"
	if parenIdx := strings.LastIndex(afterColon, "("); parenIdx > 0 {
		desc = strings.TrimSpace(afterColon[:parenIdx])
		llmName = strings.TrimSuffix(strings.TrimSpace(afterColon[parenIdx+1:]), ")")
	} else {
		desc = afterColon
	}
	ok = true
	return
}

// parseVernResultLine extracts the Vern number from an OK or FAILED line.
// Format: "    OK (llm, NB, Vern N/Total)" or "    FAILED (id, exit N, Vern N/Total)"
func parseVernResultLine(line string) int {
	// Look for "Vern N/" pattern
	idx := strings.Index(line, "Vern ")
	if idx < 0 {
		return 0
	}
	rest := line[idx+5:] // after "Vern "
	slashIdx := strings.Index(rest, "/")
	if slashIdx < 1 {
		return 0
	}
	var num int
	fmt.Sscanf(rest[:slashIdx], "%d", &num)
	return num
}

// Cancel aborts any running pipeline goroutine.
func (m DiscoveryModel) Cancel() {
	if m.vals != nil && m.vals.cancel != nil {
		m.vals.cancel()
	}
}

func (m DiscoveryModel) View() string {
	var b strings.Builder

	if m.rerun {
		b.WriteString(titleStyle.Render("Rerun Discovery"))
	} else {
		b.WriteString(titleStyle.Render("Discovery Pipeline"))
	}
	b.WriteString("\n\n")

	switch m.state {
	case discStatePathForm:
		b.WriteString(m.pathForm.View())

	case discStateProjectSelect:
		if m.projectForm == nil {
			b.WriteString(stepWarningStyle.Render("No discovery projects found in: "+m.outputBase()))
			b.WriteString("\n\n")
			b.WriteString("Create a project first using Discovery Pipeline.\n\n")
			b.WriteString(subtitleStyle.Render("Press Esc to go back"))
		} else {
			b.WriteString(m.projectForm.View())
		}

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
		} else if m.rerun {
			b.WriteString(stepOKStyle.Render("Rerunning project: "+m.vals.name))
			b.WriteString("\n\n")
			b.WriteString(fmt.Sprintf("  %s  %s\n", logHeaderStyle.Render("Path:"), dir))
			b.WriteString(fmt.Sprintf("  %s  %s\n", logHeaderStyle.Render("Prompt:"), promptPath))
			b.WriteString("\n")

			// Show prompt preview (truncated for display)
			preview := m.vals.idea
			maxPreview := 500
			if len(preview) > maxPreview {
				preview = preview[:maxPreview] + "\n..."
			}
			b.WriteString(logDimStyle.Render(preview))
			b.WriteString("\n\n")
			b.WriteString(stepWarningStyle.Render("Previous output will be cleared on run."))
			b.WriteString("\n")
			b.WriteString("You can edit prompt.md externally before pressing Enter.\n\n")
			b.WriteString(subtitleStyle.Render("Press Enter when ready to configure"))
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
		cw := contentWidth(m.width)

		// Phase-specific title
		switch m.runningPhase {
		case "vernhole":
			b.WriteString(fmt.Sprintf("%s %s\n\n", m.spinner.View(), logHeaderStyle.Render("Summoning the VernHole council...")))
		case "oracle":
			if m.oracleStep == "apply" {
				b.WriteString(fmt.Sprintf("%s %s\n\n", m.spinner.View(), logHeaderStyle.Render("Architect applying the Oracle's vision...")))
			} else {
				b.WriteString(fmt.Sprintf("%s %s\n\n", m.spinner.View(), logHeaderStyle.Render("Consulting the Oracle...")))
			}
		default:
			b.WriteString(fmt.Sprintf("%s %s\n\n", m.spinner.View(), logHeaderStyle.Render("Running discovery pipeline...")))
		}

		// Config line
		label := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render
		mode := "default"
		if m.vals.pipeline == "expanded" {
			mode = "expanded"
		}
		b.WriteString(fmt.Sprintf("  %s  %s  |  %s  %s  |  %s  %s\n",
			label("Pipeline:"), llmStyle.Render(mode),
			label("LLM:"), llmStyle.Render(m.vals.llmMode),
			label("Output:"), logDimStyle.Render(m.discoveryDir())))
		b.WriteString("  " + logDimStyle.Render(strings.Repeat("─", 50)) + "\n")

		// Ephemeral celebration (inline during running)
		if cv := m.celebration.View(); cv != "" {
			b.WriteString(cv)
			b.WriteString("\n")
		}

		// Chrome: title(2) + spinner(2) + config(1) + sep(1) + progress(2) + help(2) = 10
		totalAvail := m.height - 10 - m.celebration.Height()
		if totalAvail < 10 {
			totalAvail = 10
		}

		// Current phase's log (from phase start)
		phaseLog := m.stepLog[m.phaseLogStart:]

		if cw >= splitPanelMinWidth {
			// Split-panel layout: left=phase status panel, right=activity log
			leftW := cw * 2 / 5
			rightW := cw - leftW - 1

			var left string
			switch m.runningPhase {
			case "vernhole":
				left = m.renderVernOutlinePanel(leftW, totalAvail)
			case "oracle":
				left = m.renderOracleOutlinePanel(leftW, totalAvail)
			default:
				left = m.renderStepOutlinePanel(leftW, totalAvail)
			}
			right := renderLogPanel(phaseLog, rightW, totalAvail)

			b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, left, right))
		} else {
			// Narrow terminal: single column — just activity log
			maxLines := totalAvail
			if maxLines < 4 {
				maxLines = 4
			}
			start := 0
			if len(phaseLog) > maxLines {
				start = len(phaseLog) - maxLines
			}
			for _, line := range phaseLog[start:] {
				b.WriteString("  " + renderLogLine(line) + "\n")
			}
		}

		// Phase-specific progress bar
		b.WriteString("\n")
		switch m.runningPhase {
		case "vernhole":
			if m.totalVerns > 0 {
				pct := float64(m.vernsCompleted) / float64(m.totalVerns)
				if pct > 1 {
					pct = 1
				}
				b.WriteString(fmt.Sprintf("  %s %d/%d Verns",
					m.progress.ViewAs(pct),
					m.vernsCompleted, m.totalVerns))
			} else {
				b.WriteString(fmt.Sprintf("  %s Summoning council...", m.spinner.View()))
			}
		case "oracle":
			if m.oracleStep == "apply" {
				b.WriteString(fmt.Sprintf("  %s Architect applying changes...", m.spinner.View()))
			} else {
				b.WriteString(fmt.Sprintf("  %s Oracle working...", m.spinner.View()))
			}
		default:
			pct := float64(m.stepsCompleted) / float64(m.totalSteps)
			if pct > 1 {
				pct = 1
			}
			b.WriteString(fmt.Sprintf("  %s %d/%d steps",
				m.progress.ViewAs(pct),
				m.stepsCompleted, m.totalSteps))
		}

	case discStateDone:
		if cv := m.celebration.View(); cv != "" {
			b.WriteString(cv)
			b.WriteString("\n")
		}
		if m.statusMsg != "" {
			b.WriteString(stepOKStyle.Render(m.statusMsg) + "\n")
		}
		b.WriteString(m.viewport.View())
	}

	return b.String()
}
