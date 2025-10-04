package vkapi

import (
	"log"

	"github.com/SevereCloud/vksdk/v2/api"
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
