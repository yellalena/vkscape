package vkscape

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/yellalena/vkscape/internal/config"
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/parser"
	"github.com/yellalena/vkscape/internal/vkapi"
)

type VkScapeService struct {
	Client vkapi.VKClient
	Parser parser.VKParser
	Wg     sync.WaitGroup
}

func InitService(logger *slog.Logger) (*VkScapeService, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		output.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		output.Error("Please authenticate first using 'vkscape auth'")
		logger.Error("Failed to load config", "error", err)
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client, err := vkapi.InitClient(cfg.AccessToken, logger)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to initialize VK client: %v", err))
		logger.Error("Failed to initialize VK client", "error", err)
		return nil, fmt.Errorf("init vk client: %w", err)
	}
	parser := parser.InitParser(logger)

	return &VkScapeService{
		Client: client,
		Parser: parser,
	}, nil
}
