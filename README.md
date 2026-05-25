# Devdeck

nvim+tmux setup.

## Prerequisites

```bash
sudo dnf install -y \
  git \
  curl \
  tar \
  gzip \
  unzip \
  gcc \
  gcc-c++ \
  make \
  ripgrep \
  fd-find \
  wl-clipboard \
  fontconfig \
  tree-sitter-cli \
  nodejs \
  npm \
  python3 \
  python3-pip \
  cargo \
  rust \
  openssl-devel \
  pkgconf-pkg-config \
  go \
  dotnet-sdk-10.0
```

## Install

```bash
git clone https://github.com/jazho76/devdeck.git ~/.devdeck
cd ~/.devdeck
./scripts/install-nvim.sh
./scripts/install-tmux.sh
```

`install-nvim.sh` installs the Neovim tarball only when no `nvim` binary is available. If another Neovim is already installed, it leaves it untouched and warns when the version differs from the expected version.

`install-tmux.sh` requires tmux to be installed by your system package manager. It never installs or replaces the tmux binary.

## Uninstall

```bash
./scripts/uninstall-nvim.sh
./scripts/uninstall-tmux.sh
```

Uninstall scripts remove only Devdeck-managed symlinks, Devdeck's Neovim tarball install, and tmux plugin files. They do not remove system Neovim/tmux packages or Neovim runtime state.
