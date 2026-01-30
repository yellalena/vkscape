package tui

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellalena/vkscape/internal/progress"
)

type progressStartMsg struct {
	total int
}

type progressIncMsg struct{}

type progressStatusMsg struct {
	msg string
}

type progressDoneMsg struct{}

var (
	progressSenderMu sync.RWMutex
	progressSender   func(tea.Msg)
)

func setProgressSender(send func(tea.Msg)) {
	progressSenderMu.Lock()
	defer progressSenderMu.Unlock()
	progressSender = send
}

func getProgressSender() func(tea.Msg) {
	progressSenderMu.RLock()
	defer progressSenderMu.RUnlock()
	return progressSender
}

type tuiProgressReporter struct {
	send func(tea.Msg)
}

func newTUIProgressReporter(send func(tea.Msg)) progress.Reporter {
	if send == nil {
		return &progress.NoopReporter{}
	}
	return &tuiProgressReporter{send: send}
}

func (r *tuiProgressReporter) Start(total int) {
	r.send(progressStartMsg{total: total})
}

func (r *tuiProgressReporter) Increment() {
	r.send(progressIncMsg{})
}

func (r *tuiProgressReporter) SetStatus(msg string) {
	r.send(progressStatusMsg{msg: msg})
}

func (r *tuiProgressReporter) Done() {
	r.send(progressDoneMsg{})
}
