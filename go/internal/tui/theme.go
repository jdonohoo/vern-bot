package tui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// VernTheme returns a huh theme using the Vern brand colors.
func VernTheme() *huh.Theme {
	t := huh.ThemeBase()

	fg := lipgloss.Color("#e6edf3") // site --text

	// Focused styles
	t.Focused.Base = t.Focused.Base.BorderForeground(colorPrimary)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(colorPrimary).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(colorPrimary).Bold(true).MarginBottom(1)
	t.Focused.Description = t.Focused.Description.Foreground(colorMuted)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(colorDanger)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(colorDanger)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(colorSecondary)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(colorSecondary)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(colorSecondary)
	t.Focused.Option = t.Focused.Option.Foreground(fg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(colorSecondary)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(colorSecondary)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(colorSecondary).SetString("[x] ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(colorMuted).SetString("[ ] ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(fg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(lipgloss.Color("#000000")).Background(colorAccent).Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(fg).Background(lipgloss.Color("#374151"))
	t.Focused.Next = t.Focused.FocusedButton

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(colorSecondary)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(colorMuted)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(colorPrimary)

	// Blurred styles
	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}
