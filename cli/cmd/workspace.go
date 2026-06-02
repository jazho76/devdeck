package cmd

import (
	"errors"
	"fmt"
	"sort"
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
	Use:   "save [name]",
	Short: "Save the current tmux session as a workspace",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}
		name := workspace.CurrentSession()
		if len(args) == 1 {
			name = args[0]
		}
		if err := workspace.Save(p, name); err != nil {
			return err
		}
		ui.Info("Workspace saved: %s", name)
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

		res, err := ui.Prompt("Save workspace", workspace.CurrentSession(), validate, ui.FillHeight())
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

var workspaceRestorePopupCmd = &cobra.Command{
	Use:    "restore-popup",
	Short:  "Interactive restore prompt",
	Args:   cobra.NoArgs,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}
		if !isInteractive() {
			return errors.New("restore-popup requires a terminal")
		}

		saved, err := workspace.List(p)
		if err != nil {
			return err
		}
		if len(saved) == 0 {
			return ui.Notice("Restore workspace", "No saved workspaces yet", ui.FillHeight())
		}
		sort.Slice(saved, func(i, j int) bool {
			ti, _ := workspace.LastActivity(saved[i])
			tj, _ := workspace.LastActivity(saved[j])
			return ti.After(tj)
		})

		labels := workspace.RestoreLabels(saved, time.Now())

		choice, err := ui.SingleSelect("Restore workspace", labels, ui.FillHeight())
		if err != nil {
			return err
		}
		if choice.Cancelled {
			return nil
		}
		return workspace.Restore(p, saved[choice.Index].Slug)
	},
}

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved workspaces",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}
		saved, err := workspace.List(p)
		if err != nil {
			return err
		}
		if len(saved) == 0 {
			ui.Info("No saved workspaces")
			return nil
		}
		sort.Slice(saved, func(i, j int) bool {
			ti, _ := workspace.LastActivity(saved[i])
			tj, _ := workspace.LastActivity(saved[j])
			return ti.After(tj)
		})
		for _, label := range workspace.RestoreLabels(saved, time.Now()) {
			fmt.Println(label)
		}
		return nil
	},
}

var workspaceDeleteForce bool

var workspaceDeleteCmd = &cobra.Command{
	Use:     "delete <name>",
	Aliases: []string{"rm"},
	Short:   "Delete a saved workspace",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}
		ws, err := workspace.Load(p, args[0])
		if err != nil {
			return err
		}

		if !workspaceDeleteForce {
			if !isInteractive() {
				return errors.New("refusing to delete without confirmation; pass --force")
			}
			if !ui.Confirm(fmt.Sprintf("Delete workspace %q?", ws.Name)) {
				return nil
			}
		}

		if err := workspace.Delete(p, ws.Slug); err != nil {
			return err
		}
		ui.Info("Workspace deleted: %s", ws.Name)
		return nil
	},
}

var workspaceRenameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a saved workspace",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}
		ws, err := workspace.Rename(p, args[0], args[1])
		if err != nil {
			return err
		}
		ui.Info("Renamed to: %s", ws.Name)
		return nil
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceSaveCmd)
	workspaceCmd.AddCommand(workspaceSavePopupCmd)
	workspaceCmd.AddCommand(workspaceRestoreCmd)
	workspaceCmd.AddCommand(workspaceRestorePopupCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceDeleteCmd.Flags().BoolVarP(&workspaceDeleteForce, "force", "f", false, "skip the confirmation prompt")
	workspaceCmd.AddCommand(workspaceDeleteCmd)
	workspaceCmd.AddCommand(workspaceRenameCmd)
	rootCmd.AddCommand(workspaceCmd)
}
