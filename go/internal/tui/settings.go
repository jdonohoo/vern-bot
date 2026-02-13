package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/jdonohoo/vern-bot/go/internal/config"
)

type settingsState int

const (
	settingsStateMenu settingsState = iota
	settingsStateForm
	settingsStateSaveErr
)

type settingsAction int

const (
	settingsActionMode settingsAction = iota
	settingsActionLLMs
	settingsActionPipeline
	settingsActionDiscoveryPath
	settingsActionTimeouts
)

// settingsVals holds form-bound values on the heap so pointers survive
// bubbletea's value-copy semantics.
type settingsVals struct {
	llmMode            string
	singleLLM          string
	enabledLLMs        []string
	pipelineMode       string
	discoveryPath      string
	timeoutPipeline    string
	timeoutHistorian   string
	timeoutOracle      string
	timeoutOracleApply string
}

// SettingsModel handles the settings screen.
type SettingsModel struct {
	state       settingsState
	action      settingsAction
	cursor      int
	form        *huh.Form
	projectRoot string
	cfg         *config.Config
	width       int
	height      int
	vals        *settingsVals

	// Save result
	saveErr error
}

func NewSettingsModel(projectRoot string) SettingsModel {
	cfg := config.Load(projectRoot)
	vals := &settingsVals{
		llmMode: cfg.LLMMode,
	}
	return SettingsModel{
		state:       settingsStateMenu,
		projectRoot: projectRoot,
		cfg:         cfg,
		vals:        vals,
	}
}

// SetSize updates the terminal dimensions and resizes any active form.
func (m *SettingsModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	cw := contentWidth(w)
	fh := formHeight(h)
	if m.form != nil {
		m.form.WithWidth(cw).WithHeight(fh)
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

var settingsMenuItems = []string{
	"LLM Mode",
	"LLM Availability",
	"Pipeline Mode",
	"Default Discovery Folder",
	"Timeouts",
	"Back",
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "esc" {
			if m.state == settingsStateMenu {
				return m, backToMenu
			}
			m.state = settingsStateMenu
			m.cursor = 0
			return m, nil
		}
	}

	switch m.state {
	case settingsStateMenu:
		return m.updateMenu(msg)

	case settingsStateForm:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted {
			m.applyFormResult()
			m.saveErr = m.saveConfig()
			if m.saveErr != nil {
				m.state = settingsStateSaveErr
				return m, nil
			}
			m.state = settingsStateMenu
			m.cursor = 0
			return m, nil
		}
		if m.form.State == huh.StateAborted {
			m.state = settingsStateMenu
			m.cursor = 0
			return m, nil
		}
		return m, cmd

	case settingsStateSaveErr:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "enter" {
				m.state = settingsStateMenu
				m.cursor = 0
			}
		}
	}

	return m, nil
}

func (m SettingsModel) updateMenu(msg tea.Msg) (SettingsModel, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(settingsMenuItems)-1 {
				m.cursor++
			}
		case "enter":
			return m.executeAction(m.cursor)
		case "1":
			return m.executeAction(0)
		case "2":
			return m.executeAction(1)
		case "3":
			return m.executeAction(2)
		case "4":
			return m.executeAction(3)
		case "5":
			return m.executeAction(4)
		case "q":
			return m, backToMenu
		}
	}
	return m, nil
}

func (m SettingsModel) executeAction(index int) (SettingsModel, tea.Cmd) {
	v := m.vals
	switch index {
	case 0:
		m.action = settingsActionMode
		v.llmMode = m.cfg.LLMMode
		v.singleLLM = "claude"
		m.form = m.buildModeForm()
		m.state = settingsStateForm
		return m, m.form.Init()
	case 1:
		m.action = settingsActionLLMs
		v.enabledLLMs = m.currentEnabledLLMs()
		m.form = m.buildLLMToggleForm()
		m.state = settingsStateForm
		return m, m.form.Init()
	case 2:
		m.action = settingsActionPipeline
		v.pipelineMode = m.cfg.PipelineMode
		if v.pipelineMode == "" {
			v.pipelineMode = "default"
		}
		m.form = m.buildPipelineForm()
		m.state = settingsStateForm
		return m, m.form.Init()
	case 3:
		m.action = settingsActionDiscoveryPath
		v.discoveryPath = m.cfg.DefaultDiscoveryPath
		m.form = m.buildDiscoveryPathForm()
		m.state = settingsStateForm
		return m, m.form.Init()
	case 4:
		m.action = settingsActionTimeouts
		v.timeoutPipeline = strconv.Itoa(m.cfg.GetPipelineStepTimeout() / 60)
		v.timeoutHistorian = strconv.Itoa(m.cfg.GetHistorianTimeout() / 60)
		v.timeoutOracle = strconv.Itoa(m.cfg.GetOracleTimeout() / 60)
		v.timeoutOracleApply = strconv.Itoa(m.cfg.GetOracleApplyTimeout() / 60)
		m.form = m.buildTimeoutForm()
		m.state = settingsStateForm
		return m, m.form.Init()
	case 5:
		return m, backToMenu
	}
	return m, nil
}

func (m *SettingsModel) buildModeForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Default LLM Mode").
				Options(LLMModeOptions...).
				Height(len(LLMModeOptions)+1).
				Value(&v.llmMode),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which LLM for single mode?").
				Options(SingleLLMOptions...).
				Height(len(SingleLLMOptions)+1).
				Value(&v.singleLLM),
		).WithHideFunc(func() bool { return v.llmMode != "single_llm" }),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m *SettingsModel) buildLLMToggleForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals
	options := []huh.Option[string]{
		huh.NewOption[string]("Claude (always enabled)", "claude").Selected(true),
		huh.NewOption[string]("Codex", "codex"),
		huh.NewOption[string]("Gemini", "gemini"),
		huh.NewOption[string]("Copilot", "copilot"),
	}
	// Pre-select currently enabled LLMs
	for i, opt := range options {
		for _, enabled := range v.enabledLLMs {
			if opt.Value == enabled {
				options[i] = opt.Selected(true)
			}
		}
	}
	return huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Toggle LLM availability").
				Description("Claude is always enabled").
				Options(options...).
				Value(&v.enabledLLMs).
				Validate(func(selected []string) error {
					for _, s := range selected {
						if s == "claude" {
							return nil
						}
					}
					return fmt.Errorf("claude must be enabled")
				}),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m *SettingsModel) buildPipelineForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Default pipeline mode").
				Options(PipelineOptions...).
				Height(len(PipelineOptions)+1).
				Value(&v.pipelineMode),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m *SettingsModel) buildDiscoveryPathForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Default discovery folder").
				Description("Leave empty to use ./discovery (current directory)").
				Placeholder("~/ai-discovery").
				Value(&v.discoveryPath),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func validateMinutes(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("value is required")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("must be a number")
	}
	if n < 1 {
		return fmt.Errorf("must be at least 1 minute")
	}
	if n > 120 {
		return fmt.Errorf("max 120 minutes")
	}
	return nil
}

func (m *SettingsModel) buildTimeoutForm() *huh.Form {
	w := contentWidth(m.width)
	v := m.vals
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Pipeline step timeout (minutes)").
				Description("Per-step timeout for discovery pipeline").
				Value(&v.timeoutPipeline).
				Validate(validateMinutes),
			huh.NewInput().
				Title("Historian timeout (minutes)").
				Description("Timeout for historian indexing").
				Value(&v.timeoutHistorian).
				Validate(validateMinutes),
			huh.NewInput().
				Title("Oracle consult timeout (minutes)").
				Description("Timeout for oracle vision generation").
				Value(&v.timeoutOracle).
				Validate(validateMinutes),
			huh.NewInput().
				Title("Oracle apply timeout (minutes)").
				Description("Timeout for architect applying oracle vision").
				Value(&v.timeoutOracleApply).
				Validate(validateMinutes),
		),
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(formHeight(m.height))
}

func (m *SettingsModel) applyFormResult() {
	v := m.vals
	switch m.action {
	case settingsActionMode:
		m.cfg.LLMMode = v.llmMode
		if v.llmMode == "single_llm" {
			if mode, ok := m.cfg.LLMModes["single_llm"]; ok {
				mode.OverrideLLM = v.singleLLM
				mode.SynthesisLLM = v.singleLLM
				m.cfg.LLMModes["single_llm"] = mode
			}
		}
	case settingsActionLLMs:
		enabledSet := make(map[string]bool)
		for _, name := range v.enabledLLMs {
			enabledSet[name] = true
		}
		// Always ensure claude is enabled
		enabledSet["claude"] = true
		for _, name := range []string{"claude", "codex", "gemini", "copilot"} {
			m.cfg.LLMs[name] = enabledSet[name]
		}
	case settingsActionPipeline:
		m.cfg.PipelineMode = v.pipelineMode
	case settingsActionDiscoveryPath:
		m.cfg.DefaultDiscoveryPath = strings.TrimSpace(v.discoveryPath)
	case settingsActionTimeouts:
		if p, err := strconv.Atoi(strings.TrimSpace(v.timeoutPipeline)); err == nil {
			m.cfg.Timeouts.PipelineStep = p * 60
		}
		if h, err := strconv.Atoi(strings.TrimSpace(v.timeoutHistorian)); err == nil {
			m.cfg.Timeouts.Historian = h * 60
		}
		if o, err := strconv.Atoi(strings.TrimSpace(v.timeoutOracle)); err == nil {
			m.cfg.Timeouts.Oracle = o * 60
		}
		if a, err := strconv.Atoi(strings.TrimSpace(v.timeoutOracleApply)); err == nil {
			m.cfg.Timeouts.OracleApply = a * 60
		}
	}
}

func (m SettingsModel) currentEnabledLLMs() []string {
	var enabled []string
	for _, name := range []string{"claude", "codex", "gemini", "copilot"} {
		if e, ok := m.cfg.LLMs[name]; ok && e {
			enabled = append(enabled, name)
		}
	}
	return enabled
}

func (m SettingsModel) saveConfig() error {
	configPath := m.cfg.SourcePath

	// If config came from embedded/hardcoded (no source path) or from the
	// project default (config.default.json — shouldn't be modified), fall back
	// to the standalone user config location.
	if configPath == "" || filepath.Base(configPath) == "config.default.json" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("get home directory: %w", err)
		}
		configPath = filepath.Join(home, ".config", "vern", "config.json")
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := json.MarshalIndent(m.cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize config: %w", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}
	return nil
}

func (m SettingsModel) configSummary() string {
	label := lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render
	val := lipgloss.NewStyle().Foreground(colorSecondary).Render

	llmMode := m.cfg.LLMMode
	if llmMode == "" {
		llmMode = "mixed_claude_fallback"
	}
	pipelineMode := m.cfg.PipelineMode
	if pipelineMode == "" {
		pipelineMode = "default"
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("LLM Mode:"), val(llmMode)))

	b.WriteString(fmt.Sprintf("  %s      ", label("LLMs:")))
	for _, name := range []string{"claude", "codex", "gemini", "copilot"} {
		if enabled, ok := m.cfg.LLMs[name]; ok && enabled {
			b.WriteString(stepOKStyle.Render(name) + "  ")
		} else {
			b.WriteString(logDimStyle.Render(name) + "  ")
		}
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("Pipeline:"), val(pipelineMode)))

	discPath := m.cfg.DefaultDiscoveryPath
	if discPath == "" {
		discPath = "./discovery (default)"
	}
	b.WriteString(fmt.Sprintf("  %s  %s\n", label("Discovery:"), val(discPath)))
	b.WriteString(fmt.Sprintf("  %s  %s",
		label("Timeouts:"),
		val(fmt.Sprintf("Pipeline %dm | Historian %dm | Oracle %dm | Apply %dm",
			m.cfg.GetPipelineStepTimeout()/60,
			m.cfg.GetHistorianTimeout()/60,
			m.cfg.GetOracleTimeout()/60,
			m.cfg.GetOracleApplyTimeout()/60)),
	))

	return b.String()
}

func (m SettingsModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Settings"))
	b.WriteString("\n\n")

	switch m.state {
	case settingsStateMenu:
		// Config summary in bordered panel
		cw := contentWidth(m.width)
		b.WriteString(statusPanelStyle.Width(cw).Render(
			panelTitleStyle.Render("Current Config") + "\n" + m.configSummary(),
		))
		b.WriteString("\n\n")

		for i, item := range settingsMenuItems {
			prefix := "  "
			style := menuItemStyle
			if i == m.cursor {
				prefix = "> "
				style = menuSelectedStyle
			}

			number := fmt.Sprintf("[%d] ", i+1)
			if i == len(settingsMenuItems)-1 { // Back
				number = "[q] "
			}

			b.WriteString(style.Render(prefix + number + item))
			b.WriteString("\n")
		}

	case settingsStateForm:
		b.WriteString(m.form.View())

	case settingsStateSaveErr:
		b.WriteString(stepFailStyle.Render("Failed to save config"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("  %s\n", m.saveErr.Error()))
		b.WriteString("\n")
		b.WriteString(subtitleStyle.Render("  Press Enter to go back"))
	}

	return b.String()
}
