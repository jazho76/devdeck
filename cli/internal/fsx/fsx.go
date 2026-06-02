package fsx

import (
	"errors"
	"fmt"
	"io"
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
	Match
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

func InspectSymlink(link, expected string) (Outcome, error) {
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
	return Match, nil
}

func RemoveSymlinkIfPointsTo(link, expected string) (Outcome, error) {
	outcome, err := InspectSymlink(link, expected)
	if err != nil || outcome != Match {
		return outcome, err
	}

	if err := os.Remove(link); err != nil {
		return Absent, err
	}
	return Removed, nil
}

func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), perm); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

func CopyTree(src, dst string, skip func(rel string) bool) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if skip != nil && skip(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		target := filepath.Join(dst, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		return copyFile(path, target, info.Mode().Perm())
	})
}

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
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
