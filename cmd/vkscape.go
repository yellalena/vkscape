package cmd

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/yellalena/vkscape/internal/auth"
	"github.com/yellalena/vkscape/internal/config"
	"github.com/yellalena/vkscape/internal/models"
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/utils"
	"github.com/yellalena/vkscape/internal/vkscape"
)

func DownloadGroups(groupIDs []string, logger *slog.Logger) error {
	svc := vkscape.InitService(logger)
	output.Info(fmt.Sprintf("Processing %d group(s)...", len(groupIDs)))
	for _, groupID := range groupIDs {
		output.Info(fmt.Sprintf("📥 Downloading group: %s", groupID))
		groupDir, err := utils.CreateGroupDirectory(groupID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create directory for group %s: %v", groupID, err))
			return fmt.Errorf("failed to create directory for group %s: %w", groupID, err)
		}
		logger.Info("Created directory", "group_id", groupID, "dir", groupDir)
		posts, err := svc.Client.GetPosts(groupID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch posts for group %s: %v", groupID, err))
			return fmt.Errorf("failed to get posts for group %s: %w", groupID, err)
		}
		output.Info(fmt.Sprintf("  Found %d post(s) in group %s", len(posts), groupID))
		logger.Info("Found posts", "group_id", groupID, "count", len(posts))
		svc.Parser.ParseWallPosts(&svc.Wg, groupDir, posts)
	}

	svc.Wg.Wait()
	output.Success(fmt.Sprintf("✅ Completed downloading %d group(s)", len(groupIDs)))

	return nil
}

func DownloadAlbums(ownerID int, albumIDs []string, logger *slog.Logger) {
	svc := vkscape.InitService(logger)
	output.Info(fmt.Sprintf("Processing albums for owner %d...", ownerID))

	output.Info("Fetching available album list from VK...")
	vkAlbums := svc.Client.GetAlbums(ownerID)
	allAlbums := models.VkAlbumsToPhotoAlbums(vkAlbums)
	output.Info(fmt.Sprintf("  Found %d album(s)", len(allAlbums)))
	logger.Info("Found albums", "owner_id", ownerID, "count", len(allAlbums))

	var albums []models.PhotoAlbum
	if len(albumIDs) == 0 {
		// Use all albums
		albums = allAlbums
	} else {
		// Filter to only requested albums
		albums = models.FilterAlbumsByIDs(albumIDs, allAlbums)
		logger.Info(
			"Using provided albums",
			"owner_id",
			ownerID,
			"count",
			len(albums),
		)
	}

	output.Info(fmt.Sprintf("Downloading %d album(s)", len(albums)))

	for _, album := range albums {
		albumTitle := album.Title
		if albumTitle == "" {
			albumTitle = fmt.Sprintf("Album %d", album.ID)
		}
		output.Info(fmt.Sprintf("📷 Downloading album: %s (ID: %d)", albumTitle, album.ID))
		albumDir := utils.CreateAlbumDirectory(album)
		logger.Info("Created album directory", "album_id", album.ID, "dir", albumDir)
		photos := svc.Client.GetPhotos(ownerID, strconv.Itoa(album.ID))
		output.Info(fmt.Sprintf("  Found %d photo(s) in album '%s'", len(photos), albumTitle))
		logger.Info("Found photos", "album_id", album.ID, "count", len(photos))
		svc.Parser.ParseAlbumPhotos(&svc.Wg, albumDir, strconv.Itoa(album.ID), photos)
	}

	svc.Wg.Wait()
	output.Success(
		fmt.Sprintf("✅ Completed downloading %d album(s) for owner %d", len(albums), ownerID),
	)
}

func InteractiveAuth(logger *slog.Logger) {
	err := auth.InteractiveFlow(logger)
	if err != nil {
		output.Error(fmt.Sprintf("Authentication failed: %v", err))
		logger.Error("Authentication failed", "error", err)
		return
	}
	output.Success("Authentication successful!")
}

func AppTokenAuth(token string, logger *slog.Logger) {
	output.Info("Saving app token...")
	cfg := &config.AuthConfig{
		AuthMethod:  config.AuthMethodAppToken,
		AccessToken: token,
	}

	err := config.SaveConfig(cfg)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to save configuration: %v", err))
		logger.Error("Failed to save config", "error", err)
		return
	}
	output.Success("App token saved successfully!")
}
