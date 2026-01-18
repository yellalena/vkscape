package vkapi

import (
	"github.com/SevereCloud/vksdk/v2/api"
	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func (VK *VKClient) GetPosts(groupID string) ([]vkObject.WallWallpost, error) {
	var allPosts []vkObject.WallWallpost
	offset := 0
	count := 100 // VK API max per request

	for {
		res, err := VK.Client.WallGet(api.Params{
			"owner_id": groupID,
			"count":    count,
			"offset":   offset,
		})
		if err != nil {
			return nil, err
		}

		allPosts = append(allPosts, res.Items...)

		if len(res.Items) < count {
			break
		}
		offset += count
	}

	return allPosts, nil
}

func (VK *VKClient) GetWallPostByID(postID string) vkObject.WallWallpost {
	res, err := VK.Client.WallGetByID(api.Params{
		"posts": postID,
	})

	if err != nil {
		VK.logger.Error("Failed to get wall post by ID", "error", err, "post_id", postID)
		panic(err)
	}

	if len(res) > 0 {
		return res[0]
	}

	return vkObject.WallWallpost{}
}
