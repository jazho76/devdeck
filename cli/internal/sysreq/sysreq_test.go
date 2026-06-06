package sysreq

import (
	"os"
	"path/filepath"
	"testing"
)

func writeExecutable(t *testing.T, dir, name string) {
	t.Helper()
	writeScript(t, dir, name, "#!/bin/sh\n")
}

func writeScript(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o755); err != nil {
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

func TestMeetsMinimum(t *testing.T) {
	cases := []struct {
		current, min string
		want         bool
	}{
		{"3.6b", "3.5a", true},
		{"3.5a", "3.5a", true},
		{"3.10", "3.5a", true},
		{"3.4", "3.5a", false},
		{"3.5", "3.5a", false},
		{"4.0", "3.5a", true},
		{"3.5b", "3.5a", true},
	}
	for _, c := range cases {
		if got := meetsMinimum(c.current, c.min); got != c.want {
			t.Errorf("meetsMinimum(%q, %q) = %v, want %v", c.current, c.min, got, c.want)
		}
	}
}

func tmuxDep() Dep {
	return Dep{Name: "tmux", Binaries: []string{"tmux"}, Required: true, MinVersion: "3.5a", VersionArgs: []string{"-V"}}
}

func TestInspectVersion(t *testing.T) {
	cases := []struct {
		name    string
		version string
		want    Status
	}{
		{"newer", "3.6b", Found},
		{"equal", "3.5a", Found},
		{"older", "3.1", TooOld},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("PATH", dir)
			writeScript(t, dir, "tmux", "#!/bin/sh\necho \"tmux "+c.version+"\"\n")

			if got := Inspect(tmuxDep()); got.Status != c.want {
				t.Fatalf("status = %v, want %v (version %s)", got.Status, c.want, got.Version)
			}
		})
	}
}

func TestInspectIgnoresVersionWithoutMinVersion(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir)
	writeExecutable(t, dir, "make")

	if got := Inspect(Dep{Name: "make", Binaries: []string{"make"}}); got.Status != Found {
		t.Fatalf("status = %v, want Found", got.Status)
	}
}

func TestInspectNotFound(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	if got := Inspect(tmuxDep()); got.Status != NotFound {
		t.Fatalf("status = %v, want NotFound", got.Status)
	}
}

func TestUnsatisfiedReportsTooOld(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir)
	writeExecutable(t, dir, "git")
	writeScript(t, dir, "tmux", "#!/bin/sh\necho \"tmux 3.1\"\n")

	unsatisfied := Unsatisfied(true)
	var tmux *Result
	for i := range unsatisfied {
		if unsatisfied[i].Dep.Name == "tmux" {
			tmux = &unsatisfied[i]
		}
	}
	if tmux == nil {
		t.Fatal("expected tmux in unsatisfied required deps")
	}
	if tmux.Status != TooOld {
		t.Fatalf("tmux status = %v, want TooOld", tmux.Status)
	}
}

func TestUnsatisfiedRequiredOnly(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir)
	writeExecutable(t, dir, "git")
	writeScript(t, dir, "tmux", "#!/bin/sh\necho \"tmux 3.6b\"\n")

	if got := Unsatisfied(true); len(got) != 0 {
		t.Fatalf("required deps satisfied, want none unsatisfied, got %d", len(got))
	}

	got := Unsatisfied(false)
	if len(got) == 0 {
		t.Fatal("optional deps absent, want them reported when requiredOnly is false")
	}
	for _, r := range got {
		if r.Dep.Required {
			t.Fatalf("required dep %s reported despite being satisfied", r.Dep.Name)
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
