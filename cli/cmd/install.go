package cmd

import (
	"github.com/jazho76/devdeck/cli/internal/nvim"
	"github.com/jazho76/devdeck/cli/internal/source"
	"github.com/jazho76/devdeck/cli/internal/tmux"
	"github.com/spf13/cobra"
)

var installAll bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the devdeck environment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, log, err := setup(cmd)
		if err != nil {
			return err
		}

		log.Step("Fetching source")
		if err := source.EnsureClone(p, log); err != nil {
			return err
		}

		log.Step("Installing tmux")
		if err := tmux.Install(p, log); err != nil {
			return err
		}

		log.Step("Configuring toolsets")
		if err := configureToolsets(p, log, installAll); err != nil {
			return err
		}

		log.Step("Installing Neovim")
		if err := nvim.Install(p, log); err != nil {
			return err
		}

		log.Info("\nDone. Re-pick toolsets anytime with: devdeck toolsets")
		return nil
	},
}

func init() {
	installCmd.Flags().BoolVar(&installAll, "all", false, "enable every toolset, non-interactively")
	rootCmd.AddCommand(installCmd)
}
