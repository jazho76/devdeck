package workspace

import (
	"errors"
	"fmt"
	"os"

	"github.com/jazho76/devdeck/cli/internal/paths"
)

func Delete(p paths.Paths, name string) error {
	slug := Slugify(name)
	if slug == "" {
		return fmt.Errorf("workspace name %q has no usable characters", name)
	}
	if err := os.Remove(p.WorkspaceFile(slug)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("no such workspace: %q", name)
		}
		return err
	}
	return nil
}

func Rename(p paths.Paths, oldName, newName string) (Workspace, error) {
	if err := ValidateName(newName); err != nil {
		return Workspace{}, err
	}
	newSlug := Slugify(newName)
	if newSlug == "" {
		return Workspace{}, fmt.Errorf("workspace name %q has no usable characters", newName)
	}

	oldSlug := Slugify(oldName)
	ws, err := load(p.WorkspaceFile(oldSlug))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Workspace{}, fmt.Errorf("no such workspace: %q", oldName)
		}
		return Workspace{}, err
	}

	if newSlug == oldSlug {
		ws.Name = newName
		if err := writeWorkspace(p, ws); err != nil {
			return Workspace{}, err
		}
		return ws, nil
	}

	if _, err := os.Stat(p.WorkspaceFile(newSlug)); err == nil {
		return Workspace{}, fmt.Errorf("workspace %q already exists", newName)
	} else if !errors.Is(err, os.ErrNotExist) {
		return Workspace{}, err
	}

	ws.Name = newName
	ws.Slug = newSlug
	if err := writeWorkspace(p, ws); err != nil {
		return Workspace{}, err
	}
	if err := os.Remove(p.WorkspaceFile(oldSlug)); err != nil {
		return Workspace{}, err
	}
	return ws, nil
}
