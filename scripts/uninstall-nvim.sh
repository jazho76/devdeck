#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
NVIM_CONFIG_SOURCE="$REPO_DIR/nvim"
NVIM_CONFIG_TARGET="$HOME/.config/nvim"
NVIM_INSTALL_DIR="$HOME/.local/devdeck/nvim"
NVIM_BIN_TARGET="$HOME/.local/bin/nvim"

remove_if_symlink_to() {
  local link_path="$1"
  local expected_target="$2"

  if [ ! -L "$link_path" ]; then
    if [ -e "$link_path" ]; then
      echo "Keeping unmanaged path: $link_path"
    fi
    return
  fi

  actual_target="$(readlink -f "$link_path")"
  expected_target="$(readlink -f "$expected_target")"

  if [ "$actual_target" = "$expected_target" ]; then
    rm "$link_path"
    echo "Removed symlink: $link_path"
  else
    echo "Keeping unmanaged symlink: $link_path -> $actual_target"
  fi
}

remove_if_symlink_to "$NVIM_CONFIG_TARGET" "$NVIM_CONFIG_SOURCE"
remove_if_symlink_to "$NVIM_BIN_TARGET" "$NVIM_INSTALL_DIR/bin/nvim"

if [ -d "$NVIM_INSTALL_DIR" ]; then
  rm -rf "$NVIM_INSTALL_DIR"
  echo "Removed Devdeck Neovim: $NVIM_INSTALL_DIR"
else
  echo "No Devdeck Neovim install found: $NVIM_INSTALL_DIR"
fi

echo "Kept Neovim runtime state under ~/.local/share/nvim, ~/.local/state/nvim, and ~/.cache/nvim."
