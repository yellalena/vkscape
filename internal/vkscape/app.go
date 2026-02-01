package vkscape

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/yellalena/vkscape/internal/auth"
	"github.com/yellalena/vkscape/internal/config"
	"github.com/yellalena/vkscape/internal/models"
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/progress"
	"github.com/yellalena/vkscape/internal/utils"
)

func DownloadGroups(groupIDs []string, logger *slog.Logger) error {
	if logger == nil {
		return fmt.Errorf("logger is nil")
	}
	svc, err := InitService(logger)
	if err != nil {
		return err
	}
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

func DownloadAlbums(ownerID int, albumIDs []string, logger *slog.Logger, reporter progress.Reporter) error {
	if logger == nil {
		return fmt.Errorf("logger is nil")
	}
	svc, err := InitService(logger)
	if err != nil {
		return err
	}
	output.Info(fmt.Sprintf("Processing albums for owner %d...", ownerID))

	output.Info("Fetching available album list from VK...")
	vkAlbums, err := svc.Client.GetAlbums(ownerID)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to fetch albums: %v", err))
		logger.Error("Failed to fetch albums", "error", err, "owner_id", ownerID)
		return err
	}
	allAlbums := models.VkAlbumsToPhotoAlbums(vkAlbums)
	logger.Info("Found albums", "owner_id", ownerID, "count", len(allAlbums))

	var albums []models.PhotoAlbum
	if len(albumIDs) == 0 {
		// Use all albums
		albums = allAlbums
	} else {
		// Filter to only requested albums
		var invalidIDs []string
		albums, invalidIDs = models.FilterAlbumsByIDs(albumIDs, allAlbums)
		logger.Info(
			"Using provided albums",
			"owner_id",
			ownerID,
			"count",
			len(albums),
		)
		if len(invalidIDs) > 0 {
			output.Error(fmt.Sprintf("Warning: invalid album IDs: %s", strings.Join(invalidIDs, ", ")))
			logger.Warn("Invalid album IDs", "ids", invalidIDs, "owner_id", ownerID)
		}
	}

	output.Info(fmt.Sprintf("Found %d album(s)", len(albums)))
	reporter.Start(len(albums))

	if len(albums) == 0 {
		output.Info("No albums to download. Exiting.")
		reporter.Done()
		return nil
	}

	output.Info(fmt.Sprintf("Downloading %d album(s)", len(albums)))

	for _, album := range albums {
		albumTitle := album.Title
		if albumTitle == "" {
			albumTitle = fmt.Sprintf("Album %d", album.ID)
		}
		reporter.SetStatus(fmt.Sprintf("Downloading album %d", album.ID))
		output.Info(fmt.Sprintf("📷 Downloading album: %s (ID: %d)", albumTitle, album.ID))
		albumDir, err := utils.CreateAlbumDirectory(album)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create directory for album %d: %v", album.ID, err))
			logger.Error("Failed to create album directory", "error", err, "album_id", album.ID)
			reporter.Increment()
			continue
		}
		logger.Info("Created album directory", "album_id", album.ID, "dir", albumDir)
		photos, err := svc.Client.GetPhotos(ownerID, strconv.Itoa(album.ID))
		if err != nil {
			output.Error(fmt.Sprintf("Warning! Album %d photos download is incomplete due to an internal error.", album.ID))
			logger.Error("Failed to fetch photos", "error", err, "album_id", album.ID, "owner_id", ownerID)
		}
		output.Info(fmt.Sprintf("  Found %d photo(s) in album '%s'", len(photos), albumTitle))
		logger.Info("Found photos", "album_id", album.ID, "count", len(photos))
		if len(photos) > 0 {
			svc.Parser.ParseAlbumPhotos(&svc.Wg, albumDir, strconv.Itoa(album.ID), photos)
		}

		svc.Wg.Wait()
		reporter.Increment()
		errCount := svc.Parser.CloseErrorsAndCount()
		if errCount > 0 {
			reporter.SetStatus(fmt.Sprintf("Completed with %d errors", errCount))
		}
	}

	output.Success(
		fmt.Sprintf("✅ Completed downloading %d album(s) for owner %d", len(albums), ownerID),
	)
	reporter.Done()
	return nil
}

func InteractiveAuth(logger *slog.Logger) error {
	if logger == nil {
		return fmt.Errorf("logger is nil")
	}
	err := auth.InteractiveFlow(logger)
	if err != nil {
		output.Error(fmt.Sprintf("Authentication failed: %v", err))
		logger.Error("Authentication failed", "error", err)
		return err
	}
	output.Success("Successfully retrieved token!")
	return nil
}

func AppTokenAuth(token string, logger *slog.Logger) error {
	if logger == nil {
		return fmt.Errorf("logger is nil")
	}
	output.Info("Saving app token...")
	cfg := &config.AuthConfig{
		AuthMethod:  config.AuthMethodAppToken,
		AccessToken: token,
	}

	err := config.SaveConfig(cfg)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to save configuration: %v", err))
		logger.Error("Failed to save config", "error", err)
		return err
	}
	output.Success("App token saved successfully!")
	return nil
}
