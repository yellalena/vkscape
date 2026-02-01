package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/progress"
	"github.com/yellalena/vkscape/internal/utils"

	"github.com/yellalena/vkscape/internal/vkscape"
)

var albumDownloadCmd = &cobra.Command{
	Use:   "albums",
	Short: "Download photos from albums",
	Long: `Download photos from albums by their IDs or all albums for an owner.

		Privacy and Authentication Requirements:
		- With user token (--user): Only downloading specific albums by IDs is available.
			Albums must have public privacy settings.
		- With service token (--token): Downloading all albums by owner ID is available,
			but the profile must be open/public for everyone on the internet.

		Examples:
		# Download specific albums (works with both token types)
		vkscape albums --owner 123456 --ids 789,790

		# Download all albums (only works with service token, profile must be public)
		vkscape albums --owner 123456`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger, logFile := output.InitLogger(verbose)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		ids, err := cmd.Flags().GetString("ids")
		if err != nil {
			output.Error(fmt.Sprintf("Error getting ids flag: %v", err))
			logger.Error("Error getting ids flag", "error", err)
			return
		}
		owner, err := cmd.Flags().GetString("owner")
		if err != nil {
			output.Error(fmt.Sprintf("Error getting owner flag: %v", err))
			logger.Error("Error getting owner flag", "error", err)
			return
		}

		// Owner ID is required
		if owner == "" {
			output.Error("Please specify owner ID using --owner flag")
			return
		}

		ownerID, err := strconv.Atoi(owner)
		if err != nil {
			output.Error("Error: owner ID must be an integer")
			return
		}

		var idList []string
		if ids != "" {
			idList = utils.ParseIDList(ids)
		}

		output.Info(fmt.Sprintf("Starting download for owner %d...", ownerID))
		if len(idList) > 0 {
			output.Info(fmt.Sprintf("Downloading %d album(s)...", len(idList)))
		} else {
			output.Info("Fetching all albums for owner...")
		}

		if err := vkscape.DownloadAlbums(ownerID, idList, logger, &progress.NoopReporter{}); err != nil {
			output.Error(fmt.Sprintf("Failed to download albums: %v", err))
			logger.Error("Failed to download albums", "error", err)
			return
		}
		output.Success(fmt.Sprintf("Successfully downloaded albums for owner %d", ownerID))
	},
}

func init() {
	albumDownloadCmd.Flags().
		StringP("ids", "", "", "Comma-separated list of album IDs to download (optional). If not provided, all albums will be downloaded (requires service token).")
	albumDownloadCmd.Flags().
		StringP("owner", "", "", "ID of the user/owner to download albums from (required)")
	albumDownloadCmd.MarkFlagRequired("owner") //nolint:errcheck,gosec
	albumDownloadCmd.Flags().
		BoolP("verbose", "v", false, "Enable verbose logging (output to both file and console)")
	rootCmd.AddCommand(albumDownloadCmd)
}
