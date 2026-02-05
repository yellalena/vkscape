package tui

import (
	"context"
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
type authResultMsg struct{}
type tokenResultMsg struct{}

func downloadAlbumsCmd(ctx context.Context, ownerID int, albumIDs []string) tea.Cmd {
	return func() tea.Msg {
		logger, logFile := output.InitLogger(false)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		reporter := newTUIProgressReporter(getProgressSender())
		if err := vkscape.DownloadAlbums(ctx, ownerID, albumIDs, logger, reporter); err != nil {
			output.Error(fmt.Sprintf("Failed to download albums: %v", err))
		}
		return downloadAlbumsDoneMsg{}
	}
}

func downloadGroupsCmd(ctx context.Context, groupIDs []string) tea.Cmd {
	return func() tea.Msg {
		logger, logFile := output.InitLogger(false)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		reporter := newTUIProgressReporter(getProgressSender())
		if err := vkscape.DownloadGroups(ctx, groupIDs, logger, reporter); err != nil {
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
			return authResultMsg{}
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
		}
		return authResultMsg{}
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

func saveTokenCmd(token string) tea.Cmd {
	return func() tea.Msg {
		logger, logFile := output.InitLogger(false)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		if err := vkscape.AppTokenAuth(token, logger); err != nil {
			output.Error(fmt.Sprintf("Failed to save token: %v", err))
		}
		return tokenResultMsg{}
	}
}
