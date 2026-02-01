package vkapi

import (
	"github.com/SevereCloud/vksdk/v2/api"
	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func (VK *VKClient) GetAlbums(ownerID int) ([]vkObject.PhotosPhotoAlbumFull, error) {
	res, err := VK.Client.PhotosGetAlbums(api.Params{
		"owner_id": ownerID,
	})
	if err != nil {
		VK.logger.Error("Failed to get albums", "error", err, "owner_id", ownerID)
		return nil, err
	}

	return res.Items, nil
}

func (VK *VKClient) GetPhotos(ownerID int, albumID string) ([]vkObject.PhotosPhoto, error) {
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
			VK.logger.Error(
				"Failed to get photos",
				"error",
				err,
				"owner_id",
				ownerID,
				"album_id",
				albumID,
			)
			return allPhotos, err
		}

		allPhotos = append(allPhotos, res.Items...)

		if len(res.Items) < count {
			break // no more photos
		}

		offset += count
	}

	return allPhotos, nil
}
