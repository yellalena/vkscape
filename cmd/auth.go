package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yellalena/vkscape/internal/output"
	"github.com/yellalena/vkscape/internal/utils"
	"github.com/yellalena/vkscape/internal/vkscape"
)

var (
	tokenFlag string
	userFlag  bool
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: utils.CommandAuthDesc,
	Long: `Authenticate with VK using either an app token or user token.

		Authentication methods:
		- App token (--token): Use a pre-generated app token.
		- User auth (--user): Opens browser for interactive login.

		Examples:
		# Authenticate with app token
		vkscape auth --token YOUR_TOKEN

		# Authenticate with user flow
		vkscape auth --user`,
	Args: cobra.NoArgs,
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
			if err := vkscape.InteractiveAuth(logger); err != nil {
				return
			}
		} else {
			if err := vkscape.AppTokenAuth(tokenFlag, logger); err != nil {
				return
			}
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
