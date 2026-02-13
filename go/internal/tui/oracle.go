package tui

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	oracleStateForm          oracleState = iota // operation picker only
	oracleStatePathForm                         // pick base directory
	oracleStateProjectSelect                    // scan + pick project
	oracleStateConfigForm                       // operation-specific config + confirm
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
	operation   string
	synthDir    string
	vtsDir      string
	visionFile  string
	contextFile string
	idea        string
	council     string
	llmMode     string
	singleLLM   string
	confirm     bool
	cancel      context.CancelFunc
	logCh       chan string

	// Path picker state
	outputPath string // "default" or "custom"
	customPath string // custom path value

	// Selected project path (set by project selector)
	projectPath string

	// VernHole output dir (derived from project)
	vernholeDir string
}

// OracleModel handles the Oracle wizard.
type OracleModel struct {
	state         oracleState
	operationForm *huh.Form
	pathForm      *huh.Form
	projectForm   *huh.Form
	configForm    *huh.Form
	spinner       spinner.Model
	progress      progress.Model
	viewport      viewport.Model
	projectRoot   string
	agentsDir     string
	width         int
	height        int
	vals          *oracleVals

	// Project browser state
	projects        []projectInfo
	selectedProject string

	// Execution state
	running        bool
	stepLog        []string
	vernsCompleted int
	totalVerns     int
	err            error
	statusMsg      string
}

func NewOracleModel(projectRoot, agentsDir string) OracleModel {
	cfg := config.Load(projectRoot)

	outputPath := "default"
	customPath := ""
	if cfg.DefaultDiscoveryPath != "" {
		outputPath = "custom"
		customPath = cfg.DefaultDiscoveryPath
	}

	vals := &oracleVals{
		operation:  "consult",
		council:    "full",
		llmMode:    "mixed_claude_fallback",
		singleLLM:  "claude",
		confirm:    true,
		outputPath: outputPath,
		customPath: customPath,
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

	m.operationForm = m.buildOperationForm()

	return m
}

// SetSize updates the terminal dimensions and resizes forms.
func (m *OracleModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	cw := contentWidth(w)
	fh := formHeight(h)
	if m.operationForm != nil {
		m.operationForm.WithWidth(cw).WithHeight(fh)
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

// buildOperationForm creates a form with just the operation picker.
func (m *OracleModel) buildOperationForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Options(OracleOperationOptions...).
				Height(len(OracleOperationOptions)+1).
				Value(&v.operation),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

// buildOraclePathForm creates a form for selecting the base directory to scan.
func (m *OracleModel) buildOraclePathForm() *huh.Form {
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

// buildOracleProjectForm creates a select form for choosing a project.
func (m *OracleModel) buildOracleProjectForm() *huh.Form {
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
				Title("Select a project").
				Description(m.projectSelectDescription()).
				Options(options...).
				Height(min(len(options)+1, 15)).
				Value(&m.selectedProject),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

// projectSelectDescription returns an operation-specific description for the project selector.
func (m *OracleModel) projectSelectDescription() string {
	switch m.vals.operation {
	case "consult":
		return "Projects with vernhole/synthesis.md"
	case "apply":
		return "Projects with oracle-vision.md"
	case "vernhole":
		return "Projects with discovery output"
	}
	return "Select a project"
}

// buildOracleConfigForm creates an operation-specific config + confirm form.
func (m *OracleModel) buildOracleConfigForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals

	var groups []*huh.Group

	// Auto-detected paths note
	groups = append(groups, huh.NewGroup(
		huh.NewNote().
			Title("Auto-Detected Paths").
			Description(m.autoDetectedSummary()),
	))

	// Operation-specific fields
	switch v.operation {
	case "consult":
		lines := textareaLines(m.height)
		groups = append(groups, huh.NewGroup(
			huh.NewText().
				Title("Original idea / prompt").
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
		))
	case "vernhole":
		lines := textareaLines(m.height)
		groups = append(groups, huh.NewGroup(
			huh.NewText().
				Title("Original idea / prompt").
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
		))
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which council do you want to summon?").
				Options(CouncilOptions...).
				Height(len(CouncilOptions)+1).
				Value(&v.council),
		))
	}

	// LLM mode (all operations)
	groups = append(groups, huh.NewGroup(
		huh.NewSelect[string]().
			Title("LLM Mode").
			Options(LLMModeOptions...).
			Height(len(LLMModeOptions)+1).
			Value(&v.llmMode),
	))
	groups = append(groups, huh.NewGroup(
		huh.NewSelect[string]().
			Title("Which LLM should run all steps?").
			Options(SingleLLMOptions...).
			Height(len(SingleLLMOptions)+1).
			Value(&v.singleLLM),
	).WithHideFunc(func() bool { return v.llmMode != "single_llm" }))

	// Confirm
	groups = append(groups, huh.NewGroup(
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
	))

	return huh.NewForm(groups...).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

// autoDetectedSummary returns a formatted summary of auto-detected paths.
func (m *OracleModel) autoDetectedSummary() string {
	v := m.vals
	dim := lipgloss.NewStyle().Foreground(colorMuted).Italic(true).Render
	label := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render

	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("Project:"), dim(filepath.Base(v.projectPath))))

	switch v.operation {
	case "consult":
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("Synthesis:"), dim(v.synthDir)))
		if v.vtsDir != "" {
			b.WriteString(fmt.Sprintf("  %s  %s\n", label("VTS Dir:"), dim(v.vtsDir)))
		}
	case "apply":
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("Vision:"), dim(v.visionFile)))
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("VTS Dir:"), dim(v.vtsDir)))
	case "vernhole":
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("Context:"), dim(v.contextFile)))
		b.WriteString(fmt.Sprintf("  %s  %s\n", label("Output:"), dim(v.vernholeDir)))
	}

	return b.String()
}

// oracleOutputBase returns the base directory for project scanning.
func (m *OracleModel) oracleOutputBase() string {
	if m.vals.outputPath == "custom" && m.vals.customPath != "" {
		return expandHome(m.vals.customPath)
	}
	return "./discovery"
}

// scanOracleProjects scans the base directory for projects matching the
// operation's criteria. Returns sorted by most recently modified first.
func scanOracleProjects(base, operation string) []projectInfo {
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
		var keyTime time.Time
		eligible := false

		switch operation {
		case "consult":
			synthPath := filepath.Join(projPath, "vernhole", "synthesis.md")
			if info, err := os.Stat(synthPath); err == nil {
				eligible = true
				keyTime = info.ModTime()
			}
		case "apply":
			visionPath := filepath.Join(projPath, "oracle-vision.md")
			if info, err := os.Stat(visionPath); err == nil {
				eligible = true
				keyTime = info.ModTime()
			}
		case "vernhole":
			outputDir := filepath.Join(projPath, "output")
			outputEntries, err := os.ReadDir(outputDir)
			if err != nil {
				continue
			}
			for _, oe := range outputEntries {
				if !oe.IsDir() && strings.HasSuffix(oe.Name(), ".md") {
					eligible = true
					if info, err := oe.Info(); err == nil {
						if info.ModTime().After(keyTime) {
							keyTime = info.ModTime()
						}
					}
				}
			}
		}

		if eligible {
			projects = append(projects, projectInfo{
				name:    e.Name(),
				path:    projPath,
				modTime: keyTime,
			})
		}
	}

	// Sort newest first
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].modTime.After(projects[j].modTime)
	})

	return projects
}

// populateFromProject auto-fills vals fields from the selected project directory.
func (m *OracleModel) populateFromProject(projPath string) {
	v := m.vals
	v.projectPath = projPath

	// Try to read idea from prompt.md
	promptPath := filepath.Join(projPath, "input", "prompt.md")
	if data, err := os.ReadFile(promptPath); err == nil {
		v.idea = string(data)
	}

	switch v.operation {
	case "consult":
		v.synthDir = filepath.Join(projPath, "vernhole")
		v.vtsDir = filepath.Join(projPath, "output", "vts")
	case "apply":
		v.visionFile = filepath.Join(projPath, "oracle-vision.md")
		v.vtsDir = filepath.Join(projPath, "output", "vts")
	case "vernhole":
		v.contextFile = findConsolidationFile(filepath.Join(projPath, "output"))
		v.vernholeDir = filepath.Join(projPath, "vernhole")
	}
}

// findConsolidationFile looks for the best context file in a discovery output
// directory. Prefers files matching *consolidation*, then falls back to the
// highest-numbered step file.
func findConsolidationFile(outputDir string) string {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return ""
	}

	var consolidation string
	var highestStep string
	var highestNum int

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		lower := strings.ToLower(e.Name())
		if strings.Contains(lower, "consolidation") {
			consolidation = filepath.Join(outputDir, e.Name())
		}
		// Track highest-numbered file (format: NN-slug-step.md)
		if len(e.Name()) >= 2 && e.Name()[0] >= '0' && e.Name()[0] <= '9' {
			var num int
			fmt.Sscanf(e.Name()[:2], "%d", &num)
			if num > highestNum {
				highestNum = num
				highestStep = filepath.Join(outputDir, e.Name())
			}
		}
	}

	if consolidation != "" {
		return consolidation
	}
	return highestStep
}

func (m OracleModel) Init() tea.Cmd {
	return m.operationForm.Init()
}

func (m OracleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" && !m.running {
			switch m.state {
			case oracleStatePathForm:
				// Go back to operation picker
				m.state = oracleStateForm
				m.operationForm = m.buildOperationForm()
				return m, m.operationForm.Init()
			case oracleStateProjectSelect:
				// Go back to path form
				m.state = oracleStatePathForm
				m.pathForm = m.buildOraclePathForm()
				return m, m.pathForm.Init()
			case oracleStateConfigForm:
				// Go back to project select
				m.state = oracleStateProjectSelect
				m.projectForm = m.buildOracleProjectForm()
				if m.projectForm != nil {
					return m, m.projectForm.Init()
				}
				return m, nil
			default:
				return m, backToMenu
			}
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
		form, cmd := m.operationForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.operationForm = f
		}
		if m.operationForm.State == huh.StateCompleted {
			m.state = oracleStatePathForm
			m.pathForm = m.buildOraclePathForm()
			return m, m.pathForm.Init()
		}
		if m.operationForm.State == huh.StateAborted {
			return m, backToMenu
		}
		return m, cmd

	case oracleStatePathForm:
		form, cmd := m.pathForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.pathForm = f
		}
		if m.pathForm.State == huh.StateCompleted {
			base := m.oracleOutputBase()
			m.projects = scanOracleProjects(base, m.vals.operation)
			m.projectForm = m.buildOracleProjectForm()
			m.state = oracleStateProjectSelect
			if m.projectForm != nil {
				return m, m.projectForm.Init()
			}
			return m, nil
		}
		if m.pathForm.State == huh.StateAborted {
			// Go back to operation picker
			m.state = oracleStateForm
			m.operationForm = m.buildOperationForm()
			return m, m.operationForm.Init()
		}
		return m, cmd

	case oracleStateProjectSelect:
		if m.projectForm == nil {
			// No projects found
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				if keyMsg.String() == "esc" || keyMsg.String() == "enter" {
					m.state = oracleStatePathForm
					m.pathForm = m.buildOraclePathForm()
					return m, m.pathForm.Init()
				}
			}
			return m, nil
		}
		form, cmd := m.projectForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.projectForm = f
		}
		if m.projectForm.State == huh.StateCompleted {
			m.populateFromProject(m.selectedProject)
			m.configForm = m.buildOracleConfigForm()
			m.state = oracleStateConfigForm
			return m, m.configForm.Init()
		}
		if m.projectForm.State == huh.StateAborted {
			// Go back to path form
			m.state = oracleStatePathForm
			m.pathForm = m.buildOraclePathForm()
			return m, m.pathForm.Init()
		}
		return m, cmd

	case oracleStateConfigForm:
		form, cmd := m.configForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.configForm = f
		}
		if m.configForm.State == huh.StateCompleted {
			if !m.vals.confirm {
				return m, backToMenu
			}
			m.state = oracleStateRunning
			m.running = true
			m.vals.logCh = make(chan string, 100)
			return m, tea.Batch(m.spinner.Tick, m.startOracle(), m.waitForLog())
		}
		if m.configForm.State == huh.StateAborted {
			// Go back to project select
			m.state = oracleStateProjectSelect
			m.projectForm = m.buildOracleProjectForm()
			if m.projectForm != nil {
				return m, m.projectForm.Init()
			}
			return m, nil
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
		synthDir := v.vernholeDir
		if synthDir == "" {
			synthDir = "./vernhole/"
		}
		synthPath := filepath.Join(synthDir, "synthesis.md")
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
			vernDir := m.vals.vernholeDir
			if vernDir == "" {
				vernDir = "./vernhole/"
			}
			content.WriteString(stepOKStyle.Render("The VernHole has spoken!"))
			content.WriteString(fmt.Sprintf("\nFiles created in: %s\n\n", logDimStyle.Render(vernDir)))
			synthPath := filepath.Join(vernDir, "synthesis.md")
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
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("Project:"), dim(filepath.Base(v.projectPath))))

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
				Timeout:      cfg.GetOracleTimeout(),
				OnLog:        onLog,
			})

		case "apply":
			err = pipeline.RunOracleApply(pipeline.OracleApplyOptions{
				Ctx:          ctx,
				VisionFile:   expandHome(v.visionFile),
				VTSDir:       expandHome(v.vtsDir),
				AgentsDir:    m.agentsDir,
				SynthesisLLM: synthesisLLM,
				Timeout:      cfg.GetOracleApplyTimeout(),
				OnLog:        onLog,
			})

		case "vernhole":
			vernDir := v.vernholeDir
			if vernDir == "" {
				vernDir = "./vernhole/"
			}
			os.MkdirAll(vernDir, 0755)
			err = pipeline.RunVernHole(pipeline.VernHoleOptions{
				Ctx:          ctx,
				Idea:         v.idea,
				OutputDir:    vernDir,
				Council:      v.council,
				Context:      expandHome(v.contextFile),
				AgentsDir:    m.agentsDir,
				Timeout:      cfg.GetPipelineStepTimeout(),
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
		b.WriteString(m.operationForm.View())

	case oracleStatePathForm:
		b.WriteString(m.pathForm.View())

	case oracleStateProjectSelect:
		if m.projectForm == nil {
			opLabel := ""
			switch m.vals.operation {
			case "consult":
				opLabel = "No projects with vernhole/synthesis.md found in:"
			case "apply":
				opLabel = "No projects with oracle-vision.md found in:"
			case "vernhole":
				opLabel = "No projects with discovery output found in:"
			}
			b.WriteString(stepWarningStyle.Render(opLabel))
			b.WriteString("\n")
			b.WriteString(logDimStyle.Render(m.oracleOutputBase()))
			b.WriteString("\n\n")
			b.WriteString(subtitleStyle.Render("Press Esc to go back"))
		} else {
			b.WriteString(m.projectForm.View())
		}

	case oracleStateConfigForm:
		b.WriteString(m.configForm.View())

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
