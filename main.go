package main

import (
	"fmt"
	"log"
	"os"

	"sync"

	"github.com/joho/godotenv"
	"github.com/yellalena/vkscape/internal/parser"
	"github.com/yellalena/vkscape/internal/vkapi"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("VK_API_KEY")

	client := vkapi.InitClient(token)
	parser := parser.InitParser(client)

	var wg sync.WaitGroup

	groupID := "shantibiotic"
	// fmt.Println("Insert group ID:")
	// fmt.Scanln(&groupID)

	// todo: move
	groupDir := fmt.Sprintf("output/group_%s", groupID)
	os.MkdirAll(groupDir, 0755)

	posts := client.GetPosts(groupID, 20)
	client.GetWallPostById("-47521143_162")
	parser.ParseWallPosts(&wg, groupDir, posts)

	wg.Wait()

	log.Println("Done.")
}
