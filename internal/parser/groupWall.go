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

func (p *VKParser) ParseWallPosts(
	wg *sync.WaitGroup,
	outputDir string,
	posts []vkObject.WallWallpost,
) {
	p.errs = make(chan error, len(posts))
	for _, post := range posts {
		wg.Add(1)
		go func(post vkObject.WallWallpost) {
			defer wg.Done()
			if err := p.processPost(outputDir, post); err != nil {
				p.errs <- err
			}
		}(post)
	}
}

func (p *VKParser) processPost(outputDir string, post vkObject.WallWallpost) error {
	if post.PostType != PostTypePost || post.CopyHistory != nil {
		// Don't download reposts or non-posts
		return nil
	}

	if post.Text == "" && len(post.Attachments) == 0 {
		// Skip empty posts (non-image attachments)
		return nil
	}

	dateStr, err := convertDate(post.Date)
	if err != nil {
		p.logger.Error("Failed to convert post date", "error", err, "post_id", post.ID, "timestamp", post.Date)
		dateStr = fmt.Sprintf("unknown_%d", post.Date)
	}
	postName := fmt.Sprintf(PostFileNameTemplate, post.ID, dateStr)
	dirName, err := utils.CreateSubDirectory(outputDir, postName)
	if err != nil {
		p.logger.Error(
			"Failed to create subdirectory",
			"error",
			err,
			"post_id",
			post.ID,
			"output_dir",
			outputDir,
		)
		return err
	}

	err = utils.SaveFile(dirName, postName+".txt", []byte(post.Text))
	if err != nil {
		p.logger.Error("Failed to save post text", "error", err, "post_id", post.ID, "dir", dirName)
		return err
	}

	for _, attachment := range post.Attachments {
		if attachment.Type == "photo" {
			photo := attachment.Photo
			if len(photo.Sizes) == 0 {
				p.logger.Error(
					"Cannot download photo: no sizes available",
					"post_id",
					post.ID,
					"photo_id",
					photo.ID,
				)
				continue
			}
			filename := fmt.Sprintf(ImageFileNameTemplate, postName, photo.ID)
			err := downloadImage(photo.Sizes[len(photo.Sizes)-1].URL, dirName, filename+".jpg")
			if err != nil {
				p.logger.Error(
					"Failed to download image",
					"error",
					err,
					"post_id",
					post.ID,
					"photo_id",
					photo.ID,
					"url",
					photo.Sizes[len(photo.Sizes)-1].URL,
				)
				return err
			}
		}
	}

	return nil
}
