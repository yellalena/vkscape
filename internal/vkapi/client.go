package vkapi

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/SevereCloud/vksdk/v2/api"
)

var VK *api.VK

func InitClient() {
	_ = godotenv.Load()

	token := os.Getenv("VK_API_KEY")
	if token == "" {
		log.Fatal("VK_API_KEY not found.")
	}

	VK = api.NewVK(token)
	log.Println("VK API client initialized.")
}