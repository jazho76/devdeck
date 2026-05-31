package cmd

import (
	"fmt"
	"os"

	"github.com/jazho76/devdeck/cli/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "devdeck",
	Short:         "DevDeck control plane",
	Version:       version.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.SetVersionTemplate(version.String() + "\n")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
