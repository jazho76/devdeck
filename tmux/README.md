# Tmux Configuration

Devdeck tmux configuration.

Install it from the repository root:

```bash
./scripts/install-tmux.sh
```

The installer requires an existing `tmux` binary from your system package manager. It links `tmux/` into `~/.config/tmux`, installs or updates TPM at `~/.local/share/tmux/plugins/tpm`, and installs TPM plugins.

If the installed tmux version differs from the expected version, the installer prints a warning and continues. It never installs, replaces, or removes the tmux binary.

Uninstall Devdeck-managed tmux files:

```bash
./scripts/uninstall-tmux.sh
```

The uninstaller removes only the Devdeck-managed config symlink and tmux plugin directory. It leaves the tmux binary untouched.
