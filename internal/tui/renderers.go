package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) renderWithLogs(content string) string {
	logsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(logsColor))
	errsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(red))
	if len(m.logs) == 0 && len(m.errs) == 0 {
		return content
	}

	logs := m.logs
	if len(logs) > maxVisibleLogLines {
		logs = logs[len(logs)-maxVisibleLogLines:]
	}

	var blocks []string
	if len(m.errs) > 0 {
		blocks = append(blocks, errsStyle.Render(strings.Join(m.errs, "\n")))
	}
	if len(logs) > 0 {
		blocks = append(blocks, logsStyle.Render(strings.Join(logs, "\n")))
	}

	return content + "\n\n" + strings.Join(blocks, "\n\n")
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
