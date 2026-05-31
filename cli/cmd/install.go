package cmd

import (
	"github.com/jazho76/devdeck/cli/internal/nvim"
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/source"
	"github.com/jazho76/devdeck/cli/internal/tmux"
	"github.com/jazho76/devdeck/cli/internal/ui"
	"github.com/spf13/cobra"
)

var installAll bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the devdeck environment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}

		ui.Step("Fetching source")
		if err := source.EnsureClone(p); err != nil {
			return err
		}

		ui.Step("Installing tmux")
		if err := tmux.Install(p); err != nil {
			return err
		}

		ui.Step("Configuring toolsets")
		if err := configureToolsets(p, installAll); err != nil {
			return err
		}

		ui.Step("Installing Neovim")
		if err := nvim.Install(p); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	installCmd.Flags().BoolVar(&installAll, "all", false, "enable every toolset, non-interactively")
	rootCmd.AddCommand(installCmd)
}
