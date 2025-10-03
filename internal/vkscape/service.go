package vkscape

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/yellalena/vkscape/internal/parser"
	"github.com/yellalena/vkscape/internal/vkapi"
)

type VkScapeService struct {
	Client vkapi.VKClient
	Parser parser.VKParser
	Wg     sync.WaitGroup
}

func InitService() *VkScapeService {
	_ = godotenv.Load()
	token := os.Getenv("VK_API_KEY")

	client := vkapi.InitClient(token)
	parser := parser.InitParser()

	return &VkScapeService{
		Client: client,
		Parser: parser,
	}
}
