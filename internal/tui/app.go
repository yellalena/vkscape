package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Start() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}
