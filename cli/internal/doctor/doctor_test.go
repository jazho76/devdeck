package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/sysreq"
)

func healthyFixture(t *testing.T) paths.Paths {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)

	p, err := paths.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	writeFile(t, p.SourceNvimInit(), "-- init")
	writeFile(t, p.SourceTmuxConf(), "# tmux")
	writeFile(t, filepath.Join(p.Source, "README.md"), "# devdeck")
	writeFile(t, filepath.Join(p.Source, ".git", "HEAD"), "ref: refs/heads/main")
	writeFile(t, p.ToolsetsLocal(), "return {}")

	symlink(t, p.SourceTmux(), p.ConfigTmux)
	symlink(t, p.SourceNvim(), p.ConfigNvim)

	writeFile(t, filepath.Join(p.TPMDir(), ".git", "HEAD"), "ref: refs/heads/master")
	writeExecutable(t, filepath.Join(p.TPMDir(), "bin", "install_plugins"))

	return p
}

func TestChecksPassOnHealthyFixture(t *testing.T) {
	p := healthyFixture(t)

	cases := []struct {
		name string
		got  Result
	}{
		{"source", checkSource(p)},
		{"tmux link", checkConfigLink("tmux config link", p.ConfigTmux, p.SourceTmux(), p.SourceTmuxConf())},
		{"nvim link", checkConfigLink("Neovim config link", p.ConfigNvim, p.SourceNvim(), p.SourceNvimInit())},
		{"tpm", checkTPM(p)},
		{"toolsets", checkToolsets(p)},
	}
	for _, c := range cases {
		if c.got.Severity != OK {
			t.Errorf("%s: severity = %v (%s), want OK", c.name, c.got.Severity, c.got.Detail)
		}
	}
}

func TestCheckSourceFailsWhenMissing(t *testing.T) {
	p := healthyFixture(t)
	if err := os.RemoveAll(p.Source); err != nil {
		t.Fatal(err)
	}
	if got := checkSource(p); got.Severity != Fail {
		t.Fatalf("severity = %v, want Fail", got.Severity)
	}
}

func TestCheckSourceFailsWhenIncomplete(t *testing.T) {
	p := healthyFixture(t)
	if err := os.Remove(filepath.Join(p.Source, "README.md")); err != nil {
		t.Fatal(err)
	}
	if got := checkSource(p); got.Severity != Fail {
		t.Fatalf("severity = %v, want Fail", got.Severity)
	}
}

func TestCheckConfigLinkFailsWhenAbsent(t *testing.T) {
	p := healthyFixture(t)
	if err := os.Remove(p.ConfigNvim); err != nil {
		t.Fatal(err)
	}
	got := checkConfigLink("Neovim config link", p.ConfigNvim, p.SourceNvim(), p.SourceNvimInit())
	if got.Severity != Fail {
		t.Fatalf("severity = %v, want Fail", got.Severity)
	}
}

func TestCheckConfigLinkWarnsWhenUnmanaged(t *testing.T) {
	p := healthyFixture(t)
	if err := os.Remove(p.ConfigNvim); err != nil {
		t.Fatal(err)
	}
	symlink(t, filepath.Join(p.Home, "elsewhere"), p.ConfigNvim)

	got := checkConfigLink("Neovim config link", p.ConfigNvim, p.SourceNvim(), p.SourceNvimInit())
	if got.Severity != Warn {
		t.Fatalf("severity = %v, want Warn", got.Severity)
	}
}

func TestCheckConfigLinkFailsWhenTargetFileMissing(t *testing.T) {
	p := healthyFixture(t)
	if err := os.Remove(p.SourceNvimInit()); err != nil {
		t.Fatal(err)
	}
	got := checkConfigLink("Neovim config link", p.ConfigNvim, p.SourceNvim(), p.SourceNvimInit())
	if got.Severity != Fail {
		t.Fatalf("severity = %v, want Fail", got.Severity)
	}
}

func TestCheckTPMFailsWhenInstallerNotExecutable(t *testing.T) {
	p := healthyFixture(t)
	installer := filepath.Join(p.TPMDir(), "bin", "install_plugins")
	if err := os.Chmod(installer, 0o644); err != nil {
		t.Fatal(err)
	}
	if got := checkTPM(p); got.Severity != Fail {
		t.Fatalf("severity = %v, want Fail", got.Severity)
	}
}

func TestCheckToolsetsWarnsWhenMissing(t *testing.T) {
	p := healthyFixture(t)
	if err := os.Remove(p.ToolsetsLocal()); err != nil {
		t.Fatal(err)
	}
	if got := checkToolsets(p); got.Severity != Warn {
		t.Fatalf("severity = %v, want Warn", got.Severity)
	}
}

func TestCheckDep(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir)
	writeExecutable(t, filepath.Join(dir, "git"))
	writeScript(t, filepath.Join(dir, "tmux"), "#!/bin/sh\necho \"tmux 3.1\"\n")

	cases := []struct {
		name string
		dep  sysreq.Dep
		want Severity
	}{
		{"present", sysreq.Dep{Name: "git", Binaries: []string{"git"}, Required: true}, OK},
		{"required absent", sysreq.Dep{Name: "make", Binaries: []string{"make"}, Required: true}, Fail},
		{"optional absent", sysreq.Dep{Name: "fd", Binaries: []string{"fd", "fdfind"}}, Warn},
		{"too old", sysreq.Dep{Name: "tmux", Binaries: []string{"tmux"}, Required: true, MinVersion: "3.5a", VersionArgs: []string{"-V"}}, Fail},
	}
	for _, c := range cases {
		if got := checkDep(c.dep); got.Severity != c.want {
			t.Errorf("%s: severity = %v, want %v", c.name, got.Severity, c.want)
		}
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeExecutable(t *testing.T, path string) {
	t.Helper()
	writeScript(t, path, "#!/bin/sh\n")
}

func writeScript(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func symlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
}
