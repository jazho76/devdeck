package cmd

import (
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/ui"
	"github.com/jazho76/devdeck/cli/internal/workspace"
	"github.com/spf13/cobra"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage saved tmux workspaces",
	Args:  cobra.NoArgs,
}

var workspaceSaveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save the current tmux layout as a workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}
		if err := workspace.Save(p, args[0]); err != nil {
			return err
		}
		ui.Info("Workspace saved: %s", args[0])
		return nil
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceSaveCmd)
	rootCmd.AddCommand(workspaceCmd)
}
