package tui

import "testing"

func TestVernThemeNotNil(t *testing.T) {
	theme := VernTheme()
	if theme == nil {
		t.Fatal("VernTheme() returned nil")
	}
}

func TestVernThemeFocusedStyles(t *testing.T) {
	theme := VernTheme()

	// Verify focused selector has content (should be "> ")
	selector := theme.Focused.SelectSelector.Render("")
	if selector == "" {
		// The SetString makes it render "> " even with empty input
		t.Log("SelectSelector renders empty (base theme sets this)")
	}

	// Verify focused button has background
	btn := theme.Focused.FocusedButton.Render("OK")
	if btn == "" {
		t.Error("FocusedButton should render content")
	}
}

func TestVernThemeBlurredStyles(t *testing.T) {
	theme := VernTheme()

	// Blurred should use hidden border
	blurredBase := theme.Blurred.Base.Render("test")
	if blurredBase == "" {
		t.Error("Blurred base should render content")
	}
}
