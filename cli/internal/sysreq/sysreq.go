package sysreq

import (
	"fmt"
	"os/exec"
)

type Dep struct {
	Name     string
	Binaries []string
	Required bool
}

var Catalog = []Dep{
	{Name: "git", Binaries: []string{"git"}, Required: true},
	{Name: "tmux", Binaries: []string{"tmux"}, Required: true},
	{Name: "make", Binaries: []string{"make"}},
	{Name: "fd", Binaries: []string{"fd", "fdfind"}},
	{Name: "wl-clipboard", Binaries: []string{"wl-copy"}},
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

func Missing(requiredOnly bool) []Dep {
	var missing []Dep
	for _, d := range Catalog {
		if requiredOnly && !d.Required {
			continue
		}
		if _, found := Check(d); !found {
			missing = append(missing, d)
		}
	}
	return missing
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
