package fsx

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var ErrUnmanagedTarget = errors.New("refusing to overwrite existing path")

func EnsureSymlink(target, link string) error {
	info, err := os.Lstat(link)
	switch {
	case err == nil:
		if info.Mode()&os.ModeSymlink == 0 {
			return fmt.Errorf("%w: %s", ErrUnmanagedTarget, link)
		}
		if err := os.Remove(link); err != nil {
			return err
		}
	case !errors.Is(err, fs.ErrNotExist):
		return err
	}

	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		return err
	}
	return os.Symlink(target, link)
}

type Outcome int

const (
	Absent Outcome = iota
	Removed
	KeptUnmanaged
	KeptNotSymlink
)

func Describe(o Outcome, link string) string {
	switch o {
	case Removed:
		return "Removed symlink: " + link
	case KeptUnmanaged:
		return "Keeping unmanaged symlink: " + link
	case KeptNotSymlink:
		return "Keeping unmanaged path: " + link
	default:
		return "Nothing to remove: " + link
	}
}

func RemoveSymlinkIfPointsTo(link, expected string) (Outcome, error) {
	info, err := os.Lstat(link)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return Absent, nil
	case err != nil:
		return Absent, err
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return KeptNotSymlink, nil
	}

	dest, err := os.Readlink(link)
	if err != nil {
		return Absent, err
	}
	if filepath.Clean(dest) != filepath.Clean(expected) {
		return KeptUnmanaged, nil
	}

	if err := os.Remove(link); err != nil {
		return Absent, err
	}
	return Removed, nil
}

func RemoveDirIfExists(path string) (bool, error) {
	if _, err := os.Lstat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	if err := os.RemoveAll(path); err != nil {
		return false, err
	}
	return true, nil
}
