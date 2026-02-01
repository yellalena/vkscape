package tui

import (
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
	stateAlbumOwnerInput
	stateAlbumIDsInput
	stateAlbumDownload
)

type userInput struct {
	ownerID  string
	albumIDs string
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
	spin   spinner.Model
	progrs progressModel

	inputValues userInput

	errMsg string

	downloadDone bool
}

func initialModel() model {
	items := []list.Item{
		menuItem{title: utils.CommandAlbumsTitle, desc: utils.CommandAlbumsDesc},
		menuItem{title: utils.CommandGroupsTitle, desc: utils.CommandGroupsDesc},
		menuItem{title: utils.CommandAuthTitle, desc: utils.CommandAuthDesc},
		menuItem{title: utils.MenuQuit, desc: utils.MenuQuit},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "VKscape"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	ti := textinput.New()
	ti.Placeholder = "123456"
	ti.CharLimit = 50
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
		m.downloadDone = true
		m.clearLogs()
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
				m.state = stateAlbumOwnerInput
				m.errMsg = ""
				m.downloadDone = false
				m.clearLogs()
				m.input.SetValue("")
				m.input.Placeholder = "Owner ID"
				m.input.Focus()
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
				m.state = stateAlbumDownload
				m.errMsg = ""
				m.resetSpinner()
				m.resetProgress()
				return m, tea.Batch(downloadAlbumsCmd(ownerID, idList), m.spin.Tick)
			case "esc":
				m.state = stateMenu
			}
		}

	case stateAlbumDownload:
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

	case stateAlbumDownload:
		content := "Downloading albums...\n\nPlease wait.\n\n(esc to cancel view)"
		if !m.downloadDone {
			content = fmt.Sprintf("%s Downloading albums...\n\nPlease wait.\n\n(esc to cancel view)", m.spin.View())
		}
		if m.downloadDone {
			content = "Download complete.\n\n(esc to return to menu)"
		}
		progressBlock := m.renderProgress()
		if progressBlock != "" {
			content = content + "\n" + progressBlock
		}
		return m.renderDownloadView(content)
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

func (m *model) renderDownloadView(content string) string {
	logsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(logsGrey))
	if len(m.logs) == 0 {
		return content
	}

	logs := m.logs
	if len(logs) > maxVisibleLogLines {
		logs = logs[len(logs)-maxVisibleLogLines:]
	}

	return content + "\n\n" + logsStyle.Render(strings.Join(logs, "\n"))
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

func (m *model) renderProgress() string {
	if m.progrs.progTotal <= 0 {
		return ""
	}
	percent := float64(m.progrs.progCurrent) / float64(m.progrs.progTotal)
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	bar := m.progrs.prog.ViewAs(percent)
	if m.progrs.progStatus == "" {
		return bar
	}

	return bar + "\n" + m.progrs.progStatus
}
