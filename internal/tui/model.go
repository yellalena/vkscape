package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yellalena/vkscape/internal/utils"
)

type state int

const (
	stateMenu state = iota
	stateDownload
	stateAlbumOwnerInput
	stateAlbumIDsInput
	stateGroupIDsInput
	stateAuthRun
	stateAuthCompleting
	stateTokenInput
	stateTokenSaving
	stateHelp
)

type userInput struct {
	ownerID      string
	albumIDs     string
	authVerifier string
	appToken     string
}

type progressModel struct {
	prog        progress.Model
	progTotal   int
	progCurrent int
	progStatus  string
}

type model struct {
	state state

	menu  list.Model
	input textinput.Model

	logs   []string
	errs   []string
	spin   spinner.Model
	progrs progressModel

	inputValues userInput

	errMsg string

	actionDone bool

	cancel context.CancelFunc
}

func initialModel() model {
	items := []list.Item{
		menuItem{title: utils.CommandAlbumsTitle, desc: utils.CommandAlbumsDesc},
		menuItem{title: utils.CommandGroupsTitle, desc: utils.CommandGroupsDesc},
		menuItem{title: utils.CommandAuthTitle, desc: utils.CommandAuthDesc},
		menuItem{title: utils.CommandTokenTitle, desc: utils.CommandTokenDesc},
		menuItem{title: utils.CommandHelpTitle, desc: utils.CommandHelpDesc},
		menuItem{title: utils.MenuQuit, desc: utils.MenuQuit},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "VKscape"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	ti := textinput.New()
	ti.Placeholder = "123456"
	ti.Width = 30

	m := model{
		state: stateMenu,
		menu:  l,
		input: ti,
	}
	m.resetSpinner()
	m.resetProgress()

	return m
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Hello VKscape!")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.menu.SetSize(msg.Width, msg.Height)
	case downloadAlbumsDoneMsg:
		m.actionDone = true
		m.clearLogs()
		m.cancel = nil
	case downloadGroupsDoneMsg:
		m.actionDone = true
		m.clearLogs()
		m.cancel = nil
	case authStartMsg:
		m.inputValues.authVerifier = msg.authVerifier
		m.errMsg = ""
		m.input.SetValue("")
		m.input.Placeholder = "Redirect URL"
		m.input.Focus()
		cmds = append(cmds, openAuthBrowserCmd(msg.authURL))
	case authResultMsg:
		m.actionDone = true
	case tokenResultMsg:
		m.actionDone = true
	case progressStartMsg:
		m.progrs.progTotal = msg.total
	case progressIncMsg:
		if m.progrs.progCurrent < m.progrs.progTotal {
			m.progrs.progCurrent++
		}
	case progressStatusMsg:
		m.progrs.progStatus = msg.msg
	case progressDoneMsg:
		m.progrs.progCurrent = m.progrs.progTotal
	case logMsg:
		m.addLog(string(msg))
	case errorLogMsg:
		m.addErrorLog(string(msg))
	case spinner.TickMsg:
		m.spin, cmd = m.spin.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	switch m.state {

	case stateMenu:
		m.menu, cmd = m.menu.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			switch m.menu.SelectedItem().(menuItem).title {
			case utils.CommandAlbumsTitle:
				m.resetEverything()
				m.state = stateAlbumOwnerInput
				m.actionDone = false
				m.input.Placeholder = "Owner ID"
				m.input.Focus()
			case utils.CommandGroupsTitle:
				m.resetEverything()
				m.state = stateGroupIDsInput
				m.actionDone = false
				m.input.Placeholder = "Group IDs"
				m.input.Focus()
			case utils.CommandAuthTitle:
				m.resetEverything()
				m.state = stateAuthRun
				m.actionDone = false
				return m, tea.Batch(authCmd(), m.spin.Tick)
			case utils.CommandTokenTitle:
				m.resetEverything()
				m.state = stateTokenInput
				m.actionDone = false
				m.input.Placeholder = "App token"
				m.input.Focus()
			case utils.CommandHelpTitle:
				m.state = stateHelp
				return m, tea.ClearScreen
			case utils.MenuQuit:
				return m, tea.Quit
			}
		}

	case stateAlbumOwnerInput:
		m.input, cmd = m.input.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				m.inputValues.ownerID = m.input.Value()
				m.state = stateAlbumIDsInput
				m.errMsg = ""
				m.input.SetValue("")
				m.input.Placeholder = "Album IDs (empty = all)"
				m.input.Focus()
			case "esc":
				m.state = stateMenu
			}
		}

	case stateAlbumIDsInput:
		m.input, cmd = m.input.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				m.inputValues.albumIDs = m.input.Value()
				ownerID, err := strconv.Atoi(strings.TrimSpace(m.inputValues.ownerID))
				if err != nil {
					m.errMsg = "Owner ID must be an integer"
					return m, nil
				}

				idList := utils.ParseIDList(m.inputValues.albumIDs)
				m.state = stateDownload
				if m.cancel != nil {
					m.cancel()
				}
				ctx, cancel := context.WithCancel(context.Background())
				m.cancel = cancel
				return m, tea.Batch(downloadAlbumsCmd(ctx, ownerID, idList), m.spin.Tick)
			case "esc":
				m.state = stateMenu
			}
		}

	case stateDownload:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
			if m.cancel != nil {
				m.cancel()
				m.cancel = nil
			}
			m.state = stateMenu
			m.clearErrorLogs()
		}

	case stateGroupIDsInput:
		m.input, cmd = m.input.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				idList := utils.ParseIDList(m.input.Value())
				if len(idList) == 0 {
					m.errMsg = "Please enter at least one group ID"
					return m, nil
				}
				m.state = stateDownload
				if m.cancel != nil {
					m.cancel()
				}
				ctx, cancel := context.WithCancel(context.Background())
				m.cancel = cancel
				return m, tea.Batch(downloadGroupsCmd(ctx, idList), m.spin.Tick)
			case "esc":
				m.state = stateMenu
			}
		}

	case stateAuthRun:
		m.input, cmd = m.input.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				redirectURL := m.input.Value()
				if strings.TrimSpace(redirectURL) == "" {
					m.errMsg = "Please paste the full redirect URL"
					return m, nil
				}

				m.errMsg = ""
				m.input.SetValue("")
				m.state = stateAuthCompleting
				return m, finishAuthCmd(m.inputValues.authVerifier, redirectURL)
			case "esc":
				m.state = stateMenu
				m.clearErrorLogs()
			}
		}

	case stateAuthCompleting:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
			m.state = stateMenu
			m.clearErrorLogs()
		}

	case stateTokenInput:
		m.input, cmd = m.input.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				token := strings.TrimSpace(m.input.Value())
				if token == "" {
					m.errMsg = "Please enter a token"
					return m, nil
				}
				m.state = stateTokenSaving
				return m, tea.Batch(saveTokenCmd(token), m.spin.Tick)
			case "esc":
				m.state = stateMenu
			}
		}

	case stateTokenSaving:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
			m.state = stateMenu
			m.clearErrorLogs()
		}

	case stateHelp:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
			m.state = stateMenu
		}
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {

	case stateMenu:
		return m.menu.View()

	case stateAlbumOwnerInput:
		if m.errMsg != "" {
			return fmt.Sprintf(
				"Enter owner ID:\n\n%s\n\nError: %s\n\n(esc to cancel)",
				m.input.View(),
				m.errMsg,
			)
		}
		return fmt.Sprintf(
			"Enter owner ID:\n\n%s\n\n(esc to cancel)",
			m.input.View(),
		)

	case stateAlbumIDsInput:
		if m.errMsg != "" {
			return fmt.Sprintf(
				"Enter album IDs (comma or space separated).\nLeave empty for all:\n\n%s\n\nError: %s\n\n(esc to cancel)",
				m.input.View(),
				m.errMsg,
			)
		}
		return fmt.Sprintf(
			"Enter album IDs (comma or space separated).\nLeave empty for all:\n\n%s\n\n(esc to cancel)",
			m.input.View(),
		)

	case stateDownload:
		content := fmt.Sprintf("%s Downloading...\n\nPlease wait.\n\n(esc to cancel view)", m.spin.View())
		if m.actionDone {
			content = "Download complete.\n\n(esc to return to menu)"
		}
		progressBlock := m.renderProgress()
		if progressBlock != "" {
			content = content + "\n" + progressBlock
		}
		return m.renderWithLogs(content)

	case stateGroupIDsInput:
		if m.errMsg != "" {
			return fmt.Sprintf(
				"Enter group IDs (comma or space separated):\n\n%s\n\nError: %s\n\n(esc to cancel)",
				m.input.View(),
				m.errMsg,
			)
		}
		return fmt.Sprintf(
			"Enter group IDs (comma or space separated):\n\n%s\n\n(esc to cancel)",
			m.input.View(),
		)

	case stateAuthRun:
		if m.actionDone {
			return fmt.Sprintf("%s\n\n(esc to return to menu)")
		}
		if m.errMsg != "" {
			prompt := fmt.Sprintf(
				"Paste the FULL redirect URL from the browser:\n\n%s\n\nError: %s\n\n(esc to cancel)",
				m.input.View(),
				m.errMsg,
			)
			return m.renderWithLogs(prompt)
		}
		prompt := fmt.Sprintf(
			"Paste the FULL redirect URL from the browser:\n\n%s\n\n(esc to cancel)",
			m.input.View(),
		)
		return m.renderWithLogs(prompt)

	case stateAuthCompleting:
		content := fmt.Sprintf("%s Authenticating...\n\nPlease wait.\n\n(esc to cancel view)", m.spin.View())
		if m.actionDone {
			content = "Authentication completed.\n\n(esc to return to menu)"
		}
		return m.renderWithLogs(content)

	case stateTokenInput:
		if m.errMsg != "" {
			return fmt.Sprintf(
				"Enter app token:\n\n%s\n\nError: %s\n\n(esc to cancel)",
				m.input.View(),
				m.errMsg,
			)
		}
		return fmt.Sprintf(
			"Enter app token:\n\n%s\n\n(esc to cancel)",
			m.input.View(),
		)

	case stateTokenSaving:
		content := fmt.Sprintf("%s Saving token...\n\nPlease wait.\n\n(esc to cancel view)", m.spin.View())
		if m.actionDone {
			content = "Token saving comleted.\n\n(esc to return to menu)"
		}
		return m.renderWithLogs(content)

	case stateHelp:
		return fmt.Sprintf("%s\n\n(esc to return to menu)", utils.AppHelpText)
	}

	return ""
}

const (
	maxLogLines        = 500
	maxVisibleLogLines = 15
)

func (m *model) addLog(line string) {
	m.logs = append(m.logs, line)
	if len(m.logs) > maxLogLines {
		m.logs = m.logs[len(m.logs)-maxLogLines:]
	}
}

func (m *model) clearLogs() {
	m.logs = nil
}

func (m *model) addErrorLog(line string) {
	m.errs = append(m.errs, line)
}

func (m *model) clearErrorLogs() {
	m.errs = nil
}

func (m *model) resetEverything() {
	m.input.SetValue("")
	m.errMsg = ""
	m.clearLogs()
	m.clearErrorLogs()
	m.resetSpinner()
	m.resetProgress()
}

func (m *model) resetSpinner() {
	m.spin = spinner.New()
	m.spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(blue))
	m.spin.Spinner = spinner.Points
}

func (m *model) resetProgress() {
	m.progrs.progTotal = 0
	m.progrs.progCurrent = 0
	m.progrs.progStatus = ""

	m.progrs.prog = progress.New(
		progress.WithScaledGradient(blue, pink),
		progress.WithSpringOptions(100, 2.5),
		progress.WithWidth(60),
	)
}
