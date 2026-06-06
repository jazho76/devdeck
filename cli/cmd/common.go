package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/sysreq"
	"github.com/jazho76/devdeck/cli/internal/toolsets"
	"github.com/jazho76/devdeck/cli/internal/ui"
	"github.com/spf13/cobra"
)

func ensurePrereqs() error {
	unsatisfied := sysreq.Unsatisfied(true)
	if len(unsatisfied) == 0 {
		return nil
	}
	ui.Step("Checking prerequisites")
	for _, r := range unsatisfied {
		switch r.Status {
		case sysreq.TooOld:
			ui.StatusFail(r.Dep.Name, fmt.Sprintf("%s, minimum %s", r.Version, r.Dep.MinVersion))
		default:
			ui.StatusFail(r.Dep.Name, "not found on PATH")
		}
	}
	return fmt.Errorf("%d required dependency(ies) unsatisfied; resolve them and re-run", len(unsatisfied))
}

type toolsetSelection struct {
	all      bool
	none     bool
	explicit []string
}

func addToolsetFlags(cmd *cobra.Command, sel *toolsetSelection) {
	f := cmd.Flags()
	f.BoolVar(&sel.all, "all-toolsets", false, "enable every toolset, non-interactively")
	f.BoolVar(&sel.none, "no-toolsets", false, "disable every toolset, non-interactively")
	f.StringSliceVar(&sel.explicit, "toolsets", nil, "enable exactly these toolsets (comma-separated), non-interactively")
	cmd.MarkFlagsMutuallyExclusive("all-toolsets", "no-toolsets", "toolsets")
}

func configureToolsets(p paths.Paths, sel toolsetSelection) error {
	available, err := toolsets.Available(p)
	if err != nil {
		return err
	}

	switch {
	case sel.all:
		if err := toolsets.Write(p, available, asSet(available)); err != nil {
			return err
		}
		ui.Info("Enabled all toolsets")
		return nil
	case sel.none:
		if err := toolsets.Write(p, available, map[string]bool{}); err != nil {
			return err
		}
		ui.Info("Disabled all toolsets")
		return nil
	case len(sel.explicit) > 0:
		chosen, err := validateToolsets(sel.explicit, available)
		if err != nil {
			return err
		}
		if err := toolsets.Write(p, available, chosen); err != nil {
			return err
		}
		ui.Info("Enabled toolsets: %s", strings.Join(sel.explicit, ", "))
		return nil
	}

	if !isInteractive() {
		ui.Warn("no terminal for the toolset picker; keeping current selection (use --all-toolsets, --no-toolsets, or --toolsets)")
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

func validateToolsets(names, available []string) (map[string]bool, error) {
	valid := asSet(available)
	chosen := make(map[string]bool, len(names))
	for _, n := range names {
		if !valid[n] {
			return nil, fmt.Errorf("unknown toolset %q (available: %s)", n, strings.Join(available, ", "))
		}
		chosen[n] = true
	}
	return chosen, nil
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
