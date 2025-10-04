package vkapi

import (
	"log"

	"github.com/SevereCloud/vksdk/v2/api"
	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func (VK *VKClient) GetAlbums(ownerID int) []vkObject.PhotosPhotoAlbumFull {
	res, err := VK.Client.PhotosGetAlbums(api.Params{
		"owner_id": ownerID,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return res.Items
}

func (VK *VKClient) GetPhotos(ownerID int, albumID string) []vkObject.PhotosPhoto {
	res, err := VK.Client.PhotosGet(api.Params{
		"owner_id": ownerID,
		"album_id": albumID,
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return res.Items
}
