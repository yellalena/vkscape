package cmd

import (
	"fmt"

	"github.com/yellalena/vkscape/internal/vkscape"
)

func DownloadGroups(groupIDs []string) {
	svc := vkscape.InitService()

	for _, groupID := range groupIDs {
		groupDir := vkscape.CreateGroupDirectory(groupID)
		fmt.Println("Created dir:", groupDir)
		posts := svc.Client.GetPosts(groupID, 20) // todo: remove limit
		fmt.Println("Found posts:", len(posts))
		svc.Parser.ParseWallPosts(&svc.Wg, groupDir, posts)
	}

	svc.Wg.Wait()
}
