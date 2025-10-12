package cmd

import (
	"fmt"
	"strconv"

	"github.com/yellalena/vkscape/internal/auth"
	"github.com/yellalena/vkscape/internal/config"
	"github.com/yellalena/vkscape/internal/models"
	"github.com/yellalena/vkscape/internal/utils"
	"github.com/yellalena/vkscape/internal/vkscape"
)

// todo improve logging

func DownloadGroups(groupIDs []string) error {
	svc := vkscape.InitService()
	// todo: add some check for groupID format - number should be negative
	// maybe separate command for parsin profile, but reusing same GetPosts method
	for _, groupID := range groupIDs {
		groupDir, err := utils.CreateGroupDirectory(groupID)
		if err != nil {
			return fmt.Errorf("failed to create directory for group %s: %w", groupID, err)
		}
		fmt.Println("Created dir:", groupDir)
		posts, err := svc.Client.GetPosts(groupID, 5) // todo: remove limit
		if err != nil {
			return fmt.Errorf("failed to get posts for group %s: %w", groupID, err)
		}
		fmt.Println("Found posts:", len(posts))
		svc.Parser.ParseWallPosts(&svc.Wg, groupDir, posts)
	}

	svc.Wg.Wait()

	return nil
}

func DownloadAlbums(ownerID int, albumIDs []string) {
	svc := vkscape.InitService()
	var albums []models.PhotoAlbum

	if len(albumIDs) == 0 {
		vkAlbums := svc.Client.GetAlbums(ownerID)
		albums = models.VkAlbumsToPhotoAlbums(vkAlbums)
		fmt.Println("Found albums:", len(albumIDs))
	} else {
		albums = models.AlbumIDsToPhotoAlbums(albumIDs)
		fmt.Println("Using provided albums:", len(albumIDs)) // todo maybe get album info from VK
	}

	for _, album := range albums {
		albumDir := utils.CreateAlbumDirectory(album)
		fmt.Println("Created dir:", albumDir)
		photos := svc.Client.GetPhotos(ownerID, strconv.Itoa(album.ID))
		fmt.Println("Found photos:", len(photos))
		svc.Parser.ParseAlbumPhotos(&svc.Wg, albumDir, strconv.Itoa(album.ID), photos)
	}

	svc.Wg.Wait()
}

func InteractiveAuth() {
	err := auth.InteractiveFlow()
	if err != nil {
		fmt.Println("Authentication failed:", err)
		return
	}
}

func AppTokenAuth(token string) {
	cfg := &config.AuthConfig{
		AuthMethod:  config.AuthMethodAppToken,
		AccessToken: token,
	}

	err := config.SaveConfig(cfg)
	if err != nil {
		fmt.Println("Failed to save config:", err)
		return
	}
}
