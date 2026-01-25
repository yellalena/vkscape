package tui

import (
	"bytes"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

type logMsg string

type logWriter struct {
	send func(tea.Msg)
	mu   sync.Mutex
	buf  bytes.Buffer
}

func newLogWriter(send func(tea.Msg)) *logWriter {
	return &logWriter{send: send}
}

func (w *logWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	n, _ := w.buf.Write(p)
	for {
		b := w.buf.Bytes()
		idx := bytes.IndexByte(b, '\n')
		if idx == -1 {
			break
		}
		line := string(bytes.TrimRight(b[:idx], "\r"))
		w.buf.Next(idx + 1)
		w.send(logMsg(line))
	}

	return n, nil
}
