package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/vkscape"
)

type downloadAlbumsDoneMsg struct{}
type downloadGroupsDoneMsg struct{}

func downloadAlbumsCmd(ownerID int, albumIDs []string) tea.Cmd {
	return func() tea.Msg {
		logger, logFile := output.InitLogger(false)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		reporter := newTUIProgressReporter(getProgressSender())
		if err := vkscape.DownloadAlbums(ownerID, albumIDs, logger, reporter); err != nil {
			output.Error(fmt.Sprintf("Failed to download albums: %v", err))
		}
		return downloadAlbumsDoneMsg{}
	}
}

func downloadGroupsCmd(groupIDs []string) tea.Cmd {
	return func() tea.Msg {
		logger, logFile := output.InitLogger(false)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		if err := vkscape.DownloadGroups(groupIDs, logger); err != nil {
			output.Error(fmt.Sprintf("Failed to download groups: %v", err))
		}
		return downloadGroupsDoneMsg{}
	}
}
