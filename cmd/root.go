package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yellalena/vkscape/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "vkscape",
	Short: "vkscape is a cli tool for downloading your VK archive",
	Long:  "vkscape is a cli tool for downloading your VK archive - the ones you can't get officially",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			tui.Start()
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing VkScape '%s'\n", err)
		os.Exit(1)
	}
}
