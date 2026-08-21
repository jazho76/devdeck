package cmd

import (
	"os"

	"github.com/jazho76/devdeck/cli/internal/run"
	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin <command> [args...]",
	Short: "Run a command and pin it to the workspace",
	Long: `Run a command and pin it to the workspace, so that "workspace save" records it
and "workspace restore" relaunches it.

Commands started any other way are not restored; nvim is the only exception.
Shell syntax needs an explicit shell:

  devdeck pin pnpm run dev
  devdeck pin sh -c 'pnpm run clean && pnpm run dev'`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code, err := run.Foreground(args[0], args[1:]...)
		if err != nil {
			return err
		}
		os.Exit(code)
		return nil
	},
}

func init() {
	pinCmd.Flags().SetInterspersed(false)
	rootCmd.AddCommand(pinCmd)
}
