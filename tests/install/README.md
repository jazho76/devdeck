# Devdeck install contract tests

These tests run Devdeck inside a disposable Fedora 44 container and inspect the installed state from inside the container. They are intentionally small: one happy-path install plus fail-closed checks for existing user config.

## Run one scenario

```bash
tests/install/run.sh fedora fresh-minimal
```

## Run the matrix

```bash
tests/install/matrix.sh fast
```

`full` currently aliases `fast` until there is a concrete slower scenario worth carrying.

## Current scenarios

- `fresh-minimal`: installs from the current staged source as a non-root user with minimal toolsets and Lazy bootstrap disabled.
- `unmanaged-nvim-config`: verifies install fails closed instead of overwriting an existing user Neovim config.
- `unmanaged-tmux-config`: verifies install fails closed instead of overwriting an existing user tmux config.

## Contract under test

The scenario scripts assert that Devdeck:

- installs the CLI from the current staged source tree
- links tmux and Neovim config to Devdeck-managed source
- keeps TPM under `~/.local/share/tmux/plugins/tpm`
- preserves unmanaged user config by failing closed

Each scenario owns a fresh `HOME` and pins all XDG paths. Keep this harness boring, small, and outcome-based.
