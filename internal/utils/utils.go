package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/yellalena/vkscape/internal/models"
)

const (
	OutputDir      = "vkscape_output"
	OutputGroupDir = "group_%s"
	OutputAlbumDir = "album_%d"
)

func CreateGroupDirectory(groupID string) (string, error) {
	groupDir := filepath.Join(OutputDir, fmt.Sprintf(OutputGroupDir, groupID))
	err := os.MkdirAll(groupDir, 0755)
	return groupDir, err
}

func CreateAlbumDirectory(album models.PhotoAlbum) string {
	albumDir := filepath.Join(OutputDir, fmt.Sprintf(OutputAlbumDir, album.ID))
	_ = os.MkdirAll(albumDir, 0755)
	_ = SaveFile(albumDir, "album_info.txt", fmt.Appendf(nil, "Title: %s\nDescription: %s\nID: %d\n", album.Title, album.Description, album.ID))
	return albumDir
}

func CreateSubDirectory(parentDir, subDir string) (string, error) {
	dir := filepath.Join(parentDir, subDir)
	err := os.MkdirAll(dir, 0755)
	return dir, err
}

func SaveFile(parentDir, filename string, content []byte) error {
	filePath := filepath.Join(parentDir, filename)
	return os.WriteFile(filePath, content, 0644)
}

func SaveObject(parentDir, filename string, content io.ReadCloser) error {
	filePath := filepath.Join(parentDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, content)
	return err
}
