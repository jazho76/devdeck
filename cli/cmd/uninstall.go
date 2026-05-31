package cmd

import (
	"github.com/jazho76/devdeck/cli/internal/fsx"
	"github.com/jazho76/devdeck/cli/internal/nvim"
	"github.com/jazho76/devdeck/cli/internal/tmux"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the devdeck environment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, log, err := setup(cmd)
		if err != nil {
			return err
		}

		log.Step("Uninstalling Neovim")
		if err := nvim.Uninstall(p, log); err != nil {
			return err
		}

		log.Step("Uninstalling tmux")
		if err := tmux.Uninstall(p, log); err != nil {
			return err
		}

		log.Step("Removing source")
		removed, err := fsx.RemoveDirIfExists(p.Source)
		if err != nil {
			return err
		}
		if removed {
			log.Info("Removed source: %s", p.Source)
		} else {
			log.Info("No source to remove: %s", p.Source)
		}

		log.Info("\nDone.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
