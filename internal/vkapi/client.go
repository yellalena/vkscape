package vkapi

import (
	"log/slog"

	"github.com/SevereCloud/vksdk/v2/api"
)

type VKClient struct {
	Client *api.VK
	logger *slog.Logger
}

func InitClient(token string, logger *slog.Logger) VKClient {
	if token == "" {
		logger.Error("VK_API_KEY not found")
		panic("VK_API_KEY not found")
	}

	VK := api.NewVK(token)
	logger.Info("VK API client initialized")

	return VKClient{
		Client: VK,
		logger: logger,
	}
}

func (*VKClient) GetVersion() string {
	return api.Version
}
