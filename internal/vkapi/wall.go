package vkapi

import (
	"log"

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

func (VK *VKClient) GetWallPostById(postID string) vkObject.WallWallpost {
	res, err := VK.Client.WallGetByID(api.Params{
		"posts": postID,
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(res) > 0 {
		return res[0]
	}

	return vkObject.WallWallpost{}
}
