package parser

import (
	"fmt"
	"sync"

	vkObject "github.com/SevereCloud/vksdk/v2/object"
	"github.com/yellalena/vkscape/internal/utils"
)

const (
	PostFileNameTemplate  = "post_%d_%s"
	ImageFileNameTemplate = "%s_%d"
)

func (p *VKParser) ParseWallPosts(wg *sync.WaitGroup, outputDir string, posts []vkObject.WallWallpost) {
	for _, post := range posts {
		wg.Add(1)
		go func(post vkObject.WallWallpost) {
			defer wg.Done()
			processPost(outputDir, post)
		}(post)
	}
}

func processPost(outputDir string, post vkObject.WallWallpost) {
	if post.PostType != PostTypePost || post.CopyHistory != nil {
		// Don't download reposts or non-posts
		return
	}

	if post.Text == "" && len(post.Attachments) == 0 {
		// Skip empty posts (non-image attachments)
	}

	post_name := fmt.Sprintf(PostFileNameTemplate, post.ID, convertDate(post.Date))
	dir_name, err := utils.CreateSubDirectory(outputDir, post_name)
	if err != nil {
		// todo: log error, return
		panic(err)
	}

	err = utils.SaveFile(dir_name, post_name+".txt", []byte(post.Text))
	if err != nil {
		// todo: log error, return
		panic(err)
	}

	for _, attachment := range post.Attachments {
		switch attachment.Type {
		case "photo":
			photo := attachment.Photo
			filename := fmt.Sprintf(ImageFileNameTemplate, post_name, photo.ID)
			err := downloadImage(photo.Sizes[len(photo.Sizes)-1].BaseImage.URL, dir_name, filename+".jpg") // todo: process errors
			if err != nil {
				// todo: log error, continue
				fmt.Println("Error downloading image:", err)
			}
		}
	}
}
