package tmux

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jazho76/devdeck/cli/internal/fsx"
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/run"
	"github.com/jazho76/devdeck/cli/internal/sysreq"
	"github.com/jazho76/devdeck/cli/internal/ui"
)

const (
	ExpectedVersion = "3.5a"
	tpmURL          = "https://github.com/tmux-plugins/tpm"
)

var ErrInsideTmux = errors.New("cannot run inside an existing tmux session; exit tmux and try again")

func Install(p paths.Paths) error {
	if err := sysreq.RequireCommand("git"); err != nil {
		return err
	}
	if _, err := os.Stat(p.SourceTmuxConf()); err != nil {
		return fmt.Errorf("tmux config not found: %s", p.SourceTmuxConf())
	}
	if err := sysreq.RequireCommand("tmux"); err != nil {
		return errors.New("tmux is not installed; install it with your package manager, then try again")
	}
	if os.Getenv("TMUX") != "" {
		return ErrInsideTmux
	}

	if v, err := installedVersion(); err == nil && v != ExpectedVersion {
		ui.Warn("tmux version is %s, expected %s; continuing without modifying the binary", v, ExpectedVersion)
	}

	if err := fsx.EnsureSymlink(p.SourceTmux(), p.ConfigTmux); err != nil {
		return err
	}
	ui.Info("Linked tmux config")

	if err := ensureTPM(p); err != nil {
		return err
	}
	if err := run.Stream(filepath.Join(p.TPMDir(), "bin", "install_plugins")); err != nil {
		return err
	}

	ui.Info("Done. Start tmux with: tmux")
	return nil
}

func Update(p paths.Paths) error {
	if _, err := os.Stat(filepath.Join(p.TPMDir(), ".git")); err != nil {
		ui.Info("No TPM checkout to update: %s", p.TPMDir())
		return nil
	}
	if err := run.Stream("git", "-C", p.TPMDir(), "pull", "--ff-only"); err != nil {
		ui.Warn("could not update TPM: %v", err)
	}
	return nil
}

func Uninstall(p paths.Paths) error {
	outcome, err := fsx.RemoveSymlinkIfPointsTo(p.ConfigTmux, p.SourceTmux())
	if err != nil {
		return err
	}
	ui.Info("%s", fsx.Describe(outcome, p.ConfigTmux))

	plugins := filepath.Join(p.TmuxData, "plugins")
	removed, err := fsx.RemoveDirIfExists(plugins)
	if err != nil {
		return err
	}
	if removed {
		ui.Info("Removed tmux plugins: %s", plugins)
	}

	ui.Info("Kept tmux binary untouched.")
	return nil
}

func ensureTPM(p paths.Paths) error {
	tpm := p.TPMDir()
	if _, err := os.Stat(filepath.Join(tpm, ".git")); err == nil {
		return run.Stream("git", "-C", tpm, "pull", "--ff-only")
	}
	if _, err := os.Lstat(tpm); err == nil {
		return fmt.Errorf("refusing to overwrite existing path: %s", tpm)
	}
	if err := os.MkdirAll(filepath.Dir(tpm), 0o755); err != nil {
		return err
	}
	ui.Info("Cloning TPM into %s", tpm)
	return run.Stream("git", "clone", tpmURL, tpm)
}

func installedVersion() (string, error) {
	out, err := run.Output("tmux", "-V")
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(out, "tmux "), nil
}
