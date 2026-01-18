package vkapi

import (
	"github.com/SevereCloud/vksdk/v2/api"
	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func (VK *VKClient) GetAlbums(ownerID int) []vkObject.PhotosPhotoAlbumFull {
	res, err := VK.Client.PhotosGetAlbums(api.Params{
		"owner_id": ownerID,
	})
	if err != nil {
		VK.logger.Error("Failed to get albums", "error", err, "owner_id", ownerID)
		panic(err)
	}

	return res.Items
}

func (VK *VKClient) GetPhotos(ownerID int, albumID string) []vkObject.PhotosPhoto {
	var allPhotos []vkObject.PhotosPhoto
	offset := 0
	count := 100 // max allowed by API

	for {
		res, err := VK.Client.PhotosGet(api.Params{
			"owner_id": ownerID,
			"album_id": albumID,
			"count":    count,
			"offset":   offset,
		})
		if err != nil {
			VK.logger.Error("Failed to get photos", "error", err, "owner_id", ownerID, "album_id", albumID)
			break
		}

		allPhotos = append(allPhotos, res.Items...)

		if len(res.Items) < count {
			break // no more photos
		}

		offset += count
	}

	return allPhotos
}
