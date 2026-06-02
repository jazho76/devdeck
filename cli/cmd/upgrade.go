package cmd

import (
	"github.com/jazho76/devdeck/cli/internal/selfupdate"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the devdeck binary to the latest release",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return selfupdate.Run()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
