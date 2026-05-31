package cmd

import (
	"os"

	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/toolsets"
	"github.com/jazho76/devdeck/cli/internal/ui"
	"github.com/spf13/cobra"
)

func setup(cmd *cobra.Command) (paths.Paths, ui.Logger, error) {
	p, err := paths.Resolve()
	if err != nil {
		return paths.Paths{}, ui.Logger{}, err
	}
	return p, ui.New(cmd.OutOrStdout()), nil
}

// configureToolsets is shared by the install and toolsets commands: it enables
// everything with --all, otherwise runs the interactive picker.
func configureToolsets(p paths.Paths, log ui.Logger, all bool) error {
	available, err := toolsets.Available(p)
	if err != nil {
		return err
	}

	if all {
		if err := toolsets.Write(p, available, asSet(available)); err != nil {
			return err
		}
		log.Info("Enabled all toolsets in %s - restart nvim to apply.", p.ToolsetsLocal())
		return nil
	}

	if !isInteractive() {
		log.Warn("no terminal for the toolset picker; keeping current selection (use --all to enable everything)")
		return nil
	}

	enabled, err := toolsets.Enabled(p, available)
	if err != nil {
		return err
	}

	result, err := toolsets.RunPicker(available, enabled)
	if err != nil {
		return err
	}
	if result.Cancelled {
		log.Info("Cancelled. No changes.")
		return nil
	}

	if err := toolsets.Write(p, available, result.Chosen); err != nil {
		return err
	}
	log.Info("Saved %s - restart nvim to apply.", p.ToolsetsLocal())
	return nil
}

func isInteractive() bool {
	if info, err := os.Stdin.Stat(); err == nil && info.Mode()&os.ModeCharDevice != 0 {
		return true
	}
	return controllingTerminalAvailable()
}

func controllingTerminalAvailable() bool {
	f, err := os.Open("/dev/tty")
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func asSet(names []string) map[string]bool {
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	return set
}
