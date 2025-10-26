package vkscape

import (
	"log/slog"
	"sync"

	"github.com/yellalena/vkscape/internal/config"
	"github.com/yellalena/vkscape/internal/parser"
	"github.com/yellalena/vkscape/internal/vkapi"
)

type VkScapeService struct {
	Client vkapi.VKClient
	Parser parser.VKParser
	Wg     sync.WaitGroup
}

func InitService(logger *slog.Logger) *VkScapeService {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		panic("failed to load config: " + err.Error())
	}

	client := vkapi.InitClient(cfg.AccessToken, logger)
	parser := parser.InitParser(logger)

	return &VkScapeService{
		Client: client,
		Parser: parser,
	}
}
