package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
		if tokenFlag == "" && !userFlag {
			fmt.Println("Please provide either --token or --user flag")
			return
		}

		if userFlag {
			InteractiveAuth()
		} else {
			AppTokenAuth(tokenFlag)
		}

		fmt.Println("Authentication successful, you can now use other commands.")
	},
}

func init() {
	authCmd.Flags().StringVar(&tokenFlag, "token", "", "App token to use")
	authCmd.Flags().BoolVar(&userFlag, "user", false, "Use browser-based user authentication")

	rootCmd.AddCommand(authCmd)
}
