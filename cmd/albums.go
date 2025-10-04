package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var albumDownloadCmd = &cobra.Command{
	Use:   "albums",
	Short: "Download photos from albums",
	Long:  "Download photos from albums by their IDs",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ids, err := cmd.Flags().GetString("ids")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		owner, err := cmd.Flags().GetString("owner")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if ids == "" || owner == "" {
			fmt.Println("Please specify both owner ID using --owner flag and at least one album ID using --ids flag")
			return
		}

		ownerID, err := strconv.Atoi(owner)
		if err != nil && owner != "" {
			fmt.Println("Error: owner ID must be an integer")
			return
		}

		idList := strings.Split(ids, ",")
		DownloadAlbums(ownerID, idList)
	},
}

func init() {
	albumDownloadCmd.Flags().StringP("ids", "", "", "Comma-separated list of group IDs to download posts from")
	albumDownloadCmd.Flags().StringP("owner", "", "", "ID of the user to download albums from")
	rootCmd.AddCommand(albumDownloadCmd)
}
