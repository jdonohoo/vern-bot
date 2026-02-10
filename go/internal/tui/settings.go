package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/jdonohoo/vern-bot/go/internal/config"
)

type settingsState int

const (
	settingsStateMenu settingsState = iota
	settingsStateForm
	settingsStateSaved
	settingsStateSaveErr
)

type settingsAction int

const (
	settingsActionMode settingsAction = iota
	settingsActionLLMs
	settingsActionPipeline
)

// settingsVals holds form-bound values on the heap so pointers survive
// bubbletea's value-copy semantics.
type settingsVals struct {
	llmMode      string
	singleLLM    string
	enabledLLMs  []string
	pipelineMode string
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
	if m.form != nil {
		m.form.WithWidth(cw).WithHeight(h)
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

var settingsMenuItems = []string{
	"Change LLM Mode",
	"Toggle LLM availability",
	"Change default pipeline mode",
	"Save config",
	"Back to menu",
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

	case settingsStateSaved, settingsStateSaveErr:
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
			v := m.vals
			switch m.cursor {
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
				m.saveErr = m.saveConfig()
				if m.saveErr != nil {
					m.state = settingsStateSaveErr
				} else {
					m.state = settingsStateSaved
				}
			case 4:
				return m, backToMenu
			}
		}
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
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(m.height)
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
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(m.height)
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
	).WithTheme(VernTheme()).WithWidth(w).WithHeight(m.height)
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
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}
	configDir := filepath.Join(home, ".config", "vern")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	configPath := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(m.cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize config: %w", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}
	return nil
}

func (m SettingsModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Settings"))
	b.WriteString("\n\n")

	switch m.state {
	case settingsStateMenu:
		// Show current config summary
		b.WriteString(fmt.Sprintf("  LLM Mode: %s\n", llmStyle.Render(m.cfg.LLMMode)))
		b.WriteString("  LLMs: ")
		for _, name := range []string{"claude", "codex", "gemini", "copilot"} {
			if enabled, ok := m.cfg.LLMs[name]; ok && enabled {
				b.WriteString(stepOKStyle.Render(name+" ") + " ")
			} else {
				b.WriteString(subtitleStyle.Render(name+" ") + " ")
			}
		}
		b.WriteString(fmt.Sprintf("\n  Pipeline: %s\n", m.cfg.PipelineMode))
		b.WriteString("\n")

		for i, item := range settingsMenuItems {
			prefix := "  "
			style := menuItemStyle
			if i == m.cursor {
				prefix = "> "
				style = menuSelectedStyle
			}
			b.WriteString(style.Render(prefix + item))
			b.WriteString("\n")
		}

	case settingsStateForm:
		b.WriteString(m.form.View())

	case settingsStateSaved:
		b.WriteString(stepOKStyle.Render("Config saved!"))
		configPath := "~/.config/vern/config.json"
		if home, err := os.UserHomeDir(); err == nil {
			configPath = filepath.Join(home, ".config", "vern", "config.json")
		}
		b.WriteString(fmt.Sprintf("\n\nSaved to: %s\n", configPath))
		b.WriteString("\n")
		b.WriteString(subtitleStyle.Render("Press Enter to continue"))

	case settingsStateSaveErr:
		b.WriteString(stepFailStyle.Render("Failed to save config"))
		b.WriteString("\n\n")
		b.WriteString(fmt.Sprintf("  %s\n", m.saveErr.Error()))
		b.WriteString("\n")
		b.WriteString(subtitleStyle.Render("Press Enter to go back"))
	}

	return b.String()
}
