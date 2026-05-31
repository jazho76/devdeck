package cmd

import (
	"github.com/jazho76/devdeck/cli/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the devdeck version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cmd.Println(version.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
