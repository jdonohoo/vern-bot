package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/jdonohoo/vern-bot/go/internal/config"
)

// MenuModel is the main menu screen.
type MenuModel struct {
	cursor  int
	items   []string
	chosen  int
	cfg     *config.Config
	version string
}

func NewMenuModel(version string) MenuModel {
	return MenuModel{
		items: []string{
			"Discovery Pipeline",
			"Rerun Discovery",
			"VernHole Council",
			"Single LLM Run",
			"Post-Processing",
			"Generate Persona",
			"Historian",
			"Settings",
			"Quit",
		},
		chosen:  -1,
		cfg:     config.Load(""),
		version: version,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.chosen = m.cursor
		case "1":
			m.chosen = 0
		case "2":
			m.chosen = 1
		case "3":
			m.chosen = 2
		case "4":
			m.chosen = 3
		case "5":
			m.chosen = 4
		case "6":
			m.chosen = 5
		case "7":
			m.chosen = 6
		case "8":
			m.chosen = 7
		case "q":
			m.chosen = 8
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	var b strings.Builder

	title := "VERN-BOT"
	if m.version != "" {
		title += "  " + subtitleStyle.Render("v"+strings.TrimPrefix(m.version, "v"))
	}
	header := headerBox.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(title),
			subtitleStyle.Render("Multi-LLM Discovery Pipeline"),
		),
	)
	b.WriteString(header)
	b.WriteString("\n\n")

	for i, item := range m.items {
		prefix := "  "
		style := menuItemStyle
		if i == m.cursor {
			prefix = "> "
			style = menuSelectedStyle
		}

		number := fmt.Sprintf("[%d] ", i+1)
		if i == len(m.items)-1 {
			number = "[q] "
		}

		b.WriteString(style.Render(prefix + number + item))
		b.WriteString("\n")
	}

	// Status bar
	llmMode := m.cfg.LLMMode
	if llmMode == "" {
		llmMode = "mixed_claude_fallback"
	}

	var llms []string
	for _, name := range []string{"claude", "codex", "gemini", "copilot"} {
		if enabled, ok := m.cfg.LLMs[name]; ok && enabled {
			llms = append(llms, name)
		}
	}

	b.WriteString("\n")
	b.WriteString(statusBarStyle.Render(
		fmt.Sprintf("LLM Mode: %s  |  LLMs: %s", llmMode, strings.Join(llms, " ")),
	))

	return b.String()
}
