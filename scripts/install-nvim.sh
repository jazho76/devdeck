#!/usr/bin/env bash
set -euo pipefail

NVIM_VERSION="0.12.2"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
NVIM_CONFIG_SOURCE="$REPO_DIR/nvim"
NVIM_CONFIG_TARGET="$HOME/.config/nvim"
DEVDECK_DIR="$HOME/.local/devdeck"
NVIM_INSTALL_DIR="$DEVDECK_DIR/nvim"
NVIM_BIN_TARGET="$HOME/.local/bin/nvim"
WORK_DIR=""

cleanup() {
  if [ -n "$WORK_DIR" ]; then
    rm -rf "$WORK_DIR"
  fi
}

trap cleanup EXIT

require_command() {
  local command_name="$1"

  if ! command -v "$command_name" >/dev/null 2>&1; then
    echo "Missing required command: $command_name"
    exit 1
  fi
}

get_nvim_version() {
  local nvim_bin="$1"
  "$nvim_bin" --version | head -n 1 | sed -E 's/^NVIM v//; s/ .*//'
}

link_nvim_config() {
  if [ -e "$NVIM_CONFIG_TARGET" ] && [ ! -L "$NVIM_CONFIG_TARGET" ]; then
    echo "Refusing to overwrite existing path: $NVIM_CONFIG_TARGET"
    echo "Move it aside and run this script again."
    exit 1
  fi

  ln -sfn "$NVIM_CONFIG_SOURCE" "$NVIM_CONFIG_TARGET"
  echo "Linked Neovim config: $NVIM_CONFIG_TARGET -> $NVIM_CONFIG_SOURCE"
}

install_nvim_tarball() {
  local arch archive url

  case "$(uname -m)" in
    x86_64)        arch="x86_64" ;;
    aarch64|arm64) arch="arm64" ;;
    *)
      echo "Unsupported architecture: $(uname -m)" >&2
      echo "Devdeck installs Neovim tarballs for x86_64 and arm64 only." >&2
      exit 1
      ;;
  esac

  archive="nvim-linux-${arch}.tar.gz"
  url="https://github.com/neovim/neovim/releases/download/v${NVIM_VERSION}/${archive}"

  WORK_DIR="$(mktemp -d)"

  curl -L "$url" -o "$WORK_DIR/$archive"
  tar xzf "$WORK_DIR/$archive" -C "$WORK_DIR"

  rm -rf "$NVIM_INSTALL_DIR"
  mkdir -p "$DEVDECK_DIR" "$HOME/.local/bin"
  mv "$WORK_DIR/${archive%.tar.gz}" "$NVIM_INSTALL_DIR"
  ln -sfn "$NVIM_INSTALL_DIR/bin/nvim" "$NVIM_BIN_TARGET"
}

require_command curl
require_command tar

if [ ! -f "$NVIM_CONFIG_SOURCE/init.lua" ]; then
  echo "Neovim config not found: $NVIM_CONFIG_SOURCE"
  exit 1
fi

mkdir -p "$HOME/.config" "$HOME/.local/bin"

echo "Installing Neovim config"
echo "Config source: ${NVIM_CONFIG_SOURCE}"

if [ -x "$NVIM_INSTALL_DIR/bin/nvim" ]; then
  installed_version="$(get_nvim_version "$NVIM_INSTALL_DIR/bin/nvim")"
  if [ "$installed_version" = "$NVIM_VERSION" ]; then
    echo "Devdeck Neovim already installed: $NVIM_INSTALL_DIR/bin/nvim"
  else
    echo "Warning: Devdeck Neovim version is $installed_version, expected $NVIM_VERSION."
    echo "Leaving existing Devdeck install untouched. Run scripts/uninstall-nvim.sh first to replace it."
  fi
elif command -v nvim >/dev/null 2>&1; then
  nvim_path="$(command -v nvim)"
  installed_version="$(get_nvim_version "$nvim_path")"
  if [ "$installed_version" = "$NVIM_VERSION" ]; then
    echo "Using existing Neovim: $nvim_path"
  else
    echo "Warning: existing Neovim at $nvim_path is version $installed_version, expected $NVIM_VERSION."
    echo "Leaving it untouched. Install Neovim $NVIM_VERSION manually or remove it from PATH to let Devdeck install its tarball."
  fi
else
  echo "Installing Neovim ${NVIM_VERSION} to $NVIM_INSTALL_DIR"
  install_nvim_tarball
fi

link_nvim_config

if command -v nvim >/dev/null 2>&1; then
  nvim --headless "+Lazy! install" "+qa"
elif [ -x "$NVIM_BIN_TARGET" ]; then
  "$NVIM_BIN_TARGET" --headless "+Lazy! install" "+qa"
else
  echo "Warning: nvim not found on PATH. Add ~/.local/bin to PATH, then run:"
  echo '  nvim --headless "+Lazy! install" "+qa"'
fi

echo
echo 'Make sure ~/.local/bin is on your PATH:'
echo 'export PATH="$HOME/.local/bin:$PATH"'
echo
echo "Done. Start Neovim with: nvim"
