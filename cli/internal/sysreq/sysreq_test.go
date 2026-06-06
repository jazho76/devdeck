package sysreq

import (
	"os"
	"path/filepath"
	"testing"
)

func writeExecutable(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestLookup(t *testing.T) {
	if _, ok := Lookup("git"); !ok {
		t.Fatal("git not found in catalog")
	}
	if _, ok := Lookup("nonexistent"); ok {
		t.Fatal("unexpected catalog hit for nonexistent dep")
	}
}

func TestCheckResolvesFirstAvailableBinary(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir)
	writeExecutable(t, dir, "fdfind")

	d := Dep{Name: "fd", Binaries: []string{"fd", "fdfind"}}
	path, found := Check(d)
	if !found {
		t.Fatal("expected fdfind to satisfy the fd dep")
	}
	if path != filepath.Join(dir, "fdfind") {
		t.Fatalf("path = %s, want %s", path, filepath.Join(dir, "fdfind"))
	}
}

func TestCheckMissing(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if _, found := Check(Dep{Name: "fd", Binaries: []string{"fd", "fdfind"}}); found {
		t.Fatal("expected fd dep to be missing on empty PATH")
	}
}

func TestMissingRequiredOnly(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir)
	writeExecutable(t, dir, "git")
	writeExecutable(t, dir, "tmux")

	if missing := Missing(true); len(missing) != 0 {
		t.Fatalf("required deps satisfied, want none missing, got %d", len(missing))
	}

	missing := Missing(false)
	if len(missing) == 0 {
		t.Fatal("optional deps absent, want them reported when requiredOnly is false")
	}
	for _, d := range missing {
		if d.Required {
			t.Fatalf("required dep %s reported missing despite being on PATH", d.Name)
		}
	}
}

func TestMissingReportsAbsentRequired(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	missing := Missing(true)
	names := map[string]bool{}
	for _, d := range missing {
		names[d.Name] = true
	}
	for _, want := range []string{"git", "tmux"} {
		if !names[want] {
			t.Fatalf("expected %s in missing required deps", want)
		}
	}
}

func TestRequireCommandError(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if err := RequireCommand("git"); err == nil {
		t.Fatal("expected error for missing git")
	}

	dir := t.TempDir()
	t.Setenv("PATH", dir)
	writeExecutable(t, dir, "git")
	if err := RequireCommand("git"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
