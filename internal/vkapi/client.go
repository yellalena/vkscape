package vkapi

import (
	"fmt"
	"log/slog"

	"github.com/SevereCloud/vksdk/v2/api"
)

type VKClient struct {
	Client *api.VK
	logger *slog.Logger
}

func InitClient(token string, logger *slog.Logger) (VKClient, error) {
	if token == "" {
		logger.Error("VK Access token not found")
		return VKClient{}, fmt.Errorf("VK Access token not found")
	}

	VK := api.NewVK(token)
	logger.Info("VK API client initialized")

	return VKClient{
		Client: VK,
		logger: logger,
	}, nil
}

func (*VKClient) GetVersion() string {
	return api.Version
}
