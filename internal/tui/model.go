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
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/utils"
	"github.com/yellalena/vkscape/internal/vkscape"
)

type state int

const (
	stateMenu state = iota
	stateAlbumOwnerInput
	stateAlbumIDsInput
	stateAlbumDownload
)

type model struct {
	state state

	menu  list.Model
	input textinput.Model

	logs []string
	spin spinner.Model
	prog progress.Model

	ownerID  string
	albumIDs string

	errMsg string

	downloadDone bool

	progTotal   int
	progCurrent int
	progStatus  string
}

func initialModel() model {
	items := []list.Item{
		menuItem("Download albums"),
		menuItem("Download groups"),
		menuItem("Authenticate"),
		menuItem("Quit"),
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
		m.progTotal = msg.total
	case progressIncMsg:
		if m.progCurrent < m.progTotal {
			m.progCurrent++
		}
	case progressStatusMsg:
		m.progStatus = msg.msg
	case progressDoneMsg:
		m.progCurrent = m.progTotal
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
			switch m.menu.SelectedItem().(menuItem) {
			case "Download albums":
				m.state = stateAlbumOwnerInput
				m.errMsg = ""
				m.downloadDone = false
				m.clearLogs()
				m.input.SetValue("")
				m.input.Placeholder = "Owner ID"
				m.input.Focus()
			case "Quit":
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
				m.ownerID = m.input.Value()
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
				m.albumIDs = m.input.Value()
				ownerID, err := strconv.Atoi(strings.TrimSpace(m.ownerID))
				if err != nil {
					m.errMsg = "Owner ID must be an integer"
					return m, nil
				}

				idList := utils.ParseIDList(m.albumIDs)
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

type downloadAlbumsDoneMsg struct{}

func downloadAlbumsCmd(ownerID int, albumIDs []string) tea.Cmd {
	return func() tea.Msg {
		logger, logFile := output.InitLogger(false)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		reporter := newTUIProgressReporter(getProgressSender())
		vkscape.DownloadAlbums(ownerID, albumIDs, logger, reporter)
		return downloadAlbumsDoneMsg{}
	}
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
	logsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
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
	m.spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	m.spin.Spinner = spinner.Points
}

func (m *model) resetProgress() {
	m.progTotal = 0
	m.progCurrent = 0
	m.progStatus = ""

	m.prog = progress.New(
		progress.WithScaledGradient("#5e61b5", "#c468ac"),
		progress.WithSpringOptions(100, 2.5),
		progress.WithWidth(60),
	)
}

func (m *model) renderProgress() string {
	if m.progTotal <= 0 {
		return ""
	}
	percent := float64(m.progCurrent) / float64(m.progTotal)
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	bar := m.prog.ViewAs(percent)
	if m.progStatus == "" {
		return bar
	}

	return bar + "\n" + m.progStatus
}
