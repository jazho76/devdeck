package cmd

import (
	"os"

	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/toolsets"
	"github.com/jazho76/devdeck/cli/internal/ui"
)

func configureToolsets(p paths.Paths, all bool) error {
	available, err := toolsets.Available(p)
	if err != nil {
		return err
	}

	if all {
		if err := toolsets.Write(p, available, asSet(available)); err != nil {
			return err
		}
		ui.Info("Enabled all toolsets")
		return nil
	}

	if !isInteractive() {
		ui.Warn("no terminal for the toolset picker; keeping current selection (use --all to enable everything)")
		return nil
	}

	enabled, err := toolsets.Enabled(p, available)
	if err != nil {
		return err
	}

	result, err := ui.MultiSelect("Select toolsets to enable", available, enabled)
	if err != nil {
		return err
	}
	if result.Cancelled {
		ui.Info("Cancelled. No changes.")
		return nil
	}

	if err := toolsets.Write(p, available, result.Selected); err != nil {
		return err
	}
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
