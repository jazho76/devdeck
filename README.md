# Devdeck

nvim+tmux setup.

## Prerequisites

Install the core Fedora packages first:

```bash
sudo dnf install -y \
  git \
  curl \
  tar \
  gzip \
  unzip \
  tmux \
  fzf \
  gcc \
  gcc-c++ \
  make \
  ripgrep \
  fd-find \
  wl-clipboard \
  fontconfig \
  tree-sitter-cli
```

Toolset packages are optional. Install the packages for the toolsets you enable:

| Toolset | Fedora packages |
| --- | --- |
| `angular` | `nodejs npm` |
| `asm` | `cargo rust openssl-devel pkgconf-pkg-config` |
| `c_cpp` | `gcc gcc-c++ make` |
| `config` | `nodejs npm` |
| `csharp` | `dotnet-sdk-10.0` |
| `eslint` | `nodejs npm` |
| `go` | `go` |
| `graphql` | `nodejs npm` |
| `javascript` | `nodejs npm` |
| `javascript-debug` | `nodejs npm` |
| `markdown` | `nodejs npm` |
| `python` | `python3 python3-pip` |
| `react` | `nodejs npm` |
| `rust` | `cargo rust` |
| `typescript` | `nodejs npm` |
| `web-markup` | `nodejs npm` |

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
