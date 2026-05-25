# Neovim Configuration

Devdeck Neovim configuration.

Install it from the repository root:

```bash
./scripts/install-nvim.sh
```

The installer links `nvim/` into `~/.config/nvim`. If no `nvim` binary is available, it installs the Neovim tarball under `~/.local/devdeck/nvim` and links `~/.local/bin/nvim` to it.

Icon glyphs require a Nerd Font. From the repository root:

```bash
./scripts/install-font.sh
```

This installs DroidSansM Nerd Font into `~/.local/share/fonts/devdeck` and refreshes the font cache. Select `DroidSansM Nerd Font` in your terminal profile.

If another Neovim is already installed, the installer leaves it untouched and prints a warning when the version differs from the expected version.

Uninstall Devdeck-managed Neovim files:

```bash
./scripts/uninstall-nvim.sh
```

The uninstaller removes only the Devdeck-managed Neovim tarball install and symlinks. It does not remove system Neovim packages or Neovim runtime state.
