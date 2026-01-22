package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateMenu state = iota
	stateAlbumOwnerInput
	stateAlbumIDsInput
)

type model struct {
	state state

	menu  list.Model
	input textinput.Model

	ownerID  string
	albumIDs string
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

	return model{
		state: stateMenu,
		menu:  l,
		input: ti,
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Hello VKscape!")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.menu.SetSize(msg.Width, msg.Height)
	}

	switch m.state {

	case stateMenu:
		m.menu, cmd = m.menu.Update(msg)

		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			switch m.menu.SelectedItem().(menuItem) {
			case "Download albums":
				m.state = stateAlbumOwnerInput
				m.input.SetValue("")
				m.input.Placeholder = "Owner ID"
				m.input.Focus()
			case "Quit":
				return m, tea.Quit
			}
		}

	case stateAlbumOwnerInput:
		m.input, cmd = m.input.Update(msg)

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				m.ownerID = m.input.Value()
				m.state = stateAlbumIDsInput
				m.input.SetValue("")
				m.input.Placeholder = "Album IDs (empty = all)"
				m.input.Focus()
			case "esc":
				m.state = stateMenu
			}
		}

	case stateAlbumIDsInput:
		m.input, cmd = m.input.Update(msg)

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				m.albumIDs = m.input.Value()
				// next state will be spinner → download later
			case "esc":
				m.state = stateMenu
			}
		}
	}

	return m, cmd
}

func (m model) View() string {
	switch m.state {

	case stateMenu:
		return m.menu.View()

	case stateAlbumOwnerInput:
		return fmt.Sprintf(
			"Enter owner ID:\n\n%s\n\n(esc to cancel)",
			m.input.View(),
		)

	case stateAlbumIDsInput:
		return fmt.Sprintf(
			"Enter album IDs (comma or space separated).\nLeave empty for all:\n\n%s\n\n(esc to cancel)",
			m.input.View(),
		)
	}

	return ""
}
