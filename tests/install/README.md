# Devdeck install contract tests

These tests run Devdeck inside disposable containers and inspect the installed state from inside the container. They are intended to harden the installer/control-plane contract without touching a developer's real home directory.

## Run one scenario

```bash
tests/install/run.sh fedora fresh-minimal
```

## Run the fast matrix

```bash
tests/install/matrix.sh fast
```

## Current scenario groups

- `fast`: Fedora 44 scenarios suitable for pull-request CI.
- `full`: Fedora 44 scenarios that include all/no-toolset variants. This is reserved for slower expansion such as Ubuntu and workspace restore coverage.

## Contract under test

The scenario scripts assert that Devdeck:

- installs the CLI from the current mounted source tree
- uses `DEVDECK_SOURCE=/opt/devdeck-src` instead of cloning GitHub `main`
- links tmux and Neovim config to Devdeck-managed source
- keeps TPM under `~/.local/share/tmux/plugins/tpm` and configures `TMUX_PLUGIN_MANAGER_PATH` accordingly
- preserves unmanaged user config by failing closed
- keeps Neovim runtime state on uninstall unless `--purge` is used
- exposes the workspace command and `ws` alias

Each scenario owns a fresh `HOME` and pins all XDG paths. Keep scenarios small and outcome-based.
