package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/vkscape"
)

var (
	tokenFlag string
	userFlag  bool
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with VK",
	Long:  "Authenticate with VK using whether an app token or user token (will open browser).", // todo
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		logger, logFile := output.InitLogger(verbose)
		if logFile != nil {
			defer logFile.Close() //nolint:errcheck
		}

		if tokenFlag == "" && !userFlag {
			output.Error("Please provide either --token or --user flag")
			return
		}

		if userFlag {
			output.Info("Starting interactive authentication...")
			vkscape.InteractiveAuth(logger)
		} else {
			output.Info("Authenticating with app token...")
			vkscape.AppTokenAuth(tokenFlag, logger)
		}

		output.Success("Authentication successful! You can now use other commands.")
	},
}

func init() {
	authCmd.Flags().StringVar(&tokenFlag, "token", "", "App token to use")
	authCmd.Flags().BoolVar(&userFlag, "user", false, "Use browser-based user authentication")
	authCmd.Flags().
		BoolP("verbose", "v", false, "Enable verbose logging (output to both file and console)")

	rootCmd.AddCommand(authCmd)
}
