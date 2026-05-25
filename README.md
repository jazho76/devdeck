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
./scripts/install.sh
```

## Uninstall

```bash
./scripts/uninstall.sh
```
