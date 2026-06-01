package cmd

import (
	"errors"
	"fmt"
	"time"

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

var workspaceSavePopupCmd = &cobra.Command{
	Use:    "save-popup",
	Short:  "Interactive save prompt",
	Args:   cobra.NoArgs,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}
		if !isInteractive() {
			return errors.New("save-popup requires a terminal")
		}

		existing, err := workspace.List(p)
		if err != nil {
			return err
		}
		saved := make(map[string]time.Time, len(existing))
		for _, w := range existing {
			saved[w.Slug] = w.UpdatedAt
		}

		validate := func(name string) (string, bool) {
			if err := workspace.ValidateName(name); err != nil {
				return err.Error(), false
			}
			slug := workspace.Slugify(name)
			if slug == "" {
				return "name has no usable characters", false
			}
			if t, ok := saved[slug]; ok {
				return fmt.Sprintf("overwrites existing (saved %s)", t.Local().Format("2006-01-02")), true
			}
			return "", true
		}

		res, err := ui.Prompt("Save workspace", validate)
		if err != nil {
			return err
		}
		if res.Cancelled {
			return nil
		}
		if err := workspace.Save(p, res.Value); err != nil {
			return err
		}
		workspace.Notify("Workspace saved: " + res.Value)
		return nil
	},
}

var workspaceRestoreCmd = &cobra.Command{
	Use:   "restore <name>",
	Short: "Replace the current tmux server with a saved workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}
		return workspace.Restore(p, args[0])
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceSaveCmd)
	workspaceCmd.AddCommand(workspaceSavePopupCmd)
	workspaceCmd.AddCommand(workspaceRestoreCmd)
	rootCmd.AddCommand(workspaceCmd)
}
