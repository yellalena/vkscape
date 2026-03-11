package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/progress"
	"github.com/yellalena/vkscape/internal/utils"
	"github.com/yellalena/vkscape/internal/vkscape"
)

var conversationPhotosDownloadCmd = &cobra.Command{
	Use:   "conversation_photos",
	Short: utils.CommandConversationsPhotosDesc,
	Long: `Download photos from conversations by conversation IDs.

		Peer ID (conversation ID):
		- Numeric conversation ID
		- For group chat: 2000000000+ conversation's id
		- For user: user's id
		- For group: group's id 

		Examples:
		# Download photos from a conversation with user
		vkscape conversation_photos --peer_id 123456`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger, logFile := output.InitLogger(verbose)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		peerIdParam, err := cmd.Flags().GetString("peer_id")
		if err != nil {
			output.Error(fmt.Sprintf("Error getting peer_id flag: %v", err))
			logger.Error("Error getting peer_id flag", "error", err)
			return
		}
		if peerIdParam == "" {
			output.Error("Please provide at least one peer ID using --peer_id flag")
			return
		}
		peerId, err := strconv.Atoi(peerIdParam)
		if err != nil {
			output.Error("peer_id could not be converted to integer")
			return
		}


		output.Info(fmt.Sprintf("Starting download from conversation %d...", peerId))
		err = vkscape.DownloadConversationPhotos(context.Background(), peerId, logger, &progress.NoopReporter{})
		if err != nil {
			output.Error(fmt.Sprintf("Failed to download photos from conversation: %v", err))
			logger.Error("Error downloading photos from conversation", "error", err)
			return
		}
		output.Success(fmt.Sprintf("Successfully downloaded photos from conversation %d", peerId))
	},
}

func init() {
	conversationPhotosDownloadCmd.Flags().
		StringP("peer_id", "", "", "Numeric conversation ID to download photos from")
	conversationPhotosDownloadCmd.MarkFlagRequired("peer_id") //nolint:errcheck,gosec
	conversationPhotosDownloadCmd.Flags().
		BoolP("verbose", "v", false, "Enable verbose logging (output to both file and console)")
	rootCmd.AddCommand(conversationPhotosDownloadCmd)
}