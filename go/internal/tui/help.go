package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

// Key bindings for each screen context

type menuKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Number key.Binding
	Quit   key.Binding
}

func (k menuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Number, k.Quit}
}

func (k menuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Select}, {k.Number, k.Quit}}
}

var menuKeys = menuKeyMap{
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("j/k", "navigate")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("", "")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Number: key.NewBinding(key.WithKeys("1", "2", "3", "4", "5", "6", "7"), key.WithHelp("1-7", "jump")),
	Quit:   key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
}

type formKeyMap struct {
	Back key.Binding
}

func (k formKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back}
}

func (k formKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Back}}
}

var formKeys = formKeyMap{
	Back: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
}

type runningKeyMap struct {
	Quit key.Binding
}

func (k runningKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k runningKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit}}
}

var runningKeys = runningKeyMap{
	Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "cancel")),
}

type doneKeyMap struct {
	Back key.Binding
	Up   key.Binding
	Down key.Binding
}

func (k doneKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Back}
}

func (k doneKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Back}}
}

var doneKeys = doneKeyMap{
	Back: key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "menu")),
	Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("j/k", "scroll")),
	Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("", "")),
}

type settingsMenuKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Number key.Binding
	Save   key.Binding
	Back   key.Binding
}

func (k settingsMenuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Number, k.Save, k.Back}
}

func (k settingsMenuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Select}, {k.Number, k.Save, k.Back}}
}

var settingsMenuKeys = settingsMenuKeyMap{
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("j/k", "navigate")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("", "")),
	Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Number: key.NewBinding(key.WithKeys("1", "2", "3"), key.WithHelp("1-3", "jump")),
	Save:   key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "save")),
	Back:   key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "back")),
}

type runDoneKeyMap struct {
	Copy key.Binding
	Back key.Binding
	Up   key.Binding
	Down key.Binding
}

func (k runDoneKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Copy, k.Back}
}

func (k runDoneKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Copy, k.Back}}
}

var runDoneKeys = runDoneKeyMap{
	Copy: key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy")),
	Back: key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "menu")),
	Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("j/k", "scroll")),
	Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("", "")),
}

type runRetryKeyMap struct {
	Copy  key.Binding
	Retry key.Binding
	Back  key.Binding
	Up    key.Binding
	Down  key.Binding
}

func (k runRetryKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Copy, k.Retry, k.Back}
}

func (k runRetryKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Copy, k.Retry, k.Back}}
}

var runRetryKeys = runRetryKeyMap{
	Copy:  key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy")),
	Retry: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "retry")),
	Back:  key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "menu")),
	Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("j/k", "scroll")),
	Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("", "")),
}

type oracleDoneKeyMap struct {
	Copy key.Binding
	Back key.Binding
	Up   key.Binding
	Down key.Binding
}

func (k oracleDoneKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Copy, k.Back}
}

func (k oracleDoneKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Copy, k.Back}}
}

var oracleDoneKeys = oracleDoneKeyMap{
	Copy: key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy")),
	Back: key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q/esc", "menu")),
	Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("j/k", "scroll")),
	Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("", "")),
}

type editFilesKeyMap struct {
	Continue key.Binding
	Back     key.Binding
}

func (k editFilesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Continue, k.Back}
}

func (k editFilesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Continue, k.Back}}
}

var editFilesKeys = editFilesKeyMap{
	Continue: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "continue")),
	Back:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
}

func newHelpModel() help.Model {
	h := help.New()
	h.Styles.ShortKey = helpKeyStyle
	h.Styles.ShortDesc = helpDescStyle
	h.Styles.ShortSeparator = helpSepStyle
	return h
}
