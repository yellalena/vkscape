package parser

import (
	"fmt"
	"sync"

	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func (p *VKParser) ParseAlbumPhotos(
	wg *sync.WaitGroup,
	outputDir, albumID string,
	photos []vkObject.PhotosPhoto,
) {
	p.errs = make(chan error, len(photos))
	for _, photo := range photos {
		wg.Add(1)
		go func(photo vkObject.PhotosPhoto) {
			defer wg.Done()
			err := p.processPhoto(outputDir, albumID, photo)
			if err != nil {
				p.errs <- err
			}
		}(photo)
	}
}

func (p *VKParser) processPhoto(outputDir, albumID string, photo vkObject.PhotosPhoto) error {
	if len(photo.Sizes) == 0 {
		p.logger.Error(
			"Photo has no sizes",
			"album_id",
			albumID,
			"photo_id",
			photo.ID,
		)
		return fmt.Errorf("photo %d has no sizes", photo.ID)
	}
	filename := fmt.Sprintf(ImageFileNameTemplate, albumID, photo.ID)
	err := downloadImage(photo.Sizes[len(photo.Sizes)-1].URL, outputDir, filename+".jpg")
	if err != nil {
		p.logger.Error(
			"Failed to download photo",
			"error",
			err,
			"album_id",
			albumID,
			"photo_id",
			photo.ID,
			"url",
			photo.Sizes[len(photo.Sizes)-1].URL,
		)
		return err
	}
	return nil
}
