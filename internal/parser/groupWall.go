package parser

import (
	"os"
	"strconv"
	"sync"

	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func (p *VKParser) ParseWallPosts(wg *sync.WaitGroup, outputDir string, posts []vkObject.WallWallpost) {
	for _, post := range posts {
		wg.Add(1)
		go func(post vkObject.WallWallpost) {
			defer wg.Done()
			p.processPost(outputDir, post)
		}(post)
	}
}

func (p *VKParser) processPost(outputDir string, post vkObject.WallWallpost) {
	if post.PostType != PostTypePost || post.CopyHistory != nil {
		// Don't download reposts or non-posts
		return
	}

	if post.Text == "" && len(post.Attachments) == 0 {
		// Skip empty posts (non-image attachments)
	}

	post_name := "post_" + strconv.Itoa(post.ID) + "_" + convertDate(post.Date)
	dir_name := outputDir + "/" + post_name
	os.MkdirAll(dir_name, 0755)

	text_filename := post_name + ".txt"
	post_text := post.Text
	os.WriteFile(dir_name+"/"+text_filename, []byte(post_text), 0644)

	for _, attachment := range post.Attachments {
		switch attachment.Type {
		case "photo":
			photo := attachment.Photo
			downloadImage(photo.Sizes[len(photo.Sizes)-1].BaseImage.URL, dir_name+"/"+post_name+"_"+strconv.Itoa(photo.ID)+".jpg")
		}
	}
}
