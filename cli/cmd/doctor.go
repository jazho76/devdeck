package cmd

import (
	"fmt"

	"github.com/jazho76/devdeck/cli/internal/doctor"
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/ui"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check that the devdeck installation is healthy",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		p, err := paths.Resolve()
		if err != nil {
			return err
		}

		ui.Step("Checking devdeck installation")
		results := doctor.Run(p)
		for _, r := range results {
			switch r.Severity {
			case doctor.OK:
				ui.StatusOK(r.Name, r.Detail)
			case doctor.Warn:
				ui.StatusWarn(r.Name, r.Detail)
			case doctor.Fail:
				ui.StatusFail(r.Name, r.Detail)
			}
			if r.Severity != doctor.OK && r.Hint != "" {
				ui.Hint(r.Hint)
			}
		}

		ok, warn, fail := doctor.Summary(results)
		ui.Info("\n%d ok, %d warning(s), %d failure(s)", ok, warn, fail)
		if fail > 0 {
			return fmt.Errorf("%d check(s) failed", fail)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
