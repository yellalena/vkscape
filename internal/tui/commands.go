package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/vkscape"
)

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
