package cmd

import (
	"github.com/jazho76/devdeck/cli/internal/nvim"
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/source"
	"github.com/jazho76/devdeck/cli/internal/tmux"
	"github.com/jazho76/devdeck/cli/internal/ui"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the devdeck source and plugins",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}

		if err := ensurePrereqs(); err != nil {
			return err
		}

		ui.Step("Updating source")
		if err := source.Pull(p); err != nil {
			return err
		}

		ui.Step("Updating tmux plugins")
		if err := tmux.Update(p); err != nil {
			return err
		}

		ui.Step("Updating Neovim plugins")
		if err := nvim.Update(p); err != nil {
			return err
		}

		ui.Info("\nDone.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
