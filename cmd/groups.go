package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yellalena/vkscape/internal/output"
)

var groupDownloadCmd = &cobra.Command{
	Use:   "groups",
	Short: "Download posts from groups",
	Long:  "Download posts from groups by their IDs",
	Args:  cobra.NoArgs,
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
		if ids == "" {
			output.Error("Please provide at least one group ID using --ids flag")
			return
		}
		idList := strings.Split(ids, ",")

		output.Info(fmt.Sprintf("Starting download for %d group(s)...", len(idList)))
		err = DownloadGroups(idList, logger)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to download groups: %v", err))
			logger.Error("Error downloading groups", "error", err)
			return
		}
		output.Success(fmt.Sprintf("Successfully downloaded posts from %d group(s)", len(idList)))
	},
}

func init() {
	groupDownloadCmd.Flags().
		StringP("ids", "", "", "Comma-separated list of group IDs to download posts from")
	groupDownloadCmd.Flags().
		BoolP("verbose", "v", false, "Enable verbose logging (output to both file and console)")
	rootCmd.AddCommand(groupDownloadCmd)
}
