package vkapi

import (
	"github.com/SevereCloud/vksdk/v2/api"
	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func (VK *VKClient) GetPosts(groupID string, count int) ([]vkObject.WallWallpost, error) {
	res, err := VK.Client.WallGet(api.Params{
		"owner_id": groupID,
		"count":    count,
	})

	if err != nil {
		return nil, err
	}

	return res.Items, nil
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
