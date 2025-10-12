package vkscape

import (
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

func InitService() *VkScapeService {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error()) // todo
	}

	client := vkapi.InitClient(cfg.AccessToken)
	parser := parser.InitParser()

	return &VkScapeService{
		Client: client,
		Parser: parser,
	}
}
