package cmd

import (
	"github.com/jazho76/devdeck/cli/internal/fsx"
	"github.com/jazho76/devdeck/cli/internal/nvim"
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/tmux"
	"github.com/jazho76/devdeck/cli/internal/ui"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the devdeck environment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}

		ui.Step("Uninstalling Neovim")
		if err := nvim.Uninstall(p); err != nil {
			return err
		}

		ui.Step("Uninstalling tmux")
		if err := tmux.Uninstall(p); err != nil {
			return err
		}

		ui.Step("Removing source")
		removed, err := fsx.RemoveDirIfExists(p.Source)
		if err != nil {
			return err
		}
		if removed {
			ui.Info("Removed source: %s", p.Source)
		} else {
			ui.Info("No source to remove: %s", p.Source)
		}

		ui.Info("\nDone.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
