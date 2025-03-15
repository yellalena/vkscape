package main

import (
	"fmt"
	"log"

	"github.com/yellalena/vkscape/internal/vkapi"
)

func main() {
	vkapi.InitClient()

	var groupID string
	fmt.Println("Insert group ID:")
	fmt.Scanln(&groupID)

	var posts = vkapi.GetPosts(groupID, 5)

	for _, post := range posts {
		fmt.Println(post.Text)
		fmt.Println("––––––––––––––––––––––")
	}

	log.Println("Done.")
}