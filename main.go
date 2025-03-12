package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/SevereCloud/vksdk/v2/api"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("VK_API_KEY")
	if token == "" {
		log.Fatal("VK_API_KEY not found.")
	}

	var groupID string
	fmt.Println("Insert group ID:")
	fmt.Scanln(&groupID)

	vk := api.NewVK(token)

	res, err := vk.GroupsGetByID(api.Params{
		"group_id": groupID,
		"fields":   "description,members_count", 
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Retrieved group info:")
	fmt.Println(fmt.Sprintf("Group name: %s, \nGroup description: %s", res[0].Name, res[0].Description))
}