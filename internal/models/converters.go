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

func AlbumIDsToPhotoAlbums(albumIDs []string) []PhotoAlbum {
	albums := make([]PhotoAlbum, len(albumIDs))

	for i, id := range albumIDs {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			panic(err) // todo
		}
		albums[i] = PhotoAlbum{
			ID:          idInt,
			Title:       id,
			Description: "",
		}
	}

	return albums
}

func FilterAlbumsByIDs(albumIDs []string, allAlbums []PhotoAlbum) []PhotoAlbum {
	// Create a set of requested IDs for quick lookup
	requestedIDs := make(map[int]bool)
	for _, idStr := range albumIDs {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
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

	return albums
}
