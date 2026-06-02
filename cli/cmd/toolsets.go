package cmd

import (
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/spf13/cobra"
)

var toolsetsSelection toolsetSelection

var toolsetsCmd = &cobra.Command{
	Use:   "toolsets",
	Short: "Choose which Neovim toolsets are enabled",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}
		return configureToolsets(p, toolsetsSelection)
	},
}

func init() {
	addToolsetFlags(toolsetsCmd, &toolsetsSelection)
	rootCmd.AddCommand(toolsetsCmd)
}
