package source

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jazho76/devdeck/cli/internal/fsx"
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/run"
	"github.com/jazho76/devdeck/cli/internal/sysreq"
	"github.com/jazho76/devdeck/cli/internal/ui"
)

const (
	URL    = "https://github.com/jazho76/devdeck"
	Branch = "main"

	// OverrideEnv names a local directory to copy as the managed source instead
	// of cloning from URL. Used by the sandbox to test the local working tree.
	OverrideEnv = "DEVDECK_SOURCE"
)

func EnsureClone(p paths.Paths) error {
	if override := os.Getenv(OverrideEnv); override != "" {
		return ensureFromLocal(p, override)
	}

	if err := sysreq.RequireCommand("git"); err != nil {
		return err
	}

	switch state, err := dirState(p.Source); {
	case err != nil:
		return err
	case state == cloned:
		return Pull(p)
	case state == foreign:
		return fmt.Errorf("refusing to manage existing path: %s", p.Source)
	default:
		ui.Info("Cloning %s into %s", URL, p.Source)
		if err := os.MkdirAll(filepath.Dir(p.Source), 0o755); err != nil {
			return err
		}
		return run.Stream("git", "clone", "--branch", Branch, URL, p.Source)
	}
}

func ensureFromLocal(p paths.Paths, override string) error {
	if fi, err := os.Stat(override); err != nil || !fi.IsDir() {
		return fmt.Errorf("%s %q is not a directory", OverrideEnv, override)
	}

	switch _, err := os.Stat(p.Source); {
	case err == nil:
		ui.Info("Using existing source at %s (%s set; not re-copying)", p.Source, OverrideEnv)
		return nil
	case errors.Is(err, fs.ErrNotExist):
		ui.Info("Copying source from %s into %s", override, p.Source)
		if err := os.MkdirAll(filepath.Dir(p.Source), 0o755); err != nil {
			return err
		}
		return fsx.CopyTree(override, p.Source, func(rel string) bool { return rel == ".git" })
	default:
		return err
	}
}

func Pull(p paths.Paths) error {
	if override := os.Getenv(OverrideEnv); override != "" {
		ui.Warn("%s set; skipping source update", OverrideEnv)
		return nil
	}

	ui.Info("Updating source")
	if err := run.Stream("git", "-C", p.Source, "pull", "--ff-only"); err != nil {
		return fmt.Errorf("%w\nmanaged source could not fast-forward; remove %s and re-run devdeck install", err, p.Source)
	}
	return nil
}

type state int

const (
	missing state = iota
	cloned
	foreign
)

func dirState(dir string) (state, error) {
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		return cloned, nil
	}
	switch _, err := os.Lstat(dir); {
	case errors.Is(err, fs.ErrNotExist):
		return missing, nil
	case err != nil:
		return missing, err
	default:
		return foreign, nil
	}
}
