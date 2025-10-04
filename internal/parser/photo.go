package parser

import (
	"fmt"
	"sync"

	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func (p *VKParser) ParseAlbumPhotos(wg *sync.WaitGroup, outputDir, albumID string, photos []vkObject.PhotosPhoto) {
	for _, photo := range photos {
		wg.Add(1)
		go func(photo vkObject.PhotosPhoto) {
			defer wg.Done()
			processPhoto(outputDir, albumID, photo)
		}(photo)
	}
}

func processPhoto(outputDir, albumID string, photo vkObject.PhotosPhoto) {
	filename := fmt.Sprintf(ImageFileNameTemplate, albumID, photo.ID)
	downloadImage(photo.Sizes[len(photo.Sizes)-1].BaseImage.URL, outputDir, filename+".jpg") // todo: process errors
}
