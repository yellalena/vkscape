package cmd

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/yellalena/vkscape/internal/auth"
	"github.com/yellalena/vkscape/internal/config"
	"github.com/yellalena/vkscape/internal/models"
	"github.com/yellalena/vkscape/internal/utils"
	"github.com/yellalena/vkscape/internal/vkscape"
)

// todo improve logging

func DownloadGroups(groupIDs []string, logger *slog.Logger) error {
	svc := vkscape.InitService(logger)
	// todo: add some check for groupID format - number should be negative
	// maybe separate command for parsin profile, but reusing same GetPosts method
	for _, groupID := range groupIDs {
		groupDir, err := utils.CreateGroupDirectory(groupID)
		if err != nil {
			return fmt.Errorf("failed to create directory for group %s: %w", groupID, err)
		}
		logger.Info("Created directory", "group_id", groupID, "dir", groupDir)
		posts, err := svc.Client.GetPosts(groupID, 5) // todo: remove limit
		if err != nil {
			return fmt.Errorf("failed to get posts for group %s: %w", groupID, err)
		}
		logger.Info("Found posts", "group_id", groupID, "count", len(posts))
		svc.Parser.ParseWallPosts(&svc.Wg, groupDir, posts)
	}

	svc.Wg.Wait()

	return nil
}

func DownloadAlbums(ownerID int, albumIDs []string, logger *slog.Logger) {
	svc := vkscape.InitService(logger)
	var albums []models.PhotoAlbum

	if len(albumIDs) == 0 {
		// Not usable yet due to lack of permissions from VK side
		vkAlbums := svc.Client.GetAlbums(ownerID)
		albums = models.VkAlbumsToPhotoAlbums(vkAlbums)
		logger.Info("Found albums", "owner_id", ownerID, "count", len(albums))
	} else {
		albums = models.AlbumIDsToPhotoAlbums(albumIDs)
		logger.Info("Using provided albums", "owner_id", ownerID, "count", len(albumIDs)) // todo maybe get album info from VK
	}

	for _, album := range albums {
		albumDir := utils.CreateAlbumDirectory(album)
		logger.Info("Created album directory", "album_id", album.ID, "dir", albumDir)
		photos := svc.Client.GetPhotos(ownerID, strconv.Itoa(album.ID))
		logger.Info("Found photos", "album_id", album.ID, "count", len(photos))
		svc.Parser.ParseAlbumPhotos(&svc.Wg, albumDir, strconv.Itoa(album.ID), photos)
	}

	svc.Wg.Wait()
}

func InteractiveAuth(logger *slog.Logger) {
	err := auth.InteractiveFlow(logger)
	if err != nil {
		logger.Error("Authentication failed", "error", err)
		return
	}
	logger.Info("Authentication successful")
}

func AppTokenAuth(token string, logger *slog.Logger) {
	cfg := &config.AuthConfig{
		AuthMethod:  config.AuthMethodAppToken,
		AccessToken: token,
	}

	err := config.SaveConfig(cfg)
	if err != nil {
		logger.Error("Failed to save config", "error", err)
		return
	}
	logger.Info("App token authentication successful")
}
