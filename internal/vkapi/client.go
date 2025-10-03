package vkapi

import (
	"log"

	"github.com/SevereCloud/vksdk/v2/api"
	vkObject "github.com/SevereCloud/vksdk/v2/object"
)

type VKClient struct {
	Client *api.VK
}

func InitClient(token string) VKClient {
	if token == "" {
		log.Fatal("VK_API_KEY not found.")
	}

	VK := api.NewVK(token)
	log.Println("VK API client initialized.")

	return VKClient{
		Client: VK,
	}
}

func (*VKClient) GetVersion() string {
	return api.Version
}

func (VK *VKClient) GetPosts(groupID string, count int) []vkObject.WallWallpost {
	res, err := VK.Client.WallGet(api.Params{
		"owner_id": groupID,
		"count":    count,
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return res.Items
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

func (VK *VKClient) GetPhotos(photoIDs []string) []vkObject.PhotosPhoto {
	res, err := VK.Client.PhotosGetByID(api.Params{
		"photos": photoIDs,
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return res
}
