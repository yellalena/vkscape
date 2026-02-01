package tui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellalena/vkscape/internal/output"
)

func Start() error {
	p := tea.NewProgram(initialModel())
	output.SetWriter(newLogWriter(p.Send))
	output.SetErrorWriter(newErrorLogWriter(p.Send))
	setProgressSender(p.Send)
	defer output.SetWriter(os.Stdout)
	defer output.SetErrorWriter(os.Stdout)
	defer setProgressSender(nil)
	_, err := p.Run()
	return err
}
