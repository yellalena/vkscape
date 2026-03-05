package vkapi

import (
	"context"

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

func (VK *VKClient) StreamConversationPhotos(
	ctx context.Context,
	peerID int,
) (<-chan vkObject.PhotosPhoto, <-chan error) {

	photos := make(chan vkObject.PhotosPhoto)
	errs := make(chan error, 1)

	go func() {
		defer close(photos)
		defer close(errs)

		count := 100
		startFrom := ""

		for {
			if ctx.Err() != nil {
				errs <- ctx.Err()
				return
			}

			params := api.Params{
				"media_type": "photo",
				"peer_id":    peerID,
				"count":      count,
			}

			if startFrom != "" {
				params["start_from"] = startFrom
			}

			res, err := VK.Client.MessagesGetHistoryAttachments(params)
			if err != nil {
				errs <- err
				return
			}

			for _, item := range res.Items {
				if item.Attachment.Type == "photo" {
					select {
					case photos <- item.Attachment.Photo:
					case <-ctx.Done():
						errs <- ctx.Err()
						return
					}
				}
			}

			if res.NextFrom == "" {
				return
			}

			startFrom = res.NextFrom
		}
	}()

	return photos, errs
}
