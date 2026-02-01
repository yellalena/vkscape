package parser

import (
	"context"
	"fmt"
	"sync"

	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func (p *VKParser) ParseAlbumPhotos(
	ctx context.Context,
	wg *sync.WaitGroup,
	outputDir, albumID string,
	photos []vkObject.PhotosPhoto,
) {
	p.errs = make(chan error, len(photos))
	for _, photo := range photos {
		if ctx.Err() != nil {
			return
		}
		wg.Add(1)
		go func(photo vkObject.PhotosPhoto) {
			defer wg.Done()
			if ctx.Err() != nil {
				return
			}
			err := p.processPhoto(ctx, outputDir, albumID, photo)
			if err != nil {
				p.errs <- err
			}
		}(photo)
	}
}

func (p *VKParser) processPhoto(ctx context.Context, outputDir, albumID string, photo vkObject.PhotosPhoto) error {
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
	if ctx.Err() != nil {
		return ctx.Err()
	}
	err := downloadImage(ctx, photo.Sizes[len(photo.Sizes)-1].URL, outputDir, filename+".jpg")
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
