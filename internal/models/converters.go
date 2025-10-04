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
