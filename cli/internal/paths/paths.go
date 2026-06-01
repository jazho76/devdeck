package paths

import (
	"os"
	"path/filepath"
)

type Paths struct {
	Home        string // ~
	Share       string // ~/.local/share/devdeck
	Source      string // managed clone, source of every symlink
	NvimInstall string // devdeck-managed nvim install
	LocalBin    string // ~/.local/bin
	NvimBin     string // ~/.local/bin/nvim
	ConfigNvim  string // ~/.config/nvim
	ConfigTmux  string // ~/.config/tmux
	TmuxData    string // ~/.local/share/tmux (TPM lives here, not under devdeck)
	NvimShare   string // ~/.local/share/nvim (plugins, Mason tools)
	NvimState   string // ~/.local/state/nvim (undo, marks, sessions)
	NvimCache   string // ~/.cache/nvim
	Workspaces  string // ~/.local/share/devdeck/workspaces
}

func Resolve() (Paths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Paths{}, err
	}

	share := filepath.Join(home, ".local", "share", "devdeck")
	localBin := filepath.Join(home, ".local", "bin")

	return Paths{
		Home:        home,
		Share:       share,
		Source:      filepath.Join(share, "source"),
		NvimInstall: filepath.Join(share, "nvim"),
		LocalBin:    localBin,
		NvimBin:     filepath.Join(localBin, "nvim"),
		ConfigNvim:  filepath.Join(home, ".config", "nvim"),
		ConfigTmux:  filepath.Join(home, ".config", "tmux"),
		TmuxData:    filepath.Join(home, ".local", "share", "tmux"),
		NvimShare:   filepath.Join(home, ".local", "share", "nvim"),
		NvimState:   filepath.Join(home, ".local", "state", "nvim"),
		NvimCache:   filepath.Join(home, ".cache", "nvim"),
		Workspaces:  filepath.Join(share, "workspaces"),
	}, nil
}

func (p Paths) WorkspaceFile(slug string) string {
	return filepath.Join(p.Workspaces, slug+".json")
}

func (p Paths) SourceNvim() string     { return filepath.Join(p.Source, "nvim") }
func (p Paths) SourceTmux() string     { return filepath.Join(p.Source, "tmux") }
func (p Paths) SourceNvimInit() string { return filepath.Join(p.SourceNvim(), "init.lua") }
func (p Paths) SourceTmuxConf() string { return filepath.Join(p.SourceTmux(), "tmux.conf") }
func (p Paths) TPMDir() string         { return filepath.Join(p.TmuxData, "plugins", "tpm") }

func (p Paths) ToolsetsDir() string {
	return filepath.Join(p.SourceNvim(), "lua", "toolsets")
}

func (p Paths) ToolsetsLocal() string {
	return filepath.Join(p.SourceNvim(), "lua", "config", "toolsets-local.lua")
}

func (p Paths) ToolsetsDefault() string {
	return filepath.Join(p.SourceNvim(), "lua", "config", "toolsets-default.lua")
}
