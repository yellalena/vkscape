package tui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellalena/vkscape/internal/output"
)

func Start() error {
	p := tea.NewProgram(initialModel())
	output.SetWriter(newLogWriter(p.Send))
	defer output.SetWriter(os.Stdout)
	_, err := p.Run()
	return err
}
