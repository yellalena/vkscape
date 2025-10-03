package cmd

import (
	"fmt"

	"github.com/yellalena/vkscape/internal/utils"
	"github.com/yellalena/vkscape/internal/vkscape"
)

func DownloadGroups(groupIDs []string) {
	svc := vkscape.InitService()

	for _, groupID := range groupIDs {
		groupDir := utils.CreateGroupDirectory(groupID)
		fmt.Println("Created dir:", groupDir)
		posts := svc.Client.GetPosts(groupID, 5) // todo: remove limit
		fmt.Println("Found posts:", len(posts))
		svc.Parser.ParseWallPosts(&svc.Wg, groupDir, posts)
	}

	// utils.ZipFolder("vk_archive.zip")

	svc.Wg.Wait()
}
