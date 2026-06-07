package sysreq

import (
	"fmt"
	"os/exec"
)

type Dep struct {
	Name        string
	Binaries    []string
	Required    bool
	MinVersion  string
	VersionArgs []string
}

var Catalog = []Dep{
	{Name: "git", Binaries: []string{"git"}, Required: true},
	{Name: "tmux", Binaries: []string{"tmux"}, Required: true, MinVersion: "3.5a", VersionArgs: []string{"-V"}},
	{Name: "make", Binaries: []string{"make"}},
	{Name: "gcc", Binaries: []string{"cc", "gcc"}, Required: true},
	{Name: "ripgrep", Binaries: []string{"rg"}, Required: true},
	{Name: "tree-sitter-cli", Binaries: []string{"tree-sitter"}, Required: true},
	{Name: "fd", Binaries: []string{"fd", "fdfind"}},
	{Name: "wl-clipboard", Binaries: []string{"wl-copy"}},
}

type Status int

const (
	Found Status = iota
	NotFound
	TooOld
)

type Result struct {
	Dep     Dep
	Status  Status
	Path    string
	Version string
}

func Lookup(name string) (Dep, bool) {
	for _, d := range Catalog {
		if d.Name == name {
			return d, true
		}
	}
	return Dep{}, false
}

func Check(d Dep) (string, bool) {
	for _, bin := range d.Binaries {
		if path, err := exec.LookPath(bin); err == nil {
			return path, true
		}
	}
	return "", false
}

func HasPkgConfig(module string) bool {
	return exec.Command("pkg-config", "--exists", module).Run() == nil
}

func Inspect(d Dep) Result {
	path, found := Check(d)
	if !found {
		return Result{Dep: d, Status: NotFound}
	}
	if d.MinVersion == "" {
		return Result{Dep: d, Status: Found, Path: path}
	}

	version, err := currentVersion(path, d.VersionArgs)
	if err != nil || version == "" {
		return Result{Dep: d, Status: Found, Path: path}
	}
	if !meetsMinimum(version, d.MinVersion) {
		return Result{Dep: d, Status: TooOld, Path: path, Version: version}
	}
	return Result{Dep: d, Status: Found, Path: path, Version: version}
}

func Unsatisfied(requiredOnly bool) []Result {
	var unsatisfied []Result
	for _, d := range Catalog {
		if requiredOnly && !d.Required {
			continue
		}
		if r := Inspect(d); r.Status != Found {
			unsatisfied = append(unsatisfied, r)
		}
	}
	return unsatisfied
}

func RequireCommand(name string) error {
	d, ok := Lookup(name)
	if !ok {
		d = Dep{Name: name, Binaries: []string{name}, Required: true}
	}
	if _, found := Check(d); !found {
		return fmt.Errorf("missing required command: %s", d.Name)
	}
	return nil
}
