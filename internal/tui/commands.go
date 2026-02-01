package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yellalena/vkscape/internal/auth"
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/vkscape"
)

type downloadAlbumsDoneMsg struct{}
type downloadGroupsDoneMsg struct{}
type authStartMsg struct {
	authVerifier string
	authURL      string
}
type authResultMsg struct {
	ok bool
}

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

		reporter := newTUIProgressReporter(getProgressSender())
		if err := vkscape.DownloadGroups(groupIDs, logger, reporter); err != nil {
			output.Error(fmt.Sprintf("Failed to download groups: %v", err))
		}
		return downloadGroupsDoneMsg{}
	}
}

func authCmd() tea.Cmd {
	return func() tea.Msg {
		logger, logFile := output.InitLogger(false)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		session, err := auth.StartInteractiveFlow(logger)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to authenticate: %v", err))
			return authResultMsg{ok: false}
		}
		return authStartMsg{authVerifier: session.Verifier, authURL: session.AuthURL}
	}
}

func finishAuthCmd(verifier, redirectURL string) tea.Cmd {
	return func() tea.Msg {
		logger, logFile := output.InitLogger(false)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		if err := auth.FinishInteractiveFlow(logger, verifier, redirectURL); err != nil {
			output.Error(fmt.Sprintf("Failed to authenticate: %v", err))
			return authResultMsg{ok: false}
		}
		return authResultMsg{ok: true}
	}
}

func openAuthBrowserCmd(url string) tea.Cmd {
	return func() tea.Msg {
		logger, logFile := output.InitLogger(false)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		auth.OpenBrowser(url, logger)
		return nil
	}
}
