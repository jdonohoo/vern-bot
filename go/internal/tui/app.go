package tui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Screen identifies which screen is active.
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenDiscovery
	ScreenHole
	ScreenRun
	ScreenOracle
	ScreenGenerate
	ScreenHistorian
	ScreenSettings
)

// App is the root Bubble Tea model.
type App struct {
	screen          Screen
	menu            MenuModel
	discovery       DiscoveryModel
	hole            HoleModel
	run             RunModel
	oracle          OracleModel
	generate        GenerateModel
	historian       HistorianModel
	settings        SettingsModel
	help            help.Model
	width           int
	height          int
	projectRoot     string
	agentsDir       string
	currentVersion  string
	updateAvailable string // latest version if newer, empty otherwise
}

// NewApp creates the root TUI application.
func NewApp(projectRoot, agentsDir, version string) App {
	return App{
		screen:         ScreenMenu,
		menu:           NewMenuModel(),
		discovery:      NewDiscoveryModel(projectRoot, agentsDir),
		hole:           NewHoleModel(projectRoot, agentsDir),
		run:            NewRunModel(projectRoot, agentsDir),
		oracle:         NewOracleModel(projectRoot, agentsDir),
		settings:       NewSettingsModel(projectRoot),
		help:           newHelpModel(),
		projectRoot:    projectRoot,
		agentsDir:      agentsDir,
		currentVersion: version,
	}
}

func (a App) Init() tea.Cmd {
	return checkForUpdate(a.currentVersion)
}

type updateAvailableMsg struct {
	latestVersion string
}

func checkForUpdate(currentVersion string) tea.Cmd {
	return func() tea.Msg {
		if currentVersion == "" || currentVersion == "dev" {
			return nil
		}

		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get("https://api.github.com/repos/jdonohoo/vern-bot/releases/latest")
		if err != nil {
			return nil
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil
		}

		var release struct {
			TagName string `json:"tag_name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return nil
		}

		latest := strings.TrimPrefix(release.TagName, "v")
		current := strings.TrimPrefix(currentVersion, "v")

		if latest != "" && latest != current && isNewer(latest, current) {
			return updateAvailableMsg{latestVersion: latest}
		}
		return nil
	}
}

// isNewer returns true if a is a newer semver than b.
func isNewer(a, b string) bool {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		if aParts[i] != bParts[i] {
			// Compare numerically by padding to same length
			aLen, bLen := len(aParts[i]), len(bParts[i])
			if aLen < bLen {
				aParts[i] = strings.Repeat("0", bLen-aLen) + aParts[i]
			} else if bLen < aLen {
				bParts[i] = strings.Repeat("0", aLen-bLen) + bParts[i]
			}
			return aParts[i] > bParts[i]
		}
	}
	return len(aParts) > len(bParts)
}

// propagateSize sends the current terminal dimensions to the active sub-screen.
func (a *App) propagateSize() {
	switch a.screen {
	case ScreenDiscovery:
		a.discovery.SetSize(a.width, a.height)
	case ScreenHole:
		a.hole.SetSize(a.width, a.height)
	case ScreenRun:
		a.run.SetSize(a.width, a.height)
	case ScreenOracle:
		a.oracle.SetSize(a.width, a.height)
	case ScreenGenerate:
		a.generate.SetSize(a.width, a.height)
	case ScreenHistorian:
		a.historian.SetSize(a.width, a.height)
	case ScreenSettings:
		a.settings.SetSize(a.width, a.height)
	}
}

// activeKeyMap returns the help.KeyMap for the current screen/state.
func (a App) activeKeyMap() help.KeyMap {
	switch a.screen {
	case ScreenMenu:
		return menuKeys
	case ScreenDiscovery:
		switch a.discovery.state {
		case discStateEditFiles:
			return editFilesKeys
		case discStateRunning:
			return runningKeys
		case discStateDone:
			return runDoneKeys
		case discStateProjectSelect:
			if a.discovery.projectForm == nil {
				return editFilesKeys
			}
			return formKeys
		default:
			return formKeys
		}
	case ScreenHole:
		switch a.hole.state {
		case holeStateRunning:
			return runningKeys
		case holeStateDone:
			return runDoneKeys
		default:
			return formKeys
		}
	case ScreenRun:
		switch a.run.state {
		case runStateRunning:
			return runningKeys
		case runStateDone:
			if a.run.err != nil {
				return runRetryKeys
			}
			return runDoneKeys
		default:
			return formKeys
		}
	case ScreenOracle:
		switch a.oracle.state {
		case oracleStateRunning:
			return runningKeys
		case oracleStateDone:
			return oracleDoneKeys
		default:
			return formKeys
		}
	case ScreenGenerate:
		switch a.generate.state {
		case genStateRunning:
			return runningKeys
		case genStateDone:
			if a.generate.err != nil {
				return runRetryKeys
			}
			return doneKeys
		default:
			return formKeys
		}
	case ScreenHistorian:
		switch a.historian.state {
		case histStateRunning:
			return runningKeys
		case histStateDone:
			if a.historian.err != nil {
				return runRetryKeys
			}
			return doneKeys
		default:
			return formKeys
		}
	case ScreenSettings:
		if a.settings.state == settingsStateMenu {
			return settingsMenuKeys
		}
		return formKeys
	}
	return menuKeys
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.help.Width = contentWidth(a.width)
		a.propagateSize()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			a.cancelActiveScreen()
			return a, tea.Quit
		}

	case backToMenuMsg:
		a.screen = ScreenMenu
		a.menu = NewMenuModel()
		return a, tea.EnableMouseCellMotion

	case updateAvailableMsg:
		a.updateAvailable = msg.latestVersion
		return a, nil
	}

	var cmd tea.Cmd
	switch a.screen {
	case ScreenMenu:
		newMenu, menuCmd := a.menu.Update(msg)
		a.menu = newMenu.(MenuModel)
		cmd = menuCmd
		// Check if menu selected something
		if a.menu.chosen != -1 {
			switch a.menu.chosen {
			case 0:
				a.screen = ScreenDiscovery
				a.discovery = NewDiscoveryModel(a.projectRoot, a.agentsDir)
				a.discovery.SetSize(a.width, a.height)
				a.menu.chosen = -1
				return a, a.discovery.Init()
			case 1:
				a.screen = ScreenDiscovery
				a.discovery = NewRerunDiscoveryModel(a.projectRoot, a.agentsDir)
				a.discovery.SetSize(a.width, a.height)
				a.menu.chosen = -1
				return a, a.discovery.Init()
			case 2:
				a.screen = ScreenHole
				a.hole = NewHoleModel(a.projectRoot, a.agentsDir)
				a.hole.SetSize(a.width, a.height)
				a.menu.chosen = -1
				return a, a.hole.Init()
			case 3:
				a.screen = ScreenRun
				a.run = NewRunModel(a.projectRoot, a.agentsDir)
				a.run.SetSize(a.width, a.height)
				a.menu.chosen = -1
				return a, a.run.Init()
			case 4:
				a.screen = ScreenOracle
				a.oracle = NewOracleModel(a.projectRoot, a.agentsDir)
				a.oracle.SetSize(a.width, a.height)
				a.menu.chosen = -1
				return a, a.oracle.Init()
			case 5:
				a.screen = ScreenGenerate
				a.generate = NewGenerateModel(a.projectRoot, a.agentsDir)
				a.generate.SetSize(a.width, a.height)
				a.menu.chosen = -1
				return a, a.generate.Init()
			case 6:
				a.screen = ScreenHistorian
				a.historian = NewHistorianModel(a.projectRoot, a.agentsDir)
				a.historian.SetSize(a.width, a.height)
				a.menu.chosen = -1
				return a, a.historian.Init()
			case 7:
				a.screen = ScreenSettings
				a.settings = NewSettingsModel(a.projectRoot)
				a.settings.SetSize(a.width, a.height)
				a.menu.chosen = -1
				return a, a.settings.Init()
			case 8:
				return a, tea.Quit
			}
			a.menu.chosen = -1
		}

	case ScreenDiscovery:
		newDisc, discCmd := a.discovery.Update(msg)
		a.discovery = newDisc.(DiscoveryModel)
		cmd = discCmd

	case ScreenHole:
		newHole, holeCmd := a.hole.Update(msg)
		a.hole = newHole.(HoleModel)
		cmd = holeCmd

	case ScreenRun:
		newRun, runCmd := a.run.Update(msg)
		a.run = newRun.(RunModel)
		cmd = runCmd

	case ScreenOracle:
		newOracle, oracleCmd := a.oracle.Update(msg)
		a.oracle = newOracle.(OracleModel)
		cmd = oracleCmd

	case ScreenGenerate:
		newGen, genCmd := a.generate.Update(msg)
		a.generate = newGen.(GenerateModel)
		cmd = genCmd

	case ScreenHistorian:
		newHist, histCmd := a.historian.Update(msg)
		a.historian = newHist.(HistorianModel)
		cmd = histCmd

	case ScreenSettings:
		newSettings, settingsCmd := a.settings.Update(msg)
		a.settings = newSettings.(SettingsModel)
		cmd = settingsCmd
	}

	return a, cmd
}

func (a App) View() string {
	var content string
	switch a.screen {
	case ScreenMenu:
		content = a.menu.View()
	case ScreenDiscovery:
		content = a.discovery.View()
	case ScreenHole:
		content = a.hole.View()
	case ScreenRun:
		content = a.run.View()
	case ScreenOracle:
		content = a.oracle.View()
	case ScreenGenerate:
		content = a.generate.View()
	case ScreenHistorian:
		content = a.historian.View()
	case ScreenSettings:
		content = a.settings.View()
	default:
		content = "Unknown screen"
	}

	// Constrain content width for readability using dynamic width
	w := contentWidth(a.width)
	styled := lipgloss.NewStyle().MaxWidth(w).Render(content)

	// Help bar
	helpView := helpBarStyle.Render(a.help.View(a.activeKeyMap()))

	// Center in terminal
	if a.width > 0 && a.height > 0 {
		// Reserve lines for help bar + optional update banner
		reserved := 2
		if a.updateAvailable != "" {
			reserved++
		}
		contentHeight := a.height - reserved

		output := lipgloss.Place(a.width, contentHeight, lipgloss.Center, lipgloss.Center, styled)
		output += "\n" + lipgloss.PlaceHorizontal(a.width, lipgloss.Center, helpView)

		if a.updateAvailable != "" {
			banner := updateStyle.Render(
				fmt.Sprintf("Update available %s! Run: brew upgrade vern", a.updateAvailable),
			)
			output += "\n" + lipgloss.PlaceHorizontal(a.width, lipgloss.Right, banner)
		}
		return output
	}
	return styled + "\n" + helpView
}

type backToMenuMsg struct{}

func backToMenu() tea.Msg {
	return backToMenuMsg{}
}

// cancelActiveScreen cancels any running goroutine on the active screen.
func (a *App) cancelActiveScreen() {
	switch a.screen {
	case ScreenDiscovery:
		a.discovery.Cancel()
	case ScreenHole:
		a.hole.Cancel()
	case ScreenRun:
		a.run.Cancel()
	case ScreenOracle:
		a.oracle.Cancel()
	case ScreenGenerate:
		a.generate.Cancel()
	case ScreenHistorian:
		a.historian.Cancel()
	}
}

// Run launches the Bubble Tea TUI.
func Run(projectRoot, agentsDir, version string) error {
	p := tea.NewProgram(NewApp(projectRoot, agentsDir, version), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}
	return nil
}
