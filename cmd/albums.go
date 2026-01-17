package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yellalena/vkscape/internal/output"
)

var albumDownloadCmd = &cobra.Command{
	Use:   "albums",
	Short: "Download photos from albums",
	Long:  "Download photos from albums by their IDs",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger, logFile := output.InitLogger(verbose)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		ids, err := cmd.Flags().GetString("ids")
		if err != nil {
			logger.Error("Error getting ids flag", "error", err)
			return
		}
		owner, err := cmd.Flags().GetString("owner")
		if err != nil {
			logger.Error("Error getting owner flag", "error", err)
			return
		}

		if ids == "" || owner == "" {
			logger.Error(
				"Please specify both owner ID using --owner flag and at least one album ID using --ids flag",
			)
			return
		}

		ownerID, err := strconv.Atoi(owner)
		if err != nil && owner != "" {
			logger.Error("Error: owner ID must be an integer")
			return
		}

		idList := strings.Split(ids, ",")
		DownloadAlbums(ownerID, idList, logger)
	},
}

func init() {
	albumDownloadCmd.Flags().
		StringP("ids", "", "", "Comma-separated list of group IDs to download posts from")
	albumDownloadCmd.Flags().StringP("owner", "", "", "ID of the user to download albums from")
	albumDownloadCmd.Flags().
		BoolP("verbose", "v", false, "Enable verbose logging (output to both file and console)")
	rootCmd.AddCommand(albumDownloadCmd)
}
