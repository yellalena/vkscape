package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var groupDownloadCmd = &cobra.Command{
	Use:   "groups",
	Short: "Download posts from groups",
	Long:  "Download posts from groups by their IDs",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ids, err := cmd.Flags().GetString("ids")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if ids == "" {
			fmt.Println("Please provide at least one group ID using --ids flag")
			return
		}
		idList := strings.Split(ids, ",")

		err = DownloadGroups(idList)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	},
}

func init() {
	groupDownloadCmd.Flags().StringP("ids", "", "", "Comma-separated list of group IDs to download posts from")
	rootCmd.AddCommand(groupDownloadCmd)
}
