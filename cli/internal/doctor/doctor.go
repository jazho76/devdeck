package doctor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jazho76/devdeck/cli/internal/fsx"
	"github.com/jazho76/devdeck/cli/internal/nvim"
	"github.com/jazho76/devdeck/cli/internal/paths"
	"github.com/jazho76/devdeck/cli/internal/run"
	"github.com/jazho76/devdeck/cli/internal/source"
	"github.com/jazho76/devdeck/cli/internal/sysreq"
	"github.com/jazho76/devdeck/cli/internal/toolsets"
)

type Severity int

const (
	OK Severity = iota
	Warn
	Fail
)

type Result struct {
	Name     string
	Severity Severity
	Detail   string
	Hint     string
	Indent   bool
	Heading  bool
}

func Run(p paths.Paths) []Result {
	var results []Result
	add := func(r Result) { results = append(results, r) }

	add(checkLocalBinOnPath(p))
	for _, d := range sysreq.Catalog {
		if d.Required {
			add(checkDep(d))
		}
	}
	add(checkSource(p))
	add(checkConfigLink("tmux config link", p.ConfigTmux, p.SourceTmux(), p.SourceTmuxConf()))
	add(checkTPM(p))
	add(checkConfigLink("Neovim config link", p.ConfigNvim, p.SourceNvim(), p.SourceNvimInit()))

	bin, ok := nvim.ResolveBin(p)
	add(checkNvimBinary(bin, ok))
	add(checkNvimVersion(bin, ok))
	add(checkNvimSmoke(bin, ok))

	for _, d := range sysreq.Catalog {
		if !d.Required {
			add(checkDep(d))
		}
	}

	add(checkToolsets(p))
	for _, r := range checkToolsetReqs(p) {
		add(r)
	}

	return results
}

func checkDep(d sysreq.Dep) Result {
	r := Result{Name: d.Name}
	switch res := sysreq.Inspect(d); res.Status {
	case sysreq.Found:
		r.Severity = OK
		r.Detail = res.Path
		if res.Version != "" {
			r.Detail = fmt.Sprintf("%s (%s)", res.Path, res.Version)
		}
	case sysreq.TooOld:
		r.Severity = Fail
		r.Detail = fmt.Sprintf("%s, minimum %s", res.Version, d.MinVersion)
	default:
		if d.Required {
			r.Severity = Fail
			r.Detail = "not found on PATH"
		} else {
			r.Severity = Warn
			r.Detail = "not found on PATH (optional)"
		}
	}
	return r
}

func checkLocalBinOnPath(p paths.Paths) Result {
	r := Result{Name: "~/.local/bin on PATH"}
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if filepath.Clean(dir) == filepath.Clean(p.LocalBin) {
			r.Severity = OK
			r.Detail = p.LocalBin
			return r
		}
	}
	r.Severity = Warn
	r.Detail = p.LocalBin + " not on PATH"
	r.Hint = "add " + p.LocalBin + " to PATH so the devdeck and nvim symlinks resolve"
	return r
}

func checkSource(p paths.Paths) Result {
	r := Result{Name: "Managed source", Detail: p.Source}
	if _, err := os.Stat(p.Source); err != nil {
		r.Severity = Fail
		r.Detail = "missing: " + p.Source
		r.Hint = "run devdeck install to fetch the managed source"
		return r
	}

	_, hasGit := os.Stat(filepath.Join(p.Source, ".git"))
	if hasGit != nil && os.Getenv(source.OverrideEnv) == "" {
		r.Severity = Fail
		r.Detail = p.Source + " is not a devdeck clone"
		r.Hint = "remove it and run devdeck install"
		return r
	}

	for _, sub := range []string{"nvim", "tmux", "README.md"} {
		if _, err := os.Stat(filepath.Join(p.Source, sub)); err != nil {
			r.Severity = Fail
			r.Detail = "missing " + sub + " under " + p.Source
			r.Hint = "the managed source is incomplete; run devdeck install"
			return r
		}
	}

	r.Severity = OK
	return r
}

func checkConfigLink(name, link, target, requiredFile string) Result {
	r := Result{Name: name, Detail: link}
	outcome, err := fsx.InspectSymlink(link, target)
	if err != nil {
		r.Severity = Fail
		r.Detail = err.Error()
		return r
	}

	switch outcome {
	case fsx.Absent:
		r.Severity = Fail
		r.Detail = "not linked: " + link
		r.Hint = "run devdeck install to link it to " + target
	case fsx.KeptNotSymlink:
		r.Severity = Warn
		r.Detail = link + " is an unmanaged path, not a devdeck symlink"
	case fsx.KeptUnmanaged:
		r.Severity = Warn
		r.Detail = link + " points elsewhere (unmanaged)"
	case fsx.Match:
		if _, err := os.Stat(requiredFile); err != nil {
			r.Severity = Fail
			r.Detail = "linked, but missing " + requiredFile
			r.Hint = "the managed source is incomplete; run devdeck install"
		} else {
			r.Severity = OK
			r.Detail = link + " -> " + target
		}
	}
	return r
}

func checkTPM(p paths.Paths) Result {
	r := Result{Name: "tmux plugin manager (TPM)", Detail: p.TPMDir()}
	if _, err := os.Stat(filepath.Join(p.TPMDir(), ".git")); err != nil {
		r.Severity = Fail
		r.Detail = "missing: " + p.TPMDir()
		r.Hint = "run devdeck install to clone TPM"
		return r
	}
	installer := filepath.Join(p.TPMDir(), "bin", "install_plugins")
	if !isExecutable(installer) {
		r.Severity = Fail
		r.Detail = installer + " is not executable"
		r.Hint = "run devdeck install to repair TPM"
		return r
	}
	r.Severity = OK
	return r
}

func checkNvimBinary(bin string, ok bool) Result {
	r := Result{Name: "Neovim binary"}
	if !ok {
		r.Severity = Fail
		r.Detail = "no nvim on PATH or devdeck-managed install"
		r.Hint = "run devdeck install to install Neovim"
		return r
	}
	r.Severity = OK
	r.Detail = bin
	return r
}

func checkNvimVersion(bin string, ok bool) Result {
	r := Result{Name: "Neovim version"}
	if !ok {
		r.Severity = Warn
		r.Detail = "skipped, no nvim binary"
		return r
	}
	v, err := nvim.Version(bin)
	if err != nil {
		r.Severity = Warn
		r.Detail = "could not determine nvim version"
		return r
	}
	if v == nvim.NvimVersion {
		r.Severity = OK
		r.Detail = v
		return r
	}
	r.Severity = Warn
	r.Detail = fmt.Sprintf("%s, expected %s", v, nvim.NvimVersion)
	r.Hint = "the config is tested against Neovim " + nvim.NvimVersion
	return r
}

func checkNvimSmoke(bin string, ok bool) Result {
	r := Result{Name: "Neovim starts"}
	if !ok {
		r.Severity = Fail
		r.Detail = "skipped, no nvim binary"
		return r
	}
	if _, err := run.Query(bin, "--headless", "-u", "NONE", "+qa"); err != nil {
		r.Severity = Fail
		r.Detail = "nvim --headless -u NONE +qa failed"
		r.Hint = err.Error()
		return r
	}
	r.Severity = OK
	r.Detail = "headless launch succeeded"
	return r
}

func checkToolsets(p paths.Paths) Result {
	r := Result{Name: "Toolsets config", Detail: p.ToolsetsLocal()}
	if _, err := os.Stat(p.ToolsetsLocal()); err != nil {
		r.Severity = Warn
		r.Detail = "missing: " + p.ToolsetsLocal()
		r.Hint = "run devdeck toolsets to configure enabled toolsets"
		return r
	}
	r.Severity = OK
	return r
}

func checkToolsetReqs(p paths.Paths) []Result {
	names, err := toolsets.EnabledNames(p)
	if err != nil {
		return nil
	}
	groups, err := toolsets.Requirements(p, names)
	if err != nil {
		return []Result{{Severity: Warn, Name: "toolset requirements", Detail: err.Error()}}
	}

	var results []Result
	for _, g := range groups {
		results = append(results, Result{Heading: true, Name: g.Toolset + " toolset"})
		for _, req := range g.Reqs {
			r := Result{Name: req.Label, Detail: req.Detail, Indent: true}
			if req.Found {
				r.Severity = OK
			} else {
				r.Severity = Warn
			}
			results = append(results, r)
		}
	}
	return results
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && info.Mode()&0o111 != 0
}

func Summary(results []Result) (ok, warn, fail int) {
	for _, r := range results {
		if r.Heading {
			continue
		}
		switch r.Severity {
		case OK:
			ok++
		case Warn:
			warn++
		case Fail:
			fail++
		}
	}
	return
}
