package vkapi

import (
	"log"
	
	"github.com/SevereCloud/vksdk/v2/api"
	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

func GetPosts(groupID string, count int) ([]vkObject.WallWallpost) {
	res, err := VK.WallGet(api.Params{
		"owner_id": groupID,
		"count":    count,
	})
	
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return res.Items
}