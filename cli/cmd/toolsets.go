package cmd

import (
	"github.com/spf13/cobra"
)

var toolsetsAll bool

var toolsetsCmd = &cobra.Command{
	Use:   "toolsets",
	Short: "Choose which Neovim toolsets are enabled",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, log, err := setup(cmd)
		if err != nil {
			return err
		}
		return configureToolsets(p, log, toolsetsAll)
	},
}

func init() {
	toolsetsCmd.Flags().BoolVar(&toolsetsAll, "all", false, "enable every toolset, non-interactively")
	rootCmd.AddCommand(toolsetsCmd)
}
