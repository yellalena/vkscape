package models

import (
	"strconv"

	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func VkAlbumsToPhotoAlbums(albumObjects []vkObject.PhotosPhotoAlbumFull) []PhotoAlbum {
	albums := make([]PhotoAlbum, len(albumObjects))

	for i, album := range albumObjects {
		albums[i] = PhotoAlbum{
			ID:          album.ID,
			Title:       album.Title,
			Description: album.Description,
		}
	}

	return albums
}

func AlbumIDsToPhotoAlbums(albumIDs []string) ([]PhotoAlbum, []string) {
	albums := make([]PhotoAlbum, 0, len(albumIDs))
	var invalid []string

	for _, id := range albumIDs {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			invalid = append(invalid, id)
			continue
		}
		albums = append(albums, PhotoAlbum{
			ID:          idInt,
			Title:       id,
			Description: "",
		})
	}

	return albums, invalid
}

func FilterAlbumsByIDs(albumIDs []string, allAlbums []PhotoAlbum) ([]PhotoAlbum, []string) {
	if len(albumIDs) == 0 {
		return allAlbums, nil
	}
	// Create a set of requested IDs for quick lookup
	requestedIDs := make(map[int]bool)
	var invalid []string
	for _, idStr := range albumIDs {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			invalid = append(invalid, idStr)
			continue
		}
		requestedIDs[idInt] = true
	}

	albums := make([]PhotoAlbum, 0, len(albumIDs))
	for _, album := range allAlbums {
		if requestedIDs[album.ID] {
			albums = append(albums, album)
		}
	}

	return albums, invalid
}
