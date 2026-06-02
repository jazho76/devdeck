package fsx

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureSymlinkCreatesWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "sub", "link")

	if err := EnsureSymlink(target, link); err != nil {
		t.Fatalf("EnsureSymlink: %v", err)
	}
	if got, _ := os.Readlink(link); got != target {
		t.Fatalf("link points at %q, want %q", got, target)
	}
}

func TestEnsureSymlinkOverwritesExistingSymlink(t *testing.T) {
	dir := t.TempDir()
	link := filepath.Join(dir, "link")
	if err := os.Symlink(filepath.Join(dir, "old"), link); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "new")

	if err := EnsureSymlink(target, link); err != nil {
		t.Fatalf("EnsureSymlink: %v", err)
	}
	if got, _ := os.Readlink(link); got != target {
		t.Fatalf("link points at %q, want %q", got, target)
	}
}

func TestEnsureSymlinkRefusesNonSymlink(t *testing.T) {
	dir := t.TempDir()
	link := filepath.Join(dir, "real")
	if err := os.WriteFile(link, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := EnsureSymlink(filepath.Join(dir, "target"), link)
	if !errors.Is(err, ErrUnmanagedTarget) {
		t.Fatalf("err = %v, want ErrUnmanagedTarget", err)
	}
}

func TestRemoveSymlinkIfPointsTo(t *testing.T) {
	const expected = "/some/expected/target"

	t.Run("removes matching symlink", func(t *testing.T) {
		link := filepath.Join(t.TempDir(), "link")
		if err := os.Symlink(expected, link); err != nil {
			t.Fatal(err)
		}
		got, err := RemoveSymlinkIfPointsTo(link, expected)
		if err != nil {
			t.Fatal(err)
		}
		if got != Removed {
			t.Fatalf("outcome = %v, want Removed", got)
		}
		if _, err := os.Lstat(link); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("link still present: %v", err)
		}
	})

	t.Run("removes matching symlink even when target is gone", func(t *testing.T) {
		dir := t.TempDir()
		gone := filepath.Join(dir, "deleted-source")
		link := filepath.Join(dir, "link")
		if err := os.Symlink(gone, link); err != nil {
			t.Fatal(err)
		}
		got, err := RemoveSymlinkIfPointsTo(link, gone)
		if err != nil {
			t.Fatal(err)
		}
		if got != Removed {
			t.Fatalf("outcome = %v, want Removed", got)
		}
	})

	t.Run("keeps unmanaged symlink", func(t *testing.T) {
		link := filepath.Join(t.TempDir(), "link")
		if err := os.Symlink("/other/target", link); err != nil {
			t.Fatal(err)
		}
		got, err := RemoveSymlinkIfPointsTo(link, expected)
		if err != nil {
			t.Fatal(err)
		}
		if got != KeptUnmanaged {
			t.Fatalf("outcome = %v, want KeptUnmanaged", got)
		}
		if _, err := os.Lstat(link); err != nil {
			t.Fatalf("link was removed: %v", err)
		}
	})

	t.Run("keeps real directory", func(t *testing.T) {
		real := filepath.Join(t.TempDir(), "real")
		if err := os.Mkdir(real, 0o755); err != nil {
			t.Fatal(err)
		}
		got, err := RemoveSymlinkIfPointsTo(real, expected)
		if err != nil {
			t.Fatal(err)
		}
		if got != KeptNotSymlink {
			t.Fatalf("outcome = %v, want KeptNotSymlink", got)
		}
	})

	t.Run("absent path", func(t *testing.T) {
		got, err := RemoveSymlinkIfPointsTo(filepath.Join(t.TempDir(), "nope"), expected)
		if err != nil {
			t.Fatal(err)
		}
		if got != Absent {
			t.Fatalf("outcome = %v, want Absent", got)
		}
	})
}

func TestCopyTree(t *testing.T) {
	src := t.TempDir()
	mustWrite(t, filepath.Join(src, "top.txt"), "top", 0o644)
	mustWrite(t, filepath.Join(src, "nested", "exec.sh"), "#!/bin/sh", 0o755)
	mustWrite(t, filepath.Join(src, ".git", "config"), "secret", 0o644)

	dst := filepath.Join(t.TempDir(), "out")
	if err := CopyTree(src, dst, func(rel string) bool { return rel == ".git" }); err != nil {
		t.Fatalf("CopyTree: %v", err)
	}

	if got, _ := os.ReadFile(filepath.Join(dst, "top.txt")); string(got) != "top" {
		t.Fatalf("top.txt = %q, want %q", got, "top")
	}

	exec := filepath.Join(dst, "nested", "exec.sh")
	info, err := os.Stat(exec)
	if err != nil {
		t.Fatalf("stat exec.sh: %v", err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Fatalf("exec.sh mode = %o, want 0755", info.Mode().Perm())
	}

	if _, err := os.Stat(filepath.Join(dst, ".git")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf(".git was copied: %v", err)
	}
}

func mustWrite(t *testing.T, path, content string, perm os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveDirIfExists(t *testing.T) {
	t.Run("removes present dir", func(t *testing.T) {
		sub := filepath.Join(t.TempDir(), "sub")
		if err := os.MkdirAll(filepath.Join(sub, "nested"), 0o755); err != nil {
			t.Fatal(err)
		}
		removed, err := RemoveDirIfExists(sub)
		if err != nil {
			t.Fatal(err)
		}
		if !removed {
			t.Fatal("removed = false, want true")
		}
		if _, err := os.Lstat(sub); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("dir still present: %v", err)
		}
	})

	t.Run("absent dir", func(t *testing.T) {
		removed, err := RemoveDirIfExists(filepath.Join(t.TempDir(), "nope"))
		if err != nil {
			t.Fatal(err)
		}
		if removed {
			t.Fatal("removed = true, want false")
		}
	})
}
