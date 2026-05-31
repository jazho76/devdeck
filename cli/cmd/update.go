package cmd

import (
	"github.com/jazho76/devdeck/cli/internal/nvim"
	"github.com/jazho76/devdeck/cli/internal/source"
	"github.com/jazho76/devdeck/cli/internal/tmux"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the devdeck source and plugins",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, log, err := setup(cmd)
		if err != nil {
			return err
		}

		log.Step("Updating source")
		if err := source.Pull(p, log); err != nil {
			return err
		}

		log.Step("Updating tmux plugins")
		if err := tmux.Update(p, log); err != nil {
			return err
		}

		log.Step("Updating Neovim plugins")
		if err := nvim.Update(p, log); err != nil {
			return err
		}

		log.Info("\nDone.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
